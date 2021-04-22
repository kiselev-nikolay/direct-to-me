package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kiselev-nikolay/direct-to-me/pkg/plain_http"
	"github.com/kiselev-nikolay/direct-to-me/pkg/server"
	"github.com/kiselev-nikolay/direct-to-me/pkg/storage"
)

func main() {
	ginServer := server.GetServer()
	storage, err := storage.NewKVStorage(nil, 1*time.Minute)
	defer storage.close()
	if err != nil {
		log.Fatal(err)
	}
	ginServer.GET("/new", func(ctx *gin.Context) {
		data := map[string]string{
			"fromURI":         strings.TrimSpace(ctx.Query("from")),
			"toURL":           strings.Trim(ctx.Query("to"), "/"),
			"redirectAfter":   strings.TrimSpace(ctx.Query("after")),
			"urlTemplate":     ctx.Query("urlTemplate"),
			"methodTemplate":  ctx.Query("methodTemplate"),
			"headersTemplate": ctx.Query("headersTemplate"),
			"bodyTemplate":    ctx.Query("bodyTemplate"),
		}
		if data["toURL"] != "" {
			_, err := plain_http.MakeHTTPTemplates(data)
			if err != nil {
				ctx.JSON(400, gin.H{
					"status":   "bad template",
					"accepted": err.Error(),
				})
			}
		}
		storage.Set(data["fromURI"], data)
		ctx.JSON(200, gin.H{
			"status":   "ok",
			"accepted": ctx.Query("from"),
		})
	})
	ginServer.GET("/list", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"status": "ok",
			"keys":   storage.List(),
		})
	})
	ginServer.Any("/f/:key", func(ctx *gin.Context) {
		databaseRecord := storage.Get(ctx.Param("key"))
		data := make(map[string]interface{})
		for k, v := range ctx.Request.URL.Query() {
			if len(v) == 1 {
				data[k] = v[0]
			} else {
				data[k] = v
			}
		}
		var bodyData map[string]interface{}
		err := ctx.ShouldBind(&bodyData)
		if err != nil {
			ctx.JSON(400, gin.H{
				"status": "cannot read body. " + err.Error(),
			})
			return
		}
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
		if databaseRecord["toURL"] != "" {
			r, err := json.Marshal(data)
			if err != nil {
				ctx.JSON(400, gin.H{
					"status": "failed to process content",
				})
				return
			}
			go http.Post(databaseRecord["toURL"], "application/json", bytes.NewBuffer(r))
		} else {
			templateMask, _ := plain_http.MakeHTTPTemplates(databaseRecord)
			plainRequest, plainRequestURL, plainRequestBody, err := plain_http.BuildHTTP(templateMask, data)
			if err != nil {
				// TODO say about data was wrong
				ctx.JSON(400, gin.H{
					"status": "failed to process content",
				})
				return
			}
			go func() {
				client := &http.Client{}
				request, err := http.ReadRequest(bufio.NewReader(strings.NewReader(plainRequest)))
				if err != nil {
					fmt.Println("sub request read error:", err)
					return
				}
				newRequest, err := http.NewRequest(request.Method, plainRequestURL, strings.NewReader(plainRequestBody))
				if err != nil {
					fmt.Println("sub request go-read error:", err)
					return
				}
				newRequest.Header = request.Header
				client.Do(newRequest)
			}()
			ctx.Redirect(303, databaseRecord["redirectAfter"])
		}
		ctx.Redirect(303, databaseRecord["redirectAfter"])
	})
	ginServer.POST("/dev/print", func(ctx *gin.Context) {
		data, _ := ctx.GetRawData()
		fmt.Println(string(data))
	})
	ginServer.StaticFile("/", "../frontend/index.html")
	ginServer.StaticFile("/logo.png", "../frontend/logo.png")
	server.RunServer(ginServer)
}
