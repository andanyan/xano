package server

import (
	"fmt"
	"log"
	"time"
	"xlq-server/common"
	"xlq-server/core"
	"xlq-server/deal"
)

type Session struct {
	*core.TcpHandle
}

func NewSession(entity *core.TcpHandle) *Session {
	return &Session{
		TcpHandle: entity,
	}
}

// RPC
func (s *Session) Rpc(route string, input, output interface{}) error {
	if input == nil || output == nil {
		return fmt.Errorf("error input or output, not allow nil")
	}

	inputPacket, err := s.genPacket(route, common.MsgTypeRpc, input)
	if err != nil {
		return err
	}

	// 向网关主节点发送Rpc请求
	tcpAddr := s.Get(common.HandleKeyTcpAddr).(string)

	cli, err := core.NewTcpClient(tcpAddr)
	if err != nil {
		return err
	}
	defer cli.Close()

	cli.Send(inputPacket)

	// 收到Response包才认为已完成、其他包直接发射回去即可
	c := make(chan struct{})
	cli.SetHandle(func(h *core.TcpHandle, p *common.Packet) {
		rmsg := new(deal.Msg)
		err := common.MsgUnMarsh(common.TcpDealProtobuf, p.Data, rmsg)
		if err != nil {
			log.Println(err)
			return
		}
		mmid := s.Get(common.HandleKeyMid).(uint64)
		if rmsg.Mid == mmid {
			// 非Response的类型，直接返回给客户端
			if rmsg.MsgType != common.MsgTypeResponse {
				s.Send(p)
				return
			}
			err := common.MsgUnMarsh(common.TcpDealProtobuf, rmsg.Data, output)
			if err != nil {
				log.Println(err)
			}
			c <- struct{}{}

		} else if rmsg.Mid > mmid {
			c <- struct{}{}
		}
	})
	defer cli.SetHandle(nil)

	t := time.NewTimer(common.TcpDeadDuration)

	select {
	case <-c:

	case <-t.C:
		return fmt.Errorf("rpc timeout")
	}

	return nil
}

// Response
func (s *Session) Response(route string, input interface{}) error {
	inputPacket, err := s.genPacket(route, common.MsgTypeResponse, input)
	if err != nil {
		return err
	}
	s.Send(inputPacket)
	return nil
}

// Push
func (s *Session) Push(route string, input interface{}) error {
	// 组装包 写入连接即可
	inputPacket, err := s.genPacket(route, common.MsgTypePush, input)
	if err != nil {
		return err
	}
	s.Send(inputPacket)
	return nil
}

// 向连接写入包
func (s *Session) genPacket(route string, msgType uint32, input interface{}) (*common.Packet, error) {
	// 组装包 写入连接即可
	inputBys, err := common.MsgMarsh(common.TcpDealProtobuf, input)
	if err != nil {
		return nil, err
	}
	msg := &deal.Msg{
		Route:   route,
		Mid:     s.Get(common.HandleKeyMid).(uint64),
		MsgType: msgType,
		Deal:    common.TcpDealProtobuf,
		Data:    inputBys,
	}
	msgBys, err := common.MsgMarsh(common.TcpDealProtobuf, msg)
	if err != nil {
		return nil, err
	}
	inputPacket := &common.Packet{
		Length: uint16(len(msgBys)),
		Data:   msgBys,
	}

	return inputPacket, nil
}
