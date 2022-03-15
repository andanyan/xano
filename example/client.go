package main

import (
	"time"
	"xano"
	"xano/common"
	"xano/core"
	"xano/deal"
	"xano/example/pb"
	"xano/logger"
)

func main() {
	xano.WithConfig("./config/client.toml")

	var sid uint64 = 0

	pool := core.GetPool("0.0.0.0:12000")
	cli, err := pool.Get()
	if err != nil {
		logger.Fatal(err)
	}

	c := make(chan struct{})

	cli.Client.SetHandleFunc(func(h *core.TcpHandle, m *deal.Msg) {
		if m.Route == "SessionInit" {
			sid = m.Sid
			logger.Infof("%+v", m)
			c <- struct{}{}
			return
		}
		out := new(pb.DivResponse)
		err := common.MsgUnMarsh(m.Deal, m.Data, out)
		if err != nil {
			logger.Error(err)
			return
		}
		logger.Infof("%+v", out)
	})

	// session init
	input := &deal.SessionInitRequest{}
	inputBys, err := common.MsgMarsh(common.GetConfig().Base.TcpDeal, input)
	if err != nil {
		logger.Error(err)
		return
	}
	inputMsg := &deal.Msg{
		Route:   "SessionInit",
		Sid:     0,
		Mid:     cli.Client.GetMid(),
		MsgType: common.MsgTypeRequest,
		Deal:    common.GetConfig().Base.TcpDeal,
		Version: common.GetConfig().Base.Version,
		Data:    inputBys,
	}
	logger.Debugf("%+v", inputMsg)
	cli.Client.Send(inputMsg)

	// 等待初始化完成
	<-c

	var i int64
	for {
		input := &pb.DivRequest{
			A: i + 1,
			B: i + 3,
		}
		inputBys, err := common.MsgMarsh(common.GetConfig().Base.TcpDeal, input)
		if err != nil {
			logger.Error(err)
		}
		inputMsg := &deal.Msg{
			Route:   "Div",
			Sid:     sid,
			Mid:     cli.Client.GetMid(),
			MsgType: common.MsgTypeRequest,
			Deal:    common.GetConfig().Base.TcpDeal,
			Version: common.GetConfig().Base.Version,
			Data:    inputBys,
		}
		logger.Debugf("%+v", inputMsg)
		cli.Client.Send(inputMsg)
		i++
		time.Sleep(1000 * time.Millisecond)
	}

	<-c
}
