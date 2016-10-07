package mongo

import (
	"strings"

	"github.com/DroiTaipei/droictx"
	"github.com/DroiTaipei/droipkg"
	"github.com/DroiTaipei/mgo"
)

const (
	UnknownError       = 1030000
	APIConnectDatabase = 1030001
	// Data Logic Error
	NotFound           = 1030301
	CollectionNotFound = 1030302
	DocumentConflict   = 1030305
	CollectionConflict = 1030306
	QueryInputArray    = 1030108
	UpdateInputArray   = 1030109
	IncrementNumeric   = 1030110

	// Data Logic Error
	DroiErrMsgQueryInputArray = "query format error"
	DroiErrMsgAddArray        = "Add/AddUnique/Remove requires an array argument"
)

const (
	eofStr    = "EOF"
	closedStr = "Closed explicitly"
	// Mongo Error Msg
	MongoMsgNotFound           = "not found"
	MongoMsgDocConflict        = "document already exists"
	MongoMsgCollectionConflict = "collection already exists"
	MongoMsgNsNotFound         = "ns not found"
	MongoMsgCollectionNotFound = "no collection"
	MongoMsgNonEmpty           = "$and/$or/$nor must be a nonempty array"
	MongoMsgAndArray           = "and needs an array"
	MongoMsgOrArray            = "$or needs an array"
	MongoMsgInArray            = "$in needs an array"
	MongoMsgEachArray          = "The argument to $each"
	MongoMsgPullAllArray       = "$pullAll requires an array argument"
	MongoMsgIncrement          = "Cannot increment with non-numeric argument"
	MongoMsgE11000             = "E11000"
	MongoMsgUnknown            = "Unknown"
	MongoMsgBulk               = "multiple errors in bulk operation"
	MongoMsgTimeout            = "read tcp"
)

var updateArrayError error
var UpdateArrayError droipkg.DroiError

func init() {
	updateArrayError = droipkg.NewError(DroiErrMsgAddArray)
	UpdateArrayError = &MgoApiError{
		Code: UpdateInputArray,
		Err:  updateArrayError,
	}
}

type MgoApiError struct {
	Code int
	Err  error
}

func (ae *MgoApiError) ErrorCode() int {
	return ae.Code
}
func (ae *MgoApiError) SetErrorCode(code int) {
	ae.Code = code
}
func (ae *MgoApiError) Error() string {
	return ae.Err.Error()
}
func (ae *MgoApiError) Wrap(msg string) {
	ae.Err = droipkg.Wrap(ae.Err, msg)
}

const one = 1
const dcMsg = "API Lost Connection with Mongo; entry:"
const dotStr = "."
const nullStr = ""

func CheckDatabaseError(err error, ctx droictx.Context, entry string) droipkg.DroiError {
	if err == nil {
		return nil
	}
	errorString := err.Error()
	code := UnknownError

	if err == mgo.ErrNotFound || strings.HasPrefix(errorString, MongoMsgUnknown) {
		code = NotFound
	} else if errorString == eofStr || errorString == closedStr {
		// do reconnect
		errLog(systemCtx, dcMsg+entry+errorString)
		Reconnect()
		code = APIConnectDatabase
	} else if errorString == MongoMsgNsNotFound {
		code = CollectionNotFound
	} else if strings.HasPrefix(errorString, MongoMsgE11000) || strings.HasPrefix(errorString, MongoMsgBulk) {
		code = DocumentConflict
		errString := err.Error()
		errString = strings.Replace(errString, MongoMsgE11000, nullStr, one)
		errString = strings.Replace(errString, mgoDatabaseName+dotStr, nullStr, one)
		// dotIndex := strings.Index(errString, dotStr)
		// errString = errString[:dotIndex-appIDLength] + errString[dotIndex+one:]
		err = droipkg.NewError(errString)
		err = droipkg.Wrap(err, MongoMsgDocConflict)
	} else if strings.HasPrefix(errorString, MongoMsgCollectionConflict) {
		code = CollectionConflict
	} else if errorString == MongoMsgNonEmpty || errorString == MongoMsgAndArray ||
		errorString == MongoMsgInArray || errorString == MongoMsgOrArray {
		code = QueryInputArray
		err = droipkg.Wrap(err, DroiErrMsgQueryInputArray)
	} else if strings.HasPrefix(errorString, MongoMsgEachArray) || strings.HasPrefix(errorString, MongoMsgPullAllArray) {
		return UpdateArrayError
	} else if strings.HasPrefix(errorString, MongoMsgIncrement) {
		code = IncrementNumeric
	} else if errorString == MongoMsgCollectionNotFound {
		code = CollectionNotFound
	} else if strings.HasPrefix(errorString, MongoMsgTimeout) {
		code = APIConnectDatabase
		errLog(systemCtx, errorString)
	}
	return &MgoApiError{
		Code: code,
		Err:  err,
	}
}
