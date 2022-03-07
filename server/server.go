package server

import (
	"time"
	"xano/common"
	"xano/core"
	"xano/deal"
	"xano/logger"
	"xano/router"
)

// 微服务对象,可独立运行tcp服务
type Server struct {
	TcpClient *core.TcpClient
}

func NewServer() *Server {
	return new(Server)
}

func (s *Server) Close() {
	s.serverClose()
}

func (s *Server) Run() {
	addr := common.GetConfig().Server.TcpAddr
	if addr == "" {
		return
	}

	// 与主节点进行通信
	go s.masterHandle()

	// 启动服务
	logger.Infof("Server Start: %s", addr)
	go core.NewTcpServer(addr, s.handle)

}

func (s *Server) handle(h *core.TcpHandle, msg *deal.Msg) {
	// 解析packet
	var err error

	switch msg.MsgType {
	case common.MsgTypeNotice, common.MsgTypeRequest, common.MsgTypeRpc:
		// 创建session 提供给接口端使用
		ss := core.GetSession(h)
		// 调用路由
		if err = ss.HandleRoute(router.GetLocalRouter(), msg); err != nil {
			logger.Error(err.Error())
		}

	case common.MsgTypePush:
		h.Send(msg)

	default:

	}
}

// 与主节点进行通信
func (s *Server) masterHandle() {
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
		Server: new(ServerServer),
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
	s.TcpClient = t

	// 同步包
	s.serverStart()

	// 启动心跳
	for {
		time.Sleep(common.TcpHeartDuration)
		s.serverHeart()
	}
}

// 启动
func (s *Server) serverStart() {
	serverAddr, err := common.ParseAddr(common.GetConfig().Server.TcpAddr)
	if err != nil {
		logger.Fatal(err)
	}
	routes := router.GetLocalRouter().GetDescs()
	input := &deal.ServerStartNotice{
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
		Mid:     s.TcpClient.GetMid(),
		MsgType: common.MsgTypeNotice,
		Deal:    common.GetConfig().Base.TcpDeal,
		Data:    inputBys,
		Version: common.GetConfig().Base.Version,
	}
	s.TcpClient.Send(msg)
}

// 心跳
func (s *Server) serverHeart() {
	input := &deal.Ping{}
	inputBys, err := common.MsgMarsh(common.GetConfig().Base.TcpDeal, input)
	if err != nil {
		logger.Error(err)
		return
	}
	msg := &deal.Msg{
		Route:   "ServerHeart",
		Mid:     s.TcpClient.GetMid(),
		MsgType: common.MsgTypeRequest,
		Deal:    common.GetConfig().Base.TcpDeal,
		Data:    inputBys,
		Version: common.GetConfig().Base.Version,
	}
	s.TcpClient.Send(msg)
}

// 关闭
func (s *Server) serverClose() {
	if s.TcpClient == nil {
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
		Mid:     s.TcpClient.GetMid(),
		MsgType: common.MsgTypeNotice,
		Deal:    common.GetConfig().Base.TcpDeal,
		Data:    inputBys,
	}
	s.TcpClient.Send(msg)
}
