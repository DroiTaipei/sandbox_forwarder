package api

import (
	"sync"

	"github.com/valyala/fasthttp"
)

const (
	SERVICE_NAME_PREFIX = "sh-app-hl-sand-"
	SERVICE_NAME_SUFFIX = ".tyd.svc.cluster.local"

	GO_BUSTER_PORT = 8081

	MGO_SANDBOX_APP_COL     = "SandboxAppZoneMapping"
	MGO_SANDBOX_METRICS_COL = "SandboxAccessMetrics"

	APP_ACTIVE = "active"
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
	SlotID        int    `bson:"slotid" json:"slotid"`
	SandboxZoneID int    `bson:"szid" json:"szid"`
	Port          int    `bson:"port" json:"port"`
}

type SandboxAccessInfo struct {
	Appid                  string                 `bson:"appid" json:"appid"`
	UpdateTime             uint                   `bson:"last_update_at" json:"last_update_at"`
	LastHourAccessCount    int                    `bson:"last_hour_access" json:"last_hour_access"`
	LastQuarterAccessCount int                    `bson:"last_quarter_access" json:"last_quarter_access"`
	AppInfo                map[string]interface{} `bson:"app_info" json:"app_info"`
}
