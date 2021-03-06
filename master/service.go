package master

import (
	"sync"
	"time"
	"xano/common"
	"xano/deal"
	"xano/router"
	"xano/session"
)

// 主节点的服务对象
type MasterServer struct {
	Mux            sync.Mutex
	MemberSessions []session.Session
	ServerSessions []session.Session
	MemberMchID    uint64
}

// member start
func (s *MasterServer) MemberStart(ss session.Session, input *deal.MemberStartRequest) error {
	addr, err := common.ParseAddr(ss.GetAddr())
	if err != nil {
		return err
	}

	s.Mux.Lock()
	defer s.Mux.Unlock()

	mchID := s.MemberMchID
	s.MemberMchID++

	newNode := &deal.MemberNode{
		MchId:        mchID,
		Version:      input.Version,
		Addr:         addr.Ip + ":" + input.Port,
		LastTime:     time.Now().Unix(),
		SessionCount: 0,
		InnerAddr:    addr.Ip + ":" + input.InnerPort,
	}

	router.GetMasterNode().AddMemberNode(newNode)

	ss.Set(common.MemberNode, newNode)

	s.MemberSessions = append(s.MemberSessions, ss)

	// 告知所有server member节点发生变更
	memberNode := router.GetMasterNode().AllMemberNode()
	for _, item := range s.ServerSessions {
		item.Push("MemberNode", &deal.MemberNodePush{
			Node: memberNode,
		})
	}

	// 告知当前节点mid
	ss.Push("MemberMchID", &deal.MemberMchIDPush{
		MchID: mchID,
	})

	// 告知当前从节点 当前所有的服务节点信息
	serverNodes := router.GetMasterNode().AllServerNode()
	return ss.Response("MemberStart", &deal.MemberStartResponse{
		Node: serverNodes,
	})
}

// member stop
func (s *MasterServer) MemberStop(ss session.Session, input *deal.MemberStopNotice) error {
	node := ss.Get(common.MemberNode).(*deal.MemberNode)

	router.GetMasterNode().RemoveMemberNode(node.Addr)

	s.Mux.Lock()
	defer s.Mux.Unlock()

	index := 0
	for _, item := range s.MemberSessions {
		if item != ss {
			s.MemberSessions[index] = item
			index++
		}
	}
	s.MemberSessions = s.MemberSessions[:index]

	// 告知所有server member节点发生变更
	memberNode := router.GetMasterNode().AllMemberNode()
	for _, item := range s.ServerSessions {
		item.Push("MemberNode", &deal.MemberNodePush{
			Node: memberNode,
		})
	}

	return nil
}

// member ping
func (s *MasterServer) MemberHeart(ss session.Session, input *deal.Ping) error {
	node := ss.Get(common.MemberNode).(*deal.MemberNode)
	node.LastTime = time.Now().Unix()
	node.Psutil = input.Psutil
	return ss.Response("MemberHeart", &deal.Pong{
		Pong: node.LastTime,
	})
}

// member session
func (s *MasterServer) MemberInfo(ss session.Session, input *deal.MemberInfoNotice) error {
	node := ss.Get(common.MemberNode).(*deal.MemberNode)
	node.SessionCount = input.SessionCount
	return nil
}

// server start
func (s *MasterServer) ServerStart(ss session.Session, input *deal.ServerStartRequest) error {
	addr, err := common.ParseAddr(ss.GetAddr())
	if err != nil {
		return err
	}

	newNode := &deal.ServerNode{
		Version:  input.Version,
		Addr:     addr.Ip + ":" + input.Port,
		LastTime: time.Now().Unix(),
		Routes:   input.Routes,
	}

	router.GetMasterNode().AddServerNode(newNode)

	ss.Set(common.ServerNode, newNode)

	s.Mux.Lock()
	defer s.Mux.Unlock()
	s.ServerSessions = append(s.ServerSessions, ss)

	serverNodes := router.GetMasterNode().AllServerNode()
	// 通知所有的从节点 服务节点信息更新
	for _, item := range s.MemberSessions {
		item.Push("MemberStart", &deal.MemberStartResponse{
			Node: serverNodes,
		})
	}

	// 告知服务节点所有的网关节点
	memberNodes := router.GetMasterNode().AllMemberNode()
	return ss.Response("ServerStart", &deal.ServerStartResponse{
		Node: memberNodes,
	})
}

// member stop
func (s *MasterServer) ServerStop(ss session.Session, input *deal.ServerStopNotice) error {
	node := ss.Get(common.ServerNode).(*deal.ServerNode)

	router.GetMasterNode().RemoveServerNode(node.Addr)

	s.Mux.Lock()
	defer s.Mux.Unlock()

	index := 0
	for _, item := range s.ServerSessions {
		if item != ss {
			s.ServerSessions[index] = item
			index++
		}
	}
	s.ServerSessions = s.ServerSessions[:index]

	// 节点数据
	serverNodes := router.GetMasterNode().AllServerNode()
	// 通知所有的从节点 服务节点信息更新
	for _, item := range s.MemberSessions {
		item.Push("MemberStart", &deal.MemberStartResponse{
			Node: serverNodes,
		})
	}

	return nil
}

// member ping
func (s *MasterServer) ServerHeart(ss session.Session, input *deal.Ping) error {
	node := ss.Get(common.ServerNode).(*deal.ServerNode)
	node.LastTime = time.Now().Unix()
	node.Psutil = input.Psutil
	return ss.Response("ServerHeart", &deal.Pong{
		Pong: node.LastTime,
	})
}
