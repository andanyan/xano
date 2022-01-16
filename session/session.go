package session

type Session interface {
	// session数据
	Set(k string, v interface{})
	Get(k string) interface{}

	// Rpc
	Rpc(route string, input, output interface{}) error

	// 发送包
	Write(input interface{}) error
}
