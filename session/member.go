package session

import (
	"fmt"
	"reflect"
	"sync"
	"time"
	"xano/common"
	"xano/core"
	"xano/deal"
	"xano/logger"
	"xano/router"
)

// 用于管理member的session
type Member struct {
	sync.RWMutex
	// session集合
	Sessions map[uint64]*MemberSession
	// 当前SID
	SID uint64
}

var member *Member

func GetMember() *Member {
	if member == nil {
		member = new(Member)
		member.Sessions = make(map[uint64]*MemberSession)
	}
	return member
}

// session init
func (m *Member) SesssionInit(s *MemberSession) error {
	m.Lock()
	defer m.Unlock()

	var sid uint64
	// 遍历次数限制
	count := 0
	for {
		if count > 1 {
			err := fmt.Errorf("Session ID Not Enough")
			return err
		}
		if m.SID >= common.MaxSessionNum {
			m.SID = 0
			count++
		}
		m.SID++
		mchID := common.GetCache().GetUInt64(common.MchIDKey)
		sid = mchID*common.MaxSessionNum + m.SID
		// 一直寻找 直到找到合适的 一般情况下一次就可以找到
		if m.Sessions[sid] == nil {
			break
		}
	}
	s.SetSid(sid)
	m.Sessions[sid] = s
	return nil
}

// session close
func (m *Member) SessionClose(s *MemberSession) {
	m.Lock()
	defer m.Unlock()
	sid := s.GetSid()
	delete(m.Sessions, sid)
}

// session find
func (m *Member) SessionFindByID(sid uint64) *MemberSession {
	m.RLock()
	defer m.RUnlock()
	return m.Sessions[sid]
}

// 获取当前连接总数
func (m *Member) SessionCount() int {
	return len(m.Sessions)
}

// member session
type MemberSession struct {
	*BaseSession
}

func GetMemberSession(h *core.TcpHandle) *MemberSession {
	return &MemberSession{
		BaseSession: GetBaseSession(h),
	}
}

func (s *MemberSession) HandleRoute(r *router.Router, m *deal.Msg) error {
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

// rpc 请求
func (s *MemberSession) RpcRequest(ss *BaseSession, msg *deal.Msg) error {
	tcpAddr := router.GetMemberServerNode().GetNodeRand(msg.Route)
	if tcpAddr == "" {
		return fmt.Errorf("not found server: %s#%s", msg.Version, msg.Route)
	}

	pool := core.GetPool(tcpAddr)
	cli, err := pool.Get()
	if err != nil {
		return err
	}
	defer pool.Recycle(cli)

	c := make(chan struct{})
	cli.Client.SetHandleFunc(func(h *core.TcpHandle, m *deal.Msg) {
		if m.MsgType == common.MsgTypeResponse {
			m.Mid = ss.GetMid()
			ss.Send(m)
			c <- struct{}{}
		}
	})
	defer cli.Client.SetHandleFunc(nil)

	msg.Mid = cli.Client.GetMid()
	cli.Client.Send(msg)

	t := time.NewTimer(common.TcpDeadDuration)

	select {
	case <-c:
	case <-t.C:
		return fmt.Errorf("Rpc Request Timeout")
	}
	return nil
}

// notice 请求
func (s *MemberSession) Notice(route string, input interface{}) error {
	tcpAddr := router.GetMemberServerNode().GetNodeRand(route)
	if tcpAddr == "" {
		logger.Warnf("not found server: %s", route)
		return nil
	}
	inputBys, err := common.MsgMarsh(common.GetConfig().Base.TcpDeal, input)
	if err != nil {
		return err
	}
	inputMsg := &deal.Msg{
		Sid:     s.GetSid(),
		Route:   route,
		Mid:     0,
		MsgType: common.MsgTypePush,
		Deal:    common.GetConfig().Base.TcpDeal,
		Data:    inputBys,
		Version: common.GetConfig().Base.Version,
	}

	err = s.SendTo(tcpAddr, inputMsg)
	if err != nil {
		return err
	}
	common.PrintMsg(inputMsg, input)

	return nil
}
