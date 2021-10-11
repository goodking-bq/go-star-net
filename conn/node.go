package conn

import (
	"encoding/hex"
	"fmt"
	"github.com/spf13/viper"
	"gnet/handle"
	"log"
	"net"
	"os"
	"sync"
)

// 网络节点
// network
// 每个节点名字唯一
type Node interface {
	Name() string           // 节点名称唯一
	Address() string        // 节点的外网地址 //需要连接的用
	Port() string           // 节点的外网端口
	Listener() net.Listener // 节点的侦听
	Interface() Interfacer
	Connections() sync.Map                  // 存储连接  name:连接
	GetConnection(name string) net.Conn     // 获取连接
	Ready() bool                            // 准备就绪
	Send(data []byte)                       //发送数据
	SendFunc(name string) func(data []byte) // 发送数据函数
}

// 节点下的连接
type Connection interface {
	ID() string
	Conn() *net.Conn
	CheckData() []byte //检查数据
}

type node struct {
	name        string
	address     string
	port        string
	listener    net.Listener
	connections sync.Map
	ifc         Interfacer
}

func (nd *node) Name() string {
	return nd.name
}

func (nd *node) Address() string {
	return nd.address
}

func (nd *node) Port() string {
	return nd.port
}

func (nd *node) Listener() net.Listener {
	return nd.listener
}

func (nd *node) Connections() sync.Map {
	return nd.connections
}

func (nd *node) Interface() Interfacer {
	return nd.ifc
}

// 获取连接
func (nd *node) GetConnection(name string) net.Conn {
	return nil
}

func (nd *node) Ready() bool {
	nd.ifc.Ready(nd.Send)

	return true
}

func (nd *node) Send(data []byte) {
	Data := handle.ICMPHandle(data)
	println(hex.Dump(data))
	_, _ = nd.Interface().Write(Data)
}
func (nd *node) SendFunc(name string) func(data []byte) {
	return func(data []byte) {
		conn := nd.GetConnection(name)
		l, err := conn.Write(data)
		if err != nil {
			log.Println("send error:  ", err)
		}
		log.Println("send success:  ", l)
	}
}

func NewNode(address string) Node {
	listener, err := net.Listen("tcp", viper.GetString("host")+":"+viper.GetString("port"))
	if err != nil {
		fmt.Println("Error listening:", err)
		os.Exit(1)
	}
	ifc := NewInterface(address)
	n := &node{
		name:        "node1",
		address:     "",
		port:        "",
		listener:    listener,
		connections: sync.Map{},
		ifc:         ifc,
	}
	return n
}

type connection struct {
}
