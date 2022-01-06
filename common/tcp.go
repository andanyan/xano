package common

import "net"

type TcpClient struct {
	Conn           net.Conn
	ReadPacketChan chan *TcpPacket
	ReadPacketFunc func(packet *TcpPacket)
	SendPacketChan chan *TcpPacket
	SendFunc       func(packet *TcpPacket)
	Value          map[string]interface{}
}

func NewTcpClient(conn net.Conn, readFunc, sendFunc func(packet *TcpPacket)) *TcpClient {
	return &TcpClient{
		Conn:           conn,
		ReadPacketChan: make(chan *TcpPacket),
		ReadPacketFunc: readFunc,
		SendPacketChan: make(chan *TcpPacket),
		SendFunc:       sendFunc,
		Value:          make(map[string]interface{}),
	}
}

func (t *TcpClient) Handle() {
	go t.Send()
	t.Read()
}

func (t *TcpClient) Send() {
	for packet := range t.ReadPacketChan {
		t.Conn.Write(packet)
	}
}
