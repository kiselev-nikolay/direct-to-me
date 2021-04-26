package main

import (
	"log"

	"github.com/kiselev-nikolay/direct-to-me/pkg/api"
	"github.com/kiselev-nikolay/direct-to-me/pkg/server"
	"github.com/kiselev-nikolay/direct-to-me/pkg/storage"
)

func main() {
	strg := &storage.BitcaskStorage{}
	err := strg.Connect()
	if err != nil {
		log.Fatal(err)
	}
	ginServer := server.GetServer()
	api.ConnectAPI(ginServer, strg)
	server.RunServer(ginServer)
}
