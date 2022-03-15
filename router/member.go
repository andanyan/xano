package router

var memberRouter *Router

func GetMemberRouter() *Router {
	if memberRouter == nil {
		memberRouter = NewRouter()
		memberRouter.Name = "Member"
	}
	return memberRouter
}

// 网关端服务寻址 主要实现地址寻址
var memberNode *Node

func GetMemberNode() *Node {
	if memberNode == nil {
		memberNode = NewNode()
	}
	return memberNode
}
