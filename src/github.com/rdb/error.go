package rdb

import (
	de "github.com/DroiTaipei/droipkg"
)

var ErrUnknown de.DroiError
var ErrResourceNotFound de.DroiError
var ErrInvalidParameter de.DroiError
var ErrJsonValidation de.DroiError
var ErrPermissionDenied de.DroiError
var ErrDatabaseUnavailable de.DroiError
var ErrDatabase de.DroiError
var ErrDataNotFound de.DroiError
var ErrProcessFailed de.DroiError
var ErrPrimaryKeyDuplicated de.DroiError
var ErrCacheKeyLimit de.DroiError
var ErrCacheValueLimit de.DroiError
var ErrCacheNotFound de.DroiError
var ErrCacheUpdate de.DroiError
var ErrCacheDelete de.DroiError
var ErrCacheRecordLimit de.DroiError
var ErrDuplicatedData de.DroiError

func init() {
	ErrUnknown = de.ConstDroiError("1020001 Unknown")
	ErrResourceNotFound = de.ConstDroiError("1020002 Resource Not Found")
	ErrInvalidParameter = de.ConstDroiError("1020003 Parameter Validation Error")
	ErrJsonValidation = de.ConstDroiError("1020004 Json Validation Failed")
	ErrPermissionDenied = de.ConstDroiError("1020005 Permission Denied")
	ErrDatabaseUnavailable = de.ConstDroiError("1020006 Database Unavailable")
	ErrDatabase = de.ConstDroiError("1020007 Database Error")
	ErrProcessFailed = de.ConstDroiError("1020008 Data Process Failed")
	ErrPrimaryKeyDuplicated = de.ConstDroiError("1020009 Primary Key Duplicated")
	ErrDataNotFound = de.ConstDroiError("1020010 Data Not Found")
	ErrCacheKeyLimit = de.ConstDroiError("1020011 Out Cache Key Limit(1-32 character)")
	ErrCacheValueLimit = de.ConstDroiError("1020012 Out Cache Value Limit(1-1024 character)")
	ErrCacheNotFound = de.ConstDroiError("1020013 Cache Not Found With Key")
	ErrCacheUpdate = de.ConstDroiError("1020014 Cache Update Failed")
	ErrCacheDelete = de.ConstDroiError("1020015 Cache Delete Failed")
	ErrCacheRecordLimit = de.ConstDroiError("1020016 Out Cache Record Limit(Upper bound is 100)")
	ErrDuplicatedData = de.ConstDroiError("1020017 Duplicated Data")
}
