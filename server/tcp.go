package server

import (
	"net"
	"xlq-server/core"
	"xlq-server/log"
	"xlq-server/session"
)

func (n *Node) serveTcp() {
	if core.Options.TcpAddr == "" {
		return
	}
	l, err := net.Listen("tcp", core.Options.TcpAddr)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err.Error())
			continue
		}
		go n.tcpHandle(conn)
	}
}

// 接收协议
func (n *Node) tcpHandle(conn net.Conn) {
	// session创建
	session := session.NewTcpSession(conn)

	// 停止操作
	defer func() {
		if err := recover(); err != nil {
			log.Fatal(err)
		}
		session.Close()
	}()

	session.Handle()
}
