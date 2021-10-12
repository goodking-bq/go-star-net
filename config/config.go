package config

type Server struct {
	Bind string `json:"bind"` //server bind ip
	Port int    `json:"port"` // bind port
}

type Config struct {
	Address string `json:"address"` //tun device ip address
	Server  Server `json:"server"`
}
