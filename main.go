package main

import (
	"flag"
	"github.com/spf13/viper"
	"gnet/conn"
)

var (
	cfg = flag.String("config", "config.yaml", "config path")
)

func main() {
	flag.Parse()
	viper.AddConfigPath(".")
	viper.SetConfigFile(*cfg)
	_ = viper.ReadInConfig()

	node := conn.NewNode(viper.GetString("address"))
	node.Ready()
}
