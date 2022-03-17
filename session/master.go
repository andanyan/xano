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

type MasterSession struct {
	*BaseSession
}

func GetMasterSession(h *core.TcpHandle) *MasterSession {
	return &MasterSession{
		BaseSession: GetBaseSession(h),
	}
}

func (s *MasterSession) HandleRoute(r *router.Router, m *deal.Msg) error {
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

	// 缓存最后一次协议类型
	s.Set(common.MessageDeal, m.Deal)

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
