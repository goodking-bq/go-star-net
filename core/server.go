package core

import (
	"fmt"
	"github.com/panjf2000/gnet"
	"log"
	"net"
	"sync"
	"time"
)

type ServerConfig struct {
	Bind string `json:"bind"` //server bind ip
	Port string `json:"port"` // bind port
}
type Server struct {
	*gnet.EventServer
	codec       gnet.ICodec // 自定义
	connections sync.Map
}

func (server *Server) OnOpened(c gnet.Conn) (out []byte, action gnet.Action) {
	server.connections.Store(c.RemoteAddr(), c)
	return
}
func (server *Server) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
	server.connections.Delete(c.RemoteAddr().String())
	return
}
func (server *Server) OnInitComplete(srv gnet.Server) (action gnet.Action) {
	log.Printf("Test codec server is listening on %s (multi-cores: %t, loops: %d)\n",
		srv.Addr.String(), srv.Multicore, srv.NumEventLoop)
	return
}

// React 处理数据 frame 输入的数据  out 返回的数据
func (server *Server) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	fmt.Println("frame:", string(frame))
	p := c.Context().(StarNetProtocol)
	switch p.DataType {
	case DataTypeBeat:
		println("beat")
		item := StarNetPool.Get().(*StarNetProtocol)
		item.DataType = DataTypeBeat
		c.SetContext(item)
		out = []byte("beat")
		return
	case DataTypeSync:
		println("sync")
		item := StarNetPool.Get().(*StarNetProtocol)
		item.DataType = DataTypeSync
		c.SetContext(item)
		out = []byte("sync")
		return
	case DataTypeData:
		item := StarNetPool.Get().(*StarNetProtocol)
		item.DataType = DataTypeData
		c.SetContext(item)
		println("data")
		out = []byte("data")
		return
	default:
		return
	}
}
func (server *Server) Tick() (delay time.Duration, action gnet.Action) {
	server.connections.Range(func(key, value interface{}) bool {
		conn := value.(gnet.Conn)
		p := StarNetPool.Get().(*StarNetProtocol)
		p.DataType = DataTypeBeat
		conn.SetContext(p)
		_ = conn.AsyncWrite([]byte("beat"))
		return true
	})
	return time.Second, gnet.None
}

func (server *Server) WriterConn(dst string) *net.Conn {
	conn, ok := server.connections.Load(dst)
	if ok {
		return conn.(*net.Conn)
	}
	return nil
}

func NewServer(cfg ServerConfig) *Server {
	codec := &StarNetProtocol{}
	server := &Server{codec: codec}
	go func() {
		err := gnet.Serve(server, "tcp://"+cfg.Bind+":"+cfg.Port, gnet.WithCodec(codec))
		if err != nil {
			log.Fatalln(err)
		}
	}()
	return server
}
