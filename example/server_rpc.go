package main

import (
	"xano"
	"xano/core"
	"xano/example/pb"
	"xano/logger"
	"xano/router"
)

type B struct{}

func (b *B) Div(s *core.Session, input *pb.DivRequest) error {
	addRes := new(pb.AddResponse)
	err := s.Rpc("Add", &pb.AddRequest{
		Args: []int64{input.A, input.B},
	}, addRes)
	if err != nil {
		logger.Error(err)
		return err
	}

	res := addRes.Result * (input.A - input.B)

	return s.Response("Div", &pb.DivResponse{
		Result: res,
	})
}

func main() {
	xano.WithConfig("./config_rpc.toml")

	xano.WithRoute(&router.RouterServer{
		Name:   "",
		Server: new(B),
	})

	xano.Run()
}
