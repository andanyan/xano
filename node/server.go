package node

import (
	"log"
	"time"
	"xlq-server/common"
	"xlq-server/core"
	"xlq-server/deal"
	"xlq-server/session"
)

type Server struct{}

func NewServer() *Server {
	return new(Server)
}

func (s *Server) Run() {
	if common.GetServiceConfig().Addr == "" {
		return
	}
	log.Printf("Gate Service Address: %s \n", common.GetGateConfig().HttpAddr)
	core.NewTcpServer(common.GetGateConfig().TcpAddr, s.Handle)
}

func (s *Server) Handle(h *core.TcpHandle, p *common.TcpPacket) {
	ss := session.NewServerSession(h)
	err := ss.Handle(p)
	if err != nil {
		log.Println(err)
	}
}

// 加入到节点中 即注册
func (s *Server) AddGate() {
	if common.GetGateConfig().GateAddr == "" {
		return
	}
	addr, err := common.ParseIpAddr(common.GetGateConfig().GateAddr)
	if err != nil {
		log.Fatalln(err.Error())
	}
	req := &deal.GateRouteRequest{
		Port:   addr.Port,
		Routes: core.GetAllRoute(),
	}
	reqBys, err := common.TcpMsgMarsh(req)
	if err != nil {
		log.Fatalln(err.Error())
	}
	// 直接发包即可
	reqPacket := &common.TcpPacket{
		Length: uint16(len(reqBys)),
		Data:   reqBys,
	}

	cli, err := core.NewTcpClient(common.GetGateConfig().GateAddr)
	if err != nil {
		log.Fatalln(err.Error())
	}
	// 服务注册和回包
	cli.SetHandle(func(h *core.TcpHandle, p *common.TcpPacket) {
		// 解析包
		res := new(deal.GateRouteResponse)
		err := common.TcpMsgUnMarsh(p.Data, res)
		if err != nil {
			log.Panicln(err.Error())
			return
		}
		core.SetRoutes(res)
	})

	for {
		cli.Send(reqPacket)
		time.Sleep(common.TcpHeartDuration)
	}
}
