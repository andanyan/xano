package common

import (
	"encoding/json"
	"fmt"

	"google.golang.org/protobuf/proto"
)

// tcp的协议加密
func TcpMsgMarsh(input interface{}) ([]byte, error) {
	return MsgMarsh(GetGateConfig().TcpDeal, input)
}

// tcp的协议解密
func TcpMsgUnMarsh(msg []byte, output interface{}) error {
	return MsgUnMarsh(GetGateConfig().TcpDeal, msg, output)
}

// 加密协议
func MsgMarsh(deal uint8, input interface{}) ([]byte, error) {
	// 默认protobuf
	switch deal {
	case TcpDealProtobuf:
		pb, ok := input.(proto.Message)
		if !ok {
			return nil, fmt.Errorf("error protobuf data")
		}
		return proto.Marshal(pb)
	case TcpDealJson:
		return json.Marshal(input)
	}
	return nil, fmt.Errorf("unsupported protocols")
}

// 解密协议
func MsgUnMarsh(deal uint8, msg []byte, output interface{}) error {
	switch deal {
	case TcpDealProtobuf:
		pb, ok := output.(proto.Message)
		if !ok {
			return fmt.Errorf("error protobuf type")
		}
		return proto.Unmarshal(msg, pb)
	case TcpDealJson:
		return json.Unmarshal(msg, output)
	}
	return fmt.Errorf("unsupported protocols")
}
