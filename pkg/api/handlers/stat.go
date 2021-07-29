package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kiselev-nikolay/direct-to-me/pkg/redirectstat"
)

func MakeStatHandler(redag *redirectstat.RedirectAggregation) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"clicks": redag.Clicks,
			"fails":  redag.Fails,
		})
	}
}
