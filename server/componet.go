package server

import "xlq-server/component"

// 注册组件
func RegisterService(module interface{}) {
	component.Register(module)
}
