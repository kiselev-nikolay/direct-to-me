package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/kiselev-nikolay/direct-to-me/pkg/redirectstat"
)

func MakeStatHandler(redag *redirectstat.RedirectAggregation) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"status": "ok",
			"clicks": redag.Clicks,
			"fails":  redag.Fails,
		})
	}
}
