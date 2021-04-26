package main

import (
	"context"
	"flag"
	"log"

	"github.com/kiselev-nikolay/direct-to-me/pkg/api"
	"github.com/kiselev-nikolay/direct-to-me/pkg/redirectstat"
	"github.com/kiselev-nikolay/direct-to-me/pkg/server"
	"github.com/kiselev-nikolay/direct-to-me/pkg/storage"
)

func main() {
	var host string
	flag.StringVar(&host, "host", "", "Enables HTTPS with autoTLS for host.")
	flag.Parse()
	strg := &storage.BitcaskStorage{}
	err := strg.Connect()
	if err != nil {
		log.Fatal(err)
	}
	ginServer := server.GetServer()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	redag := &redirectstat.RedirectAggregation{}
	redag.Worker(ctx)
	api.ConnectAPI(ginServer, strg, redag)
	server.RunServer(ginServer, host)
}
