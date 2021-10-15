package core

import (
	"fmt"
	"golang.org/x/net/ipv4"
	"log"
	"net"
	"sync"
)

type OptionConnection struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

func (opt OptionConnection) String() string {
	return opt.Host + ":" + opt.Port
}

type Options struct {
	Address     string             `json:"address" yaml:"address"`
	Connections []OptionConnection `json:"connections" yaml:"connections"`
	Server      struct {
		Enable bool   `json:"enable"`
		Bind   string `json:"bind"`
		Port   string `json:"port"`
	} `json:"server"`
}

type Node struct {
	*Server
	*Interface
	conns   sync.Map
	options Options
}

//InterfaceCallBack  handle interface save data
func (node *Node) InterfaceCallBack(data []byte) {
	header, err := ipv4.ParseHeader(data)
	if err != nil {
		println(err)
		return
	}
	dst := header.Dst.String()         //to where
	if dst == node.Interface.address { // to local
		if header.Protocol == 1 {
			_, _ = node.Write(ICMPHandle(data))
		} else {
			_, err := node.Write(data)
			if err != nil {
				return
			}
		}
	} else { //to remote

	}

}

func (node *Node) StartClient() error {
	return nil
}

func (node *Node) StartServer() error {
	return nil
}

func (node *Node) Start() error {
	go node.Ready(node.InterfaceCallBack)
	return nil
}

func NewNode(opt Options) *Node {
	ifc := NewInterface(opt.Address)
	node := &Node{
		Server:    nil,
		Interface: ifc,
		options:   opt,
	}
	if opt.Server.Enable {

	}
	for _, c := range opt.Connections {
		conn, err := net.Dial("tcp", c.String())
		if err != nil {
			fmt.Printf("error connecting: %s", err)
		}
		node.connections.Store(c.String(), conn)
		go func() {
			for {
				newPackage, err := DecodeForClient(conn)
				if err != nil {
					continue
				}
				switch newPackage.DataType {
				case DataTypeBeat:
					log.Printf("beat from %s", conn.RemoteAddr())
				case DataTypeSync:
				case DataTypeData:
					_, err := nd.ifc.Write(newPackage.Data)
					if err != nil {
						return
					}
				default:

				}
			}
		}()
	}
	return node
}
