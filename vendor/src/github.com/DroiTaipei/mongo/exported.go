package mongo

import (
	"github.com/DroiTaipei/droictx"
	"github.com/DroiTaipei/droipkg"
	"github.com/DroiTaipei/mgo"
	"github.com/DroiTaipei/mgo/bson"
)

var std *Pool

// Initialize init mongo instance
func Initialize(dbi *DBInfo) (err error) {
	std, err = NewSessionPool(dbi)
	return
}

// FIXME:for gobuster...
func SetLogTopic(acc, std string) {
	mongoAccTopic = acc
	mongoStdTopic = std
}

func Close() {
	std.Close()
}

func IsAvailable() bool {
	return std.IsAvailable()
}

func Len() int {
	return std.Len()
}

func Cap() int {
	return std.Cap()
}

func Mode() mgo.Mode {
	return std.Mode()
}

func ShowConfig() map[string]interface{} {
	return std.ShowConfig()
}

func Ping(ctx droictx.Context) (err droipkg.DroiError) {
	return std.Ping(ctx)
}

// GetCollectionNames return all collection names in the db.
func GetCollectionNames(ctx droictx.Context, dbName string) (names []string, err droipkg.DroiError) {
	return std.GetCollectionNames(ctx, dbName)
}

func CollectionCount(ctx droictx.Context, dbName, collection string) (n int, err droipkg.DroiError) {
	return std.CollectionCount(ctx, dbName, collection)
}

func Run(ctx droictx.Context, cmd interface{}, result interface{}) droipkg.DroiError {
	return std.Run(ctx, cmd, result)
}

func DBRun(ctx droictx.Context, dbName string, cmd, result interface{}) droipkg.DroiError {
	return std.DBRun(ctx, dbName, cmd, result)
}

func Insert(ctx droictx.Context, dbName, collection string, doc interface{}) droipkg.DroiError {
	return std.Insert(ctx, dbName, collection, doc)
}

func Remove(ctx droictx.Context, dbName, collection string, selector interface{}) (err droipkg.DroiError) {
	return std.Remove(ctx, dbName, collection, selector)
}
func RemoveAll(ctx droictx.Context, dbName, collection string, selector interface{}) (info *mgo.ChangeInfo, err droipkg.DroiError) {
	return std.RemoveAll(ctx, dbName, collection, selector)
}
func Update(ctx droictx.Context, dbName, collection string, selector interface{}, update interface{}) (err droipkg.DroiError) {
	return std.Update(ctx, dbName, collection, selector, update)
}
func UpdateAll(ctx droictx.Context, dbName, collection string, selector interface{}, update interface{}) (info *mgo.ChangeInfo, err droipkg.DroiError) {
	return std.UpdateAll(ctx, dbName, collection, selector, update)
}

func UpdateId(ctx droictx.Context, dbName, collection string, id interface{}, update interface{}) (err droipkg.DroiError) {
	return std.UpdateId(ctx, dbName, collection, id, update)
}

func Upsert(ctx droictx.Context, dbName, collection string, selector interface{}, update interface{}) (info *mgo.ChangeInfo, err droipkg.DroiError) {
	return std.Upsert(ctx, dbName, collection, selector, update)
}

func BulkInsert(ctx droictx.Context, dbName, collection string, documents []bson.M) (err droipkg.DroiError) {
	return std.BulkInsert(ctx, dbName, collection, documents)
}

func BulkUpsert(ctx droictx.Context, dbName, collection string, selectors, documents []bson.M) (result *mgo.BulkResult, err droipkg.DroiError) {
	return std.BulkUpsert(ctx, dbName, collection, selectors, documents)
}

func BulkDelete(ctx droictx.Context, dbName, collection string, documents []bson.M) (result *mgo.BulkResult, err droipkg.DroiError) {
	return std.BulkDelete(ctx, dbName, collection, documents)
}

func GetBulk(ctx droictx.Context, dbName, collection string) (*Bulk, droipkg.DroiError) {
	return std.GetBulk(ctx, dbName, collection)
}

func QueryCount(ctx droictx.Context, dbName, collection string, selector interface{}) (n int, err droipkg.DroiError) {
	return std.QueryCount(ctx, dbName, collection, selector)
}

func QueryAll(ctx droictx.Context, dbName, collection string, result, selector, fields interface{}, skip, limit int, sort ...string) (err droipkg.DroiError) {
	return std.QueryAll(ctx, dbName, collection, result, selector, fields, skip, limit, sort...)
}

func QueryOne(ctx droictx.Context, dbName, collection string, result, selector, fields interface{}, skip, limit int, sort ...string) (err droipkg.DroiError) {
	return std.QueryOne(ctx, dbName, collection, result, selector, fields, skip, limit, sort...)
}

// FindAndModify can only update one doc
func FindAndModify(ctx droictx.Context, dbName, collection string, result, selector, update, fields interface{},
	skip, limit int, upsert, returnNew bool, sort ...string) (err droipkg.DroiError) {
	return std.FindAndModify(ctx, dbName, collection, result, selector, update, fields, skip, limit, upsert, returnNew, sort...)
}

func FindAndRemove(ctx droictx.Context, dbName, collection string, result, selector, fields interface{},
	skip, limit int, sort ...string) (err droipkg.DroiError) {
	return std.FindAndRemove(ctx, dbName, collection, result, selector, fields, skip, limit, sort...)

}

func Indexes(ctx droictx.Context, dbName, collection string) (result []mgo.Index, err droipkg.DroiError) {
	return std.Indexes(ctx, dbName, collection)
}

func CreateIndex(ctx droictx.Context, dbName, collection string, key []string, sparse, unique bool, name string) (err droipkg.DroiError) {
	return std.CreateIndex(ctx, dbName, collection, key, sparse, unique, name)
}

func EnsureIndex(ctx droictx.Context, dbName, collection string, index mgo.Index) (err droipkg.DroiError) {
	return std.EnsureIndex(ctx, dbName, collection, index)
}

func DropIndex(ctx droictx.Context, dbName, collection string, keys []string) (err droipkg.DroiError) {
	return std.DropIndex(ctx, dbName, collection, keys)
}

func DropIndexName(ctx droictx.Context, dbName, collection, name string) (err droipkg.DroiError) {
	return std.DropIndexName(ctx, dbName, collection, name)
}

func CreateCollection(ctx droictx.Context, dbName, collection string, info *mgo.CollectionInfo) (err droipkg.DroiError) {
	return std.CreateCollection(ctx, dbName, collection, info)
}

func DropCollection(ctx droictx.Context, dbName, collection string) (err droipkg.DroiError) {
	return std.DropCollection(ctx, dbName, collection)
}

func RenameCollection(ctx droictx.Context, dbName, oldName, newName string) (err droipkg.DroiError) {
	return std.RenameCollection(ctx, dbName, oldName, newName)
}

func Pipe(ctx droictx.Context, dbName, collection string, pipeline, result interface{}) (err droipkg.DroiError) {
	return std.Pipe(ctx, dbName, collection, pipeline, result)
}
