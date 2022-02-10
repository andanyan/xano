package core

import (
	"net"
	"xlq-server/logger"
)

func NewTcpServer(addr string, handleFunc TcpHandleFunc) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			logger.Error(err.Error())
			continue
		}
		h := NewTcpHandle(conn)
		h.SetHandle(handleFunc)
		go h.handle()
	}
}
