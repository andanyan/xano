package master

import (
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"xano/common"
	"xano/core"
	"xano/deal"
	"xano/logger"
	"xano/router"
	"xano/session"
)

// 主节点处理
type Master struct{}

func NewMaster() *Master {
	return new(Master)
}

func (m *Master) Run() {

	go m.runTcp()
	go m.runHttp()
}

func (m *Master) Close() {
}

func (m *Master) runTcp() {
	addr := common.GetConfig().Master.TcpAddr
	if addr == "" {
		return
	}

	// 注册主节点函数
	router.GetMasterRouter().Register(&router.RouterServer{
		Name:   "",
		Server: new(MasterServer),
	})
	router.GetMasterRouter().Register(&router.RouterServer{
		Name:   "",
		Server: new(MasterHttpServer),
	})

	// 启动服务
	logger.Infof("Gate Master Tcp Start: %s", addr)
	core.NewTcpServer(addr, func(h *core.TcpHandle) {
		h.SetHandleFunc(m.tcpHandle)
	})
}

func (m *Master) tcpHandle(h *core.TcpHandle, msg *deal.Msg) {
	ss := session.GetBaseSession(h)
	if err := ss.HandleRoute(router.GetMasterRouter(), msg); err != nil {
		logger.Error(err.Error())
	}
}

func (m *Master) runHttp() {
	addr := common.GetConfig().Master.HttpAddr
	if addr == "" {
		return
	}

	router.GetMasterHttpRouter().Register(&router.RouterServer{
		Name:   "",
		Server: new(MasterHttpServer),
	})

	httpMux := http.NewServeMux()
	httpMux.HandleFunc("/", m.httpHandle)
	httpServe := &http.Server{
		Addr:    addr,
		Handler: httpMux,
	}
	logger.Infof("Gate Master Http Start: %s", addr)
	err := httpServe.ListenAndServe()
	if err != nil {
		logger.Fatal(err.Error())
	}
}

// 要处理session的问题 http也通过sid来保持连接状态
func (m *Master) httpHandle(w http.ResponseWriter, r *http.Request) {
	routeStr := r.URL.Path
	routeArr := strings.Split(routeStr, "/")
	routeLen := len(routeArr)
	routeName := ""
	for i := 0; i < routeLen; i++ {
		routeName += strings.Title(routeArr[i])
	}
	if routeName == "" {
		w.Write([]byte("Hello World"))
		return
	}
	// 先判断是否有这个路由 如果没有 直接返回
	if router.GetMasterHttpRouter().GetRoute(routeName) == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 请求数据获取
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// http 仅支持json
	if len(body) == 0 {
		body = []byte("{}")
	}

	// 获取一个tcp连接, 进行逻辑转发
	pool := core.GetPool(common.GetConfig().Master.TcpAddr)
	cli, err := pool.Get()
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer pool.Recycle(cli)

	c := make(chan struct{})

	cli.Client.SetHandleFunc(func(h *core.TcpHandle, m *deal.Msg) {
		// 只要response的数据
		if m.MsgType != common.MsgTypeResponse {
			return
		}
		w.Write(m.Data)
		c <- struct{}{}
	})
	defer cli.Client.SetHandleFunc(nil)

	// 组装消息并发送
	msg := &deal.Msg{
		Route:   routeName,
		Sid:     0,
		Mid:     cli.Client.GetMid(),
		MsgType: common.MsgTypeRequest,
		Deal:    common.TcpDealJson,
		Data:    body,
		Version: common.GetConfig().Base.Version,
	}
	cli.Client.Send(msg)

	// 阻塞等待回包
	t := time.NewTimer(common.HttpDeadDuration)
	select {
	case <-c:
		w.Header().Add("Content-Type", "text/plain")
		//w.WriteHeader(http.StatusOK)
	case <-t.C:
		w.WriteHeader(http.StatusRequestTimeout)
	}
}
