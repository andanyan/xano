package core

import (
	"net"
)

type TcpClient struct {
	*TcpHandle
}

func NewTcpClient(addr string) (*TcpClient, error) {
	// 生成一个连接
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}
	h := NewTcpHandle(conn)
	go h.handle()
	return &TcpClient{
		TcpHandle: h,
	}, nil
}
