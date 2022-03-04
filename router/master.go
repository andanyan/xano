package router

import (
	"sync"
	"xano/deal"
)

// 主节点路由
var masterRouter *Router

func GetMasterRouter() *Router {
	if masterRouter == nil {
		masterRouter = NewRouter()
		masterRouter.Name = "Master"
	}
	return masterRouter
}

// 全节点路由信息记录
type MasterNode struct {
	sync.RWMutex

	// 从节点
	MemberNode []*deal.MemberNode

	//服务节点
	ServerNode []*deal.ServerNode
}

var masterNode *MasterNode

func GetMasterNode() *MasterNode {
	if masterNode == nil {
		masterNode = new(MasterNode)
	}
	return masterNode
}

// 获取所有从节点
func (m *MasterNode) AllMemberNode() []*deal.MemberNode {
	return m.MemberNode
}

// 增加从节点
func (m *MasterNode) AddMemberNode(node *deal.MemberNode) {
	m.Lock()
	defer m.Unlock()
	// 判断是否已存在从节点
	for key, item := range m.MemberNode {
		if item.Addr == node.Addr {
			m.MemberNode[key] = node
			return
		}
	}
	m.MemberNode = append(m.MemberNode, node)
}

// 移除从节点
func (m *MasterNode) RemoveMemberNode(addr string) {
	m.Lock()
	defer m.Unlock()
	index := 0
	for _, item := range m.MemberNode {
		if item.Addr != addr {
			m.MemberNode[index] = item
			index++
		}
	}
	m.MemberNode = m.MemberNode[:index]
}

// 获取所有的服务节点
func (m *MasterNode) AllServerNode() []*deal.ServerNode {
	return m.ServerNode
}

// 服务节点增加
func (m *MasterNode) AddServerNode(node *deal.ServerNode) {
	m.Lock()
	defer m.Unlock()
	// 判断是否已存在从节点
	for key, item := range m.ServerNode {
		if item.Addr == node.Addr {
			m.ServerNode[key] = node
			return
		}
	}
	m.ServerNode = append(m.ServerNode, node)
}

// 节点关闭
func (m *MasterNode) RemoveServerNode(addr string) {
	m.Lock()
	defer m.Unlock()
	index := 0
	for _, item := range m.ServerNode {
		if item.Addr != addr {
			m.ServerNode[index] = item
			index++
		}
	}
	m.ServerNode = m.ServerNode[:index]
}
