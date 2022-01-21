package common

import "time"

// tcp协议类型
const (
	TcpDealProtobuf uint8 = iota + 1
	TcpDealJson
)

// packet常量
const (
	TcpPacketHeadLength int = 2
	TcpReadBufLen       int = 2048
)

// tcp请求信息
const (
	// 服务注册时差
	TcpHeartDuration time.Duration = 30
	// TCP请求超时市场
	TcpDeadDuration time.Duration = 30
)

// status
const (
	StatusZero int = iota
	StatusOne
	StatusTwo
)

// 网关相关配置
type GateConfig struct {
	// 本地网关启动地址
	TcpAddr string
	// 使用的协议类型 1-protobuf 2-json
	TcpDeal uint8

	// http服务地址
	HttpAddr string
	// http协议
	HttpDeal uint8

	// bending网关监听地址
	GateAddr string
}

// tcp服务层配置
type TcpServiceConfig struct {
	// 远程网关地址
	GateAddr string
	// 本地服务地址
	Addr string
}

// tcp数据包结构
type TcpPacket struct {
	// 数据源长度
	Length uint16
	// 源数据
	Data []byte
}
