package session

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"xlq-server/common"
	"xlq-server/component"
	"xlq-server/core"
)

type HttpSession struct {
	BaseSession
	Request  *http.Request
	Response http.ResponseWriter
}

func NewHttpSession() *HttpSession {
	httpSession := &HttpSession{}
	httpSession.Id = common.GetUuid()
	httpSession.ClientType = common.ClientTypeHttp
	httpSession.Values = make(map[string]interface{})
	return httpSession
}

func (s *HttpSession) Close() {
	s.Request.Body.Close()
}

func (s *HttpSession) Write(route string, data interface{}) {
	bys, err := s.MsgMarsh(data)
	if err != nil {
		log.Panicln(err)
		return
	}
	s.Response.Write(bys)
	s.Response.WriteHeader(http.StatusOK)
}

func (s *HttpSession) Middlewares(route string) error {
	for _, f := range core.Options.HttpMiddlewares {
		err := f(route)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *HttpSession) Handle() {
	// 获取route
	route, _ := url.QueryUnescape(s.Request.URL.Path)
	route = strings.Trim(strings.ReplaceAll(route, "/", "_"), "_")
	if route == "" {
		route = "base"
	}
	// 组装协议
	deal := common.MsgTypeProtoBuf
	switch s.Request.Header.Get("deal") {
	case "json":
		deal = common.MsgTypeJson
	}
	s.MsgType = deal

	// 读取内容
	body, err := ioutil.ReadAll(s.Request.Body)
	if err != nil {
		log.Panic(err.Error())
		s.Response.WriteHeader(http.StatusInternalServerError)
		return
	}

	msg := &common.Msg{
		Route: route,
		Data:  body,
	}

	if err := component.DoneMsg(s, msg); err != nil {
		s.Response.WriteHeader(http.StatusInternalServerError)
	}
}
