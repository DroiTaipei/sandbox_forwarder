package mimir

import (
	de "github.com/DroiTaipei/droipkg"
)

var ErrUnknown de.DroiError
var ErrResourceNotFound de.DroiError
var ErrInvalidParameter de.DroiError
var ErrJsonValidation de.DroiError
var ErrPermissionDenied de.DroiError

func init() {
	ErrUnknown = de.ConstDroiError("1370001 Unknown")
	ErrResourceNotFound = de.ConstDroiError("1370002 Resource Not Found")
	ErrInvalidParameter = de.ConstDroiError("1370003 Parameter Validation Error")
	ErrJsonValidation = de.ConstDroiError("1370004 Json Validation Failed")
	ErrPermissionDenied = de.ConstDroiError("1370005 Permission Denied")
}
