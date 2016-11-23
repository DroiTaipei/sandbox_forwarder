package api

import (
	"fmt"
	"github.com/DroiTaipei/droictx"
	"github.com/DroiTaipei/mgo/bson"
	"github.com/DroiTaipei/mongo"
	"github.com/valyala/fasthttp"
	"strconv"
	"sync"
	"time"
)

const (
	SERVICE_NAME_PREFIX = "sh-app-hl-sand-"
	SERVICE_NAME_SUFFIX = ".tyd.svc.cluster.local"

	GO_BUSTER_PORT = 8081

	MGO_SANDBOX_APP_COL     = "SandboxAppZoneMapping"
	MGO_SANDBOX_METRICS_COL = "SandboxAccessMetrics"
)

var (
	proxyClient = &fasthttp.HostClient{
		Addr: "10.128.112.184:8080",
		// set other options here if required - most notably timeouts.
	}
	mutex = &sync.Mutex{}

	connInfo map[string]AppSlotMapping
)

type AppSlotMapping struct {
	AppID         string `bson:"appid" json:"appid"`
	SlotID        int    `bson:"slot" json:"slot"`
	SandboxZoneID int    `bson:"szid" json:"szid"`
	Port          int    `bson:"port" json:"port"`
}

func prepareRequest(req *fasthttp.Request) {
	// do not proxy "Connection" header.

	req.Header.Del("Connection")

	// req.SetRequestURI("/api/v1/cluster/summary")
}

func postprocessResponse(resp *fasthttp.Response) {
	// do not proxy "Connection" header
	resp.Header.Del("Connection")
}

func requestBypassGobuster(c *fasthttp.RequestCtx, redirectURL string) {

	resp := &c.Response
	req := &c.Request

	prepareRequest(req)

	appid := string(req.Header.Peek("X-Droi-AppID"))

	ctx, _ := c.UserValue(DROI_CTX_KEY).(droictx.Context)

	queryResult := AppSlotMapping{}

	err := mongo.QueryOne(ctx, MGO_SANDBOX_APP_COL, &queryResult, bson.M{"appid": appid}, nil, 0, 10)
	if err != nil {
		fmt.Printf("db query error: %s\n", err)
		WriteError(c, ErrAppNotFound)
		return
	}

	fmt.Println(SERVICE_NAME_PREFIX + strconv.Itoa(queryResult.SandboxZoneID) + SERVICE_NAME_SUFFIX + ":" + strconv.Itoa(queryResult.Port))

	proxyClient = &fasthttp.HostClient{
		Addr: SERVICE_NAME_PREFIX + strconv.Itoa(queryResult.SandboxZoneID) + SERVICE_NAME_SUFFIX + ":" + strconv.Itoa(queryResult.Port),
	}

	// proxyClient = &fasthttp.HostClient{
	// 	Addr: "tpe-db-baas-mgo.tyd.svc.cluster.local:8081",
	// }

	req.SetRequestURI(redirectURL)

	if err := proxyClient.Do(req, resp); err != nil {
		c.Logger().Printf("error when proxying the request: %s\nRequest %+v\n", err, req)
		WriteError(c, err)
		return
	}

	fmt.Printf("Request for debug: %+v\n", req)
	fmt.Printf("Response for debug: %+v\n", resp)

	// Update sandbox access metrics
	upsertDoc := bson.M{
		"appid":          appid,
		"last_update_at": uint(time.Now().Unix()),
	}

	if _, err := mongo.Upsert(ctx, MGO_SANDBOX_METRICS_COL, bson.M{"appid": appid}, upsertDoc); err != nil {
		fmt.Printf("db query error: %s\n", err)
		WriteError(c, err)
		return
	}

	proxyClient = nil

	postprocessResponse(resp)
}

func requestBypass(c *fasthttp.RequestCtx, redirectURL string) {

	resp := &c.Response
	req := &c.Request

	prepareRequest(req)

	appid := string(req.Header.Peek("X-Droi-AppID"))

	ctx, _ := c.UserValue(DROI_CTX_KEY).(droictx.Context)

	queryResult := AppSlotMapping{}

	err := mongo.QueryOne(ctx, MGO_SANDBOX_APP_COL, &queryResult, bson.M{"appid": appid}, nil, 0, 10)
	if err != nil {
		fmt.Printf("db query error: %s\n", err)
		WriteError(c, ErrAppNotFound)
		return
	}

	fmt.Println(SERVICE_NAME_PREFIX + strconv.Itoa(queryResult.SandboxZoneID) + SERVICE_NAME_SUFFIX + ":" + strconv.Itoa(GO_BUSTER_PORT))

	proxyClient = &fasthttp.HostClient{
		Addr: SERVICE_NAME_PREFIX + strconv.Itoa(queryResult.SandboxZoneID) + SERVICE_NAME_SUFFIX + ":" + strconv.Itoa(GO_BUSTER_PORT),
	}

	// proxyClient = &fasthttp.HostClient{
	// 	Addr: "tpe-db-baas-mgo.tyd.svc.cluster.local:8081",
	// }

	req.SetRequestURI(redirectURL)

	if err := proxyClient.Do(req, resp); err != nil {
		c.Logger().Printf("error when proxying the request: %s\nRequest %+v\n", err, req)
		WriteError(c, err)
		return
	}

	fmt.Printf("Request for debug: %+v\n", req)
	fmt.Printf("Response for debug: %+v\n", resp)

	// Update sandbox access metrics
	upsertDoc := bson.M{
		"appid":          appid,
		"last_update_at": uint(time.Now().Unix()),
	}

	if _, err := mongo.Upsert(ctx, MGO_SANDBOX_METRICS_COL, bson.M{"appid": appid}, upsertDoc); err != nil {
		fmt.Printf("db query error: %s\n", err)
		WriteError(c, err)
		return
	}

	proxyClient = nil

	postprocessResponse(resp)
}

func outputMetrics(c *fasthttp.RequestCtx) {

}
