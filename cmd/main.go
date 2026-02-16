// Copyright (c) 2026 Thashreef45
// Licensed under the MIT License. See LICENSE file for details.

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

	server.Start(cfg)
}
