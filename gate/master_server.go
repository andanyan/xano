package gate

import (
	"sync"
	"time"
	"xano/common"
	"xano/core"
	"xano/deal"
	"xano/router"
)

// 主节点的服务对象
type MasterServer struct {
	Mux            sync.Mutex
	MemberSessions []*core.Session
	ServerSessions []*core.Session
}

// member start
func (s *MasterServer) MemberStart(ss *core.Session, input *deal.MemberStartNotice) error {
	addr, err := common.ParseAddr(ss.GetAddr())
	if err != nil {
		return err
	}

	newNode := &deal.MemberNode{
		Version:  input.Version,
		Addr:     addr.Ip + ":" + input.Port,
		LastTime: time.Now().Unix(),
	}

	router.GetMasterNode().AddMemberNode(newNode)

	ss.Set(common.MemberNode, newNode)

	s.Mux.Lock()
	defer s.Mux.Unlock()
	s.MemberSessions = append(s.MemberSessions, ss)

	memberNodes := router.GetMasterNode().AllMemberNode()
	// 通知所有服务节点 从节点信息更新
	for _, item := range s.ServerSessions {
		item.Push("MemberNode", &deal.MemberNodePush{
			Nodes: memberNodes,
		})
	}

	// 告知当前从节点 当前所有的服务节点信息
	serverNodes := router.GetMasterNode().AllServerNode()
	return ss.Push("ServerNode", &deal.ServerNodePush{
		Nodes: serverNodes,
	})
}

// member stop
func (s *MasterServer) MemberStop(ss *core.Session, input *deal.MemberStopNotice) error {
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

	// 获取当前全部的nodes
	memberNodes := router.GetMasterNode().AllMemberNode()

	// 通知所有服务节点 从节点信息更新
	for _, item := range s.ServerSessions {
		item.Push("MemberNode", &deal.MemberNodePush{
			Nodes: memberNodes,
		})
	}

	return nil
}

// member ping
func (s *MasterServer) MemberHeart(ss *core.Session, input *deal.Ping) error {
	node := ss.Get(common.MemberNode).(*deal.MemberNode)
	node.LastTime = time.Now().Unix()
	return ss.Response("MemberHeart", &deal.Pong{
		Pong: node.LastTime,
	})
}

// server start
func (s *MasterServer) ServerStart(ss *core.Session, input *deal.ServerStartNotice) error {
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
		item.Push("ServerNode", &deal.ServerNodePush{
			Nodes: serverNodes,
		})
	}

	// 告知服务节点所有的网关节点
	memberNodes := router.GetMasterNode().AllMemberNode()
	return ss.Push("MemberNode", &deal.MemberNodePush{
		Nodes: memberNodes,
	})
}

// member stop
func (s *MasterServer) ServerStop(ss *core.Session, input *deal.ServerStopNotice) error {
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

	// 通知所有的从节点 服务节点信息更新
	for _, item := range s.MemberSessions {
		item.Push("ServerNode", &deal.ServerNodePush{
			Nodes: router.GetMasterNode().AllServerNode(),
		})
	}

	return nil
}

// member ping
func (s *MasterServer) ServerHeart(ss *core.Session, input *deal.Ping) error {
	node := ss.Get(common.ServerNode).(*deal.ServerNode)
	node.LastTime = time.Now().Unix()
	return ss.Response("ServerHeart", &deal.Pong{
		Pong: node.LastTime,
	})
}
