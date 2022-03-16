package session

import (
	"xano/deal"
	"xano/router"
)

type Session interface {
	// 获取客户端地址
	GetAddr() string

	// 获取mid
	GetMid() uint64

	// cache
	Set(k string, v interface{})
	Get(k string) interface{}

	// session id
	SetSid(sid uint64)
	GetSid() uint64

	// 远程调用
	Rpc(route string, input, output interface{}) error
	Notice(route string, input interface{}) error
	Response(route string, input interface{}) error
	RpcResponse(route string, input interface{}) error
	Push(route string, input interface{}) error
	SendTo(addr string, msg *deal.Msg) error
	HandleRoute(r *router.Router, m *deal.Msg) error
}
