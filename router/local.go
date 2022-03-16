package router

import (
	"fmt"
	"sync"
	"xano/common"
	"xano/deal"
)

// 记录本地路由
var localRouter *Router

func GetLocalRouter() *Router {
	if localRouter == nil {
		localRouter = NewRouter()
		localRouter.Name = "Local"
	}
	return localRouter
}

// 与主节点通信回包路由
var gateRouter *Router

func GetGateRouter() *Router {
	if gateRouter == nil {
		gateRouter = NewRouter()
		gateRouter.Name = "Gate"
	}
	return gateRouter
}

type LocalMemberNode struct {
	sync.RWMutex
	Nodes []*deal.MemberNode
}

var localMemberNode *LocalMemberNode

func GetLocalMemberNode() *LocalMemberNode {
	if localMemberNode == nil {
		localMemberNode = new(LocalMemberNode)
	}
	return localMemberNode
}

func (n *LocalMemberNode) SetNode(nods []*deal.MemberNode) {
	n.Lock()
	defer n.Unlock()
	n.Nodes = nods
}

func (n *LocalMemberNode) GetNodeBySid(sid uint64) (*deal.MemberNode, error) {
	if sid == 0 {
		return nil, fmt.Errorf("INVAILD SID")
	}
	n.RLock()
	defer n.RUnlock()

	mchId := sid / common.MaxSessionNum

	for _, item := range n.Nodes {
		if item.MchId == mchId {
			return item, nil
		}
	}
	return nil, fmt.Errorf("not find machine by sid")
}
