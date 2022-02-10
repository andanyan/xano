package server

import (
	"xlq-server/common"
	"xlq-server/core"
	"xlq-server/deal"
	"xlq-server/logger"
	"xlq-server/router"
)

// 微服务对象,可独立运行tcp服务
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
	logger.Infof("Server Start: %s", addr)
	core.NewTcpServer(addr, s.handle)
}

func (s *Server) handle(h *core.TcpHandle, msg *deal.Msg) {
	// 解析packet
	var err error

	switch msg.MsgType {
	case common.MsgTypeNotice, common.MsgTypeRequest, common.MsgTypeRpc:
		// 设定当前的来源地址
		h.Set(common.HandleKeyTcpAddr, h.GetAddr())

		// 创建session 提供给接口端使用
		ss := core.GetSession(h)
		// 调用路由
		if err = ss.HandleRoute(router.LocalRouter, msg); err != nil {
			logger.Error(err.Error())
		}

	case common.MsgTypePush:
		h.Send(msg)

	default:

	}
}
