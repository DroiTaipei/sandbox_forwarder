package syncdb

import (
	"github.com/DroiTaipei/droictx"
	"github.com/DroiTaipei/droipkg"
	"github.com/DroiTaipei/mgo"
	"github.com/DroiTaipei/mgo/bson"
	"github.com/DroiTaipei/mongo"
)

func Initialize(dbi *mongo.DBInfo) (err error) {
	return mongo.Initialize(dbi)
}

func Close() {
	mongo.Close()
}

func IsAvailable() bool {
	return mongo.IsAvailable()
}

func Len() int {
	return mongo.Len()
}

func Cap() int {
	return mongo.Cap()
}

func Mode() mgo.Mode {
	return mongo.Mode()
}

func SetLogTopic(acc, std string) {
	mongo.SetLogTopic(acc, std)
}

func ShowConfig() map[string]interface{} {
	return mongo.ShowConfig()
}

func NewDBInfo(name string, addrs []string, user, password, authdbName string,
	timeout, maxConn int, direct, readSecondary bool) *mongo.DBInfo {
	return mongo.NewDBInfo(name, addrs, user, password, authdbName, timeout, maxConn, direct, readSecondary)
}

func Ping(ctx droictx.Context) (err droipkg.DroiError) {
	return mongo.Ping(ctx)
}

func GetCollectionNames(ctx droictx.Context, dbName string) (names []string, err droipkg.DroiError) {
	return mongo.GetCollectionNames(ctx, dbName)
}

func CollectionCount(ctx droictx.Context, dbName, collection string) (n int, err droipkg.DroiError) {
	return mongo.CollectionCount(ctx, dbName, collection)
}

func Run(ctx droictx.Context, cmd interface{}, result interface{}) droipkg.DroiError {
	return mongo.Run(ctx, cmd, result)
}

func DBRun(ctx droictx.Context, dbName, moveDB string, cmd, result interface{}) droipkg.DroiError {
	err := mongo.DBRun(ctx, dbName, cmd, result)
	if err == nil && len(moveDB) > 0 {
		mongo.DBRun(ctx, moveDB, cmd, result)
	}
	return err
}

func Insert(ctx droictx.Context, dbName, moveDB, collection string, doc interface{}) droipkg.DroiError {
	err := mongo.Insert(ctx, dbName, collection, doc)
	if err == nil && len(moveDB) > 0 {
		mongo.Insert(ctx, moveDB, collection, doc)
	}
	return err
}

func Remove(ctx droictx.Context, dbName, moveDB, collection string, selector interface{}) (err droipkg.DroiError) {
	err = mongo.Remove(ctx, dbName, collection, selector)
	if err == nil && len(moveDB) > 0 {
		mongo.Remove(ctx, moveDB, collection, selector)
	}
	return
}
func RemoveAll(ctx droictx.Context, dbName, moveDB, collection string, selector interface{}) (info *mgo.ChangeInfo, err droipkg.DroiError) {
	info, err = mongo.RemoveAll(ctx, dbName, collection, selector)
	if err == nil && len(moveDB) > 0 {
		mongo.RemoveAll(ctx, moveDB, collection, selector)
	}
	return
}
func Update(ctx droictx.Context, dbName, moveDB, collection string, selector interface{}, update interface{}) (err droipkg.DroiError) {
	err = mongo.Update(ctx, dbName, collection, selector, update)
	if err == nil && len(moveDB) > 0 {
		mongo.Update(ctx, moveDB, collection, selector, update)
	}
	return
}
func UpdateAll(ctx droictx.Context, dbName, moveDB, collection string, selector interface{}, update interface{}) (info *mgo.ChangeInfo, err droipkg.DroiError) {
	info, err = mongo.UpdateAll(ctx, dbName, collection, selector, update)
	if err == nil && len(moveDB) > 0 {
		mongo.UpdateAll(ctx, moveDB, collection, selector, update)
	}
	return
}

func UpdateId(ctx droictx.Context, dbName, moveDB, collection string, id interface{}, update interface{}) (err droipkg.DroiError) {
	err = mongo.UpdateId(ctx, dbName, collection, id, update)
	if err == nil && len(moveDB) > 0 {
		mongo.UpdateId(ctx, moveDB, collection, id, update)
	}
	return
}

func Upsert(ctx droictx.Context, dbName, moveDB, collection string, selector interface{}, update interface{}) (info *mgo.ChangeInfo, err droipkg.DroiError) {
	info, err = mongo.Upsert(ctx, dbName, collection, selector, update)
	if err == nil && len(moveDB) > 0 {
		mongo.Upsert(ctx, moveDB, collection, selector, update)
	}
	return
}

func BulkInsert(ctx droictx.Context, dbName, moveDB, collection string, documents []bson.M) (err droipkg.DroiError) {
	err = mongo.BulkInsert(ctx, dbName, collection, documents)
	if err == nil && len(moveDB) > 0 {
		mongo.BulkInsert(ctx, moveDB, collection, documents)
	}
	return
}

func BulkUpsert(ctx droictx.Context, dbName, moveDB, collection string, selectors, documents []bson.M) (result *mgo.BulkResult, err droipkg.DroiError) {
	result, err = mongo.BulkUpsert(ctx, dbName, collection, selectors, documents)
	if err == nil && len(moveDB) > 0 {
		mongo.BulkUpsert(ctx, moveDB, collection, selectors, documents)
	}
	return
}

func BulkDelete(ctx droictx.Context, dbName, moveDB, collection string, documents []bson.M) (result *mgo.BulkResult, err droipkg.DroiError) {
	result, err = mongo.BulkDelete(ctx, dbName, collection, documents)
	if err == nil && len(moveDB) > 0 {
		mongo.BulkDelete(ctx, moveDB, collection, documents)
	}
	return
}

type SyncBulk struct {
	bulk     *mongo.Bulk
	moveBulk *mongo.Bulk
}

func GetBulk(ctx droictx.Context, dbName, moveDB, collection string) (*SyncBulk, droipkg.DroiError) {
	bulk, err := mongo.GetBulk(ctx, dbName, collection)
	if err != nil {
		return nil, err
	}
	var mvBulk *mongo.Bulk
	if len(moveDB) > 0 {
		mvBulk, _ = mongo.GetBulk(ctx, moveDB, collection)
	}
	return &SyncBulk{
		bulk:     bulk,
		moveBulk: mvBulk,
	}, nil
}

func (b *SyncBulk) Insert(docs ...interface{}) {
	b.bulk.Insert(docs...)
	if b.moveBulk != nil {
		b.moveBulk.Insert(docs...)
	}
}
func (b *SyncBulk) Update(pairs ...interface{}) {
	b.bulk.Update(pairs...)
	if b.moveBulk != nil {
		b.moveBulk.Update(pairs...)
	}
}
func (b *SyncBulk) Upsert(pairs ...interface{}) {
	b.bulk.Upsert(pairs...)
	if b.moveBulk != nil {
		b.moveBulk.Upsert(pairs...)
	}
}
func (b *SyncBulk) Remove(selectors ...interface{}) {
	b.bulk.Remove(selectors...)
	if b.moveBulk != nil {
		b.moveBulk.Remove(selectors...)
	}
}
func (b *SyncBulk) RemoveAll(selectors ...interface{}) {
	b.bulk.RemoveAll(selectors...)
	if b.moveBulk != nil {
		b.moveBulk.RemoveAll(selectors...)
	}
}
func (b *SyncBulk) Run() (result *mgo.BulkResult, err droipkg.DroiError) {
	result, err = b.bulk.Run()
	if err == nil && b.moveBulk != nil {
		b.moveBulk.Run()
	}
	return
}

func QueryCount(ctx droictx.Context, dbName, collection string, selector interface{}) (n int, err droipkg.DroiError) {
	return mongo.QueryCount(ctx, dbName, collection, selector)
}

func QueryAll(ctx droictx.Context, dbName, collection string, result, selector, fields interface{}, skip, limit int, sort ...string) (err droipkg.DroiError) {
	return mongo.QueryAll(ctx, dbName, collection, result, selector, fields, skip, limit, sort...)
}

func QueryOne(ctx droictx.Context, dbName, collection string, result, selector, fields interface{}, skip, limit int, sort ...string) (err droipkg.DroiError) {
	return mongo.QueryOne(ctx, dbName, collection, result, selector, fields, skip, limit, sort...)
}

// FindAndModify can only update one doc
func FindAndModify(ctx droictx.Context, dbName, moveDB, collection string, result, selector, update, fields interface{},
	skip, limit int, upsert, returnNew bool, sort ...string) (err droipkg.DroiError) {
	err = mongo.FindAndModify(ctx, dbName, collection, result, selector, update, fields, skip, limit, upsert, returnNew, sort...)
	if err == nil && len(moveDB) > 0 {
		if upsert {
			mongo.Upsert(ctx, moveDB, collection, selector, update)
		} else {
			mongo.Update(ctx, moveDB, collection, selector, update)
		}
	}
	return
}

func FindAndRemove(ctx droictx.Context, dbName, moveDB, collection string, result, selector, fields interface{},
	skip, limit int, sort ...string) (err droipkg.DroiError) {
	err = mongo.FindAndRemove(ctx, dbName, collection, result, selector, fields, skip, limit, sort...)
	if err == nil && len(moveDB) > 0 {
		mongo.Remove(ctx, moveDB, collection, selector)
	}
	return
}

func Indexes(ctx droictx.Context, dbName, collection string) (result []mgo.Index, err droipkg.DroiError) {
	return mongo.Indexes(ctx, dbName, collection)
}

func CreateIndex(ctx droictx.Context, dbName, moveDB, collection string, key []string, sparse, unique bool, name string) (err droipkg.DroiError) {
	err = mongo.CreateIndex(ctx, dbName, collection, key, sparse, unique, name)
	if err == nil && len(moveDB) > 0 {
		mongo.CreateIndex(ctx, moveDB, collection, key, sparse, unique, name)
	}
	return
}

func EnsureIndex(ctx droictx.Context, dbName, moveDB, collection string, index mgo.Index) (err droipkg.DroiError) {
	err = mongo.EnsureIndex(ctx, dbName, collection, index)
	if err == nil && len(moveDB) > 0 {
		mongo.EnsureIndex(ctx, moveDB, collection, index)
	}
	return
}

func DropIndex(ctx droictx.Context, dbName, moveDB, collection string, keys []string) (err droipkg.DroiError) {
	err = mongo.DropIndex(ctx, dbName, collection, keys)
	if err == nil && len(moveDB) > 0 {
		mongo.DropIndex(ctx, moveDB, collection, keys)
	}
	return
}

func DropIndexName(ctx droictx.Context, dbName, moveDB, collection, name string) (err droipkg.DroiError) {
	err = mongo.DropIndexName(ctx, dbName, collection, name)
	if err == nil && len(moveDB) > 0 {
		mongo.DropIndexName(ctx, moveDB, collection, name)
	}
	return
}

func CreateCollection(ctx droictx.Context, dbName, moveDB, collection string, info *mgo.CollectionInfo) (err droipkg.DroiError) {
	err = mongo.CreateCollection(ctx, dbName, collection, info)
	if err == nil && len(moveDB) > 0 {
		mongo.CreateCollection(ctx, moveDB, collection, info)
	}
	return
}

func DropCollection(ctx droictx.Context, dbName, moveDB, collection string) (err droipkg.DroiError) {
	err = mongo.DropCollection(ctx, dbName, collection)
	if err == nil && len(moveDB) > 0 {
		mongo.DropCollection(ctx, moveDB, collection)
	}
	return
}

func RenameCollection(ctx droictx.Context, dbName, moveDB, oldName, newName string) (err droipkg.DroiError) {
	err = mongo.RenameCollection(ctx, dbName, oldName, newName)
	if err == nil && len(moveDB) > 0 {
		mongo.RenameCollection(ctx, moveDB, oldName, newName)
	}
	return
}

func Pipe(ctx droictx.Context, dbName, collection string, pipeline, result interface{}) (err droipkg.DroiError) {
	return mongo.Pipe(ctx, dbName, collection, pipeline, result)
}
