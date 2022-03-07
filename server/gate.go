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
	TcpClient *core.TcpClient
}

func NewGate() *Gate {
	return new(Gate)
}

// 开始运行
func (g *Gate) Run() {
	if common.GetConfig().Server.MasterAddr == "" {
		return
	}
	t, err := core.NewTcpClient(common.GetConfig().Server.MasterAddr)
	if err != nil {
		logger.Fatal(err)
		return
	}

	// 设置路由
	router.GetGateRouter().Register(&router.RouterServer{
		Name:   "",
		Server: new(ServerGateServer),
	})

	// 设置回包函数
	t.SetHandle(func(h *core.TcpHandle, m *deal.Msg) {
		ss := core.GetSession(h)
		router := router.GetGateRouter()
		if err := ss.HandleRoute(router, m); err != nil {
			logger.Error(err.Error())
		}
	})
	t.SetCloseFunc(func() {
		logger.Fatal("Master Disconnected!")
	})
	g.TcpClient = t

	// 同步包
	g.serverStart()

	// 启动心跳
	for {
		time.Sleep(common.TcpHeartDuration)
		g.serverHeart()
	}
}

// 停止运行
func (g *Gate) Close() {
	g.serverClose()
}

// 启动
func (g *Gate) serverStart() {
	routes := router.GetLocalRouter().GetDescs()
	input := &deal.ServerStartNotice{
		Version: common.GetConfig().Base.Version,
		Port:    common.GetConfig().Server.Port,
		Routes:  routes,
	}
	inputBys, err := common.MsgMarsh(common.TcpDealProtobuf, input)
	if err != nil {
		logger.Error(err)
		return
	}
	msg := &deal.Msg{
		Route:   "ServerStart",
		Mid:     g.TcpClient.GetMid(),
		MsgType: common.MsgTypeNotice,
		Deal:    common.TcpDealProtobuf,
		Data:    inputBys,
		Version: common.GetConfig().Base.Version,
	}
	g.TcpClient.Send(msg)
}

// 心跳
func (g *Gate) serverHeart() {
	input := &deal.Ping{}
	inputBys, err := common.MsgMarsh(common.TcpDealProtobuf, input)
	if err != nil {
		logger.Error(err)
		return
	}
	msg := &deal.Msg{
		Route:   "ServerHeart",
		Mid:     g.TcpClient.GetMid(),
		MsgType: common.MsgTypeRequest,
		Deal:    common.TcpDealProtobuf,
		Data:    inputBys,
		Version: common.GetConfig().Base.Version,
	}
	g.TcpClient.Send(msg)
}

// 关闭
func (g *Gate) serverClose() {
	if g.TcpClient == nil {
		return
	}
	// 发送服务断开包
	input := &deal.ServerStopNotice{}
	inputBys, err := common.MsgMarsh(common.TcpDealProtobuf, input)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	msg := &deal.Msg{
		Route:   "ServerStop",
		Mid:     g.TcpClient.GetMid(),
		MsgType: common.MsgTypeNotice,
		Deal:    common.TcpDealProtobuf,
		Data:    inputBys,
	}
	g.TcpClient.Send(msg)
}
