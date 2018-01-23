package trace

import (
	"net"
	"strconv"

	"github.com/DroiTaipei/droictx"
	"github.com/DroiTaipei/droitrace"
	opentracing "github.com/DroiTaipei/opentracing-go"
	ext "github.com/DroiTaipei/opentracing-go/ext"
	"github.com/valyala/fasthttp"
)

func CreateSpanFromReqF(spanName string, c *fasthttp.Request, ctx droictx.Context) opentracing.Span {
	var sp opentracing.Span
	wireContext, err := opentracing.GlobalTracer().Extract(
		opentracing.TextMap,
		opentracing.TextMapCarrier(ctx.HeaderMap()))
	if err != nil {
		sp = opentracing.StartSpan(spanName)
	} else {
		sp = opentracing.StartSpan(
			spanName,
			ext.RPCServerOption(wireContext))
	}
	attachSpanTagsF(sp, c)
	return sp
}

func CreateRootSpanF(spanName string, c *fasthttp.Request) opentracing.Span {
	sp := opentracing.StartSpan(spanName)
	attachSpanTagsF(sp, c)
	return sp
}

func CreateChildSpanF(spanName string, parentSpan opentracing.Span, c *fasthttp.Request) opentracing.Span {
	sp := opentracing.StartSpan(
		spanName,
		opentracing.ChildOf(parentSpan.Context()))
	attachSpanTagsF(sp, c)
	return sp
}

func CreateChildSpanByContextF(spanName string, ctx droictx.Context, c *fasthttp.Request) opentracing.Span {
	parentSpanTmp := ctx.Get(droitrace.ParentSpan)
	if parentSpanTmp == nil {
		return CreateRootSpanF(spanName, c)
	}
	parentSpan, ok := parentSpanTmp.(opentracing.Span)
	if !ok {
		return CreateRootSpanF(spanName, c)
	}
	sp := opentracing.StartSpan(
		spanName,
		opentracing.ChildOf(parentSpan.Context()))

	attachSpanTagsF(sp, c)
	return sp
}

func CreateFollowFromSpanByContextF(spanName string, ctx droictx.Context, c *fasthttp.Request) opentracing.Span {
	parentSpanTmp := ctx.Get(droitrace.ParentSpan)
	if parentSpanTmp == nil {
		return CreateRootSpanF(spanName, c)
	}
	parentSpan, ok := parentSpanTmp.(opentracing.Span)
	if !ok {
		return CreateRootSpanF(spanName, c)
	}
	sp := opentracing.StartSpan(
		spanName,
		opentracing.FollowsFrom(parentSpan.Context()))
	attachSpanTagsF(sp, c)
	return sp
}

func attachSpanTagsF(sp opentracing.Span, c *fasthttp.Request) {
	ext.HTTPMethod.Set(sp, string(c.Header.Method()))
	ext.HTTPUrl.Set(sp, c.URI().String())
	host := string(c.URI().Host())
	if host, portString, err := net.SplitHostPort(host); err == nil {
		ext.PeerHostname.Set(sp, host)
		if port, err := strconv.Atoi(portString); err != nil {
			ext.PeerPort.Set(sp, uint16(port))
		}
	} else {
		ext.PeerHostname.Set(sp, host)
	}
	SetDroiTagFromHeadersF(sp, &c.Header)
	return
}

func InjectSpanF(sp opentracing.Span, req *fasthttp.Request) error {
	if err := sp.Tracer().Inject(sp.Context(), opentracing.TextMap, &req.Header); err != nil {
		return err
	}
	return nil
}

func SetDroiTagFromHeadersF(span opentracing.Span, headers *fasthttp.RequestHeader) {
	for hk, sk := range droictx.IFieldHeaderKeyMap() {
		if v := headers.Peek(hk); len(v) > 0 {
			tag := droitrace.GenDroiTag(sk)
			span.SetTag(tag, string(v))
		}
	}
}
