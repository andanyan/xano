package core

import (
	"log"
	"net"
)

func NewTcpServer(addr string, handleFunc TcpHandleFunc) {
	l, err := net.Listen("tcp", addr)
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
		h := NewTcpHandle(conn, false)
		h.SetHandle(handleFunc)
		go h.handle()
	}
}
