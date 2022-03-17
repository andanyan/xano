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

	/* 指定任意session 推送消息 因此sid作为客户端在系统内的唯一标志，业务层应与用户进行绑定
	s.PushTo(3000001, "Div", &pb.DivResponse{
		Result: res,
	})
	*/

	return s.Response("Div", &pb.DivResponse{
		Result: res,
	})
}

func main() {
	xano.WithConfig("./config/server1.toml")

	xano.WithRoute(&router.RouterServer{
		Name:   "",
		Server: new(B),
	})

	xano.Run()
}
