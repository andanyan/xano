package xano

import (
	"os"
	"os/signal"
	"syscall"
	"time"
	"xano/common"
	"xano/logger"
	"xano/router"
	"xano/server"
)

func Run() {
	// 启动master
	master := server.NewMaster()
	master.Run()
	// 启动member
	member := server.NewMember()
	member.Run()
	// 启动server
	server := server.NewServer()
	server.Run()

	// 信号监听处理
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL)
	s := <-c
	logger.Info("Receive exit signal: ", s)
	logger.Info("Stoping...")

	// 关闭节点
	go func() {
		member.Close()
		server.Close()
		master.Close()
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
