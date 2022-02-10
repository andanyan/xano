package gate

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"xlq-server/common"
	"xlq-server/core"
	"xlq-server/deal"
	"xlq-server/logger"
	"xlq-server/router"
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

func (m *Master) runTcp() {
	gConf := common.GetConfig().GateMaster
	if gConf.TcpAddr == "" {
		return
	}

	// 注册主节点函数
	router.MasterRouter.Register(&router.RouterServer{
		Name:   "",
		Server: new(MasterServer),
	})

	// 启动服务
	logger.Infof("Gate Master Tcp Start: %s", gConf.TcpAddr)
	core.NewTcpServer(gConf.TcpAddr, m.tcpHandle)
}

func (m *Master) tcpHandle(h *core.TcpHandle, msg *deal.Msg) {
	ss := core.GetSession(h)

	if err := ss.HandleRoute(router.MasterRouter, msg); err != nil {
		logger.Error(err.Error())
	}
}

func (m *Master) runHttp() {
	gConf := common.GetConfig().GateMaster
	if gConf.HttpAddr == "" {
		return
	}

	httpMux := http.NewServeMux()
	httpMux.HandleFunc("/", m.httpHandle)
	httpServe := &http.Server{
		Addr:    gConf.HttpAddr,
		Handler: httpMux,
	}
	logger.Infof("Gate Master Http Start: %s", gConf.HttpAddr)
	err := httpServe.ListenAndServe()
	if err != nil {
		logger.Fatal(err.Error())
	}
}

func (m *Master) httpHandle(w http.ResponseWriter, r *http.Request) {
	routeStr := r.URL.Path
	routeArr := strings.Split(routeStr, "/")
	routeLen := len(routeArr)
	routeName := ""
	for i := 0; i < routeLen; i++ {
		routeName += strings.Title(routeArr[i])
	}
	if routeName == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	d := common.TcpDealJson
	if r.Header.Get("Deal") == fmt.Sprintf("%d", common.TcpDealProtobuf) {
		d = common.TcpDealProtobuf
	}

	// 获取一个tcp连接, 进行逻辑转发
	pool := core.GetPool(common.GetConfig().GateMaster.TcpAddr)
	cli, err := pool.Get()
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer pool.Recycle(cli)

	c := make(chan struct{})

	cli.Client.SetHandle(func(h *core.TcpHandle, m *deal.Msg) {
		if m.MsgType != common.MsgTypeResponse {
			return
		}
		w.Write(m.Data)
		c <- struct{}{}
	})
	defer cli.Client.SetHandle(nil)

	// 组装消息并发送
	msg := &deal.Msg{
		Route:   routeName,
		Mid:     cli.Client.GetMid(),
		MsgType: common.MsgTypeRequest,
		Deal:    d,
		Data:    body,
		Version: common.GetConfig().Base.Version,
	}
	cli.Client.Send(msg)

	// 阻塞等待回包
	t := time.NewTimer(common.HttpDeadDuration)
	select {
	case <-c:
		w.WriteHeader(http.StatusOK)
	case <-t.C:
		w.WriteHeader(http.StatusRequestTimeout)
	}

}
