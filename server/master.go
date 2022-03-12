package server

import (
	"time"
	"xano/common"
	"xano/core"
	"xano/deal"
	"xano/logger"
	"xano/router"
	"xano/session"
)

type ServerMaster struct {
	MasterClient *core.TcpClient
}

var serverMaster *ServerMaster

func GetServerMaster() *ServerMaster {
	if serverMaster == nil {
		serverMaster = new(ServerMaster)
	}
	return serverMaster
}

// 与主节点进行通信
func (s *ServerMaster) masterHandle() {
	if common.GetConfig().Server.MasterAddr == "" {
		return
	}
	t, err := core.NewTcpClient(common.GetConfig().Server.MasterAddr)
	if err != nil {
		logger.Error("MASTER DISCONNECT:", err)
		s.MasterClient = nil
		time.Sleep(common.DelayDuration)
		s.masterHandle()
		return
	}

	// 设置回包函数
	t.SetHandleFunc(func(h *core.TcpHandle, m *deal.Msg) {
		ss := session.GetSession(h)
		ss.SID = m.Sid
		router := router.GetGateRouter()
		if err := ss.HandleRoute(router, m); err != nil {
			logger.Error(err.Error())
		}
	})
	t.SetCloseFunc(func(h *core.TcpHandle) {
		logger.Error("MASTER DISCONNECT:", err)
		s.MasterClient = nil
		time.Sleep(common.DelayDuration)
		s.masterHandle()
		return
	})
	s.MasterClient = t

	// 同步包
	s.serverStart()

	// 启动心跳
	for {
		s.serverHeart()
		time.Sleep(common.TcpHeartDuration)
	}
}

// 启动
func (s *ServerMaster) serverStart() {
	if s.MasterClient == nil {
		return
	}
	serverAddr, err := common.ParseAddr(common.GetConfig().Server.TcpAddr)
	if err != nil {
		logger.Fatal(err)
	}
	routes := router.GetLocalRouter().GetDescs()
	input := &deal.ServerStartRequest{
		Version: common.GetConfig().Base.Version,
		Port:    serverAddr.Port,
		Routes:  routes,
	}
	inputBys, err := common.MsgMarsh(common.GetConfig().Base.TcpDeal, input)
	if err != nil {
		logger.Error(err)
		return
	}
	msg := &deal.Msg{
		Route:   "ServerStart",
		Sid:     0,
		Mid:     s.MasterClient.GetMid(),
		MsgType: common.MsgTypeRequest,
		Deal:    common.GetConfig().Base.TcpDeal,
		Data:    inputBys,
		Version: common.GetConfig().Base.Version,
	}
	s.MasterClient.Send(msg)
}

// 心跳
func (s *ServerMaster) serverHeart() {
	if s.MasterClient == nil {
		return
	}
	input := &deal.Ping{
		Psutil: common.GetPsutil(),
	}
	inputBys, err := common.MsgMarsh(common.GetConfig().Base.TcpDeal, input)
	if err != nil {
		logger.Error(err)
		return
	}
	msg := &deal.Msg{
		Route:   "ServerHeart",
		Sid:     0,
		Mid:     s.MasterClient.GetMid(),
		MsgType: common.MsgTypeRequest,
		Deal:    common.GetConfig().Base.TcpDeal,
		Data:    inputBys,
		Version: common.GetConfig().Base.Version,
	}
	s.MasterClient.Send(msg)
}

// 关闭
func (s *ServerMaster) serverClose() {
	if s.MasterClient == nil {
		return
	}
	// 发送服务断开包
	input := &deal.ServerStopNotice{}
	inputBys, err := common.MsgMarsh(common.GetConfig().Base.TcpDeal, input)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	msg := &deal.Msg{
		Route:   "ServerStop",
		Sid:     0,
		Mid:     s.MasterClient.GetMid(),
		MsgType: common.MsgTypeNotice,
		Deal:    common.GetConfig().Base.TcpDeal,
		Data:    inputBys,
	}
	s.MasterClient.Send(msg)
}
