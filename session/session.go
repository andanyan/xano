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

type Session struct {
	*core.TcpHandle
	SID uint64
}

// 获取Session
func GetSession(entity *core.TcpHandle) *Session {
	ss := entity.Get(common.HandleKeySession)
	if ss != nil {
		return ss.(*Session)
	}
	ns := &Session{
		TcpHandle: entity,
	}
	entity.Set(common.HandleKeySession, ns)
	return ns
}

// 获取sid
func (s *Session) GetSid() uint64 {
	return s.SID
}

// RPC
func (s *Session) Rpc(route string, input, output interface{}) error {
	// 向网关主节点发送Rpc请求
	tcpAddr := router.GetLocalNode().GetNodeRand(route)
	//tcpAddr := router.GetMemberInfo().GetNodeRand()
	if tcpAddr == "" {
		return fmt.Errorf("has no member node valid")
	}
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

	net := new(NetService)

	resMsgs, err := net.Request(tcpAddr, msg)
	if err != nil {
		return err
	}

	// 遍历resMsg 非Response数据，直接返回给客户端
	for _, item := range resMsgs {
		if item.MsgType != common.MsgTypeResponse {
			item.Mid = s.GetMid()
			s.Send(item)
		} else {
			err := common.MsgUnMarsh(common.GetConfig().Base.TcpDeal, item.Data, output)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// NOTICE
func (s *Session) Notice(route string, input interface{}) error {
	// 向网关主节点发送Rpc请求
	tcpAddr := router.GetLocalNode().GetNodeRand(route)
	if tcpAddr == "" {
		return fmt.Errorf("has no member node valid")
	}

	net := new(NetService)

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

	return net.Notice(tcpAddr, msg)
}

// Response
func (s *Session) Response(route string, input interface{}) error {
	msg, err := s.genMsg(route, common.MsgTypeResponse, input)
	if err != nil {
		return err
	}
	s.Send(msg)
	return nil
}

// Push
func (s *Session) Push(route string, input interface{}) error {
	// 组装包 写入连接即可
	msg, err := s.genMsg(route, common.MsgTypePush, input)
	if err != nil {
		return err
	}
	s.Send(msg)
	return nil
}

// Push To Other Session
func (s *Session) PushSession(route string, sid uint64, input interface{}) error {
	memberNode, err := router.GetLocalMemberNode().GetNodeBySid(sid)
	if err != nil {
		return err
	}

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

	// 获取连接
	pool := core.GetPool(memberNode.Addr)
	cli, err := pool.Get()
	if err != nil {
		return err
	}
	defer pool.Recycle(cli)

	// 发送消息
	cli.Client.Send(msg)
	return nil
}

// 向连接写入包
func (s *Session) genMsg(route string, msgType uint32, input interface{}) (*deal.Msg, error) {
	// 组装包 写入连接即可
	inputBys, err := common.MsgMarsh(common.GetConfig().Base.TcpDeal, input)
	if err != nil {
		return nil, err
	}
	msg := &deal.Msg{
		Route:   route,
		Sid:     s.GetSid(),
		Mid:     s.GetMid(),
		MsgType: msgType,
		Deal:    common.GetConfig().Base.TcpDeal,
		Version: common.GetConfig().Base.Version,
		Data:    inputBys,
	}
	return msg, nil
}

// 处理路由
func (s *Session) HandleRoute(r *router.Router, m *deal.Msg) error {
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
