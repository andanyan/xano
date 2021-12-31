package server

import (
	"net/http"
	"xlq-server/core"
	"xlq-server/log"
	"xlq-server/session"
)

func (n *Node) serveHttp() {
	if core.Options.HttpAddr == "" {
		return
	}

	var err error

	// 请求
	handleMux := http.NewServeMux()
	handleMux.HandleFunc("/", n.httlHandle)
	httpServer := &http.Server{
		Addr:    core.Options.HttpAddr,
		Handler: handleMux,
	}

	// 加密
	if core.Options.HttpTlsKey != "" {
		err = httpServer.ListenAndServeTLS(core.Options.HttpTlsCert, core.Options.HttpTlsKey)
	} else {
		err = httpServer.ListenAndServe()
	}

	if err != nil {
		log.Fatal(err.Error())
	}
}

func (n *Node) httlHandle(w http.ResponseWriter, r *http.Request) {
	httpSession := session.NewHttpSession()
	defer httpSession.Close()

	httpSession.Response = w
	httpSession.Request = r

	httpSession.Handle()

}
