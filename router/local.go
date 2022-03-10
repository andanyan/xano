package router

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

// 节点信息
var localNode *Node

func GetLocalNode() *Node {
	if localNode == nil {
		localNode = NewNode()
	}
	return localNode
}
