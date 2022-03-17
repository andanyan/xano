package server

import (
	"xano/common"
	"xano/core"
	"xano/deal"
	"xano/logger"
	"xano/router"
	"xano/session"
)

// 微服务对象,可独立运行tcp服务
type Server struct{}

func NewServer() *Server {
	return new(Server)
}

func (s *Server) Close() {
	GetServerMaster().serverClose()
}

func (s *Server) Run() {
	// 设置路由
	router.GetGateRouter().Register(&router.RouterServer{
		Name:   "",
		Server: new(ServerServer),
	})

	// 与主节点进行通信
	go GetServerMaster().masterHandle()

	// 启动tcp
	go s.runTcp()
}

func (s *Server) runTcp() {
	addr := common.GetConfig().Server.TcpAddr
	if addr == "" {
		return
	}
	// 启动服务
	logger.Infof("Server Start: %s", addr)
	core.NewTcpServer(addr, func(h *core.TcpHandle) {
		h.SetHandleFunc(s.tcpHandle)
	})
}

func (s *Server) tcpHandle(h *core.TcpHandle, msg *deal.Msg) {
	// 解析packet
	var err error

	switch msg.MsgType {
	case common.MsgTypeNotice, common.MsgTypeRequest:
		// 创建session 提供给接口端使用
		ss := session.GetServerSession(h)
		ss.SetSid(msg.Sid)
		// 调用路由
		if err = ss.HandleRoute(router.GetLocalRouter(), msg); err != nil {
			logger.Error(err.Error())
		}

	default:

	}
}
