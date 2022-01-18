package common

var (
	gateConf      *GateConfig
	serviceConfig *TcpServiceConfig
)

func init() {
	gateConf = new(GateConfig)
	serviceConfig = new(TcpServiceConfig)
}

func SetGateConfig(conf *GateConfig) {
	if conf == nil {
		return
	}
	gateConf = conf
}

func SetServiceConfig(conf *TcpServiceConfig) {
	if conf == nil {
		return
	}
	serviceConfig = conf
}

func GetGateConfig() *GateConfig {
	return gateConf
}

func GetServiceConfig() *TcpServiceConfig {
	return serviceConfig
}
