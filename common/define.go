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
	TcpHeartDuration time.Duration = time.Duration(60) * time.Second
	// TCP请求超时市场
	TcpDeadDuration time.Duration = time.Duration(60) * time.Second
	// HTTP请求超时时间
	HttpDeadDuration time.Duration = time.Duration(60) * time.Second
	// 连接池
	TcpPoolIdMin int = 0
	TcpPoolIdMax int = 10
	// 获取对象最长等待时长 单位毫秒
	TcpPoolMaxWaitTime int = 1000
	// 30分钟
	TcpPoolLifeTime int64 = 1800
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
	HandleKeySession string = "Session"
	HandleKeyMid     string = "Mid"
	HandleKeyTcpAddr string = "TcpAddr"
)
