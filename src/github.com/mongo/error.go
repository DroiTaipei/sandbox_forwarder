package mongo

import "github.com/DroiTaipei/droipkg"

var ErrUnknown droipkg.DroiError

// ErrAPIConnectDatabase is constant error for connecting mongodb fail
var ErrAPIConnectDatabase droipkg.DroiError

// ErrAPIFullResource is constant error for not enough mongo connections
var ErrAPIFullResource droipkg.DroiError

// ErrMongoPoolClosed is constant error for pool closed
var ErrMongoPoolClosed droipkg.DroiError

func init() {
	ErrUnknown = droipkg.NewCarrierDroiError(UnknownError, "unknown error")
	ErrAPIConnectDatabase = droipkg.NewCarrierDroiError(APIConnectDatabase, "cannot connect database, please retry later")
	ErrAPIFullResource = droipkg.NewCarrierDroiError(APIFullResource, "mongo resource not enough")
	ErrMongoPoolClosed = droipkg.NewCarrierDroiError(MongoPoolClosed, "mongo connection pool closed")
}
