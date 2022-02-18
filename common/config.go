package common

import (
	"xano/logger"

	"github.com/BurntSushi/toml"
)

var config *Config

// 加载配置文件
func SetConfig(file string) {
	if file == "" {
		logger.Fatal("config is nil")
	}
	if config == nil {
		config = new(Config)
	}
	_, err := toml.DecodeFile(file, config)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Debug("Config:", config)
}

// 获取配置
func GetConfig() *Config {
	if config == nil {
		logger.Fatal("config is nil")
	}
	return config
}

type Config struct {
	Base Base

	Server ServerConfig

	GateMaster GateMaster

	GateMember GateMember
}

type Base struct {
	// 版本号
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
	TcpAddr  string
	HttpAddr string
}

type GateMember struct {
	Host string
	Port string
	// 主节点地址
	MasterAddr string
}
