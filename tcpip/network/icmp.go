package network

import (
	"encoding/binary"
	"golang.org/x/net/ipv4"
)

type Error struct {
	msg string

	ignoreStats bool
}

type ICMPHeader struct {
	Type       ipv4.ICMPType
	Code       uint8
	Checksum   uint16
	Identifier uint16
	Sequence   uint16
	Timestamp  uint64
}

func (im *ICMPHeader) Len() int {
	return 16
}

func (im *ICMPHeader) Marshal() ([]byte, error) {
	b := make([]byte, im.Len())
	b[0] = byte(im.Type)
	b[1] = byte(0)
	binary.BigEndian.PutUint16(b[2:4], uint16(0))
	binary.BigEndian.PutUint16(b[4:6], uint16(im.Identifier))
	binary.BigEndian.PutUint16(b[6:8], uint16(im.Sequence))
	binary.BigEndian.PutUint64(b[8:16], im.Timestamp)
	return b, nil
}

func checkSum(data []byte) uint16 {
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

func ParseHeader(data []byte) (*ICMPHeader, error) {
	im := &ICMPHeader{
		Type:       ipv4.ICMPType(data[0]),
		Code:       uint8(data[1]),
		Checksum:   uint16(0), //binary.BigEndian.Uint16(data[2:4]),
		Identifier: binary.BigEndian.Uint16(data[4:6]),
		Sequence:   binary.BigEndian.Uint16(data[6:8]),
		Timestamp:  binary.BigEndian.Uint64(data[8:16]),
	}
	buf, _ := im.Marshal()
	cs := checkSum(buf)
	if cs == binary.BigEndian.Uint16(data[2:4]) {
		im.Checksum = binary.BigEndian.Uint16(data[2:4])
		return im, nil
	} else {
		return nil, nil
	}
}
