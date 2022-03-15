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
		//h.SetInitFunc(m.tcpInit)
	}, func(h *core.TcpHandle) {
		h.SetCloseFunc(m.tcpClose)
	})
}

// 转发逻辑
func (m *Member) tcpHandle(h *core.TcpHandle, msg *deal.Msg) {
	// 获取当前session
	ss := session.GetSession(h)
	sid := ss.GetSid()
	// 如果mid为0,sid不为0且和当前不等,消息类型为推送 意味着是内部转发包
	if msg.Mid == 0 && msg.Sid > 0 && msg.Sid != sid && msg.MsgType == common.MsgTypePush {
		ss := session.GetMember().SessionFindByID(msg.Sid)
		msg.Mid = ss.GetMid()
		ss.Send(msg)
		return
	}

	// 如果是客户端初始化请求
	if msg.Route == common.SessionInitKey {
		// session初始化
		m.tcpInit(h)
	}

	// 赋值Sid
	msg.Sid = sid

	// 正常客户端消息处理
	switch msg.MsgType {
	case common.MsgTypeRequest:
		// 从连接池中拿到连接转发出去即可，拿到response之后释放连接
		tcpAddr := router.GetMemberNode().GetNodeRand(msg.Route)
		if tcpAddr == "" {
			logger.Errorf("not found server: %s#%s", msg.Version, msg.Route)
			return
		}
		net := session.NewNetService()
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
		// 从连接池中拿到连接转发出去即可，拿到response之后释放连接
		tcpAddr := router.GetMemberNode().GetNodeRand(msg.Route)
		if tcpAddr == "" {
			logger.Errorf("not found server: %s#%s", msg.Version, msg.Route)
			return
		}
		net := session.NewNetService()
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
	ss := session.GetSession(h)
	sid := ss.GetSid()
	if sid > 0 {
		return
	}
	session.GetMember().SesssionInit(ss)
	GetMemberMaster().memberInfo()
	ss.Response("SessionInit", &deal.SessionInitResponse{
		Sid: sid,
	})
	logger.Debug("tcp session init: ", sid)
}

// 关闭逻辑
func (m *Member) tcpClose(h *core.TcpHandle) {
	ss := session.GetSession(h)
	session.GetMember().SessionClose(ss)
	GetMemberMaster().memberInfo()
	logger.Debug("tcp session close: ", ss.GetSid())
}
