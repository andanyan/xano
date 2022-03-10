package master

import (
	"xano/deal"
	"xano/router"
	"xano/session"
)

type MasterHttpServer struct{}

// 获取所有的node
func (s *MasterHttpServer) ServerNodes(ss *session.Session, input *deal.ServerNodesRequest) error {
	serverNodes := router.GetMasterNode().AllServerNode()
	return ss.Response("ServerNodes", &deal.ServerNodesResponse{
		Nodes: serverNodes,
	})
}

// 获取所有的node
func (s *MasterHttpServer) MemberNodes(ss *session.Session, input *deal.MemberNodesRequest) error {
	memberNodes := router.GetMasterNode().AllMemberNode()
	return ss.Response("ServerNodes", &deal.MemberNodesResponse{
		Nodes: memberNodes,
	})
}
