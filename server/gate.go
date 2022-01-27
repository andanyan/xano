package server

import (
	"log"
	"time"
	"xlq-server/common"
	"xlq-server/core"
	"xlq-server/deal"
	"xlq-server/router"
)

// 与网关之间的通信
type Gate struct {
	*core.TcpClient
	IsClose bool
	Conf    common.ServerConfig
}

func NewGate() *Gate {
	sConf := common.GetConfig().Server
	g := new(Gate)
	g.Conf = sConf
	return g
}

// 开始运行
func (g *Gate) Run() {
	if g.Conf.GateAddr == "" {
		return
	}
	t, err := core.NewTcpClient(g.Conf.GateAddr)
	if err != nil {
		log.Println(err)
		return
	}
	// 不需要设置回包函数
	g.TcpClient = t

	g.Start()

	// 开启心跳
	go g.Heartbert()
}

// 发送同步路由包
func (g *Gate) Start() {
	// 同步路由
	routes := router.LocalRouter.GetDescs()
	input := &deal.ServerStartNotice{
		Version: common.GetConfig().Base.Version,
		Port:    g.Conf.Port,
		Routes:  routes,
	}
	inputBys, err := common.MsgMarsh(common.TcpDealProtobuf, input)
	if err != nil {
		log.Println(err)
		return
	}
	msg := &deal.Msg{
		Route:   "ServerStart",
		Mid:     g.GetMid(),
		MsgType: common.MsgTypeNotice,
		Deal:    common.TcpDealProtobuf,
		Data:    inputBys,
	}
	g.Send(msg)
}

// 停止运行
func (g *Gate) Close() {
	g.IsClose = true

	// 发送服务断开包
	input := &deal.ServerCloseNotice{}
	inputBys, err := common.MsgMarsh(common.TcpDealProtobuf, input)
	if err != nil {
		log.Println(err)
		return
	}
	msg := &deal.Msg{
		Route:   "ServerClose",
		Mid:     g.GetMid(),
		MsgType: common.MsgTypeNotice,
		Deal:    common.TcpDealProtobuf,
		Data:    inputBys,
	}
	g.Send(msg)
}

// 定时发送心跳包
func (g *Gate) Heartbert() {
	for {
		time.Sleep(common.TcpHeartDuration)

		ping := &deal.Ping{
			Ping: time.Now().Unix(),
		}
		pingBys, err := common.MsgMarsh(common.TcpDealProtobuf, ping)
		if err != nil {
			log.Println(err)
			continue
		}
		msg := &deal.Msg{
			Route:   "Heartbert",
			Mid:     g.GetMid(),
			MsgType: common.MsgTypeRequest,
			Deal:    common.TcpDealProtobuf,
			Data:    pingBys,
		}
		g.Send(msg)
	}
}
