package file

import (
	derr "github.com/DroiTaipei/droipkg"
)

var (
	// Global
	OkCode = 0
	// Server error
	ErrInternal            = derr.CarrierDroiError("1180000 Server internal error")
	ErrNewLogHandler       = derr.CarrierDroiError("1180001 Initialize log handler error")
	ErrInvalidHash         = derr.CarrierDroiError("1180002 Hash code is wrong")
	ErrInputForamt         = derr.CarrierDroiError("1180003 Input format is wrong")
	ErrInvalidSessionToken = derr.CarrierDroiError("1180004 Session token is invalid")
	ErrDroiFileNotFound    = derr.CarrierDroiError("1180005 Droi File not exist")
	ErrFunctionDisabled    = derr.CarrierDroiError("1180006 Server functions disabled")
	ErrInvalidAdmin        = derr.CarrierDroiError("1180007 Invalid admin")
	ErrRequestSize         = derr.CarrierDroiError("1180008 Request size is too big")
	ErrAppIDLib            = derr.CarrierDroiError("1180009 AppIDLib error, not valid appid")
	ErrSessionExist        = derr.CarrierDroiError("1180010 Targit session exist")
	ErrSessionBusy         = derr.CarrierDroiError("1180011 Sessions busy")
	ErrAppUnavailable      = derr.CarrierDroiError("1180012 Target app unavailable, please contact droibaas")
	ErrInvalidCallbackSize = derr.CarrierDroiError("1180013 Upload session not match the size")
	ErrInvalidFileStatus   = derr.CarrierDroiError("1180014 Invalid file to be query")
	ErrBusyBulkDelete      = derr.CarrierDroiError("1180015 Delete Busy, please wait")
	// Qiniu
	ErrQiniuUnknown        = derr.CarrierDroiError("1180100 Unknown Qiniu error")
	ErrQiniuNoPermission   = derr.CarrierDroiError("1180101 No Permission")
	ErrQiniuAccountControl = derr.CarrierDroiError("1180102 Account control internal error")
	ErrQiniuSDK            = derr.CarrierDroiError("1180103 Qiniu SDK error")
	// File Majesty
	ErrMajestyUnknown = derr.CarrierDroiError("1180200 Unknown File Majesty error")
	ErrMajestyRequest = derr.CarrierDroiError("1180201 Request is not valid")
	// RDB API
	ErrRDBUnknown                = derr.CarrierDroiError("1180300 Unknown RDB API error")
	ErrRDBJsonFormat             = derr.CarrierDroiError("1180301 JSON format error")
	ErrRDBRecordNotFound         = derr.CarrierDroiError("1180302 No Result")
	ErrRDBQiniuAccountJsonFormat = derr.CarrierDroiError("1180303 Qiniu account is not JSON format string")
	// GOB API
	ErrGOBUnknown        = derr.CarrierDroiError("1180400 Unknown GOB API error")
	ErrGOBJsonFormat     = derr.CarrierDroiError("1180401 JSON format error")
	ErrGOBRecordNotFound = derr.CarrierDroiError("1180402 No Result")
	ErrGOBObjectExist    = derr.CarrierDroiError("1180403 Object exist")
	// Scan Errors
	ErrScanImage = derr.CarrierDroiError("1181001 The resource contains inappropriate content related porn/violent")
	ErrScanFile  = derr.CarrierDroiError("1181002 The resource contains inappropriate content related virus attack")
	ErrScanText  = derr.CarrierDroiError("1181003 The resource contains inappropriate content")
)
