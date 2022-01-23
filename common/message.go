package common

import (
	"encoding/json"
	"fmt"

	"google.golang.org/protobuf/proto"
)

// 加密协议
func MsgMarsh(deal uint32, input interface{}) ([]byte, error) {
	// 默认protobuf
	if deal < TcpDealProtobuf || deal > TcpDealJson {
		deal = TcpDealProtobuf
	}
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
func MsgUnMarsh(deal uint32, msg []byte, output interface{}) error {
	if deal < TcpDealProtobuf || deal > TcpDealJson {
		deal = TcpDealProtobuf
	}
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
