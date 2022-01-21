package node

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"xlq-server/common"
	"xlq-server/core"
	"xlq-server/deal"
	"xlq-server/session"
)

type Gate struct{}

func NewGate() *Gate {
	return new(Gate)
}

// tcp gate接口
func (g *Gate) RunTcp() {
	if common.GetGateConfig().TcpAddr == "" {
		return
	}
	log.Printf("Gate Tcp Address: %s \n", common.GetGateConfig().TcpAddr)
	core.NewTcpServer(common.GetGateConfig().TcpAddr, g.TcpHandle)
}

// 接收到包
func (g *Gate) TcpHandle(h *core.TcpHandle, p *common.TcpPacket) {
	// 获取Session, 包含用户的数据
	s := session.NewGateSession(h)
	err := s.HandleTcp(p)
	if err != nil {
		log.Println(err)
	}
}

// 运行http
func (g *Gate) RunHttp() {
	if common.GetGateConfig().HttpAddr == "" {
		return
	}

	// 请求
	handleMux := http.NewServeMux()
	handleMux.HandleFunc("/", g.HttpHandle)
	httpServer := &http.Server{
		Addr:    common.GetGateConfig().HttpAddr,
		Handler: handleMux,
	}

	log.Printf("Gate Http Address: %s \n", common.GetGateConfig().HttpAddr)
	log.Fatal(httpServer.ListenAndServe())
}

// http逻辑处理
func (g *Gate) HttpHandle(w http.ResponseWriter, r *http.Request) {
	// 路由解析
	route, err := url.QueryUnescape(r.URL.Path)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Panicln(err.Error())
		return
	}

	// 路由转化
	routeArr := strings.Split(route, "/")
	routeLen := len(routeArr)
	if routeLen < 3 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	newRoute := ""
	for i := 1; i < routeLen; i++ {
		newRoute += strings.Title(routeArr[i])
		if i == 1 {
			newRoute += "_"
		}
	}

	// 读取内容
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Panic(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s := session.NewGateSession(nil)
	res, err := s.HandleHttp(newRoute, body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// 提供内部注册服务
func (g *Gate) RunInner() {
	if common.GetGateConfig().GateAddr == "" {
		return
	}

	log.Printf("Gate Inner Address: %s \n", common.GetGateConfig().GateAddr)
	core.NewTcpServer(common.GetGateConfig().GateAddr, g.InnerHandle)
}

func (g *Gate) InnerHandle(h *core.TcpHandle, p *common.TcpPacket) {
	// 解析包
	msg := new(deal.GateRouteRequest)
	err := common.TcpMsgUnMarsh(p.Data, msg)
	if err != nil {
		log.Panicln(err.Error())
		return
	}
	remoteAddr, err := common.ParseIpAddr(h.GetAddr())
	if err != nil {
		log.Panicln(err.Error())
		return
	}
	serverAddr := remoteAddr.Ip + ":" + msg.Port

	// 节点退出
	if int(msg.Status) == common.StatusOne {
		core.RemoveRoute(serverAddr)
		return
	}

	// 节点加入
	core.AddRoute(serverAddr, msg.Routes)
	// 组装回包
	res := &deal.GateRouteResponse{
		Routes: core.GetAllRoute(),
	}
	resBys, err := common.TcpMsgMarsh(res)
	if err != nil {
		log.Println(err.Error())
		return
	}
	resPacket := &common.TcpPacket{
		Length: uint16(len(resBys)),
		Data:   resBys,
	}
	h.Send(resPacket)
}
