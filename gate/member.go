package gate

import (
	"time"
	"xano/common"
	"xano/core"
	"xano/deal"
	"xano/logger"
	"xano/router"
)

type Member struct {
	MasterClient *core.TcpClient
}

func NewMember() *Member {
	return new(Member)
}

func (m *Member) Close() {
	m.memberClose()
}

func (m *Member) Run() {
	gConf := common.GetConfig().GateMember
	addr := common.GenAddr(gConf.Host, gConf.Port)
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
	router := router.GetGateInfo()
	tcpAddr := router.GetNodeRand(msg.Route)
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

	c := make(chan struct{})
	cli.Client.SetHandle(func(_ *core.TcpHandle, rm *deal.Msg) {
		nMid := h.GetMid()
		rm.Mid = nMid
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

	// 与主节点建立连接
	cli, err := core.NewTcpClient(gConf.MasterAddr)
	if err != nil {
		logger.Fatal(err)
		return
	}
	defer cli.Close()
	cli.SetHandle(func(h *core.TcpHandle, m *deal.Msg) {
		ss := core.GetSession(h)
		if err := ss.HandleRoute(router.GetMemberRouter(), m); err != nil {
			logger.Error(err.Error())
		}
	})
	cli.SetCloseFunc(func() {
		logger.Fatal("Master Disconnected!")
	})
	m.MasterClient = cli

	// 发起Start通信
	m.memberStart()

	// 启动心跳
	for {
		time.Sleep(common.TcpHeartDuration)
		m.memberHeart()
	}
}

// notice master member start
func (m *Member) memberStart() {
	input := &deal.MemberStartNotice{
		Version: common.GetConfig().Base.Version,
		Port:    common.GetConfig().GateMember.Port,
	}
	inputBys, err := common.MsgMarsh(common.TcpDealProtobuf, input)
	if err != nil {
		logger.Error(err)
		return
	}
	msg := &deal.Msg{
		Route:   "MemberStart",
		Mid:     m.MasterClient.GetMid(),
		MsgType: common.MsgTypeNotice,
		Deal:    common.TcpDealProtobuf,
		Data:    inputBys,
		Version: common.GetConfig().Base.Version,
	}
	m.MasterClient.Send(msg)
}

// notice master member close
func (m *Member) memberClose() {
	if m.MasterClient == nil {
		return
	}
	input := &deal.MemberStopNotice{}
	inputBys, err := common.MsgMarsh(common.TcpDealProtobuf, input)
	if err != nil {
		logger.Error(err)
		return
	}
	msg := &deal.Msg{
		Route:   "MemberStop",
		Mid:     m.MasterClient.GetMid(),
		MsgType: common.MsgTypeNotice,
		Deal:    common.TcpDealProtobuf,
		Data:    inputBys,
		Version: common.GetConfig().Base.Version,
	}
	m.MasterClient.Send(msg)
}

// heart master
func (m *Member) memberHeart() {
	input := &deal.Ping{}
	inputBys, err := common.MsgMarsh(common.TcpDealProtobuf, input)
	if err != nil {
		logger.Error(err)
		return
	}
	msg := &deal.Msg{
		Route:   "MemberHeart",
		Mid:     m.MasterClient.GetMid(),
		MsgType: common.MsgTypeRequest,
		Deal:    common.TcpDealProtobuf,
		Data:    inputBys,
		Version: common.GetConfig().Base.Version,
	}
	m.MasterClient.Send(msg)
}
