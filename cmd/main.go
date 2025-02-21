package main

import (
	"flag"
	"fmt"

	"connor.run/deepcog/api"
	"connor.run/deepcog/pkg/config"
)

func main() {
	cfgPath := flag.String("config", "config.toml", "path to config file")
	flag.Parse()

	cfg, err := config.LoadConfig(*cfgPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	e := api.NewServer()
	e.Start(fmt.Sprintf("0.0.0.0:%s", cfg.Server.Port))
}
