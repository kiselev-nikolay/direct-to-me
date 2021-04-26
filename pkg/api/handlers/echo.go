package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func MakeEchoHandler() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		data, _ := ctx.GetRawData()
		fmt.Println(string(data))
	}
}
