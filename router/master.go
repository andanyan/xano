package router

import (
	"sync"
	"time"
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

	// addr -> MasterNodeItem
	Data map[string]*deal.NodeItem
}

var masterNode *MasterNode

func GetMasterNode() *MasterNode {
	if masterNode == nil {
		masterNode = new(MasterNode)
		masterNode.Data = make(map[string]*deal.NodeItem)
	}
	return masterNode
}

// 增加节点
func (m *MasterNode) AddNode(addr, version string, routes []string) bool {
	m.Lock()
	defer m.Unlock()

	if m.Data[addr] != nil {
		m.Data[addr].LastTime = time.Now().Unix()
		return false
	}

	m.Data[addr] = &deal.NodeItem{
		LastTime: time.Now().Unix(),
		Status:   true,
		Version:  version,
		Addr:     addr,
		Routes:   routes,
	}
	return true
}

// 节点关闭
func (m *MasterNode) RemoveNode(addr string) {
	m.Lock()
	defer m.Unlock()

	if m.Data[addr] != nil {
		return
	}

	m.Data[addr].Status = false
}

// 获取全部节点
func (m *MasterNode) GetAllNode() []*deal.NodeItem {
	m.RLock()
	defer m.RUnlock()
	res := make([]*deal.NodeItem, 0)
	for _, node := range m.Data {
		res = append(res, node)
	}
	return res
}
