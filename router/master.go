package router

import (
	"sync"
	"time"
)

// 主节点路由
var MasterRouter = new(Router)

// 全节点路由信息记录
type MasterNode struct {
	sync.RWMutex

	// addr -> MasterNodeItem
	Data map[string]*MasterNodeItem
}
type MasterNodeItem struct {
	LastTime int64
	Status   bool
	Version  string
	Addr     string
	Routes   []string
}

var masterNode *MasterNode

func GetMasterNode() *MasterNode {
	if masterNode == nil {
		masterNode = new(MasterNode)
		masterNode.Data = make(map[string]*MasterNodeItem)
	}
	return masterNode
}

// 增加节点
func (m *MasterNode) AddNode(addr, version string, routes []string) bool {
	m.Lock()
	defer m.Unlock()

	if m.Data[addr] != nil {
		return false
	}

	m.Data[addr] = &MasterNodeItem{
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

// 更改时间
func (m *MasterNode) UpdateTime(addr string, t int64) {
	m.Lock()
	defer m.Unlock()

	if m.Data[addr] == nil {
		return
	}

	m.Data[addr].LastTime = t
}

// 获取全部节点
func (m *MasterNode) GetAllNode() []*MasterNodeItem {
	m.RLock()
	defer m.RUnlock()
	res := make([]*MasterNodeItem, 0)
	for _, node := range m.Data {
		res = append(res, node)
	}
	return res
}