package api

import (
	"github.com/gin-gonic/gin"
	"github.com/kiselev-nikolay/direct-to-me/pkg/api/handlers"
	"github.com/kiselev-nikolay/direct-to-me/pkg/redirectstat"
	"github.com/kiselev-nikolay/direct-to-me/pkg/storage"
)

func ConnectAPI(ginServer *gin.Engine, authKey string, strg storage.Storage, redag *redirectstat.RedirectAggregation) {
	authorized := ginServer.Group("/api")
	authorized.Use(Auth(authKey))
	authorized.GET("/stats", handlers.MakeStatHandler(redag))
	authorized.POST("/new", handlers.MakeNewRedirectHandler(strg))
	authorized.GET("/list", handlers.MakeListRedirectsHandler(strg))
	authorized.POST("/dev/print", handlers.MakeEchoHandler())

	ginServer.StaticFile("/", "./assets/public/index.html")
	ginServer.GET("/bundle.js", authBundle(authKey))
	ginServer.Static("/static", "./assets/public/")
	redirectHandler := handlers.MakeRedirectHandler(strg)
	ginServer.NoMethod(redirectHandler)
	ginServer.NoRoute(redirectHandler)
}
