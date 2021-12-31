package common

type Session interface {
	GetValue(key string) interface{}
	SetValue(key string, val interface{})

	Rpc(route string, input interface{}, output interface{}) error

	Write(route string, data interface{})
	MsgMarsh(data interface{}) ([]byte, error)
	MsgUnMarsh(data []byte, v interface{}) error

	Middlewares(route string) error
}

const (
	ClientTypeTcp uint8 = iota + 1
	ClientTypeHttp
)
