package api

import (
	"fmt"
	"github.com/DroiTaipei/droictx"
	"github.com/DroiTaipei/mgo/bson"
	"github.com/DroiTaipei/mongo"
	"github.com/valyala/fasthttp"
	"strconv"
	"sync"
)

const (
	SERVICE_NAME_PREFIX = "sh-sand-zone-"
	SERVICE_NAME_SUFFIX = ".tyd.svc.cluster.local"

	MGO_SANDBOX_ZONE_COL = "SandboxZonePodMapping"
	MGO_SANDBOX_APP_COL  = "SandboxAppZoneMapping"
	MGO_METRICS_COL      = "SandboxAccessMetrics"
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

func requestBypass(c *fasthttp.RequestCtx) {

	resp := &c.Response
	req := &c.Request

	prepareRequest(req)

	appid := string(req.Header.Peek("X-Droi-AppID"))

	ctx, _ := c.UserValue(DROI_CTX_KEY).(droictx.Context)

	queryResult := AppSlotMapping{}

	err := mongo.QueryOne(ctx, MGO_SANDBOX_APP_COL, &queryResult, bson.M{"appid": appid}, nil, 0, 10)
	if err != nil {
		fmt.Printf("db query error: %s\n", err)
		WriteError(c, err)
	}
	// fmt.Printf("Query result : %+v\n", queryResult)
	// fmt.Println(queryResult.AppID)

	// mutex.Lock()

	proxyClient = &fasthttp.HostClient{
		Addr: SERVICE_NAME_PREFIX + strconv.Itoa(queryResult.SandboxZoneID) + SERVICE_NAME_SUFFIX + ":" + strconv.Itoa(queryResult.Port),
	}

	// proxyClient.Addr = SERVICE_NAME_PREFIX + strconv.Itoa(queryResult.SandboxZoneID) + SERVICE_NAME_SUFFIX + ":" + strconv.Itoa(queryResult.Port)
	fmt.Println(SERVICE_NAME_PREFIX + strconv.Itoa(queryResult.SandboxZoneID) + SERVICE_NAME_SUFFIX + ":" + strconv.Itoa(queryResult.Port))

	if err := proxyClient.Do(req, resp); err != nil {
		c.Logger().Printf("error when proxying the request: %s\nRequest %+v\n", err, req)
		WriteError(c, err)
	}

	proxyClient = nil

	postprocessResponse(resp)
}
