package session

import (
	"fmt"
	"reflect"
	"xano/common"
	"xano/core"
	"xano/deal"
	"xano/logger"
	"xano/router"
)

type BaseSession struct {
	*core.TcpHandle
}

func GetBaseSession(h *core.TcpHandle) *BaseSession {
	v := h.Get(common.HandleKeySession)
	if v != nil {
		return v.(*BaseSession)
	}
	return &BaseSession{
		TcpHandle: h,
	}
}

func (b *BaseSession) SetSid(sid uint64) {
	b.Set(common.SessionIDKey, sid)
}

func (b *BaseSession) GetSid() uint64 {
	v := b.Get(common.SessionIDKey)
	if v == nil {
		return 0
	}
	return v.(uint64)
}

func (b *BaseSession) Rpc(route string, input, output interface{}) error {
	return fmt.Errorf("Not Support Rpc")
}

func (b *BaseSession) Notice(route string, input interface{}) error {
	inputBys, err := common.MsgMarsh(common.GetConfig().Base.TcpDeal, input)
	if err != nil {
		return err
	}
	inputMsg := &deal.Msg{
		Sid:     b.GetSid(),
		Route:   route,
		Mid:     b.GetMid(),
		MsgType: common.MsgTypeNotice,
		Deal:    common.GetConfig().Base.TcpDeal,
		Data:    inputBys,
		Version: common.GetConfig().Base.Version,
	}
	common.PrintMsg(inputMsg, input)
	b.Send(inputMsg)
	return nil
}

func (b *BaseSession) Response(route string, input interface{}) error {
	inputBys, err := common.MsgMarsh(common.GetConfig().Base.TcpDeal, input)
	if err != nil {
		return err
	}
	inputMsg := &deal.Msg{
		Sid:     b.GetSid(),
		Route:   route,
		Mid:     b.GetMid(),
		MsgType: common.MsgTypeResponse,
		Deal:    common.GetConfig().Base.TcpDeal,
		Data:    inputBys,
		Version: common.GetConfig().Base.Version,
	}
	common.PrintMsg(inputMsg, input)
	b.Send(inputMsg)
	return nil
}

func (b *BaseSession) RpcResponse(route string, input interface{}) error {
	return fmt.Errorf("Not Support Rpc Response")
}

func (b *BaseSession) Push(route string, input interface{}) error {
	inputBys, err := common.MsgMarsh(common.GetConfig().Base.TcpDeal, input)
	if err != nil {
		return err
	}
	inputMsg := &deal.Msg{
		Sid:     b.GetSid(),
		Route:   route,
		Mid:     b.GetMid(),
		MsgType: common.MsgTypePush,
		Deal:    common.GetConfig().Base.TcpDeal,
		Data:    inputBys,
		Version: common.GetConfig().Base.Version,
	}
	b.Send(inputMsg)
	common.PrintMsg(inputMsg, input)
	return nil
}

func (b *BaseSession) PushTo(sid uint64, route string, input interface{}) error {
	return fmt.Errorf("Not Support PushTo")
}

func (b *BaseSession) SendTo(addr string, msg *deal.Msg) error {
	pool := core.GetPool(addr)
	cli, err := pool.Get()
	if err != nil {
		logger.Error(err)
		return err
	}
	defer pool.Recycle(cli)
	cli.Client.SetHandleFunc(nil)
	msg.Mid = cli.Client.GetMid()
	cli.Client.Send(msg)
	return nil
}

func (b *BaseSession) HandleRoute(r *router.Router, m *deal.Msg) error {
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
	arg := []reflect.Value{reflect.ValueOf(b), reflect.ValueOf(input)}
	res := route.Method.Call(arg)

	if len(res) == 0 {
		return nil
	}
	if err := res[0].Interface(); err != nil {
		return fmt.Errorf("%+v", err)
	}
	return nil
}
