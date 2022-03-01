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
	xano.WithConfig("./config.toml")

	pool := core.GetPool("0.0.0.0:12000")
	cli, err := pool.Get()
	if err != nil {
		logger.Fatal(err)
	}
	cli.Client.SetHandle(func(h *core.TcpHandle, m *deal.Msg) {
		out := new(pb.DivResponse)
		err := common.MsgUnMarsh(m.Deal, m.Data, out)
		if err != nil {
			logger.Error(err)
			return
		}
		logger.Infof("%+v", out)
	})

	var i int64
	for {
		input := &pb.DivRequest{
			A: i + 1,
			B: i + 3,
		}
		logger.Debugf("%+v", input)
		inputBys, err := common.MsgMarsh(common.TcpDealProtobuf, input)
		if err != nil {
			logger.Error(err)
		}
		inputMsg := &deal.Msg{
			Route:   "Div",
			Mid:     cli.Client.GetMid(),
			MsgType: common.MsgTypeRequest,
			Deal:    common.TcpDealProtobuf,
			Version: common.GetConfig().Base.Version,
			Data:    inputBys,
		}
		cli.Client.Send(inputMsg)
		i++
		time.Sleep(time.Second)
	}

	c := make(chan struct{})
	<-c
}
