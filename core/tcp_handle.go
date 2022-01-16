package core

import (
	"io"
	"log"
	"net"
	"sync"
	"time"
	"xlq-server/common"
)

type TcpHandleFunc func(h *TcpHandle, p *common.TcpPacket)

type TcpHandle struct {
	Value map[string]interface{}

	// tcp状态
	sync.RWMutex
	status     bool
	conn       net.Conn
	sendChan   chan *common.TcpPacket
	readChan   chan *common.TcpPacket
	handleFunc TcpHandleFunc
	mid        uint64
	// 连接类型
	isClient bool
}

func NewTcpHandle(conn net.Conn, isClient bool) *TcpHandle {
	return &TcpHandle{
		Value:    make(map[string]interface{}),
		status:   true,
		conn:     conn,
		sendChan: make(chan *common.TcpPacket),
		readChan: make(chan *common.TcpPacket),
		// handleFunc: handleFunc,
		isClient: isClient,
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
func (h *TcpHandle) Send(p *common.TcpPacket) {
	if !h.status {
		return
	}
	h.sendChan <- p
}

// 设置值
func (h *TcpHandle) Set(k string, v interface{}) {
	h.Lock()
	defer h.Unlock()
	h.Value[k] = v
}

// 获取值
func (h *TcpHandle) Get(k string) interface{} {
	h.RLock()
	defer h.RUnlock()
	return h.Value[k]
}

// 获取新的消息id
func (h *TcpHandle) GetMid() uint64 {
	h.Lock()
	defer h.Unlock()
	mid := h.mid
	h.mid++
	return mid
}

// 处理执行
func (h *TcpHandle) handle() {
	go h.runSend()
	go h.runRead()
	if h.isClient {
		go h.runHeart()
	}
	h.handleRead()
}

// 心跳处理
func (h *TcpHandle) runHeart() {
	heartPacket := &common.TcpPacket{
		Length: common.HeartPacketLength,
		Data:   []byte(common.HeartPacketRequest),
	}
	for {
		if !h.status {
			break
		}
		time.Sleep(10 * time.Second)
		h.Send(heartPacket)
	}
}

// 包处理
func (h *TcpHandle) runRead() {
	heartPacket := &common.TcpPacket{
		Length: common.HeartPacketLength,
		Data:   []byte(common.HeartPacketResponse),
	}
	for p := range h.readChan {
		if !h.status {
			break
		}
		// 心跳包辨别
		if p.Length == common.HeartPacketLength && string(p.Data) == common.HeartPacketRequest {
			h.Send(heartPacket)
			return
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
