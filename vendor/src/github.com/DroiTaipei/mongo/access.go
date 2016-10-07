package mongo

import (
	"fmt"
	"strings"
	"time"

	"github.com/DroiTaipei/droictx"
	"github.com/DroiTaipei/droipkg"
	"github.com/DroiTaipei/mgo"
	"github.com/DroiTaipei/mgo/bson"
)

const (
	METHOD_CREATE_DOC  = "CreateDocument"
	METHOD_UPDATE_DOC  = "UpdateDocument"
	METHOD_READ_DOC    = "ReadDocument"
	METHOD_DELETE_DOC  = "DeleteDocument"
	METHOD_CREATE_COL  = "CreateCollection"
	METHOD_DROP_COL    = "DropCollection"
	METHOD_CREATE_BULK = "BulkCreate"
)

func Ping(ctx droictx.Context) (err droipkg.DroiError) {
	s, entry := getWriteSess()
	defer s.Close()
	err = CheckDatabaseError(s.Ping(), ctx, entry)
	return
}

func CollectionCount(ctx droictx.Context, collection string) (n int, err droipkg.DroiError) {
	s, entry := getWriteSess()
	defer s.Close()
	// start := time.Now()
	col := getMongoCollection(s, collection)
	n, dberr := col.Count()
	err = CheckDatabaseError(dberr, ctx, entry)
	// accessLog(ctx, METHOD_CREATE_DOC, fmt.Sprintf("Collection:%s,Stuff:%v", collection, doc), start)
	return
}

func Insert(ctx droictx.Context, collection string, doc interface{}) droipkg.DroiError {
	s, entry := getWriteSess()
	defer s.Close()

	start := time.Now()
	col := getMongoCollection(s, collection)

	// Insert document to collection
	dberr := col.Insert(doc)
	err := CheckDatabaseError(dberr, ctx, entry)
	accessLog(ctx, METHOD_CREATE_DOC, fmt.Sprintf("Collection:%s,Stuff:%v", collection, doc), start)
	return err
}

func Remove(ctx droictx.Context, collection string, selector interface{}) (err droipkg.DroiError) {
	s, entry := getWriteSess()
	defer s.Close()

	start := time.Now()
	col := getMongoCollection(s, collection)
	dberr := col.Remove(selector)
	err = CheckDatabaseError(dberr, ctx, entry)
	if err != nil {
		return
	}
	accessLog(ctx, METHOD_DELETE_DOC, fmt.Sprintf("Collection:%s,Selector:%v", collection, selector), start)
	return
}
func RemoveAll(ctx droictx.Context, collection string, selector interface{}) (info *mgo.ChangeInfo, err droipkg.DroiError) {
	s, entry := getWriteSess()
	defer s.Close()

	start := time.Now()
	col := getMongoCollection(s, collection)
	info, dberr := col.RemoveAll(selector)
	err = CheckDatabaseError(dberr, ctx, entry)
	if err != nil {
		return
	}
	accessLog(ctx, METHOD_DELETE_DOC, fmt.Sprintf("Collection:%s,Selector:%v", collection, selector), start)
	return
}
func Update(ctx droictx.Context, collection string, selector interface{}, update interface{}) (err droipkg.DroiError) {
	s, entry := getWriteSess()
	defer s.Close()

	start := time.Now()
	col := getMongoCollection(s, collection)
	dberr := col.Update(selector, update)
	err = CheckDatabaseError(dberr, ctx, entry)
	if err != nil {
		return
	}
	accessLog(ctx, METHOD_UPDATE_DOC, fmt.Sprintf("Collection:%s,Selector:%v,Stuff:%v", collection, selector, update), start)
	return
}
func UpdateAll(ctx droictx.Context, collection string, selector interface{}, update interface{}) (info *mgo.ChangeInfo, err droipkg.DroiError) {
	s, entry := getWriteSess()
	defer s.Close()

	start := time.Now()
	col := getMongoCollection(s, collection)
	info, dberr := col.UpdateAll(selector, update)
	err = CheckDatabaseError(dberr, ctx, entry)
	if err != nil {
		return
	}
	accessLog(ctx, METHOD_UPDATE_DOC, fmt.Sprintf("Collection:%s,Selector:%v,Stuff:%v", collection, selector, update), start)
	return
}
func UpdateId(ctx droictx.Context, collection string, id interface{}, update interface{}) (err droipkg.DroiError) {
	s, entry := getWriteSess()
	defer s.Close()

	start := time.Now()
	col := getMongoCollection(s, collection)
	dberr := col.UpdateId(id, update)
	err = CheckDatabaseError(dberr, ctx, entry)
	if err != nil {
		return
	}
	accessLog(ctx, METHOD_UPDATE_DOC, fmt.Sprintf("Collection:%s,Selector:%v,Stuff:%v", collection, id, update), start)
	return
}
func Upsert(ctx droictx.Context, collection string, selector interface{}, update interface{}) (info *mgo.ChangeInfo, err droipkg.DroiError) {
	s, entry := getWriteSess()
	defer s.Close()

	start := time.Now()
	col := getMongoCollection(s, collection)
	info, dberr := col.Upsert(selector, update)
	err = CheckDatabaseError(dberr, ctx, entry)
	if err != nil {
		return
	}
	if err == nil && info.UpsertedId != nil {
		accessLog(ctx, METHOD_CREATE_DOC, fmt.Sprintf("Collection:%s,Selector:%v,Stuff:%v", collection, selector, update), start)
	} else {
		accessLog(ctx, METHOD_UPDATE_DOC, fmt.Sprintf("Collection:%s,Selector:%v,Stuff:%v", collection, selector, update), start)
	}
	return
}

func BulkInsert(ctx droictx.Context, collection string, documents []interface{}) (err droipkg.DroiError) {
	s, entry := getWriteSess()
	defer s.Close()

	start := time.Now()
	col := getMongoCollection(s, collection)
	bulk := col.Bulk()
	bulk.Unordered()
	b := len(documents)
	for i := 0; i < b; i++ {
		// NOTE: bson.NewObjectId would fail with goroutine, even there is only one goroutine worker.
		bulk.Insert(documents[i])
	}
	// Set document _id if not set
	// Insert document to collection
	_, dberr := bulk.Run()
	err = CheckDatabaseError(dberr, ctx, entry)
	accessLog(ctx, METHOD_CREATE_BULK, collection, start)
	return
}

type Bulk struct {
	Bulk       *mgo.Bulk
	Collection string
	Entry      string
	Ctx        droictx.Context
	Session    *mgo.Session
}

func GetBulk(ctx droictx.Context, collection string) *Bulk {
	s, entry := getWriteSess()

	col := getMongoCollection(s, collection)
	bulk := col.Bulk()
	bulk.Unordered()
	return &Bulk{
		Bulk:       bulk,
		Collection: collection,
		Entry:      entry,
		Ctx:        ctx,
		Session:    s,
	}
}

func (b *Bulk) Insert(docs ...interface{}) {
	b.Bulk.Insert(docs...)
}
func (b *Bulk) Update(pairs ...interface{}) {
	b.Bulk.Update(pairs...)
}
func (b *Bulk) Remove(selectors ...interface{}) {
	b.Bulk.Remove(selectors...)
}
func (b *Bulk) RemoveAll(selectors ...interface{}) {
	b.Bulk.RemoveAll(selectors...)
}
func (b *Bulk) Run() (result *mgo.BulkResult, err droipkg.DroiError) {
	defer b.Session.Close()
	start := time.Now()
	result, dberr := b.Bulk.Run()
	err = CheckDatabaseError(dberr, b.Ctx, b.Entry)
	accessLog(b.Ctx, METHOD_CREATE_BULK, b.Collection, start)
	return
}

// TO DISCUSSION
// func Find(collection string, query interface{}) *mgo.Query {
// 	s, entry := getWriteSess()
// 	defer s.Close()
// 	col := getMongoCollection(s, collection)
// 	return col.Find(query)
// }
// func One(ctx droictx.Context, query *mgo.Query, result interface{}) (err droipkg.DroiError) {
// 	s, entry := getWriteSess()
// 	defer s.Close()
// 	start := time.Now()
// 	dberr := query.One(result)
// 	err = CheckDatabaseError(dberr, ctx, entry)
// 	logMsg := fmt.Sprintf("Collection:%s,Query:%v,Projection::%v,Sort:%v,Skip:%d,Limit:%d",
// 		collection,
// 		selector,
// 		fields,
// 		sort,
// 		skip,
// 		limit,
// 	)
// 	accessLog(ctx, METHOD_READ_DOC, logMsg, start)
// 	return nil
// }
// func All(ctx droictx.Context, query *mgo.Query, result interface{}) (err droipkg.DroiError) {
// 	start := time.Now()
// 	dberr := q.Query.All(result)
// 	err := CheckDatabaseError(dberr, q.Ctx, q.Entry, q.ReconnectHandler)
// 	accessLog(ctx, METHOD_READ_DOC, logMsg, start)
// 	return nil
// }

func QueryCount(ctx droictx.Context, collection string, selector interface{}) (n int, err droipkg.DroiError) {
	s, entry := getReadSess()
	defer s.Close()
	start := time.Now()
	query := getMongoCollection(s, collection).Find(selector)
	n, dberr := query.Count()
	err = CheckDatabaseError(dberr, ctx, entry)
	logMsg := fmt.Sprintf("Collection:%s,Query:%v",
		collection,
		selector,
	)
	accessLog(ctx, METHOD_READ_DOC, logMsg, start)
	return
}

func QueryAll(ctx droictx.Context, collection string, result, selector, fields interface{}, skip, limit int, sort ...string) (err droipkg.DroiError) {
	s, entry := getReadSess()
	defer s.Close()
	start := time.Now()
	query := getMongoCollection(s, collection).Find(selector)
	if fields != nil {
		query.Select(fields)
	}
	query.Skip(skip)
	query.Limit(limit)
	if len(sort) > 0 {
		query.Sort(sort...)
	}

	dberr := query.All(result)
	err = CheckDatabaseError(dberr, ctx, entry)
	logMsg := fmt.Sprintf("Collection:%s,Query:%v,Projection::%v,Sort:%v,Skip:%d,Limit:%d",
		collection,
		selector,
		fields,
		sort,
		skip,
		limit,
	)
	accessLog(ctx, METHOD_READ_DOC, logMsg, start)
	return
}

func QueryOne(ctx droictx.Context, collection string, result, selector, fields interface{}, skip, limit int, sort ...string) (err droipkg.DroiError) {
	s, entry := getReadSess()
	defer s.Close()
	start := time.Now()
	query := getMongoCollection(s, collection).Find(selector)
	if fields != nil {
		query.Select(fields)
	}
	query.Skip(skip)
	query.Limit(limit)
	if len(sort) > 0 {
		query.Sort(sort...)
	}
	dberr := query.One(result)

	err = CheckDatabaseError(dberr, ctx, entry)
	logMsg := fmt.Sprintf("Collection:%s,Filter:%v,Projection::%v,Sort:%v,Skip:%d,Limit:%d",
		collection,
		selector,
		fields,
		sort,
		skip,
		limit,
	)
	accessLog(ctx, METHOD_READ_DOC, logMsg, start)
	return
}

func Indexes(ctx droictx.Context, collection string) (result []mgo.Index, err droipkg.DroiError) {
	s, entry := getWriteSess()
	defer s.Close()
	result, dberr := getMongoCollection(s, collection).Indexes()
	err = CheckDatabaseError(dberr, ctx, entry)
	return
}

func CreateIndex(ctx droictx.Context, collection string, key []string, sparse, unique bool, name string) (err droipkg.DroiError) {
	s, entry := getWriteSess()
	defer s.Close()

	index := mgo.Index{
		Key:        key,
		Background: true,
		Sparse:     sparse,
		Unique:     unique,
		Name:       name,
	}
	dberr := getMongoCollection(s, collection).EnsureIndex(index)
	err = CheckDatabaseError(dberr, ctx, entry)
	if err != nil {
		errLog(ctx, collection+" add index:"+strings.Join(index.Key, ",")+" err:"+err.Error())
	}
	return
}

func EnsureIndex(ctx droictx.Context, collection string, index mgo.Index) (err droipkg.DroiError) {
	s, entry := getWriteSess()
	defer s.Close()
	dberr := getMongoCollection(s, collection).EnsureIndex(index)
	err = CheckDatabaseError(dberr, ctx, entry)
	if err != nil {
		errLog(ctx, collection+" add index:"+strings.Join(index.Key, ",")+" err:"+err.Error())
	}
	return
}

func DropIndex(ctx droictx.Context, collection string, keys []string) (err droipkg.DroiError) {
	s, entry := getWriteSess()
	defer s.Close()

	dberr := getMongoCollection(s, collection).DropIndex(keys...)
	err = CheckDatabaseError(dberr, ctx, entry)
	return
}

func CreateCollection(ctx droictx.Context, collection string, info *mgo.CollectionInfo) (err droipkg.DroiError) {
	start := time.Now()
	s, entry := getWriteSess()
	defer s.Close()

	dberr := getMongoCollection(s, collection).Create(info)

	err = CheckDatabaseError(dberr, ctx, entry)
	if err != nil {
		return
	}
	accessLog(ctx, METHOD_CREATE_COL, collection, start)
	return
}

func DropCollection(ctx droictx.Context, collection string) (err droipkg.DroiError) {
	start := time.Now()
	s, entry := getWriteSess()
	defer s.Close()

	dberr := getMongoCollection(s, collection).DropCollection()
	err = CheckDatabaseError(dberr, ctx, entry)
	if err != nil {
		return
	}
	accessLog(ctx, METHOD_DROP_COL, collection, start)
	return
}

// TODO: maybe removed in the future
func SetSharding(ctx droictx.Context, collection string) (err droipkg.DroiError) {
	start := time.Now()
	s := getShardSess()
	defer s.Close()

	result := bson.M{}
	fullName := mgoDatabaseName + "." + collection
	dberr := s.Run(bson.D{{"shardCollection", fullName},
		{"key", bson.M{mgoDefaultPK: 1}}}, &result)
	err = CheckDatabaseError(dberr, ctx, "shard")
	if err != nil {
		return
	}
	accessLog(ctx, "Sharding", collection, start)
	return
}
