package core

import (
	"fmt"
	"github.com/panjf2000/gnet"
	"io"
	"log"
	"net"
	"sync"
)

type ServerConfig struct {
	Bind string `json:"bind"` //server bind ip
	Port string `json:"port"` // bind port
}
type Server struct {
	*gnet.EventServer
	connections sync.Map
}

func (server *Server) OnOpened(c gnet.Conn) (out []byte, action gnet.Action) {
	server.connections.Store(c.RemoteAddr(), c)
	return
}

func (server *Server) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	out = frame
	return
}

func (server *Server) WriterConn(dst string) *net.Conn {
	conn, ok := server.connections.Load(dst)
	if ok {
		return conn.(*net.Conn)
	}
	return nil
}

func NewServer(cfg ServerConfig) *Server {
	server := &Server{}
	go func() {
		err := gnet.Serve(server, "tcp://"+cfg.Bind+":"+cfg.Port)
		if err != nil {
			log.Fatalln(err)
		}
	}()
	return server
}

func handleCheck(conn net.Conn) bool {
	fmt.Println("pair up ...")
	conn.Write([]byte("report"))
	buf := make([]byte, 1024)
	reqLen, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error to read message because of ", err)
		return false
	}
	fmt.Println(string(buf[:reqLen-1]))
	return true
}

func handleRequest(conn net.Conn) {
	defer conn.Close()
	for {
		io.Copy(conn, conn)
	}
}
