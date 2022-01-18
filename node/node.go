package node

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"xlq-server/common"
)

type Node struct{}

func NewNode() *Node {
	return new(Node)
}

func (n *Node) WithGate(conf *common.GateConfig) {
	common.SetGateConfig(conf)
}

func (n *Node) WithService(conf *common.TcpServiceConfig) {
	common.SetServiceConfig(conf)
}

// 运行
func (n *Node) Run() {
	// 网关启动
	gate := NewGate()
	go gate.RunTcp()
	go gate.RunHttp()

	// 服务启动
	server := NewServer()
	go server.Run()

	// 监听信号
	sg := make(chan os.Signal, 1)
	signal.Notify(sg, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	// 阻塞进程
	<-sg

	// 关闭
	log.Panicln("server is stopping...")
	n.Close()
	log.Panicln("server stopped")
}

func (n *Node) Close() {
	// 数据维护和存储

}
