package session

import (
	"fmt"
	"time"
	"xano/common"
	"xano/core"
	"xano/deal"
	"xano/logger"
)

type NetService struct{}

// 所有请求均针对网关层
func NewNetService() *NetService {
	return new(NetService)
}

// Request
func (n *NetService) Request(addr string, msg *deal.Msg) ([]*deal.Msg, error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
		}
	}()

	// 获取连接
	pool := core.GetPool(addr)
	cli, err := pool.Get()
	if err != nil {
		return nil, err
	}
	defer pool.Recycle(cli)

	// 设置回包
	resMsgs := make([]*deal.Msg, 0)
	c := make(chan struct{})
	cli.Client.SetHandleFunc(func(h *core.TcpHandle, m *deal.Msg) {
		resMsgs = append(resMsgs, m)
		if m.MsgType == common.MsgTypeResponse {
			c <- struct{}{}
		}
	})
	defer cli.Client.SetHandleFunc(nil)

	// 发送消息
	msg.Mid = cli.Client.GetMid()
	cli.Client.Send(msg)

	t := time.NewTimer(common.TcpDeadDuration)

	select {
	case <-c:

	case <-t.C:
		return nil, fmt.Errorf("handle timeout")
	}

	return resMsgs, nil
}

// Notice
func (n *NetService) Notice(addr string, msg *deal.Msg) error {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
		}
	}()

	// 获取连接
	pool := core.GetPool(addr)
	cli, err := pool.Get()
	if err != nil {
		return err
	}
	defer pool.Recycle(cli)

	// 发送消息
	msg.Mid = cli.Client.GetMid()
	cli.Client.Send(msg)
	return nil
}
