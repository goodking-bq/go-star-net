package handler

import "golang.org/x/net/ipv4"

func Handle(data []byte) []byte {
	header, err := ipv4.ParseHeader(data)
	if err != nil {
		println(err)
		return nil
	}
	switch header.Protocol {
	case 1:
		return ICMPHandle(data)
	case 6: //TCP
		println("TCP protocol")
		return nil
	case 17: //UDP
		println("UDP protocol")
		return nil
	default:
		println("error protocol")
		return nil
	}
	return nil
}
