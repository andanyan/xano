package main

import (
	"xano"
	"xano/example/pb"
	"xano/logger"
	"xano/router"
	"xano/session"
)

type B struct{}

func (b *B) Div(s session.Session, input *pb.DivRequest) error {
	addRes := new(pb.AddResponse)
	err := s.Rpc("Add", &pb.AddRequest{
		Args: []int64{input.A, input.B},
	}, addRes)
	if err != nil {
		logger.Error(err)
		return err
	}

	res := addRes.Result * (input.B - input.A)
	logger.Info("div res: ", res)
	return s.Response("Div", &pb.DivResponse{
		Result: res,
	})
}

func main() {
	xano.WithConfig("./config/server2.toml")

	xano.WithRoute(&router.RouterServer{
		Name:   "",
		Server: new(B),
	})

	xano.Run()
}
