package api

import (
	"github.com/DroiTaipei/droipkg"
)

const (
	ErrUnknown          = droipkg.ConstDroiError("1160000 Unknown Error")
	ErrPanic            = droipkg.ConstDroiError("1160001 Panic Error")
	ErrAppNotFound      = droipkg.ConstDroiError("1160002 This app has been paused or not deployed yet. Please redeploy sandbox app again.")
	ErrForwardRequest   = droipkg.ConstDroiError("1160003 Forward Request Error")
	ErrForwardTimeout   = droipkg.ConstDroiError("1160004 Forward Request Timeout")
	ErrDatabase         = droipkg.ConstDroiError("1160005 Database Error")
	ErrAccessRestrictrd = droipkg.ConstDroiError("1160006 This app has been suspended, please contact to your account manager")
)

func getDroiErrorCode(err error) int {
	cause := droipkg.Cause(err)
	if de, ok := cause.(droipkg.DroiError); ok {
		return de.ErrorCode()
	}

	return ErrUnknown.ErrorCode()
}
