package main

import (
	"flag"
	"github.com/goodking-bq/go-star-net/core"
	"github.com/spf13/viper"
)

var (
	cfg    = flag.String("config", "config.yaml", "config path")
	config core.Config
)

func init() {
	flag.Parse()
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetConfigFile(*cfg)
	_ = viper.ReadInConfig()
	println(viper.GetBool("with_server"))
}

func main() {
	err := viper.Unmarshal(&config)
	if err != nil {
		return
	}
	node := core.NewNode(config)
	node.Ready()
}
