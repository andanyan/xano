package server

import "xlq-server/core"

func WithTcp(addr string) {
	core.Options.TcpAddr = addr
}

func WithHttp(addr, tlsKey, tlsCert string) {
	core.Options.HttpAddr = addr
	core.Options.HttpTlsKey = tlsKey
	core.Options.HttpTlsCert = tlsCert
}

// 远程节点
func WithMember(memberAddr string, masterAddr string, isMaster bool) {
	core.Options.MasterBool = isMaster
	core.Options.MasterAddr = masterAddr
	core.Options.MemberAddr = memberAddr
}
