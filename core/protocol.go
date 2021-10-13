package core

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/panjf2000/gnet"
	"log"
	"sync"
)

const (
	DefaultVersion    = 0x1001
	DefaultHeadLength = 8
	DataTypeBeat      = 0x0001
	DataTypeSync      = 0x0002
	DataTypeData      = 0x0003
)

type StarNetProtocol struct {
	Version    uint16
	DataType   uint16
	DataLength uint32
	Data       []byte
}

// Encode data
func (sp *StarNetProtocol) Encode(c gnet.Conn, buf []byte) ([]byte, error) {
	// take out the param
	item := c.Context().(*StarNetProtocol)
	item.Data = buf
	return item.Pack()
}

//Pack ...
func (sp *StarNetProtocol) Pack() ([]byte, error) {
	result := make([]byte, 0)

	buffer := bytes.NewBuffer(result)

	if err := binary.Write(buffer, binary.BigEndian, sp.Version); err != nil {
		s := fmt.Sprintf("Pack version error , %v", err)
		return nil, errors.New(s)
	}

	if err := binary.Write(buffer, binary.BigEndian, sp.DataType); err != nil {
		s := fmt.Sprintf("Pack type error , %v", err)
		return nil, errors.New(s)
	}
	dataLen := uint32(len(sp.Data))
	if err := binary.Write(buffer, binary.BigEndian, dataLen); err != nil {
		s := fmt.Sprintf("Pack datalength error , %v", err)
		return nil, errors.New(s)
	}
	if dataLen > 0 {
		if err := binary.Write(buffer, binary.BigEndian, sp.Data); err != nil {
			s := fmt.Sprintf("Pack data error , %v", err)
			return nil, errors.New(s)
		}
	}

	return buffer.Bytes(), nil
}

// Decode data
func (sp *StarNetProtocol) Decode(c gnet.Conn) ([]byte, error) {
	// parse header
	headerLen := DefaultHeadLength // uint16+uint16+uint32
	if size, header := c.ReadN(headerLen); size == headerLen {
		byteBuffer := bytes.NewBuffer(header)
		var pbVersion, actionType uint16
		var dataLength uint32
		_ = binary.Read(byteBuffer, binary.BigEndian, &pbVersion)
		_ = binary.Read(byteBuffer, binary.BigEndian, &actionType)
		_ = binary.Read(byteBuffer, binary.BigEndian, &dataLength)
		// to check the protocol version and actionType,
		// reset buffer if the version or actionType is not correct
		if pbVersion != DefaultVersion || isCorrectType(actionType) == false {
			c.ResetBuffer()
			log.Println("not normal protocol:", pbVersion, DefaultVersion, actionType, dataLength)
			return nil, errors.New("not normal protocol")
		}
		// parse payload
		dataLen := int(dataLength) // max int32 can contain 210MB payload
		protocolLen := headerLen + dataLen
		if dataSize, data := c.ReadN(protocolLen); dataSize == protocolLen {
			c.ShiftN(protocolLen)
			// log.Println("parse success:", data, dataSize)

			// return the payload of the data
			return data[headerLen:], nil
		}
		// log.Println("not enough payload data:", dataLen, protocolLen, dataSize)
		return nil, errors.New("not enough payload data")

	}
	// log.Println("not enough header data:", size)
	return nil, errors.New("not enough header data")
}

//UnPack ...
func (sp *StarNetProtocol) UnPack(data []byte) error {
	// parse header
	headerLen := DefaultHeadLength // uint16+uint16+uint32
	if len(data) < headerLen {
		return errors.New("not enough header data")
	}
	byteBuffer := bytes.NewBuffer(data[:headerLen])
	var pbVersion, actionType uint16
	var dataLength uint32
	_ = binary.Read(byteBuffer, binary.BigEndian, &pbVersion)
	_ = binary.Read(byteBuffer, binary.BigEndian, &actionType)
	_ = binary.Read(byteBuffer, binary.BigEndian, &dataLength)
	// to check the protocol version and actionType,
	// reset buffer if the version or actionType is not correct
	if pbVersion != DefaultVersion || isCorrectType(actionType) == false {
		log.Println("not normal protocol:", pbVersion, DefaultVersion, actionType, dataLength)
		return errors.New("not normal protocol")
	}
	// parse payload
	dataLen := int(dataLength) // max int32 can contain 210MB payload
	protocolLen := headerLen + dataLen
	if len(data) == protocolLen {
		sp.DataLength = dataLength
		sp.Version = pbVersion
		sp.DataType = actionType
		sp.Data = data[headerLen:]
		return nil
	}
	// log.Println("not enough payload data:", dataLen, protocolLen, dataSize)
	return errors.New("data length error")
}

func isCorrectType(t uint16) bool {
	switch t {
	case DataTypeBeat, DataTypeSync, DataTypeData:
		return true
	default:
		return false
	}
}

var StarNetPool = sync.Pool{
	New: func() interface{} {
		return &StarNetProtocol{
			Version:    DefaultVersion,
			DataType:   DataTypeBeat,
			DataLength: 0,
			Data:       nil,
		}
	},
}
