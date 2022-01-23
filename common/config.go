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
	Server ServerConfig
}

type ServerConfig struct {
	// 服务器地址
	Host string
	Port string

	// 通信的网关地址
	GateAddr string

	// 版本号
	Version string
}
