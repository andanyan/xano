package server

import (
	"xano/core"
	"xano/deal"
	"xano/router"
)

type ServerServer struct{}

// 心跳返回
func (s *ServerServer) ServerHeart(ss *core.Session, input *deal.Pong) error {
	return nil
}

// 节点推送
func (m *ServerServer) MemberNode(ss *core.Session, input *deal.MemberNodePush) error {
	router.GetMemberInfo().SetNode(input.Nodes)
	return nil
}
