package session

import (
	"fmt"
	"reflect"
	"xlq-server/common"
	"xlq-server/core"
	"xlq-server/deal"
)

type ServerSession struct {
	BaseSession
}

func NewServerSession(entity *core.TcpHandle) *ServerSession {
	s := new(ServerSession)
	s.TcpHandle = entity
	return s
}

func (s *ServerSession) Handle(packet *common.TcpPacket) error {
	var err error
	// 解析请求
	msg := new(deal.Msg)
	err = common.TcpMsgUnMarsh(packet.Data, msg)
	if err != nil {
		return err
	}

	// 获取服务
	componet := core.GetComponent(msg.Route)
	if componet == nil {
		return fmt.Errorf("not found service: " + msg.Route)
	}

	input := reflect.New(componet.Input)
	err = common.TcpMsgUnMarsh(msg.Data, input.Interface())
	if err != nil {
		return err
	}

	args := []reflect.Value{reflect.ValueOf(s), input}
	results := componet.Method.Func.Call(args)
	resLen := len(results)
	if resLen != 2 {
		return fmt.Errorf("error output param nums")
	}

	// 如果不需要回包
	if componet.Ouput == nil {
		return results[0].Interface().(error)
	}

	// 需要回包
	err = results[1].Interface().(error)
	if err != nil {
		return err
	}

	// 生成回包
	outputPacket, err := s.EncodePacket(msg.Route, msg.Mid, results[0].Interface())
	if err != nil {
		return err
	}

	s.Send(outputPacket)

	return nil
}
