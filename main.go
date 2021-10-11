package main

import (
	"flag"
	"github.com/goodking-bq/go-star-net/conn"
	"github.com/spf13/viper"
)

var (
	cfg = flag.String("config", "config.yaml", "config path")
)

func init() {
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetConfigFile(*cfg)
	_ = viper.ReadInConfig()
	println("env address is: ", viper.GetString("address"))
}

func main() {
	flag.Parse()

	node := conn.NewNode(viper.GetString("address"))
	node.Ready()
}
