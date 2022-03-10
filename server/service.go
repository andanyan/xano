package server

import (
	"xano/deal"
	"xano/router"
	"xano/session"
)

type ServerServer struct{}

// 心跳返回
func (s *ServerServer) ServerHeart(ss *session.Session, input *deal.Pong) error {
	return nil
}

// 节点推送
func (m *ServerServer) ServerStart(ss *session.Session, input *deal.ServerStartResponse) error {
	router.GetLocalNode().SetNode(input.Nodes)
	return nil
}