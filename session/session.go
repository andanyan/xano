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
	Del(k string)
	GetInt(k string) int
	GetUInt(k string) uint
	GetInt32(k string) int32
	GetUInt32(k string) uint32
	GetInt64(k string) int64
	GetUInt64(k string) uint64
	GetFloat(k string) float64
	GetBool(k string) bool
	GetString(k string) string

	// session id
	SetSid(sid uint64)
	GetSid() uint64

	// 远程调用
	Rpc(route string, input, output interface{}) error
	Notice(route string, input interface{}) error
	Response(route string, input interface{}) error
	RpcResponse(route string, input interface{}) error
	Push(route string, input interface{}) error
	PushTo(sid uint64, route string, input interface{}) error
	SendTo(addr string, msg *deal.Msg) error
	HandleRoute(r *router.Router, m *deal.Msg) error
}
