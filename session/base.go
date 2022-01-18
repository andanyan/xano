package session

import (
	"fmt"
	"log"
	"time"
	"xlq-server/common"
	"xlq-server/core"
	"xlq-server/deal"
)

type BaseSession struct {
	*core.TcpHandle

	// LastMsg *deal.Msg
}

func NewBaseSession(entity *core.TcpHandle) *BaseSession {
	s := new(BaseSession)
	s.TcpHandle = entity
	return s
}

func (s *BaseSession) Set(k string, v interface{}) {
	s.TcpHandle.Set(k, v)
}

func (s *BaseSession) Get(k string) interface{} {
	return s.TcpHandle.Get(k)
}

// rpc请求
func (s *BaseSession) Rpc(route string, input, output interface{}) error {
	// 先用本地地址
	addr := common.GetGateConfig().TcpAddr
	cli, err := core.NewTcpClient(addr)
	if err != nil {
		return err
	}
	defer cli.Close()

	// 封装成packet
	inputPacket, err := s.EncodePacket(route, cli.TcpHandle.GetMid(), common.TcpDealProtobuf, input)
	if err != nil {
		return err
	}
	cli.TcpHandle.Send(inputPacket)

	// 不需要回包的处理
	if output == nil {
		return nil
	}

	// 要求回包
	c := make(chan struct{})
	// 设置收包
	cli.TcpHandle.SetHandle(func(h *core.TcpHandle, p *common.TcpPacket) {
		// 只要第一个包
		cli.TcpHandle.SetHandle(nil)

		err := s.DecodePacket(p, output)
		if err != nil {
			log.Panicln(err.Error())
		}
		c <- struct{}{}
	})

	// 读取一个包
	t := time.NewTimer(common.TcpDeadDuration * time.Second)
	select {
	case <-c:
		// 处理完了

	case <-t.C:
		// 超时了 连接池写好后直接回收连接
		return fmt.Errorf("request limit")
	}

	return nil
}

// msg请求包生成
func (s *BaseSession) EncodePacket(route string, mid uint64, doneDeal uint8, input interface{}) (*common.TcpPacket, error) {
	inputBys, err := common.TcpMsgMarsh(input)
	if err != nil {
		return nil, err
	}

	msg := &deal.Msg{
		Route: route,
		Mid:   mid,
		Deal:  uint32(doneDeal),
		Data:  inputBys,
	}

	msgBys, err := common.TcpMsgMarsh(msg)
	if err != nil {
		return nil, err
	}
	inputPacket := &common.TcpPacket{
		Length: uint16(len(msgBys)),
		Data:   msgBys,
	}
	return inputPacket, nil
}

// msg解析一个包
func (s *BaseSession) DecodePacket(packet *common.TcpPacket, output interface{}) error {
	var err error
	msg := new(deal.Msg)
	err = common.TcpMsgUnMarsh(packet.Data, msg)
	if err != nil {
		return err
	}
	err = common.MsgUnMarsh(uint8(msg.Deal), msg.Data, output)
	if err != nil {
		return err
	}
	return nil
}
