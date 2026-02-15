package main

import (
	"log"

	"github.com/Thashreef45/proxy-server/boot"
)

func main() {

	cfg, err := boot.InitConfig()
	if err != nil {
		log.Printf(err.Error())
		return
	}

	server, err := boot.NewServer(cfg)
	if err != nil {
		log.Printf(err.Error())
		return
	}

	server.Start(cfg.ListenPort)
}
