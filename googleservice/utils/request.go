package utils

import (
	"github.com/valyala/fasthttp"
)

func BuildRequest(req *fasthttp.Request, method string, body []byte, apiKey string, url string) {
	req.SetBody(body)
	req.Header.SetMethod(method)
	req.Header.Set("API_KEY", apiKey)
	req.Header.SetContentType("application/json")
	req.SetRequestURI(url)
}
