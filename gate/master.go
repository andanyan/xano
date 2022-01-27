package gate

import (
	"log"
	"xlq-server/common"
	"xlq-server/core"
	"xlq-server/deal"
	"xlq-server/router"
)

// 主节点处理
type Master struct{}

func NewMaster() *Master {
	return new(Master)
}

func (m *Master) Run() {
	gConf := common.GetConfig().GateMaster
	addr := gConf.Host + ":" + gConf.Port
	if addr == "" {
		return
	}

	// 注册主节点函数
	router.MasterRouter.Register(&router.RouterServer{
		Name:   "",
		Server: new(MasterServer),
	})

	// 启动服务
	log.Printf("Gate Master Start: %s \n", addr)
	core.NewTcpServer(addr, m.handle)
}

func (m *Master) handle(h *core.TcpHandle, msg *deal.Msg) {
	ss := core.GetSession(h)

	if err := ss.HandleRoute(router.MasterRouter, msg); err != nil {
		log.Println(err)
	}
}
