package server

import (
	"log"
	"xlq-server/common"
	"xlq-server/core"
)

// tcp服务
type Server struct{}

func NewServer() *Server {
	return new(Server)
}

func (s *Server) Run() {
	sConf := common.GetConfig().Server
	addr := common.GenAddr(sConf.Host, sConf.Port)
	if addr == "" {
		return
	}
	log.Printf("Gate Service Address: %s \n", addr)
	core.NewTcpServer(addr, s.handle)
}

func (s *Server) handle(h *core.TcpHandle, p *common.Packet) {
	ss := core.NewSession(h)

}
