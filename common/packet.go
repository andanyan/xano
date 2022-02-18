package common

import (
	"bytes"
	"encoding/binary"
	"xano/logger"
)

// tcp数据包结构
type Packet struct {
	// 数据源长度
	Length uint16
	// 源数据
	Data []byte
}

// 包转码
func PacketMarsh(p *Packet) ([]byte, error) {
	var err error
	buf := bytes.NewBuffer(nil)
	err = binary.Write(buf, binary.LittleEndian, p.Length)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.LittleEndian, p.Data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// 包解析 返回占用字节数、包、错误
func PacketUnMarsh(bys []byte) (int, []*Packet) {
	var index int
	var packets []*Packet
	blen := len(bys)
	for {
		p := new(Packet)
		headStart := index
		headEnd := index + PacketHeadLength
		if headEnd >= blen {
			break
		}
		binary.Read(bytes.NewReader(bys[headStart:headEnd]), binary.LittleEndian, &p.Length)
		dataStart := headEnd
		dataEnd := dataStart + int(p.Length)
		if dataEnd > blen {
			break
		}
		p.Data = make([]byte, p.Length)
		err := binary.Read(bytes.NewReader(bys[dataStart:dataEnd]), binary.LittleEndian, &(p.Data))
		if err != nil {
			logger.Error(err)
		}
		packets = append(packets, p)
		index = dataEnd
	}
	return index, packets
}
