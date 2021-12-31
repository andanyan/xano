package session

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"xlq-server/log"

	"xlq-server/common"
	"xlq-server/core"

	"google.golang.org/protobuf/proto"
)

type BaseSession struct {
	Id         string
	ClientType uint8
	MsgType    uint8

	Values map[string]interface{}
}

// 回调
func (s *BaseSession) Write(route string, data interface{}) {

}

// Rpc
func (s *BaseSession) Rpc(route string, input interface{}, output interface{}) error {
	// 获取一个地址, 发起请求, 返回数据
	conn, err := net.Dial("tcp", core.Options.TcpAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	// 组包
	inputBys, err := s.MsgMarsh(input)
	if err != nil {
		return err
	}
	inputMsg := &common.Msg{
		Route: route,
		Data:  inputBys,
	}
	inputMsgBys, err := s.MsgMarsh(inputMsg)
	if err != nil {
		return err
	}
	inputPacket := &common.Packet{
		Type:    s.MsgType,
		DataLen: uint16(len(inputMsgBys)),
		Data:    inputMsgBys,
	}
	inputPacketBys, err := s.MsgMarsh(inputPacket)
	if err != nil {
		return err
	}
	// 发送请求
	_, err = conn.Write(inputPacketBys)
	if err != nil {
		return err
	}

	// 不要回复
	if output == nil {
		return nil
	}

	var outPacket *common.Packet
	packetBuf := make([]byte, 0)
	buf := make([]byte, common.ConnReadBufLen)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Println("Read [%s], session will be closed immediately", err.Error())
				return err
			}
			continue
		}
		if n <= 0 {
			continue
		}
		packetBuf = append(packetBuf, buf[:n]...)

		// 解析包
		packets, index := s.ParsePacket(packetBuf)
		packetLen := len(packets)
		if packetLen > 0 {
			packetBuf = packetBuf[index:]
			outPacket = packets[packetLen-1]
			break
		}
	}

	if outPacket == nil {
		return nil
	}

	// 解析出msg
	s.MsgType = outPacket.Type
	msg := new(common.Msg)
	err = s.MsgUnMarsh(outPacket.Data, msg)
	if err != nil {
		return err
	}

	err = s.MsgUnMarsh(msg.Data, output)
	if err != nil {
		return err
	}

	return nil
}

// 关闭session
func (s *BaseSession) Close() {
}

// 获取值
func (s *BaseSession) GetValue(key string) interface{} {
	return s.Values[key]
}

// 设置值
func (s *BaseSession) SetValue(key string, val interface{}) {
	s.Values[key] = val
}

// 解析packet
func (s *BaseSession) ParsePacket(buf []byte) ([]*common.Packet, int) {
	var packets []*common.Packet
	var index int
	var err error
	bufLen := len(buf)
	for index < bufLen {
		if bufLen-index < common.MsgTypeLen+common.MsgDataLen {
			break
		}
		packet := new(common.Packet)
		// 读取dataLen
		dataLenStart := index + common.MsgTypeLen
		dataLenEnd := dataLenStart + common.MsgDataLen
		err = binary.Read(bytes.NewReader(buf[dataLenStart:dataLenEnd]), binary.LittleEndian, &packet.DataLen)
		if err != nil {
			log.Println(err)
			break
		}
		if bufLen-index < common.MsgTypeLen+common.MsgDataLen+int(packet.DataLen) {
			break
		}
		// 读取type
		typeStart := index
		typeEnd := typeStart + common.MsgTypeLen
		err = binary.Read(bytes.NewReader(buf[typeStart:typeEnd]), binary.LittleEndian, &packet.Type)
		if err != nil {
			log.Println(err)
			break
		}
		// 包内容
		bufStart := index + common.MsgTypeLen + common.MsgDataLen
		bufEnd := bufStart + int(packet.DataLen)
		packet.Data = buf[bufStart:bufEnd]
		packets = append(packets, packet)
		index = bufEnd
	}
	return packets, index
}

// 解析数据包
func (s *BaseSession) MsgMarsh(data interface{}) ([]byte, error) {
	switch s.MsgType {
	case common.MsgTypeProtoBuf:
		pb, ok := data.(proto.Message)
		if !ok {
			return nil, fmt.Errorf("protobuf: convert on wrong type value")
		}
		return proto.Marshal(pb)
	case common.MsgTypeJson:
		return json.Marshal(data)
	}
	return nil, fmt.Errorf("UnSupport Msg Type, %d", s.MsgType)
}

// 加密数据包
func (s *BaseSession) MsgUnMarsh(data []byte, v interface{}) error {
	switch s.MsgType {
	case common.MsgTypeProtoBuf:
		pb, ok := v.(proto.Message)
		if !ok {
			return fmt.Errorf("protobuf: convert on wrong type value")
		}
		return proto.Unmarshal(data, pb)
	case common.MsgTypeJson:
		return json.Unmarshal(data, v)
	}
	return fmt.Errorf("UnSupport Msg Type, %d", s.MsgType)
}

// 中间件
func (s *BaseSession) Middlewares(msg *common.Msg) error {
	return nil
}
