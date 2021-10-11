package main

import (
	"flag"
	"github.com/goodking-bq/go-star-net/conn"
	"github.com/spf13/viper"
)

var (
	cfg = flag.String("config", "config.yaml", "config path")
)

func main() {
	flag.Parse()
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetConfigFile(*cfg)
	_ = viper.ReadInConfig()

	node := conn.NewNode(viper.GetString("address"))
	node.Ready()
}
