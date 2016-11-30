package api

import (
	"bytes"
	"fmt"
	"runtime"

	ur "util/request"

	"github.com/DroiTaipei/droictx"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

const (
	API_VERSION  = "v1"
	DROI_CTX_KEY = "DroiCtx"
)

type Route struct {
	Method      string
	Pattern     string
	HandlerFunc fasthttprouter.Handle
}
type Routes []Route

var (
	routes = Routes{
		Route{
			"GET",
			"/health",
			HealthCheckHandler,
		},

		Route{
			"POST",
			"/sandbox/*url",
			ReceiveRequestBypassGobuster,
		},

		Route{
			"POST",
			"/gobuster/*url",
			ReceiveRequest,
		},

		Route{
			"GET",
			"/metrics",
			MetricsHandler,
		},
	}

	requestRoutes = Routes{
		Route{
			"POST",
			"/*url",
			RequestHandler,
		},
	}
)

func logStackOnRecover(ctx *fasthttp.RequestCtx, r interface{}) {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("[Panic] recover from panic situation: - %v\r\n", r))
	for i := 2; ; i += 1 {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		buffer.WriteString(fmt.Sprintf("    %s:%d\r\n", file, line))

	}
	logHeaders := ur.FetchFromRequest(&ctx.Request, ctx.Method())
	ctxLog(logHeaders, buffer.String())
	WriteError(ctx, ErrPanic)
	// TODO: set an error page and more info to kafka
	return
}

func gloablMiddleware(h fasthttprouter.Handle) fasthttprouter.Handle {
	return fasthttprouter.Handle(func(c *fasthttp.RequestCtx, ps fasthttprouter.Params) {
		var v []byte
		var ctx droictx.Context
		for headerKey, fieldkey := range ur.KeyMap {
			v = c.Request.Header.Peek(headerKey)
			if len(v) > 0 {
				ctx.Set(fieldkey, string(v))
			}
		}
		c.SetUserValue(DROI_CTX_KEY, ctx)
		HTTPAccessLog(ctx, string(c.Method()), string(c.Path()), c.RemoteAddr().String(), c.Request.Header.ContentLength())
		h(c, ps)
		return
	})
}

func ApiRegist() *fasthttprouter.Router {
	r := fasthttprouter.New()
	r.PanicHandler = logStackOnRecover
	var routingPath string
	for _, route := range routes {
		routingPath = fmt.Sprintf("/%s%s", API_VERSION, route.Pattern)
		debugf("%s : %s", route.Method, routingPath)
		r.Handle(route.Method, routingPath, gloablMiddleware(route.HandlerFunc))
	}

	return r
}

func ForwarderRegist() *fasthttprouter.Router {
	r := fasthttprouter.New()
	r.PanicHandler = logStackOnRecover
	var routingPath string
	for _, route := range requestRoutes {
		routingPath = fmt.Sprintf("%s", route.Pattern)
		debugf("%s : %s", route.Method, routingPath)
		r.Handle(route.Method, routingPath, gloablMiddleware(route.HandlerFunc))
	}
	return r
}
