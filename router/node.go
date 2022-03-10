package router

import (
	"math/rand"
	"sync"
	"xano/common"
	"xano/deal"
)

type Node struct {
	sync.RWMutex
	Nodes []*deal.ServerNode
}

func NewNode() *Node {
	return new(Node)
}

// 设置值
func (n *Node) SetNode(nodes []*deal.ServerNode) {
	n.Lock()
	defer n.Unlock()
	n.Nodes = nodes
}

// 获取一个地址
func (n *Node) GetNodeRand(route string) string {
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
