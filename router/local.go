package router

import (
	"math/rand"
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

type MemberInfo struct {
	sync.RWMutex
	Nodes []*deal.MemberNode
}

var memberInfo *MemberInfo

func GetMemberInfo() *MemberInfo {
	if memberInfo == nil {
		memberInfo = new(MemberInfo)
	}
	return memberInfo
}

// 设置值
func (m *MemberInfo) SetNode(nodes []*deal.MemberNode) {
	m.Lock()
	defer m.Unlock()
	m.Nodes = nodes
}

// 获取一个地址
func (m *MemberInfo) GetNodeRand() string {
	m.RLock()
	defer m.RUnlock()
	nodes := make([]*deal.MemberNode, 0)
	for _, item := range m.Nodes {
		if item.Version == common.GetConfig().Base.Version {
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
