package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kiselev-nikolay/direct-to-me/pkg/server"
	"github.com/kiselev-nikolay/direct-to-me/pkg/storage"
)

func ConnectAPI(ginServer *gin.Engine, fs *storage.FireStoreStorage) {
	ginServer.GET("/new", func(ctx *gin.Context) {
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
	})
	ginServer.GET("/list", func(ctx *gin.Context) {
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
	})
	ginServer.Any("/f/:key", func(ctx *gin.Context) {
		redirect, err := fs.GetRedirect(ctx.Param("key"))
		if err != nil {
			ctx.JSON(500, gin.H{
				"status": "database unreachable",
			})
			return
		}
		data := make(map[string]interface{})
		for k, v := range ctx.Request.URL.Query() {
			if len(v) == 1 {
				data[k] = v[0]
			} else {
				data[k] = v
			}
		}
		var bodyData map[string]interface{}
		for k, v := range bodyData {
			data[k] = v
		}
		for k, v := range ctx.Request.PostForm {
			if len(v) == 1 {
				data[k] = v[0]
			} else {
				data[k] = v
			}
		}
		if redirect.ToURL != "" {
			r, err := json.Marshal(data)
			if err != nil {
				ctx.JSON(400, gin.H{
					"status": "failed to process content",
				})
				return
			}
			go http.Post(redirect.ToURL, "application/json", bytes.NewBuffer(r))
		} else {
			req, err := processTemplate(redirect, &data)
			if err != nil {
				log.Print(err)
				ctx.JSON(400, gin.H{
					"status": "failed to process content",
				})
				return
			}
			go func() {
				HTTPClient := http.Client{}
				_, err := HTTPClient.Do(req)
				if err != nil {
					log.Print(err)
				}
			}()
		}
		ctx.Redirect(303, redirect.RedirectAfter)
	})
	ginServer.POST("/dev/print", func(ctx *gin.Context) {
		data, _ := ctx.GetRawData()
		fmt.Println(string(data))
	})
	ginServer.StaticFile("/", "./frontend/index.html")
	ginServer.StaticFile("/logo.png", "./frontend/logo.png")
	server.RunServer(ginServer)
}
