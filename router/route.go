package router

import (
	"reflect"
	"xano/logger"
)

type Router struct {
	Name string
	// route -> func
	Data map[string]*RouterItem
	// []route
	Desc []string
}

type RouterItem struct {
	Method reflect.Value
	Input  reflect.Type
}

type RouterServer struct {
	Name   string
	Server interface{}
}

func NewRouter() *Router {
	r := new(Router)
	r.Data = make(map[string]*RouterItem)
	return r
}

// 服务对象, 服务模块名
func (r *Router) Register(obj *RouterServer) {
	var vName, mName string

	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("error type func %s.%s: %s\n", vName, mName, err)
		}
	}()

	if obj.Server == nil {
		return
	}

	fv := reflect.ValueOf(obj.Server)
	ft := fv.Type()

	vName = obj.Name

	numMethod := fv.NumMethod()

	for i := 0; i < numMethod; i++ {
		method := fv.Method(i)

		mName = ft.Method(i).Name

		rName := vName + mName

		if _, ok := r.Data[rName]; ok {
			panic("duplicate route names")
		}

		inLen := method.Type().NumIn()
		if inLen < 2 {
			return
		}

		r.Data[rName] = &RouterItem{
			Method: method,
			Input:  method.Type().In(1),
		}
		r.Desc = append(r.Desc, rName)
	}
}

func (r *Router) GetRoute(route string) *RouterItem {
	return r.Data[route]
}

func (r *Router) GetDescs() []string {
	return r.Desc
}
