package server

import (
	"flag"
	"net/http"
	"sync"

	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
)

// GetServer return server
func GetServer() *gin.Engine {
	r := gin.Default()

	// Staff handlers
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	r.GET("/favicon.ico", func(c *gin.Context) {
		c.String(http.StatusOK, "")
	})

	return r
}

// RunServer run server
func RunServer(r *gin.Engine) {
	var serverWaitGroup sync.WaitGroup
	var production bool
	flag.BoolVar(&production, "production", false, "Enables production mode, HTTPS with autoTLS.")
	flag.Parse()
	if production {
		gin.SetMode(gin.ReleaseMode)
		serverWaitGroup.Add(1)
		go func() {
			defer serverWaitGroup.Done()
			autotls.Run(r, "direct-to-me.com")
		}()
	}
	serverWaitGroup.Add(1)
	go func() {
		defer serverWaitGroup.Done()
		r.Run()
	}()
	serverWaitGroup.Wait()
}
