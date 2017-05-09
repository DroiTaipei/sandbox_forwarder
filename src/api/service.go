package api

import (
	"errors"
	"fmt"
	"strconv"
	"time"

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

func outputMetrics(c *fasthttp.RequestCtx) {
	// resp := &c.Response
	req := &c.Request

	c.Logger().Printf("%+v", req)
}

func updateAppMapping(c *fasthttp.RequestCtx) {
}

func requestToGobuster(c *fasthttp.RequestCtx, redirectURL string) {
	resp := &c.Response
	req := &c.Request

	prepareRequest(req)

	appid := string(req.Header.Peek("X-Droi-Service-AppID"))

	ctx, _ := c.UserValue(DROI_CTX_KEY).(droictx.Context)

	queryResult := AppSlotMapping{}

	err := mongo.QueryOne(ctx, MGO_SANDBOX_APP_COL, &queryResult, bson.M{"appid": appid, "status": APP_ACTIVE}, nil, 0, 10)
	if err != nil {
		c.Logger().Printf("query app failed: %s\n", err)
		WriteError(c, ErrAppNotFound)
		return
	}

	req.Header.Set("X-Droi-SlotID", strconv.Itoa(queryResult.SlotID))

	proxyClient := &fasthttp.HostClient{
		Addr: SERVICE_NAME_PREFIX + fmt.Sprintf("%04d", queryResult.SandboxZoneID) + SERVICE_NAME_SUFFIX + ":" + strconv.Itoa(GO_BUSTER_PORT),
	}

	req.SetRequestURI(getFullURI(redirectURL, c.URI().QueryString()))

	if err := proxyClient.Do(req, resp); err != nil {
		c.Logger().Printf("error when proxying the request: %s\nRequest: %+v\nResponse: %+v\n", err, req, resp)
		errorLog(errors.New(fmt.Sprintf("error when proxying the request: %s\nRequest: %+v\nResponse: %+v\n", err, req, resp)))
		WriteError(c, ErrForwardRequest)
		return
	}

	// Update sandbox access metrics
	upsertDoc := SandboxAccessInfo{
		Appid:      appid,
		UpdateTime: uint(time.Now().Unix()),
	}

	if _, err := mongo.Upsert(ctx, MGO_SANDBOX_METRICS_COL, bson.M{"appid": appid}, upsertDoc); err != nil {
		c.Logger().Printf("upsert metric failed: %s", err)
		WriteError(c, ErrDatabase)
		return
	}

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
