package api

import (
	"github.com/gin-gonic/gin"
	"github.com/kiselev-nikolay/direct-to-me/pkg/api/handlers"
	"github.com/kiselev-nikolay/direct-to-me/pkg/server"
	"github.com/kiselev-nikolay/direct-to-me/pkg/storage"
)

func ConnectAPI(ginServer *gin.Engine, fs *storage.FireStoreStorage) {
	ginServer.POST("/api/new", handlers.MakeNewRedirectHandler(fs))
	ginServer.GET("/api/list", handlers.MakeListRedirectsHandler(fs))
	ginServer.POST("/dev/print", handlers.MakeEchoHandler(fs))
	ginServer.StaticFile("/", "./assets/public/index.html")
	ginServer.Static("/static", "./assets/public/")
	redirectHandler := handlers.MakeRedirectHandler(fs)
	ginServer.NoMethod(redirectHandler)
	ginServer.NoRoute(redirectHandler)
	server.RunServer(ginServer)
}
