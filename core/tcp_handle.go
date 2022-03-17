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
type TcpInitFunc func(h *TcpHandle)
type TcpCloseFunc func(h *TcpHandle)

type TcpHandle struct {
	//value map[string]interface{}
	*common.Cache

	// tcp状态
	sync.RWMutex
	status     bool
	conn       net.Conn
	sendChan   chan *common.Packet
	readChan   chan *common.Packet
	handleFunc TcpHandleFunc
	initFunc   TcpInitFunc
	closeFunc  TcpCloseFunc
	// 消息id
	mid uint64
}

func NewTcpHandle(conn net.Conn) *TcpHandle {
	return &TcpHandle{
		Cache:    common.NewCache(),
		status:   true,
		conn:     conn,
		sendChan: make(chan *common.Packet, 100),
		readChan: make(chan *common.Packet, 100),
	}
}

// 设置初始化函数
func (h *TcpHandle) SetInitFunc(f TcpInitFunc) {
	h.initFunc = f
}

// 设置断开函数
func (h *TcpHandle) SetCloseFunc(f TcpCloseFunc) {
	h.closeFunc = f
}

// 设置包处理函数
func (h *TcpHandle) SetHandleFunc(handleFunc TcpHandleFunc) {
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
func (h *TcpHandle) Send(msg *deal.Msg) {
	if !h.status {
		logger.Warn("TCP IS DISCONNECT")
		return
	}
	msgBys, err := common.MsgMarsh(common.GetConfig().Base.TcpDeal, msg)
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
	if h.mid == common.MaxUint64 {
		h.mid = 0
	}
	return atomic.AddUint64(&h.mid, 1)
}

// 处理执行
func (h *TcpHandle) handle() {
	// 连接初始化
	if h.initFunc != nil {
		h.initFunc(h)
	}

	// 连接断开
	defer func() {
		h.Close()
		if h.closeFunc != nil {
			h.closeFunc(h)
		}
	}()

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
			err := common.MsgUnMarsh(common.GetConfig().Base.TcpDeal, p.Data, msg)
			if err != nil {
				logger.Error(err.Error())
				continue
			}
			// 消息id序号校验
			mmid := h.Get(common.HandleKeyMid)
			var nmmid uint64 = 0
			if mmid != nil {
				nmmid = mmid.(uint64)
			}
			if nmmid == common.MaxUint64 {
				nmmid = 0
			}
			if nmmid != msg.Mid-1 {
				logger.Error("Fatal MsgId: ", nmmid, msg.Mid)
				return
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
			logger.Debug(err)
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
