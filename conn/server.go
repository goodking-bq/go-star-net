package conn

import (
	"fmt"
	"io"
	"net"
	"os"
)

type server struct {
	listener    net.Listener
	ifc         Interfacer
	connections []byte
}

func InitServer(host, port string) {
	listener, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		fmt.Println("Error listening:", err)
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Println("Listening on " + host + ":" + port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err)
			os.Exit(1)
		}
		//logs an incoming message
		fmt.Printf("Received message %s -> %s \n", conn.RemoteAddr(), conn.LocalAddr())
		// Handle connections in a new goroutine.
		ok := handleCheck(conn)
		if ok {

		} else {
			conn.Close()
		}
	}
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
