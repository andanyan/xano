package member

import (
	"xano/common"
	"xano/core"
	"xano/deal"
	"xano/logger"
	"xano/router"
	"xano/session"
)

type Member struct {
}

func NewMember() *Member {
	return new(Member)
}

func (m *Member) Close() {
	GetMemberMaster().memberClose()
}

func (m *Member) Run() {
	// 注册回包服务
	router.GetMemberRouter().Register(&router.RouterServer{
		Name:   "",
		Server: new(Service),
	})

	// 启动与主节点的连接
	go GetMemberMaster().masterHandle()

	// 启动tcp服务
	go m.runTcp()

	// 启动内部访问地址
	go m.innerTcp()
}

// 运行tcp
func (m *Member) runTcp() {
	addr := common.GetConfig().Member.TcpAddr
	if addr == "" {
		return
	}
	logger.Infof("Gate Member TCP Server Start: %s", addr)
	core.NewTcpServer(addr, func(h *core.TcpHandle) {
		h.SetHandleFunc(m.tcpHandle)
	}, func(h *core.TcpHandle) {
		h.SetInitFunc(m.tcpInit)
	}, func(h *core.TcpHandle) {
		h.SetCloseFunc(m.tcpClose)
	})
}

// 转发逻辑
func (m *Member) tcpHandle(h *core.TcpHandle, msg *deal.Msg) {
	logger.Infof("%+v", msg)
	switch msg.MsgType {
	case common.MsgTypeRequest, common.MsgTypeNotice:
		ss := session.GetMemberSession(h)
		sid := ss.GetSid()
		tcpAddr := router.GetMemberServerNode().GetNodeRand(msg.Route)
		if tcpAddr == "" {
			logger.Errorf("not found server: %s#%s", msg.Version, msg.Route)
			return
		}
		msg.Sid = sid
		ss.SendTo(tcpAddr, msg)
	}
}

// 初始逻辑
func (m *Member) tcpInit(h *core.TcpHandle) {
	ss := session.GetMemberSession(h)
	err := session.GetMember().SesssionInit(ss)
	if err != nil {
		logger.Error(err)
		ss.Close()
		return
	}
	GetMemberMaster().memberInfo()
	logger.Debug("tcp session init: ", ss.GetSid())
}

// 关闭逻辑
func (m *Member) tcpClose(h *core.TcpHandle) {
	ss := session.GetMemberSession(h)
	session.GetMember().SessionClose(ss)
	GetMemberMaster().memberInfo()
	logger.Debug("tcp session close: ", ss.GetSid())
}

// **********************************
// 内部tcp
// **********************************
// 运行tcp
func (m *Member) innerTcp() {
	addr := common.GetConfig().Member.InnerAddr
	if addr == "" {
		return
	}
	logger.Infof("Gate Member Inner TCP Server Start: %s", addr)
	core.NewTcpServer(addr, func(h *core.TcpHandle) {
		h.SetHandleFunc(m.innerHandle)
	})
}

// 转发逻辑
func (m *Member) innerHandle(h *core.TcpHandle, msg *deal.Msg) {
	logger.Debugf("Inner %+v", msg)
	sid := msg.Sid
	if sid <= 0 {
		logger.Warnf("No Session Msg: %+v", msg)
		return
	}
	switch msg.MsgType {
	case common.MsgTypeRequest:
		// 同步请求 用于rpc
		ss := session.GetMember().SessionFindByID(sid)
		if ss == nil {
			logger.Errorf("Session Invaild: %d", sid)
			return
		}
		err := ss.RpcRequest(session.GetBaseSession(h), msg)
		if err != nil {
			logger.Error(err)
		}

	case common.MsgTypeNotice:
		// notice请求 用于事件通知
		ss := session.GetMember().SessionFindByID(sid)
		if ss == nil {
			logger.Errorf("Session Invaild: %d", sid)
			return
		}
		tcpAddr := router.GetMemberServerNode().GetNodeRand(msg.Route)
		if tcpAddr == "" {
			logger.Errorf("not found server: %s#%s", msg.Version, msg.Route)
			return
		}
		ss.SendTo(tcpAddr, msg)
	case common.MsgTypeResponse, common.MsgTypePush:
		ss := session.GetMember().SessionFindByID(sid)
		if ss == nil {
			logger.Errorf("Session Invaild: %d", sid)
			return
		}
		msg.Mid = ss.GetMid()
		ss.Send(msg)
	}
}
