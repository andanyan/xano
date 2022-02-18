package router

import (
	"math/rand"
	"sync"
	"time"
	"xano/deal"
)

var MemberRouter = NewRouter()

// 普通网关节点
// 主要实现地址寻址
type GateInfo struct {
	sync.RWMutex
	Data []*deal.NodeItem
	// version#route = addr
	Route map[string][]string
}

var gateInfo *GateInfo

func init() {
	rand.Seed(time.Now().Unix())
}

func GetGateInfo() *GateInfo {
	if gateInfo == nil {
		gateInfo = new(GateInfo)
		gateInfo.Route = make(map[string][]string)
	}
	return gateInfo
}

// 设置值
func (g *GateInfo) SetData(nodes []*deal.NodeItem) {
	g.Lock()
	defer g.Unlock()

	g.Data = nodes
	g.Route = make(map[string][]string)

	for _, node := range g.Data {
		for _, r := range node.Routes {
			k := node.Version + "#" + r
			tmp, ok := g.Route[k]
			if !ok {
				tmp = make([]string, 0)
			}
			tmp = append(tmp, node.Addr)
			g.Route[node.Version+"#"+r] = tmp
		}
	}
}

// 获取节点
func (g *GateInfo) GetNodeAddr(route, version string) string {
	g.RLock()
	defer g.RUnlock()

	// 找出可以使用的节点，然后随机一个即可
	res := g.Route[version+"#"+route]
	rlen := len(res)

	if rlen == 0 {
		return ""
	}

	i := rand.Intn(rlen)

	return res[i]
}
