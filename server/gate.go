package server

import (
	"time"
	"xano/common"
	"xano/core"
	"xano/deal"
	"xano/logger"
	"xano/router"
)

// 与网关之间的通信
type Gate struct {
	*core.TcpClient
	// 配置加载到缓存
	Conf common.ServerConfig
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
		logger.Error(err.Error())
		return
	}
	// 不需要设置回包函数
	g.TcpClient = t

	// 开启心跳
	go g.Start()
}

// 发送同步路由包
func (g *Gate) Start() {
	for {
		routes := router.GetLocalRouter().GetDescs()
		input := &deal.ServerStartNotice{
			Version: common.GetConfig().Base.Version,
			Port:    g.Conf.Port,
			Routes:  routes,
		}
		inputBys, err := common.MsgMarsh(common.TcpDealProtobuf, input)
		if err != nil {
			logger.Error(err.Error())
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

		time.Sleep(common.TcpHeartDuration)
	}
}

// 停止运行
func (g *Gate) Close() {
	if g.Conf.GateAddr == "" {
		return
	}
	// 发送服务断开包
	input := &deal.ServerCloseNotice{}
	inputBys, err := common.MsgMarsh(common.TcpDealProtobuf, input)
	if err != nil {
		logger.Error(err.Error())
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
