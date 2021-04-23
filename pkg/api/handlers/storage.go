package handlers

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kiselev-nikolay/direct-to-me/pkg/storage"
)

func MakeNewRedirectHandler(fs *storage.FireStoreStorage) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		redirect := storage.Redirect{
			FromURI:         strings.TrimSpace(ctx.Query("from")),
			ToURL:           strings.Trim(ctx.Query("to"), "/"),
			RedirectAfter:   strings.TrimSpace(ctx.Query("after")),
			URLTemplate:     ctx.Query("urlTemplate"),
			MethodTemplate:  ctx.Query("methodTemplate"),
			HeadersTemplate: ctx.Query("headersTemplate"),
			BodyTemplate:    ctx.Query("bodyTemplate"),
		}
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
			"status": "ok",
			"keys":   redirects,
		})
	}
}
