package router

import (
	"math/rand"
	"sync"
	"xano/common"
	"xano/deal"
)

var memberRouter *Router

func GetMemberRouter() *Router {
	if memberRouter == nil {
		memberRouter = NewRouter()
		memberRouter.Name = "Member"
	}
	return memberRouter
}

// 网关端服务寻址 主要实现地址寻址
type MemberServerNode struct {
	sync.RWMutex
	Nodes []*deal.ServerNode
}

var memberServerNode *MemberServerNode

func GetMemberServerNode() *MemberServerNode {
	if memberServerNode == nil {
		memberServerNode = new(MemberServerNode)
	}
	return memberServerNode
}

// 设置值
func (n *MemberServerNode) SetNode(nodes []*deal.ServerNode) {
	n.Lock()
	defer n.Unlock()
	n.Nodes = nodes
}

// 获取一个地址
func (n *MemberServerNode) GetNodeRand(route string) string {
	n.RLock()
	defer n.RUnlock()
	nodes := make([]*deal.ServerNode, 0)
	for _, item := range n.Nodes {
		if item.Version == common.GetConfig().Base.Version && common.InStringArr(route, item.Routes) {
			nodes = append(nodes, item)
		}
	}
	nodeLen := len(nodes)
	if nodeLen == 0 {
		return ""
	}
	index := rand.Intn(nodeLen)
	return nodes[index].Addr
}
