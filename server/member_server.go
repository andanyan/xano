package server

import (
	"xano/core"
	"xano/deal"
	"xano/router"
)

type MemberServer struct{}

// 心跳返回
func (m *MemberServer) MemberHeart(ss *core.Session, input *deal.Pong) error {
	return nil
}

// 节点推送
func (m *MemberServer) ServerNode(ss *core.Session, input *deal.ServerNodePush) error {
	router.GetGateInfo().SetNode(input.Nodes)
	return nil
}
