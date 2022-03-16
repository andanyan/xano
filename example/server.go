package main

import (
	"xano"
	"xano/example/pb"
	"xano/router"
	"xano/session"
)

type A struct{}

func (a *A) Add(s session.Session, input *pb.AddRequest) error {
	var res int64
	for _, item := range input.Args {
		res += item
	}
	return s.RpcResponse("Add", &pb.AddResponse{
		Result: res,
	})
}

func main() {
	xano.WithConfig("./config/server.toml")

	xano.WithRoute(&router.RouterServer{
		Name:   "",
		Server: new(A),
	})

	xano.Run()
}
