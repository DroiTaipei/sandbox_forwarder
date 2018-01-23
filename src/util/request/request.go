package request

import (
	"github.com/DroiTaipei/droictx"
	"github.com/valyala/fasthttp"
)

const (
	KEY_APP_ID           = "X-Droi-AppID"
	KEY_DEVICE_ID        = "X-Droi-DeviceID"
	KEY_USER_ID          = "X-Droi-UserID"
	KEY_APP_FILTER       = "X-Droi-AppFilter"
	KEY_REQ_ID           = "X-Droi-ReqID"
	KEY_PLATFORM_KEY     = "X-Droi-Platform-Key"
	KEY_APP_ID_MODE      = "X-Droi-AidMode"
	LOG_KEY_APP_ID       = "Aid"
	LOG_KEY_APP_ID_MODE  = "Aidm"
	LOG_KEY_DEVICE_ID    = "Did"
	LOG_KEY_USER_ID      = "Uid"
	LOG_KEY_REQ_ID       = "Rid"
	LOG_KEY_PLATFORM_KEY = "XPk"
	LOG_KEY_URL          = "URL"
)

var KeyMap map[string]string

func init() {
	KeyMap = map[string]string{
		KEY_APP_ID:       LOG_KEY_APP_ID,
		KEY_APP_ID_MODE:  LOG_KEY_APP_ID_MODE,
		KEY_DEVICE_ID:    LOG_KEY_DEVICE_ID,
		KEY_USER_ID:      LOG_KEY_USER_ID,
		KEY_REQ_ID:       LOG_KEY_REQ_ID,
		KEY_PLATFORM_KEY: LOG_KEY_PLATFORM_KEY,
	}

}

func fetchFromHeader(c *fasthttp.RequestHeader) droictx.Context {
	r := &droictx.DoneContext{}
	keyMap := map[string]string{
		KEY_APP_ID:       LOG_KEY_APP_ID,
		KEY_APP_ID_MODE:  LOG_KEY_APP_ID_MODE,
		KEY_DEVICE_ID:    LOG_KEY_DEVICE_ID,
		KEY_USER_ID:      LOG_KEY_USER_ID,
		KEY_REQ_ID:       LOG_KEY_REQ_ID,
		KEY_PLATFORM_KEY: LOG_KEY_PLATFORM_KEY,
	}
	r.Set(LOG_KEY_REQ_ID, "")
	var v string
	for headerKey, fieldkey := range keyMap {
		v = string(c.Peek(headerKey))
		if v != "" {
			r.Set(fieldkey, v)
		}
	}
	return r
}

func FetchFromRequest(c *fasthttp.Request, method []byte) droictx.Context {
	ctx := fetchFromHeader(&c.Header)
	ctx.Set(LOG_KEY_URL, string(append(method, c.URI().FullURI()[:]...)))
	return ctx
}
