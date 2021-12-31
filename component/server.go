package component

import (
	"fmt"
	"reflect"
	"xlq-server/common"
)

// 处理请求
func DoneMsg(s common.Session, msg *common.Msg) error {
	var err error
	// 解析路由
	entry, ok := service.Routes[msg.Route]
	if !ok {
		return fmt.Errorf("not found services: %s", "aaa")
	}

	// 执行中间件
	if err := s.Middlewares(msg.Route); err != nil {
		return err
	}

	// 输入数据参数化
	input := reflect.New(entry.Input.Elem()).Interface()
	err = s.MsgUnMarsh(msg.Data, input)
	if err != nil {
		return fmt.Errorf("error input data: %s", entry.Input.Name())
	}

	// 调用函数 session, input
	args := []reflect.Value{reflect.ValueOf(s), reflect.ValueOf(input)}
	result := entry.Method.Func.Call(args)

	if len(result) != 2 {
		return fmt.Errorf("error out data: %s", msg.Route)
	}

	err = result[1].Interface().(error)
	if err != nil {
		return fmt.Errorf("error out data: %s", err.Error())
	}

	// 依次写入session
	s.Write(msg.Route, result[0].Interface())
	return nil
}
