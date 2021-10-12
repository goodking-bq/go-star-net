package handler

import (
	"encoding/binary"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"time"
)

func CheckSum(data []byte) uint16 {
	var (
		sum    uint32
		length int = len(data)
		index  int
	)
	for length > 1 {
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		index += 2
		length -= 2
	}
	if length > 0 {
		sum += uint32(data[index])
	}
	sum += (sum >> 16)

	return uint16(^sum)
}

func timeToBytes(t time.Time) []byte {
	nsec := t.UnixNano()
	b := make([]byte, 8)
	for i := uint8(0); i < 8; i++ {
		b[i] = byte((nsec >> ((7 - i) * 8)) & 0xff)
	}
	return b
}

func ICMPHandle(data []byte) []byte {
	a, _ := ipv4.ParseHeader(data)
	Data := make([]byte, len(data))
	a.Src, a.Dst = a.Dst, a.Src
	a.Checksum = 0
	a.TotalLen = 84
	headerData, err := a.Marshal()
	if err != nil {
		println(err)
		return nil
	}
	binary.BigEndian.PutUint16(headerData[2:4], uint16(a.TotalLen))
	binary.BigEndian.PutUint16(headerData[10:12], CheckSum(headerData))
	copy(Data, headerData)
	replyData := data[a.Len:]
	replyData[0] = 0
	binary.BigEndian.PutUint32(replyData[8:16], uint32(time.Now().Unix()))
	message, _ := icmp.ParseMessage(1, replyData)
	e, _ := message.Marshal(nil)
	copy(Data[a.Len:], e)
	return Data
}
