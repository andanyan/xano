package gate

import (
	"fmt"
	"time"
	"xlq-server/common"
	"xlq-server/core"
	"xlq-server/deal"
	"xlq-server/router"
)

// 主节点的服务对象
type MasterServer struct{}

// 新节点启动
func (s *MasterServer) ServerStart(ss *core.Session, input *deal.ServerStartNotice) error {
	addr, err := common.ParseAddr(ss.GetAddr())
	if err != nil {
		return err
	}
	remoteAddr := addr.Ip + ":" + input.Port
	ss.Set(common.HandleKeyTcpAddr, remoteAddr)
	router.GetMasterNode().AddNode(addr.Ip+":"+input.Port, input.Version, input.Routes)
	return nil
}

// 节点关闭
func (s *MasterServer) ServerClose(ss *core.Session, input *deal.ServerCloseNotice) error {
	remoteAddr := ss.Get(common.HandleKeyTcpAddr).(string)
	if remoteAddr == "" {
		return fmt.Errorf("unkoned remote addr")
	}
	router.GetMasterNode().RemoveNode(remoteAddr)
	return nil
}

// 获取全部可用节点
func (s *MasterServer) AllNode(ss *core.Session, input *deal.AllNodeRequest) error {
	nodes := router.GetMasterNode().GetAllNode()

	timeNow := time.Now().Unix()
	res := make([]*deal.NodeItem, 0)
	for _, node := range nodes {
		if node.Status && timeNow < node.LastTime+2*int64(common.TcpHeartDuration) {
			res = append(res, node)
		}
	}

	return ss.Response("AllNode", &deal.AllNodeResponse{
		Nodes: res,
	})
}

// 获取全节点信息
func (s *MasterServer) GetMasterInfo(ss *core.Session, input *deal.MasterInfoRequest) error {
	nodes := router.GetMasterNode().GetAllNode()
	return ss.Response("MasterInfo", &deal.MasterInfoResponse{
		Nodes: nodes,
	})
}
