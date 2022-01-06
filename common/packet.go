package common

import (
	"encoding/json"
	"fmt"

	"google.golang.org/protobuf/proto"
)

func (p *TcpPacket) Marsh() ([]byte, error) {
	switch p.Deal {
	case TcpDealProtobuf:
		pb, ok := p.(proto.Message)
		if !ok {
			return nil, fmt.Errorf("protobuf: convert on wrong type value")
		}
		return proto.Marshal(pb)
	case TcpDealJson:
		return json.Marshal(p)
	}
	return nil, fmt.Errorf("unsupported protocol type")
}
