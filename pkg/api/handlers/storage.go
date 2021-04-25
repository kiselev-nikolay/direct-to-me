package handlers

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/kiselev-nikolay/direct-to-me/pkg/storage"
)

func MakeNewRedirectHandler(strg storage.Storage) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		redirect := storage.Redirect{}
		ctx.Bind(&redirect)
		if redirect.ToURL != "" {
			// FIXIT validate template
			log.Print("FIXIT validate template")
		}
		err := strg.SetRedirect(redirect.FromURI, &redirect)
		if err != nil {
			log.Println(err)
			ctx.JSON(500, gin.H{
				"status": "database unreachable",
			})
			return
		}
		ctx.JSON(200, gin.H{
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
				ctx.JSON(200, gin.H{
					"status":    "ok",
					"redirects": []bool{},
				})
				return
			}
			log.Println(err)
			ctx.JSON(500, gin.H{
				"status": "database unreachable",
			})
			return
		}
		ctx.JSON(200, gin.H{
			"status":    "ok",
			"redirects": redirects,
		})
	}
}
