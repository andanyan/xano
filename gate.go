package core

import (
	"sync"
	"time"
	"xlq-server/common"
	"xlq-server/deal"
)

// 路由数据
type GateRoute struct {
	sync.RWMutex
	// addr -> item
	Routes map[string]*deal.GateRouteItem
}

//
var gateRoute *GateRoute

func init() {
	if gateRoute == nil {
		gateRoute = new(GateRoute)
		gateRoute.Routes = make(map[string]*deal.GateRouteItem)
	}
}

// 增加路由
func AddRoute(addr string, routes []string) {
	gateRoute.Lock()
	defer gateRoute.Unlock()

	item := &deal.GateRouteItem{
		Addr:     addr,
		Routes:   routes,
		LastTime: time.Now().Unix(),
	}
	gateRoute.Routes[addr] = item
}

func RemoveRoute(addr string) {
	gateRoute.Lock()
	defer gateRoute.Unlock()

	delete(gateRoute.Routes, addr)
}

// 获取全部路由
func GetAllRoute() []*deal.GateRouteItem {
	gateRoute.RLock()
	defer gateRoute.RUnlock()

	timeNow := time.Now().Unix()
	res := make([]*deal.GateRouteItem, 0)
	for k, v := range gateRoute.Routes {
		if timeNow-v.LastTime < int64(3*common.TcpHeartDuration) {
			res = append(res, v)
		} else {
			delete(gateRoute.Routes, k)
		}
	}
	return res
}
