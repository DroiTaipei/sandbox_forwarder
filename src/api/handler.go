package api

import (
	"github.com/valyala/fasthttp"
)

const (
	MAX_KEY_LENGTH   = 100
	MAX_VALUE_LENGTH = 1024
	MAX_TTL          = 86400
)

func HealthCheckHandler(c *fasthttp.RequestCtx) {

	resp := NewResponse()
	resp.Write(c)
}

func RequestHandler(c *fasthttp.RequestCtx) {

	url := string(c.URI().Path())

	requestToGobuster(c, url)
}
