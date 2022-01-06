package common

// tcp协议类型
const (
	TcpDealProtobuf uint8 = iota + 1
	TcpDealJson
)

// 网关相关配置
type GateConfig struct {
	TcpAddr string

	TcpDeal uint8
}

// tcp数据包结构
type TcpPacket struct {
	// 数据源长度
	Length uint16
	// 源数据
	Data []byte
}
