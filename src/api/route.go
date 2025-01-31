package api

import (
	"bytes"
	"fmt"
	"runtime"
	"time"

	ur "util/request"
	"util/trace"

	"github.com/DroiTaipei/droictx"
	"github.com/DroiTaipei/droitrace"
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
	HandlerFunc fasthttp.RequestHandler
}
type Routes []Route

var (
	routes = Routes{
		Route{
			"GET",
			"/health",
			HealthCheckHandler,
		},
	}

	requestRoutes = Routes{
		Route{
			"GET",
			"/*url",
			RequestHandler,
		},

		Route{
			"POST",
			"/*url",
			RequestHandler,
		},

		Route{
			"PATCH",
			"/*url",
			RequestHandler,
		},

		Route{
			"DELETE",
			"/*url",
			RequestHandler,
		},

		Route{
			"PUT",
			"/*url",
			RequestHandler,
		},

		Route{
			"OPTIONS",
			"/*url",
			RequestHandler,
		},
	}
)

func recv(ctx *fasthttp.RequestCtx) {
	if rcv := recover(); rcv != nil {
		logStackOnRecover(ctx, rcv)
		return
	}
}

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

func globalMiddleware(h fasthttp.RequestHandler, timeout int) fasthttp.RequestHandler {
	return func(c *fasthttp.RequestCtx) {
		var v []byte
		ctx := &droictx.DoneContext{}
		for headerKey, fieldkey := range ur.KeyMap {
			v = c.Request.Header.Peek(headerKey)
			if len(v) > 0 {
				ctx.Set(fieldkey, string(v))
			}
		}
		c.SetUserValue(DROI_CTX_KEY, ctx)
		HTTPAccessLog(ctx, string(c.Method()), string(c.Path()), c.RemoteAddr().String(), c.Request.Header.ContentLength())

		// Add trace
		sp := trace.CreateSpanFromReqF(droictx.ComponentForwarder, &c.Request, ctx)
		defer sp.Finish()
		ctx.Set(droitrace.ParentSpan, sp)

		doneCh := make(chan struct{})
		go func() {

			defer close(doneCh)
			defer recv(c)

			h(c)
		}()

		select {
		case <-doneCh:

			debugf("Request detail: request: %+v, with response: %+v", &c.Request, &c.Response)

		case <-time.After(time.Duration(timeout) * time.Second):

			WriteError(c, ErrForwardTimeout)

		}

		return
	}
}

func ApiRegist(timeout int) *fasthttprouter.Router {
	r := fasthttprouter.New()
	r.PanicHandler = logStackOnRecover
	var routingPath string
	for _, route := range routes {
		routingPath = fmt.Sprintf("/%s%s", API_VERSION, route.Pattern)
		r.Handle(route.Method, routingPath, globalMiddleware(route.HandlerFunc, timeout))
	}

	return r
}

func ForwarderRegist(timeout int) *fasthttprouter.Router {
	r := fasthttprouter.New()
	r.PanicHandler = logStackOnRecover
	var routingPath string
	for _, route := range requestRoutes {
		routingPath = fmt.Sprintf("%s", route.Pattern)
		debugf("%s : %s", route.Method, routingPath)
		r.Handle(route.Method, routingPath, globalMiddleware(route.HandlerFunc, timeout))
	}
	return r
}
