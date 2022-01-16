package core

import "reflect"

// 所有服务处理器
type Service struct {
	// route => ServiceComponent
	Components map[string]*ServiceComponent
}

// 服务组件
type ServiceComponent struct {
	Method reflect.Method
	Input  reflect.Type
	Ouput  reflect.Type
}

var service *Service

// 获取服务
func GetComponent(route string) *ServiceComponent {
	if service == nil {
		return nil
	}
	return service.Components[route]
}

// 注册服务
