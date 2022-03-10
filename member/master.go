package member

import (
	"time"
	"xano/common"
	"xano/core"
	"xano/deal"
	"xano/logger"
	"xano/router"
	"xano/session"
)

// 与主节点的通信
type MemberMaster struct {
	MasterClient *core.TcpClient
}

var memberMaster *MemberMaster

func GetMemberMaster() *MemberMaster {
	if memberMaster == nil {
		memberMaster = new(MemberMaster)
	}
	return memberMaster
}

// 与主节点通信
func (m *MemberMaster) masterHandle() {
	addr := common.GetConfig().Member.MasterAddr
	if addr == "" {
		return
	}

	// 与主节点建立连接
	cli, err := core.NewTcpClient(addr)
	if err != nil {
		logger.Fatal(err)
		return
	}
	defer cli.Close()
	cli.SetHandleFunc(func(h *core.TcpHandle, m *deal.Msg) {
		ss := session.GetSession(h)
		ss.SID = m.Sid
		if err := ss.HandleRoute(router.GetMemberRouter(), m); err != nil {
			logger.Error(err.Error())
		}
	})
	cli.SetCloseFunc(func(h *core.TcpHandle) {
		// 连接断开，即服务停止
		logger.Fatal("MASTER DISCONNECT")
	})
	m.MasterClient = cli

	// 发起Start通信
	m.memberStart()

	// 获取sid
	m.memberSid()

	// 启动心跳
	for {
		m.memberHeart()
		time.Sleep(common.TcpHeartDuration)
	}
}

// notice master member start
func (m *MemberMaster) memberStart() {
	memberAddr, err := common.ParseAddr(common.GetConfig().Member.TcpAddr)
	if err != nil {
		logger.Fatal(err)
	}
	input := &deal.MemberStartRequest{
		Version: common.GetConfig().Base.Version,
		Port:    memberAddr.Port,
	}
	inputBys, err := common.MsgMarsh(common.GetConfig().Base.TcpDeal, input)
	if err != nil {
		logger.Error(err)
		return
	}
	msg := &deal.Msg{
		Route:   "MemberStart",
		Sid:     0,
		Mid:     m.MasterClient.GetMid(),
		MsgType: common.MsgTypeNotice,
		Deal:    common.GetConfig().Base.TcpDeal,
		Data:    inputBys,
		Version: common.GetConfig().Base.Version,
	}
	m.MasterClient.Send(msg)
}

// notice master member close
func (m *MemberMaster) memberClose() {
	if m.MasterClient == nil {
		return
	}
	input := &deal.MemberStopNotice{}
	inputBys, err := common.MsgMarsh(common.GetConfig().Base.TcpDeal, input)
	if err != nil {
		logger.Error(err)
		return
	}
	msg := &deal.Msg{
		Route:   "MemberStop",
		Sid:     0,
		Mid:     m.MasterClient.GetMid(),
		MsgType: common.MsgTypeNotice,
		Deal:    common.GetConfig().Base.TcpDeal,
		Data:    inputBys,
		Version: common.GetConfig().Base.Version,
	}
	m.MasterClient.Send(msg)
}

// heart master
func (m *MemberMaster) memberHeart() {
	input := &deal.Ping{
		Psutil: common.GetPsutil(),
	}
	inputBys, err := common.MsgMarsh(common.GetConfig().Base.TcpDeal, input)
	if err != nil {
		logger.Error(err)
		return
	}
	msg := &deal.Msg{
		Route:   "MemberHeart",
		Sid:     0,
		Mid:     m.MasterClient.GetMid(),
		MsgType: common.MsgTypeRequest,
		Deal:    common.GetConfig().Base.TcpDeal,
		Data:    inputBys,
		Version: common.GetConfig().Base.Version,
	}
	m.MasterClient.Send(msg)
	// 每次心跳是上报
	GetMemberMaster().memberInfo()
}

// 发起sid取值
func (m *MemberMaster) memberSid() {
	if m.MasterClient == nil {
		return
	}
	input := &deal.MemberGetSidRequest{}
	inputBys, err := common.MsgMarsh(common.GetConfig().Base.TcpDeal, input)
	if err != nil {
		logger.Error(err)
		return
	}
	msg := &deal.Msg{
		Route:   "MemberGetSid",
		Sid:     0,
		Mid:     m.MasterClient.GetMid(),
		MsgType: common.MsgTypeNotice,
		Deal:    common.GetConfig().Base.TcpDeal,
		Data:    inputBys,
		Version: common.GetConfig().Base.Version,
	}
	m.MasterClient.Send(msg)
}

// 同步session
func (m *MemberMaster) memberInfo() {
	if m.MasterClient == nil {
		return
	}
	input := &deal.MemberInfoNotice{
		SessionCount: session.GetConnect().GetCount(),
	}
	inputBys, err := common.MsgMarsh(common.GetConfig().Base.TcpDeal, input)
	if err != nil {
		logger.Error(err)
		return
	}
	msg := &deal.Msg{
		Route:   "MemberInfo",
		Sid:     0,
		Mid:     m.MasterClient.GetMid(),
		MsgType: common.MsgTypeNotice,
		Deal:    common.GetConfig().Base.TcpDeal,
		Data:    inputBys,
		Version: common.GetConfig().Base.Version,
	}
	m.MasterClient.Send(msg)
}
