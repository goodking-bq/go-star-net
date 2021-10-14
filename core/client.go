package core

import (
	"fmt"
	"log"
	"net"
	"sync"
)

type Connect struct {
	Host string
	Port string
}

func (c Connect) String() string {
	return c.Host + ":" + c.Port
}

type ClientConfig struct {
	Connections []Connect `json:"connections" yaml:"connections"`
}

type Client struct {
	connections sync.Map //链接
}

func (c *Client) WriteTo(data []byte, to string) {

}

func (c *Client) WriteToLocal(data []byte) {

}

func NewClient(cfg ClientConfig) *Client {
	client := &Client{}
	for _, c := range cfg.Connections {
		conn, err := net.Dial("tcp", c.String())
		if err != nil {
			fmt.Println("Error connecting:", err)
			continue
		}
		client.connections.Store(c.String(), conn)
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

				default:

				}
			}
		}()
	}
	return client
}

// ConnToDevice transport  client conn to device
func ConnToDevice(conn net.Conn, tun Interface) {

}
