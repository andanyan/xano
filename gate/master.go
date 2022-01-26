package gate

import (
	"log"
	"xlq-server/common"
	"xlq-server/core"
	"xlq-server/inner"
	"xlq-server/router"
)

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
		Server: new(inner.Master),
	})

	// 启动服务
	log.Printf("Gate Master Start: %s \n", addr)
	core.NewTcpServer(addr, m.handle)
}

func (m *Master) handle(h *core.TcpHandle, p *common.Packet) {

}
