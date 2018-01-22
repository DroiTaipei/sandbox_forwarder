package da

import (
	de "github.com/DroiTaipei/droipkg"
)

var ErrUnknown de.DroiError
var ErrResourceNotFound de.DroiError
var ErrInvalidParameter de.DroiError
var ErrJsonValidation de.DroiError
var ErrPermissionDenied de.DroiError

var ErrMongoConn de.DroiError
func init() {
	ErrUnknown = de.ConstDroiError("1380001 Unknown")
	ErrResourceNotFound = de.ConstDroiError("1380002 Resource Not Found")
	ErrInvalidParameter = de.ConstDroiError("1380003 Parameter Validation Error")
	ErrJsonValidation = de.ConstDroiError("1380004 Json Validation Failed")
	ErrPermissionDenied = de.ConstDroiError("1380005 Permission Denied")
	ErrMongoConn = de.ConstDroiError("1380006 Mongo Connection Error")
}
