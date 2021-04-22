package main

import (
	"log"

	"github.com/kiselev-nikolay/direct-to-me/pkg/api"
	"github.com/kiselev-nikolay/direct-to-me/pkg/conf"
	"github.com/kiselev-nikolay/direct-to-me/pkg/server"
	"github.com/kiselev-nikolay/direct-to-me/pkg/storage"
)

func main() {
	conf := conf.ReadConfig("./conf.yaml")
	projectID := "decent-genius-311507"
	fs := &storage.FireStoreStorage{}
	err := fs.Connect(storage.FireStoreStorageConf{
		ProjectID:       projectID,
		CredentialsPath: conf.Google.Application.Credentials.Storage,
	})
	if err != nil {
		log.Fatal(err)
	}
	ginServer := server.GetServer()
	api.ConnectAPI(ginServer, fs)
	server.RunServer(ginServer)
}
