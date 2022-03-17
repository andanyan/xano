package session

import (
	"fmt"
	"reflect"
	"time"
	"xano/common"
	"xano/core"
	"xano/deal"
	"xano/logger"
	"xano/router"
)

type ServerSession struct {
	*BaseSession
}

func GetServerSession(h *core.TcpHandle) *ServerSession {
	return &ServerSession{
		BaseSession: GetBaseSession(h),
	}
}

func (s *ServerSession) Rpc(route string, input, output interface{}) error {
	if input == nil || output == nil {
		return fmt.Errorf("input or output is null")
	}

	// 组装消息
	inputBys, err := common.MsgMarsh(common.GetConfig().Base.TcpDeal, input)
	if err != nil {
		return err
	}
	msg := &deal.Msg{
		Route:   route,
		Sid:     s.GetSid(),
		Mid:     0, //会在下个环节重新赋值
		MsgType: common.MsgTypeRequest,
		Deal:    common.GetConfig().Base.TcpDeal,
		Version: common.GetConfig().Base.Version,
		Data:    inputBys,
	}

	// 向网关主节点发送Rpc请求
	rpcNode, err := router.GetLocalMemberNode().GetNodeBySid(s.GetSid())
	if err != nil {
		return fmt.Errorf("has no member node valid")
	}

	// 获取连接
	pool := core.GetPool(rpcNode.InnerAddr)
	cli, err := pool.Get()
	if err != nil {
		logger.Error(err)
		return err
	}
	defer pool.Recycle(cli)

	c := make(chan struct{})
	cli.Client.SetHandleFunc(func(h *core.TcpHandle, m *deal.Msg) {
		if m.MsgType == common.MsgTypeResponse {
			err := common.MsgUnMarsh(common.GetConfig().Base.TcpDeal, m.Data, output)
			if err != nil {
				logger.Error(err)
			}
			c <- struct{}{}
		}
	})
	defer cli.Client.SetHandleFunc(nil)

	// 发送消息
	msg.Mid = cli.Client.GetMid()
	cli.Client.Send(msg)

	common.PrintMsg(msg, input)

	t := time.NewTimer(common.TcpDeadDuration)

	select {
	case <-c:
		// 正常返回
	case <-t.C:
		return fmt.Errorf("rpc timeout")
	}

	return nil
}

func (s *ServerSession) Notice(route string, input interface{}) error {
	// 组装消息
	inputBys, err := common.MsgMarsh(common.GetConfig().Base.TcpDeal, input)
	if err != nil {
		return err
	}
	msg := &deal.Msg{
		Route:   route,
		Sid:     s.GetSid(),
		Mid:     0, //会在下个环节重新赋值
		MsgType: common.MsgTypeNotice,
		Deal:    common.GetConfig().Base.TcpDeal,
		Version: common.GetConfig().Base.Version,
		Data:    inputBys,
	}

	// 向网关主节点发送Rpc请求
	rpcNode, err := router.GetLocalMemberNode().GetNodeBySid(s.GetSid())
	if err != nil {
		return fmt.Errorf("has no member node valid")
	}
	tcpAddr := rpcNode.InnerAddr

	err = s.SendTo(tcpAddr, msg)
	if err != nil {
		return nil
	}
	common.PrintMsg(msg, input)

	return err
}

func (s *ServerSession) Response(route string, input interface{}) error {
	// 组装消息
	inputBys, err := common.MsgMarsh(common.GetConfig().Base.TcpDeal, input)
	if err != nil {
		return err
	}
	msg := &deal.Msg{
		Route:   route,
		Sid:     s.GetSid(),
		Mid:     0, //会在下个环节重新赋值
		MsgType: common.MsgTypeResponse,
		Deal:    common.GetConfig().Base.TcpDeal,
		Version: common.GetConfig().Base.Version,
		Data:    inputBys,
	}

	// 向网关主节点发送Rpc请求
	rpcNode, err := router.GetLocalMemberNode().GetNodeBySid(s.GetSid())
	if err != nil {
		return fmt.Errorf("has no member node valid")
	}
	tcpAddr := rpcNode.InnerAddr

	err = s.SendTo(tcpAddr, msg)
	if err != nil {
		return nil
	}
	common.PrintMsg(msg, input)

	return err
}

func (s *ServerSession) RpcResponse(route string, input interface{}) error {
	// 组装消息
	inputBys, err := common.MsgMarsh(common.GetConfig().Base.TcpDeal, input)
	if err != nil {
		return err
	}
	msg := &deal.Msg{
		Route:   route,
		Sid:     s.GetSid(),
		Mid:     s.GetMid(), //会在下个环节重新赋值
		MsgType: common.MsgTypeResponse,
		Deal:    common.GetConfig().Base.TcpDeal,
		Version: common.GetConfig().Base.Version,
		Data:    inputBys,
	}

	// 向网关主节点发送Rpc请求
	s.Send(msg)
	common.PrintMsg(msg, input)
	return nil
}

func (s *ServerSession) Push(route string, input interface{}) error {
	// 组装消息
	inputBys, err := common.MsgMarsh(common.GetConfig().Base.TcpDeal, input)
	if err != nil {
		return err
	}
	msg := &deal.Msg{
		Route:   route,
		Sid:     s.GetSid(),
		Mid:     0,
		MsgType: common.MsgTypePush,
		Deal:    common.GetConfig().Base.TcpDeal,
		Version: common.GetConfig().Base.Version,
		Data:    inputBys,
	}

	// 向网关主节点发送Rpc请求
	rpcNode, err := router.GetLocalMemberNode().GetNodeBySid(s.GetSid())
	if err != nil {
		return fmt.Errorf("has no member node valid")
	}
	tcpAddr := rpcNode.InnerAddr

	err = s.SendTo(tcpAddr, msg)
	if err != nil {
		return nil
	}
	common.PrintMsg(msg, input)

	return err
}

func (s *ServerSession) PushTo(sid uint64, route string, input interface{}) error {
	// 组装消息
	inputBys, err := common.MsgMarsh(common.GetConfig().Base.TcpDeal, input)
	if err != nil {
		return err
	}
	msg := &deal.Msg{
		Route:   route,
		Sid:     sid,
		Mid:     0,
		MsgType: common.MsgTypePush,
		Deal:    common.GetConfig().Base.TcpDeal,
		Version: common.GetConfig().Base.Version,
		Data:    inputBys,
	}

	// 向网关主节点发送Rpc请求
	rpcNode, err := router.GetLocalMemberNode().GetNodeBySid(sid)
	if err != nil {
		return fmt.Errorf("has no member node valid")
	}
	tcpAddr := rpcNode.InnerAddr

	err = s.SendTo(tcpAddr, msg)
	if err != nil {
		return nil
	}
	common.PrintMsg(msg, input)

	return err
}

func (s *ServerSession) HandleRoute(r *router.Router, m *deal.Msg) error {
	// 获取路由
	route := r.GetRoute(m.Route)
	if route == nil {
		return fmt.Errorf("error route " + m.Route)
	}

	// 解析输入
	input := reflect.New(route.Input.Elem()).Interface()
	err := common.MsgUnMarsh(m.Deal, m.Data, input)
	if err != nil {
		logger.Error(err)
		return err
	}
	common.PrintMsg(m, input)

	// 调用函数
	arg := []reflect.Value{reflect.ValueOf(s), reflect.ValueOf(input)}
	res := route.Method.Call(arg)

	if len(res) == 0 {
		return nil
	}
	if err := res[0].Interface(); err != nil {
		return fmt.Errorf("%+v", err)
	}
	return nil
}
