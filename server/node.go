package server

// 当前节点数据
type Node struct{}

var node *Node

func init() {
	node = new(Node)
}

func (n *Node) Start() {
	go n.serveTcp()
	go n.serveHttp()

}
