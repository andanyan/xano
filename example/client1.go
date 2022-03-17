package main

import (
	"xano"
	"xano/common"
	"xano/core"
	"xano/deal"
	"xano/example/pb"
	"xano/logger"
)

func main() {
	xano.WithConfig("./config/client.toml")

	//var sid uint64 = 0

	pool := core.GetPool("0.0.0.0:11000")
	cli, err := pool.Get()
	if err != nil {
		logger.Fatal(err)
	}
	defer pool.Recycle(cli)

	c := make(chan struct{})

	cli.Client.SetHandleFunc(func(h *core.TcpHandle, m *deal.Msg) {
		out := new(pb.DivResponse)
		err := common.MsgUnMarsh(m.Deal, m.Data, out)
		if err != nil {
			logger.Error(err)
			return
		}
		logger.Infof("%+v", out)
	})

	<-c
}
