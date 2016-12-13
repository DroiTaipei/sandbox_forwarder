package api

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

const (
	MAX_KEY_LENGTH   = 100
	MAX_VALUE_LENGTH = 1024
	MAX_TTL          = 86400
)

func HealthCheckHandler(c *fasthttp.RequestCtx, ps fasthttprouter.Params) {

	resp := NewResponse()
	resp.Write(c)
}

func ReceiveRequestBypassGobuster(c *fasthttp.RequestCtx, ps fasthttprouter.Params) {

	key := ps.ByName("url")

	requestBypassGobuster(c, key)
}

func ReceiveRequest(c *fasthttp.RequestCtx, ps fasthttprouter.Params) {

	key := ps.ByName("url")

	requestBypass(c, key)
}

func MetricsHandler(c *fasthttp.RequestCtx, ps fasthttprouter.Params) {

	outputMetrics(c)

	resp := NewResponse()
	resp.Write(c)
}

func RequestHandler(c *fasthttp.RequestCtx, ps fasthttprouter.Params) {

	url := ps.ByName("url")

	requestToGobuster(c, url)
}
