package common

var config *Config

// 加载配置文件
func SetConfig(conf *Config) {
	config = conf
}

// 获取配置
func GetConfig() *Config {
	return config
}

type Config struct {
	Base Base

	Server ServerConfig

	GateMaster GateMaster

	GateMember GateMember
}

type Base struct {
	Version string
}

type ServerConfig struct {
	// 服务器地址
	Host string
	Port string

	// 通信的网关地址
	GateAddr string
}

type GateMaster struct {
	Host string
	Port string
}

type GateMember struct {
	Host string
	Port string
	// 主节点地址
	MasterAddr string
	// 版本号
}
