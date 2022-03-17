package master

import (
	"xano/deal"
	"xano/router"
	"xano/session"
)

type MasterHttpServer struct{}

// 获取所有的node
func (s *MasterHttpServer) ServerNode(ss session.Session, input *deal.ServerNodeRequest) error {
	serverNodes := router.GetMasterNode().AllServerNode()
	return ss.Response("ServerNode", &deal.ServerNodeResponse{
		Node: serverNodes,
	})
}

// 获取所有的node
func (s *MasterHttpServer) MemberNode(ss session.Session, input *deal.MemberNodeRequest) error {
	memberNodes := router.GetMasterNode().AllMemberNode()
	return ss.Response("MemberNode", &deal.MemberNodeResponse{
		Node: memberNodes,
	})
}
