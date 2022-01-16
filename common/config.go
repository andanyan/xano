package common

var (
	gateConf      *GateConfig
	serviceConfig *TcpServiceConfig
)

func init() {
	gateConf = new(GateConfig)
	serviceConfig = new(TcpServiceConfig)
}

func GetGateConfig() *GateConfig {
	return gateConf
}

func GetServiceConfig() *TcpServiceConfig {
	return serviceConfig
}
