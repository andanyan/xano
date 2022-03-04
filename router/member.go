package router

import (
	"math/rand"
	"sync"
	"time"
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

// 普通网关节点
// 主要实现地址寻址
type GateInfo struct {
	sync.RWMutex
	Nodes []*deal.ServerNode
}

var gateInfo *GateInfo

func init() {
	rand.Seed(time.Now().Unix())
}

func GetGateInfo() *GateInfo {
	if gateInfo == nil {
		gateInfo = new(GateInfo)
	}
	return gateInfo
}

// 设置值
func (g *GateInfo) SetNode(nodes []*deal.ServerNode) {
	g.Lock()
	defer g.Unlock()
	g.Nodes = nodes
}

// 获取一个地址
func (g *GateInfo) GetNodeRand(route string) string {
	g.RLock()
	defer g.RUnlock()
	nodes := make([]*deal.ServerNode, 0)
	for _, item := range g.Nodes {
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
