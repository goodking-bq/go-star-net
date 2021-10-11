package conn

import (
	"fmt"
	"github.com/songgao/water"
	"golang.org/x/net/icmp"
	"log"
	"os/exec"
	"runtime"
)

type Interfacer interface {
	Address() string                  // 网卡地址
	Name() string                     // 网卡名
	Init() bool                       // 初始化网卡
	IsReady() bool                    // 是否准备好
	Ready(callBack func(data []byte)) // callback
	Write(data []byte) (n int, err error)
}

type Interface struct {
	ifc     *water.Interface
	config  string
	address string
	netmask string
	name    string
}

func (ifc *Interface) Address() string {
	return ifc.address
}

func (ifc *Interface) Name() string {
	return ifc.name
}
func (ifc *Interface) IsReady() bool {
	return true
}
func (ifc *Interface) Write(data []byte) (n int, err error) {
	n, err = ifc.ifc.Write(data)
	println("write ", n)
	return n, err
}
func (ifc *Interface) Init() bool {
	ifce, err := water.New(water.Config{
		DeviceType: water.TUN,
	})
	if err != nil {
		log.Fatal(err)
		return false
	}
	log.Printf("os name: %s\n", runtime.GOOS)
	var command string
	switch runtime.GOOS {
	case "darwin":
		command = fmt.Sprintf("sudo ifconfig %s 10.3.0.10 10.3.0.20 up", ifce.Name())
	case "linux":
		command = fmt.Sprintf("sudo ifconfig %s 10.3.0.10 netmask 255.255.255.0", ifce.Name())
	}
	log.Printf("Interface Name: %s\n", ifce.Name())
	cmd := exec.Command("/bin/bash", "-c", command)
	if _, err := cmd.Output(); err != nil {
		log.Println(err)
		return false
	}
	ifc.ifc = ifce
	ifc.name = ifce.Name()
	_, err = icmp.ListenPacket("ip4:icmp", "10.3.0.10")
	if err != nil {
		log.Fatal(err)
	}
	return true
}

func (ifc *Interface) Ready(callBack func(data []byte)) {
	packet := make([]byte, 2000)

	for {
		n, err := ifc.ifc.Read(packet)
		if err != nil {
			log.Fatal(err)
		}
		callBack(packet[:n])
	}

}

func NewInterface(address string) Interfacer {
	ifc := &Interface{
		ifc:     nil,
		config:  "",
		address: address,
		netmask: "",
		name:    "",
	}
	ifc.Init()
	return ifc
}
