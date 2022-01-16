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

	HeartPacketLength   uint16 = 4
	HeartPacketRequest  string = "Ping"
	HeartPacketResponse string = "Pong"
)

// tcp请求信息
const (
	// TCP请求超时市场
	TcpDeadDuration time.Duration = 30
)

// 网关相关配置
type GateConfig struct {
	// 本地网关启动地址
	TcpAddr string
	// 使用的协议类型
	TcpDeal uint8
}

// tcp服务层配置
type TcpServiceConfig struct {
	LocalAddr string
	TcpAddr   string
}

// tcp数据包结构
type TcpPacket struct {
	// 数据源长度
	Length uint16
	// 源数据
	Data []byte
}
