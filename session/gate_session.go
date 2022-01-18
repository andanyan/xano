package session

import (
	"fmt"
	"time"
	"xlq-server/common"
	"xlq-server/core"
	"xlq-server/deal"
)

type GateSession struct {
	BaseSession
}

func NewGateSession(entity *core.TcpHandle) *GateSession {
	s := &GateSession{
		BaseSession{
			entity,
		},
	}
	return s
}

func (s *GateSession) HandleTcp(packet *common.TcpPacket) error {
	var err error

	/*msg := new(deal.Msg)
	err = common.TcpMsgUnMarsh(packet.Data, msg)
	if err != nil {
		return err
	}*/

	// 转发逻辑
	addr := common.GetGateConfig().TcpAddr
	cli, err := core.NewTcpClient(addr)
	if err != nil {
		return err
	}
	cli.SetHandle(func(h *core.TcpHandle, p *common.TcpPacket) {
		// 收到回包
		s.Send(p)
	})
	cli.Send(packet)
	return nil
}

// http请求
func (s *GateSession) HandleHttp(route string, body []byte) ([]byte, error) {
	// 暂时使用本地
	addr := common.GetGateConfig().TcpAddr
	cli, err := core.NewTcpClient(addr)
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	// 封装成packet
	msg := &deal.Msg{
		Route: route,
		Mid:   cli.TcpHandle.GetMid(),
		Deal:  uint32(common.GetGateConfig().HttpDeal),
		Data:  body,
	}
	msgBys, err := common.TcpMsgMarsh(msg)
	if err != nil {
		return nil, err
	}
	inputPacket := &common.TcpPacket{
		Length: uint16(len(msgBys)),
		Data:   msgBys,
	}
	cli.TcpHandle.Send(inputPacket)

	// 指定返回值
	var res []byte
	// 要求回包
	c := make(chan struct{})
	// 设置收包
	cli.TcpHandle.SetHandle(func(h *core.TcpHandle, p *common.TcpPacket) {
		// 只要第一个包
		cli.TcpHandle.SetHandle(nil)

		// 读取到回包
		msg := new(deal.Msg)
		err = common.TcpMsgUnMarsh(p.Data, msg)
		if err != nil {
			return
		}
		res = msg.Data
		// 通知回包成功
		c <- struct{}{}
	})

	// 读取一个包
	t := time.NewTimer(common.TcpDeadDuration * time.Second)
	select {
	case <-c:
		// 处理完了

	case <-t.C:
		// 超时了 连接池写好后直接回收连接
		return nil, fmt.Errorf("request limit")
	}

	return res, nil
}
