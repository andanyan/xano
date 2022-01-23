package router

import (
	"log"
	"reflect"
)

// 本地服务器 服务记录
type LocalRoute struct {
	// route -> func
	Data map[string]*LocalRouteItem
	// []route
	Desc []string
}

type LocalRouteItem struct {
	Method reflect.Value
	Input  reflect.Type
}

var localRoute *LocalRoute

// 服务注册
func Register(v interface{}) {
	if localRoute == nil {
		localRoute = new(LocalRoute)
		localRoute.Data = make(map[string]*LocalRouteItem)
	}

	var vName, mName string

	defer func() {
		if err := recover(); err != nil {
			log.Printf("error type func %s.%s: %s\n", vName, mName, err)
		}
	}()

	fv := reflect.ValueOf(v)
	ft := fv.Type()

	vName = ft.Elem().Name()

	numMethod := fv.NumMethod()

	for i := 0; i < numMethod; i++ {
		method := fv.Method(i)

		mName = ft.Method(i).Name

		rName := vName + mName

		if _, ok := localRoute.Data[rName]; ok {
			panic("duplicate route names")
		}

		localRoute.Data[rName] = &LocalRouteItem{
			Method: method,
			Input:  method.Type().In(1),
		}
		localRoute.Desc = append(localRoute.Desc, rName)
		log.Panicln("route:", rName)
	}
}

// 获取路由
func GetLocalRoute(route string) *LocalRouteItem {
	return localRoute.Data[route]
}

// 获取全部路由地址
func GetLocalRoutes() []string {
	return localRoute.Desc
}
