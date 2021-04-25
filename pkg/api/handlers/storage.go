package handlers

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/kiselev-nikolay/direct-to-me/pkg/storage"
)

func MakeNewRedirectHandler(fs *storage.FireStoreStorage) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		redirect := storage.Redirect{}
		ctx.Bind(&redirect)
		if redirect.ToURL != "" {
			// FIXIT validate template
			log.Print("FIXIT validate template")
		}
		err := fs.SetRedirect(redirect.FromURI, &redirect)
		if err != nil {
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

func MakeListRedirectsHandler(fs *storage.FireStoreStorage) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		redirects, err := fs.ListRedirects()
		if err != nil {
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
