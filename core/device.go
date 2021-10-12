package core

import (
	"fmt"
	"github.com/songgao/water"
	"golang.org/x/net/icmp"
	"log"
	"os/exec"
	"runtime"
)

// Interface tun device interface
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
	ifc.ifc = ifce
	ifc.name = ifce.Name()
	switch runtime.GOOS {
	case "darwin":
		ifc.initDarwin()
	case "linux":
		ifc.initLinux()
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

func NewInterface(address string) *Interface {
	ifc := &Interface{
		ifc:     nil,
		config:  "",
		address: address,
		netmask: "255.255.255.0",
		name:    "",
	}
	ifc.Init()
	return ifc
}

func (ifc *Interface) initDarwin() {
	command := fmt.Sprintf("sudo ifconfig %s 10.3.0.10 %s up", ifc.Name(), ifc.address)
	log.Printf("Interface Name: %s\n", ifc.Name())
	cmd := exec.Command("/bin/bash", "-c", command)
	if _, err := cmd.Output(); err != nil {
		log.Println(err)
		log.Fatal(err)
	}
	_, err := icmp.ListenPacket("ip4:icmp", "10.3.0.10")
	if err != nil {
		log.Fatal(err)
	}
}

func (ifc *Interface) initLinux() {
	command := fmt.Sprintf("sudo ifconfig %s %s netmask %s", ifc.Name(), ifc.address, ifc.netmask)
	log.Printf("Interface Name: %s\n", ifc.Name())
	cmd := exec.Command("/bin/bash", "-c", command)
	if _, err := cmd.Output(); err != nil {
		log.Println(err)
		log.Fatal(err)
	}
	_, err := icmp.ListenPacket("ip4:icmp", ifc.address)
	if err != nil {
		log.Fatal(err)
	}
}
