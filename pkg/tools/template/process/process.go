package process

import (
	"bufio"
	"bytes"
	"errors"
	"net/http"
	"net/textproto"
	"strings"
	"text/template"

	"github.com/kiselev-nikolay/direct-to-me/pkg/storage"
)

type parsedRedirect struct {
	method     string
	rawURL     string
	rawHeaders string
	body       string
}

func processTemplatePart(t string, data map[string]interface{}) (string, error) {
	temp, err := template.New("XTemplate").Parse(t)
	if err != nil {
		return "", err
	}
	b := bytes.Buffer{}
	err = temp.Execute(&b, data)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

func readRedirectTemplate(r *storage.Redirect, data map[string]interface{}) (*parsedRedirect, error) {
	method, err := processTemplatePart(r.MethodTemplate, data)
	if err != nil {
		return nil, errors.New("Wrong template for Method: " + err.Error())
	}
	rawURL, err := processTemplatePart(r.URLTemplate, data)
	if err != nil {
		return nil, errors.New("Wrong template for URL: " + err.Error())
	}
	rawHeaders, err := processTemplatePart(r.HeadersTemplate, data)
	if err != nil {
		return nil, errors.New("Wrong template for Headers: " + err.Error())
	}
	body, err := processTemplatePart(r.BodyTemplate, data)
	if err != nil {
		return nil, errors.New("Wrong template for Body: " + err.Error())
	}
	return &parsedRedirect{method, rawURL, rawHeaders, body}, nil
}

func ProcessTemplate(r *storage.Redirect, data map[string]interface{}) (*http.Request, error) {
	parsed, err := readRedirectTemplate(r, data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(parsed.method, parsed.rawURL, strings.NewReader(parsed.body))
	if err != nil {
		return nil, errors.New("Cannot build request: " + err.Error())
	}
	if strings.Count(parsed.rawHeaders, ":") > 0 {
		headersMIME, err := textproto.NewReader(bufio.NewReader(strings.NewReader(parsed.rawHeaders))).ReadMIMEHeader()
		if err != nil {
			return nil, errors.New("Cannot read headers: " + err.Error())
		}
		for k, v := range headersMIME {
			req.Header[k] = v
		}
	}
	return req, nil
}

func ValidateRedirectTemplate(r *storage.Redirect, data map[string]interface{}) error {
	_, err := readRedirectTemplate(r, data)
	return err
}
