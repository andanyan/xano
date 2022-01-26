package server

import (
	"fmt"
	"log"
	"reflect"
	"xlq-server/common"
	"xlq-server/core"
	"xlq-server/deal"
	"xlq-server/router"
)

// 微服务对象,可独立运行tcp服务
type Server struct{}

func NewServer() *Server {
	return new(Server)
}

func (s *Server) Run() {
	sConf := common.GetConfig().Server
	addr := common.GenAddr(sConf.Host, sConf.Port)
	if addr == "" {
		return
	}
	log.Printf("Server Start: %s \n", addr)
	core.NewTcpServer(addr, s.handle)
}

func (s *Server) handle(h *core.TcpHandle, p *common.Packet) {
	// 解析packet
	var err error
	msg := new(deal.Msg)
	err = common.MsgUnMarsh(common.TcpDealProtobuf, p.Data, msg)
	if err != nil {
		log.Panicln(err)
		return
	}

	switch msg.MsgType {
	case common.MsgTypeNotice, common.MsgTypeRequest, common.MsgTypeRpc:
		// 设置当前mid
		h.Set(common.HandleKeyMid, msg.Mid)
		// 设定当前的来源地址
		h.Set(common.HandleKeyTcpAddr, h.GetAddr())

		// 创建session 提供给接口端使用
		ss := GetSession(h)
		// 调用路由
		if err = s.handleRoute(ss, msg); err != nil {
			log.Println(err)
		}

	case common.MsgTypePush:
		h.Send(p)

	default:

	}
}

func (s *Server) handleRoute(ss *Session, msg *deal.Msg) error {
	// 获取路由
	route := router.LocalRouter.GetRoute(msg.Route)
	if route == nil {
		return fmt.Errorf("error route " + msg.Route)
	}

	// 解析输入
	input := reflect.New(route.Input).Interface()
	err := common.MsgUnMarsh(msg.Deal, msg.Data, input)
	if err != nil {
		return err
	}

	// 调用函数
	arg := []reflect.Value{reflect.ValueOf(ss), reflect.ValueOf(input)}
	res := route.Method.Call(arg)

	if len(res) == 0 {
		return nil
	}
	if err := res[0].Interface(); err != nil {
		return fmt.Errorf("%+v", err)
	}
	return nil
}
