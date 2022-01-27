package gate

import (
	"xlq-server/core"
	"xlq-server/deal"
	"xlq-server/router"
)

type MemberServer struct{}

// 收到路由回包
func (s *MemberServer) AllNode(ss *core.Session, input *deal.AllNodeResponse) error {
	router.GetGateInfo().SetData(input.Nodes)
	return nil
}
