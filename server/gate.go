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

	// 同步路由
	routes := router.GetLocalRoutes()
	input := &deal.ServerRunNotice{
		Version: g.Conf.Version,
		Port:    g.Conf.Port,
		Routes:  routes,
	}
	inputBys, err := common.MsgMarsh(common.TcpDealProtobuf, input)
	if err != nil {
		log.Println(err)
		return
	}
	msg := &deal.Msg{
		Route:   "ServerRun",
		Mid:     g.GetMid(),
		MsgType: common.MsgTypeNotice,
		Deal:    common.TcpDealProtobuf,
		Data:    inputBys,
	}
	msgBys, err := common.MsgMarsh(common.TcpDealProtobuf, msg)
	if err != nil {
		log.Panicln(err)
		return
	}
	packet := &common.Packet{
		Length: uint16(len(msgBys)),
		Data:   msgBys,
	}
	g.Send(packet)

	// 心跳处理
	go g.heart()
}

// 停止运行
func (g *Gate) Close() {
	g.IsClose = true

	// 发送服务断开包
	input := &deal.ServerCloseNotice{
		Version: g.Conf.Version,
		Port:    g.Conf.Port,
	}
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
	msgBys, err := common.MsgMarsh(common.TcpDealProtobuf, msg)
	if err != nil {
		log.Println(err)
		return
	}
	packet := &common.Packet{
		Length: uint16(len(msgBys)),
		Data:   msgBys,
	}
	g.Send(packet)
}

// 心跳处理
func (g *Gate) heart() {
	input := &deal.ServerHeartRequest{}
	inputBys, err := common.MsgMarsh(common.TcpDealProtobuf, input)
	if err != nil {
		log.Println(err)
		return
	}
	msg := &deal.Msg{
		Route:   "ServerHeart",
		Mid:     g.GetMid(),
		MsgType: common.MsgTypeRequest,
		Deal:    common.TcpDealProtobuf,
		Data:    inputBys,
	}
	msgBys, err := common.MsgMarsh(common.TcpDealProtobuf, msg)
	if err != nil {
		log.Println(err)
		return
	}
	packet := &common.Packet{
		Length: uint16(len(msgBys)),
		Data:   msgBys,
	}

	// 持续发送心跳包
	for g.Status() && !g.IsClose {
		g.Send(packet)
		time.Sleep(common.TcpHeartDuration)
	}

	// 如果断开了,则不断尝试重启
	for !g.Status() && !g.IsClose {
		g.Run()
	}
}
