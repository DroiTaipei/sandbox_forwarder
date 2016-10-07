package api

import (
	"github.com/DroiTaipei/droipkg"
)

const (
	ErrUnknown             = droipkg.ConstDroiError("1130000 Unknown Error")
	ErrPanic               = droipkg.ConstDroiError("1130001 Panic Error")
	ErrResourceNotFound    = droipkg.ConstDroiError("1130002 Resouce Not Found")
	ErrParameterValidation = droipkg.ConstDroiError("1130003 Parameter Validation Failed")
	ErrJsonValidation      = droipkg.ConstDroiError("1130004 Json Validation Failed")
	ErrPermissionDenied    = droipkg.ConstDroiError("1130005 Permission Denied")
	ErrDatabaseUnavailable = droipkg.ConstDroiError("1130006 Database Unavailable")
	ErrDatabase            = droipkg.ConstDroiError("1130007 Database Error")
	ErrDataProcessFailed   = droipkg.ConstDroiError("1130008 Data Process Failed")
	ErrDataNotFound        = droipkg.ConstDroiError("1130010 Data Not Found")
	ErrResourceBound       = droipkg.ConstDroiError("1130011 Resource Bound")
)

func getDroiErrorCode(err error) int {
	cause := droipkg.Cause(err)
	if de, ok := cause.(droipkg.DroiError); ok {
		return de.ErrorCode()
	}

	return ErrUnknown.ErrorCode()
}
