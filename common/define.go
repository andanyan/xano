package common

import "time"

// tcp协议类型
const (
	TcpDealProtobuf uint8 = iota + 1
	TcpDealJson
)

// packet常量
const (
	PacketHeadLength int = 2
	TcpReadBufLen    int = 2048
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
