package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kiselev-nikolay/direct-to-me/pkg/storage"
	template_process "github.com/kiselev-nikolay/direct-to-me/pkg/tools/template/process"
)

func MakeRedirectHandler(strg storage.Storage) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		requestURI := strings.Trim(ctx.Request.URL.Path, "/")
		redirect, err := strg.GetRedirect(requestURI)
		if err != nil {
			if storage.IsNotFoundError(err) {
				ctx.JSON(404, gin.H{
					"status": "not found",
				})
				return
			}
			log.Println(err)
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
				log.Println(err)
				ctx.JSON(400, gin.H{
					"status": "failed to process content",
				})
				return
			}
			go http.Post(redirect.ToURL, "application/json", bytes.NewBuffer(r))
		} else {
			req, err := template_process.ProcessTemplate(redirect, &data)
			if err != nil {
				log.Println(err)
				ctx.JSON(400, gin.H{
					"status": "failed to process content",
				})
				return
			}
			go func() {
				HTTPClient := http.Client{}
				_, err := HTTPClient.Do(req)
				if err != nil {
					log.Println(err)
				}
			}()
		}
		ctx.Redirect(303, redirect.RedirectAfter)
	}
}
