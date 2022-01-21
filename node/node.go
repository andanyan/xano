package node

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"xlq-server/common"
)

type Node struct {
	Gate   *Gate
	Server *Server
}

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
	go gate.RunInner()

	// 服务启动
	server := NewServer()
	go server.Run()
	// 加入到gate中
	go server.AddGate()

	// 赋值
	n.Gate = gate
	n.Server = server

	// 监听信号
	sg := make(chan os.Signal, 1)
	signal.Notify(sg, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	// 阻塞进程
	<-sg

	// 关闭
	log.Panicln("server is stopping...")
	n.Close()
	time.Sleep(3 * time.Second)
	log.Panicln("server stopped")
}

func (n *Node) Close() {
	n.Server.CloseGate()
}
