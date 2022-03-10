package core

import (
	"net"
	"xano/logger"
)

func NewTcpServer(addr string, opsFunc ...func(h *TcpHandle)) {
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
		for _, f := range opsFunc {
			f(h)
		}
		//h.SetHandle(handleFunc)
		go h.handle()
	}
}
