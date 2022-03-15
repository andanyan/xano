package member

import (
	"xano/common"
	"xano/deal"
	"xano/router"
	"xano/session"
)

type Service struct{}

// 心跳返回
func (m *Service) MemberHeart(ss *session.Session, input *deal.Ping) error {
	return nil
}

// 启动回包
func (m *Service) MemberStart(ss *session.Session, input *deal.MemberStartResponse) error {
	router.GetMemberNode().SetNode(input.Nodes)
	return nil
}

// 获取sid
func (m *Service) MemberMchID(ss *session.Session, input *deal.MemberMchIDPush) error {
	common.GetCache().Set(common.MchIDKey, input.MchID)
	return nil
}
