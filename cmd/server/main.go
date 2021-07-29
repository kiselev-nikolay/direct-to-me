package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"flag"
	"log"

	"github.com/kiselev-nikolay/direct-to-me/pkg/api"
	"github.com/kiselev-nikolay/direct-to-me/pkg/redirectstat"
	"github.com/kiselev-nikolay/direct-to-me/pkg/server"
	"github.com/kiselev-nikolay/direct-to-me/pkg/storage"
)

const authKeyRandomLen = 32

func main() {
	var host string
	var authKey string
	flag.StringVar(&host, "host", "", "Enables HTTPS with autoTLS for host.")
	flag.StringVar(&authKey, "authkey", "", "Set authorization header key for api endpoints. Default is random string.")
	flag.Parse()
	if authKey == "" {
		buf := make([]byte, authKeyRandomLen)
		if _, err := rand.Read(buf); err != nil {
			log.Fatal(err)
		}
		authKey = base64.RawStdEncoding.EncodeToString(buf)
	}
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
	api.ConnectAPI(ginServer, authKey, strg, redag)
	server.RunServer(ginServer, host)
}
