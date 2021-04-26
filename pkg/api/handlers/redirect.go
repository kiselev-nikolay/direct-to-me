package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	rs "github.com/kiselev-nikolay/direct-to-me/pkg/redirectstat"
	"github.com/kiselev-nikolay/direct-to-me/pkg/storage"
	template_process "github.com/kiselev-nikolay/direct-to-me/pkg/tools/template/process"
)

func IsSocialReferer(h string) bool {
	v, err := url.Parse(h)
	if err != nil {
		return false
	}
	domain := strings.Split(v.Hostname(), ".")
	if len(domain) < 2 {
		return false
	}
	switch strings.ToLower(domain[len(domain)-2]) {
	case "facebook", "instagram", "youtube", "twitter", "tiktok", "pinterest", "snapchat", "whatsapp":
		return true
	}
	return false
}

func MakeRedirectHandler(strg storage.Storage) func(ctx *gin.Context) {
	statCh := rs.GetStatChannels()
	return func(ctx *gin.Context) {
		requestURI := strings.Trim(ctx.Request.URL.Path, "/")
		redirect, err := strg.GetRedirect(requestURI)
		if err != nil {
			if storage.IsNotFoundError(err) {
				statCh.FailsChannel <- &rs.Fail{RedirectKey: requestURI, NotFound: 1}
				ctx.JSON(404, gin.H{
					"status": "not found",
				})
				return
			}
			statCh.FailsChannel <- &rs.Fail{RedirectKey: requestURI, DatabaseUnreachable: 1}
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
				statCh.FailsChannel <- &rs.Fail{RedirectKey: requestURI, ClientContentProcessFailed: 1}
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
				statCh.FailsChannel <- &rs.Fail{RedirectKey: requestURI, TemplateProcessFailed: 1}
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
					statCh.FailsChannel <- &rs.Fail{RedirectKey: requestURI, TemplateProcessFailed: 1}
					log.Println(err)
				}
			}()
		}
		if IsSocialReferer(ctx.Request.Header.Get("Referer")) {
			statCh.ClicksChannel <- &rs.Click{RedirectKey: requestURI, Social: 1}
		} else {
			statCh.ClicksChannel <- &rs.Click{RedirectKey: requestURI, Direct: 1}
		}
		ctx.Redirect(303, redirect.RedirectAfter)
	}
}
