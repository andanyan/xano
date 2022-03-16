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
	logger.Debugf("Config: %+v", config)
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

	Server Server

	Master Master

	Member Member
}

type Base struct {
	// 版本号
	Version string
	// Msg协议
	TcpDeal uint32
}

type Server struct {
	// 服务器地址
	TcpAddr string

	// 通信的网关地址
	MasterAddr string
}

type Master struct {
	TcpAddr  string
	HttpAddr string
}

type Member struct {
	// tcp对外地址
	TcpAddr string
	// tcp对内地址
	InnerAddr string
	// 主节点地址
	MasterAddr string
}
