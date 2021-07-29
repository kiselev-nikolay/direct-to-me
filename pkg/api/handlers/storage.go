package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kiselev-nikolay/direct-to-me/pkg/storage"
)

func MakeNewRedirectHandler(strg storage.Storage) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		redirect := storage.Redirect{}
		if err := ctx.Bind(&redirect); err != nil {
			log.Print(err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": err.Error(),
			})
			return
		}
		if redirect.ToURL != "" {
			// FIXIT validate template
			log.Print("FIXIT validate template")
		}
		err := strg.SetRedirect(redirect.FromURI, &redirect)
		if err != nil {
			log.Println(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status": "database unreachable",
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"status":   "ok",
			"accepted": ctx.Query("from"),
		})
	}
}

func MakeListRedirectsHandler(strg storage.Storage) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		redirects, err := strg.ListRedirects()
		if err != nil {
			if storage.IsNotFoundError(err) {
				ctx.JSON(http.StatusOK, gin.H{
					"status":    "ok",
					"redirects": []bool{},
				})
				return
			}
			log.Println(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status": "database unreachable",
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"redirects": redirects,
		})
	}
}
