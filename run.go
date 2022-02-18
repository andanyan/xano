package xano

import (
	"os"
	"xano/common"
	"xano/gate"
	"xano/logger"
	"xano/router"
	"xano/server"
)

func Run() {
	// 启动网关主节点
	gateMaster := gate.NewMaster()
	gateMaster.Run()
	// 启动member
	gateMember := gate.NewMember()
	gateMember.Run()
	// 启动网关层客户端
	serverGate := server.NewGate()
	serverGate.Run()
	// 启动服务层
	serverServer := server.NewServer()
	serverServer.Run()
}

// 设置配置
func WithConfig(file string) {
	common.SetConfig(file)
}

// 注册路由
func WithRoute(obj *router.RouterServer) {
	if obj == nil {
		return
	}
	router.LocalRouter.Register(obj)
}

// 日志设定
func WithLog(out *os.File, level int) {
	logger.SetLoggerLevel(level)
	logger.SetOutput(out)
}
