package common

import "time"

// tcp协议类型
const (
	TcpDealProtobuf uint32 = iota + 1
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
	TcpHeartDuration time.Duration = 60
	// TCP请求超时市场
	TcpDeadDuration time.Duration = 60
)

// status
const (
	StatusZero int = iota
	StatusOne
	StatusTwo
)

// msg type
const (
	MsgTypeNone uint32 = iota
	MsgTypeRequest
	MsgTypeResponse
	MsgTypeNotice
	MsgTypePush
	MsgTypeRpc
)

// handle key
const (
	HandleKeyMid     string = "Mid"
	HandleKeyTcpAddr string = "TcpAddr"
)
