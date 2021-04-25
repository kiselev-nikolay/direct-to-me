package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/kiselev-nikolay/direct-to-me/pkg/storage"
)

func MakeEchoHandler(strg storage.Storage) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		data, _ := ctx.GetRawData()
		fmt.Println(string(data))
	}
}
