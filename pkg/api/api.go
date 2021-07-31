package api

import (
	"github.com/gin-gonic/gin"
	"github.com/kiselev-nikolay/direct-to-me/pkg/api/handlers"
	"github.com/kiselev-nikolay/direct-to-me/pkg/redirectstat"
	"github.com/kiselev-nikolay/direct-to-me/pkg/storage"
)

func ConnectAPI(ginServer *gin.Engine, authKey string, strg storage.Storage, redag *redirectstat.RedirectAggregation) {
	limiter := GetLimiter()
	authorized := ginServer.Group("/api")
	authorized.Use(Auth(authKey))
	authorized.Use(limiter)
	authorized.GET("/stats", handlers.MakeStatHandler(redag))
	authorized.POST("/new", handlers.MakeNewRedirectHandler(strg))
	authorized.GET("/list", handlers.MakeListRedirectsHandler(strg))
	authorized.POST("/dev/print", handlers.MakeEchoHandler())

	frontend := ginServer.Group("control", gin.BasicAuth(gin.Accounts{
		"staff": authKey,
	}))
	frontend.Use(limiter)
	frontend.StaticFile("/", "./assets/public/index.html")
	frontend.GET("/bundle.js", authBundle(authKey))
	frontend.Static("/static", "./assets/public/")
	redirectHandler := handlers.MakeRedirectHandler(strg)
	ginServer.NoMethod(redirectHandler)
	ginServer.NoRoute(redirectHandler)
}
