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
	NodeMap map[uint64]*deal.MemberNode
}

var localMemberNode *LocalMemberNode

func GetLocalMemberNode() *LocalMemberNode {
	if localMemberNode == nil {
		localMemberNode = new(LocalMemberNode)
		localMemberNode.NodeMap = make(map[uint64]*deal.MemberNode)
	}
	return localMemberNode
}

func (n *LocalMemberNode) SetNode(nodes []*deal.MemberNode) {
	n.Lock()
	defer n.Unlock()

	n.NodeMap = make(map[uint64]*deal.MemberNode)

	for _, node := range nodes {
		n.NodeMap[node.MchId] = node
	}
}

func (n *LocalMemberNode) GetNodeBySid(sid uint64) (*deal.MemberNode, error) {
	if sid == 0 {
		return nil, fmt.Errorf("INVAILD SID")
	}
	n.RLock()
	defer n.RUnlock()

	mchId := sid / common.MaxSessionNum

	node, ok := n.NodeMap[mchId]
	if !ok {
		return nil, fmt.Errorf("not find machine by sid")
	}

	return node, nil
}
