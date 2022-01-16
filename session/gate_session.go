package session

import (
	"xlq-server/common"
	"xlq-server/core"
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

func (s *GateSession) Handle(packet *common.TcpPacket) error {
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
		s.Send(p)
	})

	cli.Send(packet)
	return nil
}
