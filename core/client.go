package core

import (
	"fmt"
	"net"
	"os"
	"strconv"
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
	connections sync.Map
}

func NewClient(cfg ClientConfig) *Client {
	client := &Client{}
	for _, c := range cfg.Connections {
		conn, err := net.Dial("tcp", c.String())
		if err != nil {
			fmt.Println("Error connecting:", err)
			os.Exit(1)
		}
		client.connections.Store(c.String(), conn)

	}
	return client
}

func handleWrite(conn net.Conn, done chan string) {
	for i := 10; i > 0; i-- {
		_, e := conn.Write([]byte("hello " + strconv.Itoa(i) + "\r\n"))
		if e != nil {
			fmt.Println("Error to send message because of ", e.Error())
			break
		}
	}
	done <- "Sent"
}
func handleRead(conn net.Conn, done chan string) {
	buf := make([]byte, 1024)
	reqLen, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error to read message because of ", err)
		return
	}
	fmt.Println(string(buf[:reqLen-1]))
	done <- "Read"
}
