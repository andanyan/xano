package xano

import (
	"os"
	"os/signal"
	"syscall"
	"time"
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
	go serverGate.Run()
	// 启动服务层
	serverServer := server.NewServer()
	serverServer.Run()

	// 信号监听处理
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL)
	s := <-c
	logger.Info("Receive exit signal: ", s)
	logger.Info("Stoping...")

	// 关闭节点
	go func() {
		serverServer.Close()
		serverGate.Close()
		gateMember.Close()
		gateMaster.Close()
	}()

	time.Sleep(3 * time.Second)
	logger.Info("Stoped")
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
	router.GetLocalRouter().Register(obj)
}

// 日志设定
func WithLog(out *os.File, level int) {
	logger.SetLoggerLevel(level)
	logger.SetOutput(out)
}
