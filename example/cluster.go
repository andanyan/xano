package main

import (
	"os"
	"xano"
	"xano/example/pb"
	"xano/logger"
	"xano/router"
	"xano/session"
)

type A struct{}

func (a *A) Add(s *session.Session, input *pb.AddRequest) error {
	var res int64
	for _, val := range input.Args {
		res += val
	}
	return s.Response("Add", &pb.AddResponse{
		Result: res,
	})
}

type B struct{}

func (b *B) Div(s *session.Session, input *pb.DivRequest) error {
	addRes := new(pb.AddResponse)
	err := s.Rpc("Add", &pb.AddRequest{
		Args: []int64{input.A, input.B},
	}, addRes)
	if err != nil {
		logger.Error(err)
		return err
	}

	res := addRes.Result * (input.B - input.A)

	return s.Response("Div", &pb.DivResponse{
		Result: res,
	})
}

func main() {
	xano.WithConfig("./config/cluster.toml")

	xano.WithRoute(&router.RouterServer{
		Name:   "",
		Server: new(A),
	})
	xano.WithRoute(&router.RouterServer{
		Name:   "",
		Server: new(B),
	})
	xano.WithLog(os.Stdout, logger.LoggerLevelInfo)

	xano.Run()
}
