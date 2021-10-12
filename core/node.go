package core

import (
	"encoding/hex"
	"golang.org/x/net/ipv4"
	"io"
	"net"
)

type Config struct {
	Address  string       `json:"address" yaml:"address"` //tun device ip address
	Server   ServerConfig `json:"server" yaml:"server"`
	Client   ClientConfig `json:"client" yaml:"client"`
	Leaf     bool         `json:"leaf" yaml:"leaf"`
	IsServer bool         `json:"is_server" yaml:"is_server"`
}

// Connection 节点下的连接
type Connection interface {
	ID() string
	Conn() *net.Conn
	CheckData() []byte //检查数据
}

// Node 网络节点
// network
// 每个节点名字唯一
type Node struct {
	Name    string
	Address string
	Port    string
	server  *Server
	client  *Client
	ifc     *Interface
}

func (nd *Node) Interface() *Interface {
	return nd.ifc
}

func (nd *Node) Ready() bool {
	nd.ifc.Ready(nd.Send)

	return true
}

func (nd *Node) Send(data []byte) {
	var Data []byte
	header, err := ipv4.ParseHeader(data)
	if err != nil {
		println(err)
		return
	}
	dst := header.Dst.String()
	var writer io.Writer
	if dst == nd.Address {
		writer = nd.Interface()
	}
	println(dst)
	switch header.Protocol {
	case 1:
		copy(Data, ICMPHandle(data))
	case 6: //TCP
		println("TCP protocol")
		return
	case 17: //UDP
		println("UDP protocol")
		return
	default:
		println("error protocol")
		return
	}
	if Data != nil {
		println(hex.Dump(data))
		_, _ = writer.Write(Data)
	}
}

//func (nd *Node) SendFunc(name string) func(data []byte) {
//	return func(data []byte) {
//		conn := nd.GetConnection(name)
//		l, err := conn.Write(data)
//		if err != nil {
//			log.Println("send error:  ", err)
//		}
//		log.Println("send success:  ", l)
//	}
//}

// NewNode 新建节点
func NewNode(cfg Config) *Node {
	ifc := NewInterface(cfg.Address)
	n := &Node{
		Name:    cfg.Address,
		Address: cfg.Address,
		Port:    "",
		ifc:     ifc,
	}
	if cfg.Leaf {
		server := NewServer(cfg.Server)
		n.server = server
	} else {
		client := NewClient(cfg.Client)
		n.client = client
	}
	return n
}
