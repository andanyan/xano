package common

var (
	gateConf *GateConfig
)

func init() {
	gateConf = new(GateConfig)
}

func GetGateConfig() *GateConfig {
	return gateConf
}
