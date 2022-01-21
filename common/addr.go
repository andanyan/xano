package common

import (
	"errors"
	"strings"
)

type Addr struct {
	Ip   string
	Port string
}

// 组装地址
func GenAddr(host, port string) string {
	return host + ":" + port
}

// 解析地址
func ParseAddr(s string) (*Addr, error) {
	sArr := strings.Split(s, ":")
	if len(sArr) != 2 {
		return nil, errors.New("error addr: " + s)
	}
	addr := &Addr{
		Ip:   sArr[0],
		Port: sArr[1],
	}
	return addr, nil
}
