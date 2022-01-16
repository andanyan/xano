package node

import (
	"log"
	"xlq-server/common"
	"xlq-server/core"
	"xlq-server/session"
)

type Server struct{}

func NewServer() *Server {
	return new(Server)
}

func (s *Server) Run() {
	if common.GetServiceConfig().LocalAddr == "" {
		return
	}
	core.NewTcpServer(common.GetGateConfig().TcpAddr, s.Handle)
}

func (s *Server) Handle(h *core.TcpHandle, p *common.TcpPacket) {
	ss := session.NewServerSession(h)
	err := ss.Handle(p)
	if err != nil {
		log.Println(err)
	}
}
