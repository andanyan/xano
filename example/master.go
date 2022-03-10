package main

import (
	"xano"
	"xano/example/pb"
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

func main() {
	xano.WithConfig("./config/master.toml")

	xano.WithRoute(&router.RouterServer{
		Name:   "",
		Server: new(A),
	})

	xano.Run()
}
