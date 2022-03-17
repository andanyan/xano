package common

import (
	"strconv"
	"strings"
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

	var cpuCores int32 = 0
	cpuModelName := ""
	var cpuUsePercent float64 = 0
	if len(c) > 0 {
		cpuCores = c[0].Cores
		cpuModelName = c[0].ModelName
	}
	if len(cp) > 0 {
		cpuUsePercent = cp[0]
	}

	return &deal.Psutil{
		MemTotal:      m.Total,
		MemAvailable:  m.Available,
		MemUsed:       m.Used,
		MemUsePercent: m.UsedPercent,
		MemFree:       m.Free,
		// cpu
		CpuCores:      cpuCores,
		CpuModelName:  cpuModelName,
		CpuUsePercent: cpuUsePercent,
		// host
		HostName:      h.Hostname,
		HostBoostTime: h.BootTime,
		HostOs:        h.OS,
		HostPlatform:  h.Platform,
	}

}

// version compare
func VersionCompare(v1, v2 string) bool {
	arr1 := strings.Split(v1, ".")
	arr2 := strings.Split(v2, ".")
	len1 := len(arr1)
	len2 := len(arr2)
	if len1 != len2 {
		return len1 > len2
	}
	for i := 0; i < len1; i++ {
		t1, _ := strconv.Atoi(arr1[i])
		t2, _ := strconv.Atoi(arr2[i])
		if t1 != t2 {
			return t1 > t2
		}
	}
	return false
}
