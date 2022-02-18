package gate

import (
	"xano/core"
	"xano/deal"
	"xano/router"
)

type MemberServer struct{}

// 收到路由回包
func (s *MemberServer) AllNode(ss *core.Session, input *deal.AllNodeResponse) error {
	router.GetGateInfo().SetData(input.Nodes)
	return nil
}
