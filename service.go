package core

import (
	"log"
	"reflect"
	"sync"
	"xlq-server/deal"
)

// 所有服务处理器
type Service struct {
	sync.RWMutex
	// route => ServiceComponent
	Components map[string]*ServiceComponent
	// route => []addr
	RemoteRoutes map[string][]string
}

// 服务组件
type ServiceComponent struct {
	Func  reflect.Value
	Input reflect.Type
	Ouput reflect.Type
}

var service *Service

func init() {
	if service != nil {
		return
	}
	service = new(Service)
	service.Components = make(map[string]*ServiceComponent)
}

// 获取服务
func GetComponent(route string) *ServiceComponent {
	if service == nil {
		return nil
	}
	return service.Components[route]
}

// 注册服务
func RegisterComponent(comp interface{}) {
	fv := reflect.Indirect(reflect.ValueOf(service))
	ft := fv.Type()

	// 获取模块名
	fName := ft.Name()

	// 获取方法
	methodNum := ft.NumMethod()
	if methodNum == 0 {
		log.Panicln("comp has no method: " + fName)
		return
	}

	for i := 0; i < methodNum; i++ {
		method := ft.Method(i)

		mName := method.Name

		numIn := method.Type.NumIn()
		if numIn != 2 {
			log.Fatalln("error input param name", fName, mName)
		}

		numOut := method.Type.NumOut()
		if numOut != 2 {
			log.Fatalln("error output param name", fName, mName)
		}

		// 组件
		c := &ServiceComponent{
			Func:  method.Func,
			Input: method.Type.In(1),
			Ouput: method.Type.In(0),
		}

		// 路由名字转为
		r := fName + "_" + mName

		service.Components[r] = c
	}
}

// 获取全部的路由服务
func GetLocalRoute() []string {
	var res []string
	for r, _ := range service.Components {
		res = append(res, r)
	}
	return res
}

// 设置路由位置
func SetRoutes(res *deal.GateRouteResponse) {
	service.Lock()
	defer service.Unlock()

	// route = []addr
	service.RemoteRoutes = make(map[string][]string)
	for _, item := range res.Routes {
		for _, r := range item.Routes {
			addrs, ok := service.RemoteRoutes[r]
			if !ok {
				addrs = []string{}
			}
			service.RemoteRoutes[r] = append(addrs, item.Addr)
		}
	}
}
