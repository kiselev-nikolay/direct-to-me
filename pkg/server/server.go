package server

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
)

func GetServer() *gin.Engine {
	r := gin.Default()

	// Staff handlers
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	r.GET("/favicon.ico", func(c *gin.Context) {
		c.String(http.StatusNotFound, "Chrome why?")
	})

	return r
}

func RunServer(r *gin.Engine, host string) {
	var serverWaitGroup sync.WaitGroup
	if host != "" {
		gin.SetMode(gin.ReleaseMode)
		serverWaitGroup.Add(1)
		go func() {
			defer serverWaitGroup.Done()
			autotls.Run(r, host)
		}()
	}
	serverWaitGroup.Add(1)
	go func() {
		defer serverWaitGroup.Done()
		r.Run()
	}()
	serverWaitGroup.Wait()
}
