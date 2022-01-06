package gate

import (
	"log"
	"net"
	"xlq-server/common"
)

type GateTcp struct {
	Conn              net.Conn
	ReceivePacketChan chan *common.TcpPacket
	SendPacketChan    chan *common.TcpPacket
}

func (g *Gate) RunTcp() {
	if common.GetGateConfig().TcpAddr == "" {
		return
	}

	l, err := net.Listen("tcp", common.GetGateConfig().TcpAddr)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Panicln(err.Error())
			continue
		}
		go g.tcpHandle(conn)
	}
}

// tcp处理
func (g *Gate) tcpHandle(conn net.Conn) {
	tcpClient := &common.TcpClient{
		Conn:           conn,
		ReadPacketChan: make(chan *common.TcpPacket),

		SendPacketChan: make(chan *common.TcpPacket),
		Value:          make(map[string]interface{}),
	}

}
