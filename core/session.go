package core

import (
	"fmt"
	"log"
	"time"
)

type Session struct {
	*TcpHandle
}

func NewSession(entity *TcpHandle) *Session {
	return &Session{
		TcpHandle: entity,
	}
}

// RPC
func (s *Session) Rpc(route string, input, output interface{}) error {
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
	cli.TcpHandle.SetHandle(func(h *core.TcpHandle, p *common.Packet) {
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

// Response
func (s *Session) Response(route string, input interface{}) error {
	return nil
}

// Push
func (s *Session) Push(route string, input interface{}) error {
	return nil
}
