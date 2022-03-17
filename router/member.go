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
	//Nodes    []*deal.ServerNode
	RouteMap map[string]*MemberServerRoute
}

type MemberServerRoute struct {
	LastVersion string
	Addr        []string
	Version     []string
}

var memberServerNode *MemberServerNode

func GetMemberServerNode() *MemberServerNode {
	if memberServerNode == nil {
		memberServerNode = new(MemberServerNode)
		memberServerNode.RouteMap = make(map[string]*MemberServerRoute)
	}
	return memberServerNode
}

// 设置值
func (n *MemberServerNode) SetNode(nodes []*deal.ServerNode) {
	n.Lock()
	defer n.Unlock()

	n.RouteMap = make(map[string]*MemberServerRoute)

	for _, node := range nodes {
		for _, route := range node.Routes {
			item, ok := n.RouteMap[route]
			if !ok {
				item = new(MemberServerRoute)
			}
			item.Addr = append(item.Addr, node.Addr)
			item.Version = append(item.Version, node.Version)
			if common.VersionCompare(node.Version, item.LastVersion) {
				item.LastVersion = node.Version
			}
			n.RouteMap[route] = item
		}
	}
}

// 获取一个地址
func (n *MemberServerNode) GetNodeRand(version, route string) string {
	n.RLock()
	defer n.RUnlock()

	// 获取地址
	rs, ok := n.RouteMap[route]
	if !ok {
		return ""
	}
	if len(rs.Addr) == 0 {
		return ""
	}
	if version == "" {
		version = rs.LastVersion
	}
	ids := make([]int, 0)
	for i := 0; i < len(rs.Version); i++ {
		if rs.Version[i] == version {
			ids = append(ids, i)
		}
	}
	idsLen := len(ids)
	if idsLen == 0 {
		return ""
	}
	index := rand.Intn(idsLen)
	return rs.Addr[index]
}
