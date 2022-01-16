package common

import (
	"bytes"
	"encoding/binary"
)

// 包转码
func PacketMarsh(p *TcpPacket) ([]byte, error) {
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
func PacketUnMarsh(bys []byte) (int, []*TcpPacket) {
	var index int
	var packets []*TcpPacket
	blen := len(bys)
	for {
		p := new(TcpPacket)
		headStart := index
		headEnd := index + TcpPacketHeadLength
		if headEnd >= blen {
			break
		}
		binary.Read(bytes.NewReader(bys[headStart:headEnd]), binary.LittleEndian, &p.Length)
		dataStart := headEnd
		dataEnd := dataStart + int(p.Length)
		if dataEnd >= blen {
			break
		}
		binary.Read(bytes.NewReader(bys[dataStart:dataEnd]), binary.LittleEndian, &p.Data)
		packets = append(packets, p)
		index = dataEnd
	}
	return index, packets
}
