package common

import (
	"xano/deal"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

// 判断对象是否在字符串数组中
func InStringArr(s string, arr []string) bool {
	for _, val := range arr {
		if val == s {
			return true
		}
	}
	return false
}

// 获取服务配置信息
func GetPsutil() *deal.Psutil {
	m, _ := mem.VirtualMemory()
	c, _ := cpu.Info()
	cp, _ := cpu.Percent(0, true)
	h, _ := host.Info()
	return &deal.Psutil{
		MemTotal:      m.Total,
		MemAvailable:  m.Available,
		MemUsed:       m.Used,
		MemUsePercent: m.UsedPercent,
		MemFree:       m.Free,
		// cpu
		CpuCores:      c[0].Cores,
		CpuModelName:  c[0].ModelName,
		CpuUsePercent: cp[0],
		// host
		HostName:      h.Hostname,
		HostBoostTime: h.BootTime,
		HostOs:        h.OS,
		HostPlatform:  h.Platform,
	}

}
