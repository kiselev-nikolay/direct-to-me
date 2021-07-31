package api

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	limit "github.com/yangxikun/gin-limit-by-key"

	"golang.org/x/time/rate"
)

type authHeader struct {
	Authorization string `header:"Authorization" binding:"required"`
}

const lenBearerAndToken = 2

func Auth(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := new(authHeader)
		if err := c.ShouldBindHeader(authorization); err != nil {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		tokens := strings.Split(authorization.Authorization, " ")
		if len(tokens) != lenBearerAndToken || tokens[0] != "Bearer" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if tokens[1] != key {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		c.Next()
	}
}

const tokenPlaceholder = "{{API_TOKEN}}"

func authBundle(authKey string) func(c *gin.Context) {
	return func(c *gin.Context) {
		file, err := os.Open("./assets/public/bundle.js")
		defer file.Close()
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusNotFound)
		}
		data, err := io.ReadAll(file)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusNotFound)
		}
		bundle := strings.ReplaceAll(string(data), tokenPlaceholder, authKey)
		bundleReader := strings.NewReader(bundle)
		c.DataFromReader(http.StatusOK, int64(len(bundle)), "text/javascript", bundleReader, nil)
	}
}

func GetLimiter() gin.HandlerFunc {
	return limit.NewRateLimiter(func(c *gin.Context) string {
		return c.ClientIP()
	}, func(c *gin.Context) (*rate.Limiter, time.Duration) {
		return rate.NewLimiter(rate.Every(time.Minute), 30), 15 * time.Minute
	}, func(c *gin.Context) {
		c.AbortWithStatus(429)
	})
}
