package api

import (
	"fmt"
	"net"
	"strconv"
	"time"
	"util/trace"

	"github.com/DroiTaipei/droictx"
	"github.com/DroiTaipei/mgo/bson"
	"github.com/DroiTaipei/mongo"
	"github.com/valyala/fasthttp"
)

const qMark = "?"

func prepareRequest(req *fasthttp.Request) {
	// do not proxy "Connection" header.
	req.Header.Del("Connection")
}

func postprocessResponse(resp *fasthttp.Response) {
	// do not proxy "Connection" header
	resp.Header.Del("Connection")
}

func requestToGobuster(c *fasthttp.RequestCtx, redirectURL string) {
	resp := &c.Response
	req := &c.Request

	prepareRequest(req)

	appid := string(req.Header.Peek("X-Droi-Service-AppID"))

	ctx, _ := c.UserValue(DROI_CTX_KEY).(droictx.Context)

	queryResult := AppSlotMapping{}

	err := mongo.QueryOne(ctx, MGO_SANDBOX_DB_NAME, MGO_SANDBOX_APP_COL, &queryResult, bson.M{"appid": appid, "status": APP_ACTIVE}, nil, 0, 10)
	if err != nil {
		wrapErr := ErrAppNotFound.Wrap(fmt.Sprintf("App ID %s not found: %+v", appid, err))
		ctxLog(ctx, wrapErr)
		WriteError(c, ErrAppNotFound)
		return
	}

	if queryResult.SandboxZoneID == 3000 {
		wrapErr := ErrAppNotFound.Wrap(fmt.Sprintf("This App:%s has been suspended", appid))
		ctxLog(ctx, wrapErr)
		WriteError(c, ErrAccessRestrictrd)
		return
	}

	req.Header.Set("X-Droi-SlotID", strconv.Itoa(queryResult.SlotID))

	proxyClient := &fasthttp.HostClient{
		Addr: SERVICE_NAME_PREFIX + fmt.Sprintf("%04d", queryResult.SandboxZoneID) + SERVICE_NAME_SUFFIX + ":" + strconv.Itoa(GO_BUSTER_PORT),
	}

	ipList, lookErr := net.LookupIP(SERVICE_NAME_PREFIX + fmt.Sprintf("%04d", queryResult.SandboxZoneID) + SERVICE_NAME_SUFFIX)
	if lookErr != nil {
		wrapErr := ErrForwardRequest.Wrap(fmt.Sprintf("Lookup domain error: %+v", lookErr))
		ctxLog(ctx, wrapErr)
		WriteError(c, ErrForwardRequest)
		return
	}
	debugf("Lookup IP: %+v", ipList)

	req.SetRequestURI(getFullURI(redirectURL, c.URI().QueryString()))
	sp := trace.CreateChildSpanByContextF(droictx.ComponentForwarder, ctx, req)
	defer sp.Finish()
	trace.InjectSpanF(sp, req)

	if err := proxyClient.Do(req, resp); err != nil {
		wrapErr := ErrForwardRequest.Wrap(fmt.Sprintf("Proxying request error: %s\nRequest: %#v\nResponse: %#v\n", err.Error(), req, resp))
		ctxLog(ctx, wrapErr)
		WriteError(c, ErrForwardRequest)
		return
	}

	// Update sandbox access metrics
	upsertDoc := SandboxAccessInfo{
		Appid:      appid,
		UpdateTime: uint(time.Now().Unix()),
	}

	if _, err := mongo.Upsert(ctx, MGO_SANDBOX_DB_NAME, MGO_SANDBOX_METRICS_COL, bson.M{"appid": appid}, upsertDoc); err != nil {
		wrapErr := ErrDatabase.Wrap(fmt.Sprintf("Update mongo error: %#v", err))
		ctxLog(ctx, wrapErr)
		WriteError(c, ErrDatabase)
		return
	}

	RequestTotal.WithLabelValues(appid).Inc()

	proxyClient = nil

	postprocessResponse(resp)
}

func getFullURI(URL string, queryBuf []byte) string {
	n := len(URL) + len(queryBuf) + 1
	bs := make([]byte, n)
	bl := 0
	bl += copy(bs[bl:], URL)
	bl += copy(bs[bl:], qMark)
	bl += copy(bs[bl:], queryBuf)
	return string(bs)
}
