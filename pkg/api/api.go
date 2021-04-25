package api

import (
	"github.com/gin-gonic/gin"
	"github.com/kiselev-nikolay/direct-to-me/pkg/api/handlers"
	"github.com/kiselev-nikolay/direct-to-me/pkg/server"
	"github.com/kiselev-nikolay/direct-to-me/pkg/storage"
)

func ConnectAPI(ginServer *gin.Engine, strg storage.Storage) {
	ginServer.POST("/api/new", handlers.MakeNewRedirectHandler(strg))
	ginServer.GET("/api/list", handlers.MakeListRedirectsHandler(strg))
	ginServer.POST("/dev/print", handlers.MakeEchoHandler(strg))
	ginServer.StaticFile("/", "./assets/public/index.html")
	ginServer.Static("/static", "./assets/public/")
	redirectHandler := handlers.MakeRedirectHandler(strg)
	ginServer.NoMethod(redirectHandler)
	ginServer.NoRoute(redirectHandler)
	server.RunServer(ginServer)
}
