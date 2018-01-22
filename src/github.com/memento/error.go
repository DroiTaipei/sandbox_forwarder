package memento

import (
	"github.com/DroiTaipei/droipkg"
	"github.com/DroiTaipei/droipkg/mongo"
)

const (
	// ErrorMemento defines memento error code preifx
	ErrorMemento = 1420000
	// 1420001 is preserved for mongo connect error
	// 1420013 db resource not enough
	// 1420014 db pool closed

	// ErrorMementoGRPC is the error code of gRPC error
	ErrorMementoGRPC = ErrorMemento + 2
	// ErrorMementoPanic is the error code of internal service error
	ErrorMementoPanic = ErrorMemento + 3
	// ErrorMementoDataBroken is the error code of bson marshal error
	ErrorMementoDataBroken = ErrorMemento + 4

	// ErrorMementoData defines memento error code preifx of data logic part
	ErrorMementoData = ErrorMemento + 100
	// ErrorMementoKey is the error code of key too long error
	ErrorMementoKey = ErrorMementoData + 1
	// ErrorMementoValue is the error code of value too big error
	ErrorMementoValue = ErrorMementoData + 2
	// ErrorMementoBadRequest is the error code of request parse error
	ErrorMementoBadRequest = ErrorMementoData + 3

	// 1420300 preserved for mongo error

	// ErrorMongoAddend for transforming mongo 103xxxx error to 142xxxx
	ErrorMongoAddend = ErrorMemento - mongo.UnknownError
)

// ErrMementoPanic is the error of internal service error
var ErrMementoPanic = droipkg.NewCarrierDroiError(ErrorMementoPanic, "internal service error")

// ErrMementoKey is the error of key too long error
var ErrMementoKey = droipkg.NewCarrierDroiError(ErrorMementoKey, "key too long > 1000")

// ErrMementoValue is the error of value too big error
var ErrMementoValue = droipkg.NewCarrierDroiError(ErrorMementoValue, "value too big > 1MB")
