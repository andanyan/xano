package core

import (
	"fmt"
	"reflect"
	"time"
	"xano/common"
	"xano/deal"
	"xano/logger"
	"xano/router"
)

type Session struct {
	*TcpHandle
}

func GetSession(entity *TcpHandle) *Session {
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

// RPC
func (s *Session) Rpc(route string, input, output interface{}) error {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
		}
	}()
	if input == nil || output == nil {
		return fmt.Errorf("error input or output, not allow nil")
	}

	// 发送远程包
	msg, err := s.genMsg(route, common.MsgTypeRpc, input)
	if err != nil {
		return err
	}

	// 向网关主节点发送Rpc请求
	tcpAddr := s.Get(common.HandleKeyTcpAddr).(string)
	if tcpAddr == "" {
		return fmt.Errorf("has no server: %s", route)
	}

	// 获取连接
	pool := GetPool(tcpAddr)
	poolObj, err := pool.Get()
	if err != nil {
		return err
	}
	defer pool.Recycle(poolObj)

	// 发送包
	poolObj.Client.Send(msg)

	// 收到Response包才认为已完成、其他包直接发射回去即可
	c := make(chan struct{})
	poolObj.Client.SetHandle(func(h *TcpHandle, m *deal.Msg) {
		mmid := s.Get(common.HandleKeyMid).(uint64)
		if m.Mid == mmid {
			// 非Response的类型，直接返回给客户端
			if m.MsgType != common.MsgTypeResponse {
				s.Send(m)
				return
			}
			err := common.MsgUnMarsh(common.TcpDealProtobuf, m.Data, output)
			if err != nil {
				logger.Error(err.Error())
			}
			c <- struct{}{}

		} else if m.Mid > mmid {
			c <- struct{}{}
		}
	})
	defer poolObj.Client.SetHandle(nil)

	t := time.NewTimer(common.TcpDeadDuration)

	select {
	case <-c:

	case <-t.C:
		return fmt.Errorf("rpc timeout")
	}

	return nil
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

// 向连接写入包
func (s *Session) genMsg(route string, msgType uint32, input interface{}) (*deal.Msg, error) {
	// 组装包 写入连接即可
	inputBys, err := common.MsgMarsh(common.TcpDealProtobuf, input)
	if err != nil {
		return nil, err
	}
	msg := &deal.Msg{
		Route:   route,
		Mid:     s.Get(common.HandleKeyMid).(uint64),
		MsgType: msgType,
		Deal:    common.TcpDealProtobuf,
		Data:    inputBys,
	}
	logger.Infof("Route: %s, Mid: %d, MsgType: %d, deal: %d, date: [%+v]", msg.Route, msg.Mid, msg.MsgType, msg.Deal, input)
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

	logger.Infof("Route: %s, Mid: %d, MsgType: %d, deal: %d, date: [%+v]", m.Route, m.Mid, m.MsgType, m.Deal, input)

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
