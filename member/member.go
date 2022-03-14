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
	MasterClient *core.TcpClient
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
	// sid拦截
	if msg.Sid <= 0 {
		logger.Warn("Invaild Msg No Sid")
		return
	}

	// 从连接池中拿到连接转发出去即可，拿到response之后释放连接
	tcpAddr := router.GetMemberNode().GetNodeRand(msg.Route)
	if tcpAddr == "" {
		logger.Errorf("not found server: %s#%s", msg.Version, msg.Route)
		return
	}

	net := session.NewNetService()

	switch msg.MsgType {
	case common.MsgTypeRequest:
		resMsgs, err := net.Request(tcpAddr, msg)
		if err != nil {
			logger.Error(err)
			return
		}
		// 所有包全部返回
		for _, item := range resMsgs {
			item.Mid = h.GetMid()
			h.Send(item)
		}

	case common.MsgTypeNotice:
		err := net.Notice(tcpAddr, msg)
		if err != nil {
			logger.Error(err)
			return
		}

	default:
		// 搞不清楚的全部发回客户端
		msg.Mid = h.GetMid()
		h.Send(msg)
	}
}

// 初始逻辑
func (m *Member) tcpInit(h *core.TcpHandle) {
	// 获取sid
	if !session.GetConnect().IsEnough() {
		GetMemberMaster().memberSid()
	}
	sid, err := session.GetConnect().GetSid()
	if err != nil {
		logger.Error(err)
		return
	}
	ss := session.GetSession(h)
	ss.SID = sid
	ss.Push("SessionInit", &deal.SessionInitPush{
		Sid: ss.GetSid(),
	})
	GetMemberMaster().memberInfo()
}

// 关闭逻辑
func (m *Member) tcpClose(h *core.TcpHandle) {
	ss := session.GetSession(h)
	sid := ss.GetSid()
	session.GetConnect().DelSid(sid)
	err := ss.Notice("SessionClose", &deal.SessionCloseNotice{
		Sid: sid,
	})
	if err != nil {
		logger.Warn(err)
	}
	GetMemberMaster().memberInfo()
}
