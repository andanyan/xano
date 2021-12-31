package common

// 配置数据
type Options struct {
	// 本机公网信息
	MasterBool bool   // 是否master节点
	MasterAddr string // master节点的情况下, 用于启动tcp服务器
	MemberAddr string // 本机rpc服务对应的地址

	// 本机tcp
	TcpAddr        string
	TcpMiddlewares []func(route string) error

	// 本机http
	HttpAddr        string
	HttpTlsKey      string
	HttpTlsCert     string
	HttpMiddlewares []func(route string) error
}
