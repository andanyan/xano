package common

import (
	"encoding/json"
	"fmt"
	"xano/deal"
	"xano/logger"

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

// 打印消息
var msgTypeDesc = []string{"request", "response", "notice", "push"}
var dealDesc = []string{"protobuf", "json"}

func PrintMsg(msg *deal.Msg, data ...interface{}) {
	if len(data) == 0 {
		logger.Printf("sid: %d, route: %s, mid: %d, msgType: %s, deal: %s, version: %s, dataLen: %d",
			msg.Sid, msg.Route, msg.Mid, msgTypeDesc[msg.MsgType-1], dealDesc[msg.Deal-1], msg.Version, len(msg.Data))
		return
	}
	logger.Printf("sid: %d, route: %s, mid: %d, msgType: %s, deal: %s, version: %s, dataLen: %d, data: %+v",
		msg.Sid, msg.Route, msg.Mid, msgTypeDesc[msg.MsgType-1], dealDesc[msg.Deal-1], msg.Version, len(msg.Data), data)
}
