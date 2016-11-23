package api

import (
	// "encoding/json"
	// "github.com/DroiTaipei/droictx"
	// "github.com/DroiTaipei/droipkg"
	"fmt"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

const (
	MAX_KEY_LENGTH   = 100
	MAX_VALUE_LENGTH = 1024
	MAX_TTL          = 86400
)

// func ReqGetMemcache(c *fasthttp.RequestCtx, ps fasthttprouter.Params) {
// 	ctx, _ := c.UserValue(DROI_CTX_KEY).(droictx.Context)
// 	resp := getMemcache(ctx, ps.ByName("key"))
// 	resp.Write(c)
// }
func ReceiveRequestBypassGobuster(c *fasthttp.RequestCtx, ps fasthttprouter.Params) {

	fmt.Println("Test receive request!")

	key := ps.ByName("url")

	fmt.Println(key)

	requestBypassGobuster(c, key)
}

func ReceiveRequest(c *fasthttp.RequestCtx, ps fasthttprouter.Params) {

	fmt.Println("Test receive request!")

	key := ps.ByName("url")

	fmt.Println(key)

	requestBypass(c, key)
}

func MetricsHandler(c *fasthttp.RequestCtx, ps fasthttprouter.Params) {

	outputMetrics(c)

	resp := NewResponse()
	resp.Write(c)
}

func HealthCheckHandler(c *fasthttp.RequestCtx, ps fasthttprouter.Params) {

	resp := NewResponse()
	resp.Write(c)
}
