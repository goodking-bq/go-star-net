package tcpip

type Process interface {
	Pack(data []byte) ([]byte, error)
	UnPack(data []byte) ([]byte, error)
}
