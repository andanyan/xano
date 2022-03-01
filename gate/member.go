package gate

import (
	"time"
	"xano/common"
	"xano/core"
	"xano/deal"
	"xano/logger"
	"xano/router"
)

type Member struct{}

func NewMember() *Member {
	return new(Member)
}

func (m *Member) Close() {
	return
}

func (m *Member) Run() {
	gConf := common.GetConfig().GateMember
	addr := gConf.Host + ":" + gConf.Port
	if addr == "" {
		return
	}

	// 注册回包服务
	router.GetMemberRouter().Register(&router.RouterServer{
		Name:   "",
		Server: new(MemberServer),
	})

	// 处理与主节点的通信
	go m.masterHandle()

	// 启动服务
	logger.Infof("Gate Member Start: %s", addr)
	go core.NewTcpServer(addr, m.handle)
}

// 转发逻辑
func (m *Member) handle(h *core.TcpHandle, msg *deal.Msg) {
	// 从连接池中拿到连接转发出去即可，拿到response之后释放连接
	tcpAddr := router.GetGateInfo().GetNodeAddr(msg.Route, msg.Version)
	if tcpAddr == "" {
		logger.Errorf("not found server: %s#%s", msg.Version, msg.Route)
		return
	}

	pool := core.GetPool(tcpAddr)
	cli, err := pool.Get()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	defer pool.Recycle(cli)

	curMid := msg.Mid

	c := make(chan struct{})
	cli.Client.SetHandle(func(_ *core.TcpHandle, rm *deal.Msg) {
		rm.Mid = curMid
		h.Send(rm)
		if rm.MsgType == common.MsgTypeResponse {
			c <- struct{}{}
		}
	})
	defer cli.Client.SetHandle(nil)

	// 发送消息
	msg.Mid = cli.Client.GetMid()
	cli.Client.Send(msg)
	<-c
}

// 与主节点通信
func (m *Member) masterHandle() {
	gConf := common.GetConfig().GateMember
	if gConf.MasterAddr == "" {
		return
	}

	cli, err := core.NewTcpClient(gConf.MasterAddr)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	cli.SetHandle(func(h *core.TcpHandle, m *deal.Msg) {
		ss := core.GetSession(h)
		if err := ss.HandleRoute(router.GetMemberRouter(), m); err != nil {
			logger.Error(err.Error())
		}
	})

	input := &deal.AllNodeRequest{
		Version: common.GetConfig().Base.Version,
	}
	inputBys, err := common.MsgMarsh(common.TcpDealProtobuf, input)
	if err != nil {
		return
	}

	// 定时获取网关接口
	time.Sleep(common.DelayDuration)
	for {
		msg := &deal.Msg{
			Route:   "AllNode",
			Mid:     cli.GetMid(),
			MsgType: common.MsgTypeRequest,
			Deal:    common.TcpDealProtobuf,
			Data:    inputBys,
			Version: common.GetConfig().Base.Version,
		}
		cli.Send(msg)
		time.Sleep(common.TcpHeartDuration)
	}
}
