package common

import "reflect"

// 原始消息包
type Packet struct {
	Type    uint8
	DataLen uint16
	Data    []byte
}

// 数据包
type Msg struct {
	// 消息id
	Mid uint64
	// route 格式为 Aa_Bb
	Route string
	// 数据
	Data []byte
}

// 路由
type Route struct {
	Method reflect.Method
	Input  reflect.Type
	OutPut reflect.Type
}
