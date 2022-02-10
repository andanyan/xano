package gate

import (
	"time"
	"xlq-server/common"
	"xlq-server/core"
	"xlq-server/deal"
	"xlq-server/logger"
	"xlq-server/router"
)

type Member struct{}

func NewMember() *Member {
	return new(Member)
}

func (m *Member) Run() {
	gConf := common.GetConfig().GateMember
	addr := gConf.Host + ":" + gConf.Port
	if addr == "" {
		return
	}

	// 注册回包服务
	router.MasterRouter.Register(&router.RouterServer{
		Name:   "",
		Server: new(MemberServer),
	})

	// 处理与主节点的通信
	go m.masterHandle()

	// 启动服务
	logger.Infof("Gate Member Start: %s", addr)
	core.NewTcpServer(addr, m.handle)
}

// 转发逻辑
func (m *Member) handle(h *core.TcpHandle, msg *deal.Msg) {
	// 从连接池中拿到连接转发出去即可，拿到response之后释放连接
	tcpAddr := router.GetGateInfo().GetNodeAddr(msg.Route, msg.Version)
	if tcpAddr == "" {
		logger.Error("not found server:", msg.Version, msg.Route)
		return
	}

	pool := core.GetPool(tcpAddr)
	cli, err := pool.Get()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	defer pool.Recycle(cli)
	c := make(chan struct{})
	cli.Client.SetHandle(func(ch *core.TcpHandle, cm *deal.Msg) {
		ch.Send(cm)
		if cm.MsgType != common.MsgTypeResponse {
			c <- struct{}{}
		}
	})
	defer cli.Client.SetHandle(nil)
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
		if err := ss.HandleRoute(router.MemberRouter, m); err != nil {
			logger.Error(err.Error())
		}
	})

	input := &deal.AllNodeRequest{
		Version: common.GetConfig().Base.Version,
	}
	inputBys, err := common.MsgMarsh(common.TcpDealProtobuf, input)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	msg := &deal.Msg{
		Route:   "AllNode",
		Mid:     cli.GetMid(),
		MsgType: common.MsgTypeRequest,
		Deal:    common.TcpDealProtobuf,
		Data:    inputBys,
		Version: "",
	}

	// 定时获取网关接口
	for {
		cli.Send(msg)
		time.Sleep(common.TcpHeartDuration)
	}
}
