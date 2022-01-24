package core

import (
	"io"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"
	"xlq-server/common"
	"xlq-server/deal"
)

type TcpHandleFunc func(h *TcpHandle, p *common.Packet)

type TcpHandle struct {
	value map[string]interface{}

	// tcp状态
	sync.RWMutex
	status     bool
	conn       net.Conn
	sendChan   chan *common.Packet
	readChan   chan *common.Packet
	handleFunc TcpHandleFunc
	// 消息id
	mid uint64
}

func NewTcpHandle(conn net.Conn) *TcpHandle {
	return &TcpHandle{
		value:    make(map[string]interface{}),
		status:   true,
		conn:     conn,
		sendChan: make(chan *common.Packet),
		readChan: make(chan *common.Packet),
	}
}

// 添加处理函数
func (h *TcpHandle) SetHandle(handleFunc TcpHandleFunc) {
	h.handleFunc = handleFunc
}

// 处理关闭
func (h *TcpHandle) Close() {
	h.status = false
	close(h.readChan)
	close(h.sendChan)
	h.conn.Close()
}

// 包入列
func (h *TcpHandle) Send(p *common.Packet) {
	if !h.status {
		return
	}
	h.sendChan <- p
}

// 设置值
func (h *TcpHandle) Set(k string, v interface{}) {
	h.Lock()
	defer h.Unlock()
	h.value[k] = v
}

// 获取值
func (h *TcpHandle) Get(k string) interface{} {
	h.RLock()
	defer h.RUnlock()
	return h.value[k]
}

// 删除值
func (h *TcpHandle) Del(k string) {
	h.Lock()
	defer h.Unlock()
	delete(h.value, k)
}

// 获取地址
func (h *TcpHandle) GetAddr() string {
	return h.conn.RemoteAddr().String()
}

// 获取连接状态
func (h *TcpHandle) Status() bool {
	return h.status
}

func (h *TcpHandle) GetMid() uint64 {
	return atomic.AddUint64(&h.mid, 1)
}

// 处理执行
func (h *TcpHandle) handle() {
	go h.runSend()
	go h.runRead()
	h.handleRead()
}

// 心跳包 仅做保连
func (h *TcpHandle) ping() {
	for {
		time.Sleep(common.TcpHeartDuration)

		ping := &deal.Ping{
			Ping: time.Now().Unix(),
		}
		pingBys, err := common.MsgMarsh(common.TcpDealProtobuf, ping)
		if err != nil {
			log.Println(err)
			continue
		}

		msg := &deal.Msg{
			Route:   "Ping",
			Mid:     h.GetMid(),
			MsgType: common.MsgTypeRequest,
			Deal:    common.TcpDealProtobuf,
			Data:    pingBys,
		}
		msgBys, err := common.MsgMarsh(common.TcpDealProtobuf, msg)
		if err != nil {
			log.Println(err)
			continue
		}
		packet := &common.Packet{
			Length: uint16(len(msgBys)),
			Data:   msgBys,
		}
		h.Send(packet)
	}
}
func (h *TcpHandle) pong(p *common.Packet) bool {
	msg := new(deal.Msg)
	err := common.MsgUnMarsh(msg.Deal, p.Data, msg)
	if err != nil {
		log.Println(err)
		return true
	}

	if msg.Route == "Ping" {
		pong := &deal.Pong{
			Pong: time.Now().Unix(),
		}
		pongBys, err := common.MsgMarsh(common.TcpDealProtobuf, pong)
		if err != nil {
			log.Println(err)
			return true
		}

		rmsg := &deal.Msg{
			Route:   "Pong",
			Mid:     msg.Mid,
			MsgType: common.MsgTypeResponse,
			Deal:    common.TcpDealProtobuf,
			Data:    pongBys,
		}
		rmsgBys, err := common.MsgMarsh(common.TcpDealProtobuf, rmsg)
		if err != nil {
			log.Println(err)
			return true
		}
		packet := &common.Packet{
			Length: uint16(len(rmsgBys)),
			Data:   rmsgBys,
		}
		h.Send(packet)
		return true
	}

	return false
}

// 包处理
func (h *TcpHandle) runRead() {
	for p := range h.readChan {
		if !h.status {
			break
		}

		if h.pong(p) {
			continue
		}

		if h.handleFunc != nil {
			h.handleFunc(h, p)
		}
	}
}

// 包发送
func (h *TcpHandle) runSend() {
	for packet := range h.sendChan {
		if !h.status {
			break
		}
		bys, err := common.PacketMarsh(packet)
		if err != nil {
			log.Println(err)
			continue
		}
		h.conn.Write(bys)
	}
}

// 包读取
func (h *TcpHandle) handleRead() {
	defer h.Close()

	// 读包
	cacheBuf := make([]byte, 0)
	buf := make([]byte, common.TcpReadBufLen)
	for {
		n, err := h.conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				// 连接已断开
				log.Println("tcp connection is disconnected")
				h.Close()
				break
			}
			log.Println("tcp error: " + err.Error())
			continue
		}
		if n <= 0 {
			continue
		}
		cacheBuf = append(cacheBuf, buf[:n]...)

		index, packets := common.PacketUnMarsh(cacheBuf)
		if index > 0 {
			cacheBuf = cacheBuf[index:]
			for _, p := range packets {
				h.readChan <- p
			}
		}
	}
}
