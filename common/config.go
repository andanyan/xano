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
	// 服务器数据
	Host string
	Port string
}
