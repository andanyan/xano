package session

import (
	"io"
	"net"
	"xlq-server/common"
	"xlq-server/component"
	"xlq-server/core"
	"xlq-server/log"
)

type TcpSession struct {
	BaseSession

	LastMid   uint64
	WriteChan chan []byte
	Conn      net.Conn
}

// 建立tcp Session
func NewTcpSession(conn net.Conn) *TcpSession {
	tcpSession := &TcpSession{}
	tcpSession.Id = common.GetUuid()
	tcpSession.Conn = conn
	tcpSession.ClientType = common.ClientTypeTcp
	tcpSession.Values = make(map[string]interface{})
	tcpSession.WriteChan = make(chan []byte)

	return tcpSession
}

// 回调
func (s *TcpSession) Write(route string, data interface{}) {
	//数据封装
	resBys, err := s.MsgMarsh(data)
	if err != nil {
		log.Println("Error Result Data:", route)
		return
	}
	resMsg := common.Msg{
		Mid:   s.LastMid,
		Route: route,
		Data:  resBys,
	}
	msgBys, err := s.MsgMarsh(resMsg)
	if err != nil {
		log.Println("Error Msg Data:", route)
		return
	}
	resPacket := common.Packet{
		Type:    s.MsgType,
		DataLen: uint16(len(msgBys)),
		Data:    msgBys,
	}
	packetBys, err := s.MsgMarsh(resPacket)
	if err != nil {
		log.Println("Error Package Data:", route)
		return
	}
	s.WriteChan <- packetBys
}

// 关闭session
func (s *TcpSession) Close() {
	close(s.WriteChan)
	s.Conn.Close()
	// 注销Session
}

func (s *TcpSession) Middlewares(route string) error {
	for _, f := range core.Options.TcpMiddlewares {
		err := f(route)
		if err != nil {
			return err
		}
	}
	return nil
}

// tcp读取信息
func (s *TcpSession) Handle() {
	// 发送
	go func() {
		for bys := range s.WriteChan {
			s.Conn.Write(bys)
		}
	}()

	packetBuf := make([]byte, 0)
	buf := make([]byte, common.ConnReadBufLen)
	for {
		n, err := s.Conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Println("Read [%s], session will be closed immediately", err.Error())
			}
			return
		}
		if n <= 0 {
			continue
		}
		packetBuf = append(packetBuf, buf[:n]...)

		packets, index := s.ParsePacket(packetBuf)
		packetLen := len(packets)
		if packetLen > 0 {
			packetBuf = packetBuf[index:]
			for _, packet := range packets {
				s.MsgType = packet.Type
				msg := new(common.Msg)
				err := s.MsgUnMarsh(packet.Data, msg)
				if err != nil {
					log.Println(err.Error())
					return
				}
				s.LastMid = msg.Mid
				if err := component.DoneMsg(s, msg); err != nil {
					log.Println(err.Error())
				}
			}
		}
	}

}
