package component

import (
	"reflect"
	"xlq-server/common"
	"xlq-server/log"
)

// 服务
type Service struct {
	Routes map[string]*common.Route
}

var service Service

// 组件注册
func Register(module interface{}) {
	mfv := reflect.Indirect(reflect.ValueOf(module))
	mft := mfv.Type()

	// 模块名
	mName := mft.Name()

	// 获取全部函数
	methodNum := mft.NumMethod()

	for i := 0; i < methodNum; i++ {
		method := mft.Method(i)

		routeName := mName + "_" + method.Name

		methodFunc := method.Func
		methodFuncType := methodFunc.Type()

		if methodFuncType.NumIn() != 2 {
			log.Fatal("func error input " + routeName)
		}
		if methodFuncType.NumOut() != 2 {
			log.Fatal("func error output " + routeName)
		}

		route := &common.Route{
			Method: method,
			Input:  methodFuncType.In(1),
			OutPut: methodFuncType.Out(1),
		}
		service.Routes[routeName] = route
	}

}
