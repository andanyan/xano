package node

import (
	"log"
	"xlq-server/common"
	"xlq-server/core"
	"xlq-server/session"
)

type Gate struct{}

func NewGate() *Gate {
	return new(Gate)
}

// tcp gate接口
func (g *Gate) RunTcp() {
	if common.GetGateConfig().TcpAddr == "" {
		return
	}
	core.NewTcpServer(common.GetGateConfig().TcpAddr, g.TcpHandle)
}

// 接收到包
func (g *Gate) TcpHandle(h *core.TcpHandle, p *common.TcpPacket) {
	// 获取Session, 包含用户的数据
	s := session.NewGateSession(h)
	err := s.Handle(p)
	if err != nil {
		log.Println(err)
	}
}

//
func (g *Gate) RunHttp() {

}
