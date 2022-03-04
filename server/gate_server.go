package server

import (
	"xano/core"
	"xano/deal"
	"xano/router"
)

type ServerGateServer struct{}

// 心跳返回
func (s *ServerGateServer) ServerHeart(ss *core.Session, input *deal.Pong) error {
	return nil
}

// 节点推送
func (m *ServerGateServer) MemberNode(ss *core.Session, input *deal.MemberNodePush) error {
	router.GetMemberInfo().SetNode(input.Nodes)
	return nil
}
