package core

import (
	"net"
	"sync"
	"sync/atomic"
	"xano/common"
	"xano/deal"
	"xano/logger"
)

type TcpHandleFunc func(h *TcpHandle, m *deal.Msg)

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
		sendChan: make(chan *common.Packet, 100),
		readChan: make(chan *common.Packet, 100),
	}
}

// 添加处理函数
func (h *TcpHandle) SetHandle(handleFunc TcpHandleFunc) {
	h.handleFunc = handleFunc
}

// 处理关闭
func (h *TcpHandle) Close() {
	if !h.status {
		return
	}
	h.status = false
	close(h.readChan)
	close(h.sendChan)
	h.conn.Close()
}

// 包入列
func (h *TcpHandle) Send(m *deal.Msg) {
	if !h.status {
		return
	}
	msgBys, err := common.MsgMarsh(common.TcpDealProtobuf, m)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	packet := &common.Packet{
		Length: uint16(len(msgBys)),
		Data:   msgBys,
	}
	h.sendChan <- packet
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

// 获取消息id
func (h *TcpHandle) GetMid() uint64 {
	return atomic.AddUint64(&h.mid, 1)
}

// 处理执行
func (h *TcpHandle) handle() {
	go h.runSend()
	go h.runRead()
	h.handleRead()
}

// 包处理
func (h *TcpHandle) runRead() {
	for p := range h.readChan {
		if !h.status {
			break
		}

		if h.handleFunc != nil {
			msg := new(deal.Msg)
			err := common.MsgUnMarsh(common.TcpDealProtobuf, p.Data, msg)
			if err != nil {
				logger.Error(err.Error())
				continue
			}
			// 消息id序号校验
			mmid := h.Get(common.HandleKeyMid)
			if mmid != nil && mmid.(uint64) != msg.Mid-1 {
				continue
			}
			// 设置当前消息id
			h.Set(common.HandleKeyMid, msg.Mid)
			h.handleFunc(h, msg)
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
			logger.Error(err.Error())
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
			logger.Error(err.Error())
			h.Close()
			break
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
