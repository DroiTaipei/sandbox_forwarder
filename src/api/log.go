package api

import (
	"github.com/DroiTaipei/droictx"
	"github.com/DroiTaipei/droipkg"
)

const (
	METHOD_FIELD   = "Md"
	URI_FIELD      = "Uri"
	CLIENT_ADDRESS = "CI"
	REQUEST_LENGTH = "Rl"
	REQUEST_ID     = "Rid"
)

func errorLog(err error) {
	droipkg.GetLogger().Error(err.Error())
}

func ctxLog(ctx droictx.Context, v ...interface{}) {
	droipkg.GetLogger().WithMap(ctx.Map()).Error(v...)
}

func HTTPAccessLog(ctx droictx.Context, method, uri, addr string, length int) {
	droipkg.GetLogger().WithMap(ctx.Map()).
		WithField(METHOD_FIELD, method).
		WithField(URI_FIELD, uri).
		WithField(CLIENT_ADDRESS, addr).
		WithField(REQUEST_LENGTH, length).
		Info("HTTP Access")
}

func debugf(format string, args ...interface{}) {
	droipkg.GetLogger().Debugf(format, args...)
}
