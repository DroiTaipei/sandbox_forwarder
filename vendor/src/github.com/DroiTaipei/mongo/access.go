package mongo

import (
	"fmt"
	"strings"
	"time"

	"github.com/DroiTaipei/droictx"
	"github.com/DroiTaipei/droipkg"
	me "github.com/DroiTaipei/droipkg/mongo"
	"github.com/DroiTaipei/mgo"
	"github.com/DroiTaipei/mgo/bson"
)

const (
	METHOD_CREATE_DOC  = "CreateDocument"
	METHOD_UPDATE_DOC  = "UpdateDocument"
	METHOD_UPDATE_DOCS = "UpdateDocuments"
	METHOD_READ_DOC    = "ReadDocument"
	METHOD_DELETE_DOC  = "DeleteDocument"
	METHOD_DELETE_DOCS = "DeleteDocuments"
	METHOD_CREATE_COL  = "CreateCollection"
	METHOD_DROP_COL    = "DropCollection"
	METHOD_RENAME_COL  = "RenameCollection"
	METHOD_CREATE_BULK = "CreateBulker"
	METHOD_BULK_INSERT = "BulkInsert"
	METHOD_BULK_UPSERT = "BulkUpsert"
	METHOD_BULK_DELETE = "BulkDelete"
)

func needReconnect(s string) bool {
	switch s {
	case me.MongoMsgNoReachableServers, me.MongoMsgEOF, me.MongoMsgKernelEOF, me.MongoMsgClose:
		return true
	}
	if strings.HasPrefix(s, me.MongoMsgNoHost) || strings.HasPrefix(s, me.MongoMsgNoHost2) || strings.HasPrefix(s, me.MongoMsgWriteUnavailable) {
		return true
	}
	return false
}

func getMongoCollection(s *Session, dbName, colName string) *mgo.Collection {
	return s.Session().DB(dbName).C(colName)
}

func (p *Pool) checkDatabaseError(err error, ctx droictx.Context, s *Session) droipkg.DroiError {
	if err == nil {
		p.put(s)
		return nil
	}
	errorString := err.Error()
	code := me.UnknownError

	switch {
	case needReconnect(errorString):
		// do reconnect
		errLog(ctx, fmt.Sprintf("API Lost Connection with Mongos: %s; error:%s", s.Addr(), errorString))
		go p.backgroundReconnect(s)
		code = me.APIConnectDatabase
		// special case, we won't put session right now
		return droipkg.NewCarrierDroiError(code, errorString)
	case err == mgo.ErrNotFound || strings.HasPrefix(errorString, me.MongoMsgUnknown):
		code = me.NotFound
	case errorString == me.MongoMsgNsNotFound || errorString == me.MongoMsgCollectionNotFound || errorString == me.MongoMsgNamespaceNotFound:
		code = me.CollectionNotFound
	case strings.HasPrefix(errorString, me.MongoMsgE11000) || strings.HasPrefix(errorString, me.MongoMsgBulk):
		code = me.DocumentConflict
		errorString = "document already exists:" + strings.Replace(errorString, me.MongoMsgE11000, "", 1)
	case strings.HasSuffix(errorString, me.MongoMsgTimeout) ||
		strings.HasPrefix(errorString, me.MongoMsgReadTCP) || strings.HasPrefix(errorString, me.MongoMsgWriteTCP):
		infoLog(ctx, errorString)
		code = me.Timeout
		s.Session().Refresh()
	case strings.HasSuffix(errorString, me.MongoMsgCollectionConflict):
		code = me.CollectionConflict
	case strings.HasSuffix(errorString, me.MongoMsgArray):
		code = me.QueryInputArray
		errorString = "query format error:" + errorString
	case strings.HasPrefix(errorString, me.MongoMsgEachArray) || strings.HasPrefix(errorString, me.MongoMsgPullAllArray):
		code = me.UpdateInputArray
		errorString = "Add/AddUnique/Remove requires an array argument"
	case strings.HasPrefix(errorString, me.MongoMsgIncrement):
		code = me.IncrementNumeric
	case errorString == me.MongoMsgRegexString:
		code = me.RegexString
	case strings.HasPrefix(errorString, me.MongoMsgDotField):
		code = me.DotField
	case strings.HasPrefix(errorString, me.MongoMsgwiredTigerIndex):
		code = me.StringIndexTooLong
	}
	p.put(s)
	return droipkg.NewCarrierDroiError(code, errorString)
}

func (p *Pool) Ping(ctx droictx.Context) (err droipkg.DroiError) {
	s, err := p.get(ctx)
	if err != nil {
		return
	}
	err = p.checkDatabaseError(s.Session().Ping(), ctx, s)
	return
}

// GetCollectionNames return all collection names in the db.
func (p *Pool) GetCollectionNames(ctx droictx.Context, dbName string) (names []string, err droipkg.DroiError) {
	s, err := p.get(ctx)
	if err != nil {
		return
	}
	names, dberr := s.Session().DB(dbName).CollectionNames()
	err = p.checkDatabaseError(dberr, ctx, s)
	return
}

func (p *Pool) CollectionCount(ctx droictx.Context, dbName, collection string) (n int, err droipkg.DroiError) {
	s, err := p.get(ctx)
	if err != nil {
		return
	}

	// start := time.Now()
	col := getMongoCollection(s, dbName, collection)
	n, dberr := col.Count()
	err = p.checkDatabaseError(dberr, ctx, s)
	// accessLog(ctx, METHOD_CREATE_DOC, fmt.Sprintf("Collection:%s,Stuff:%v", collection, doc), start)
	return
}

func (p *Pool) Run(ctx droictx.Context, cmd interface{}, result interface{}) droipkg.DroiError {
	s, err := p.get(ctx)
	if err != nil {
		return err
	}

	dberr := s.Session().Run(cmd, result)
	err = p.checkDatabaseError(dberr, ctx, s)
	infoLog(ctx, fmt.Sprintf("DB Operations:%v ,Result:%v", cmd, result))
	return err
}

func (p *Pool) DBRun(ctx droictx.Context, dbName string, cmd, result interface{}) droipkg.DroiError {
	s, err := p.get(ctx)
	if err != nil {
		return err
	}

	dberr := s.Session().DB(dbName).Run(cmd, result)
	err = p.checkDatabaseError(dberr, ctx, s)
	infoLog(ctx, fmt.Sprintf("DB Operations:%v ,Result:%v", cmd, result))
	return err
}

func (p *Pool) Insert(ctx droictx.Context, dbName, collection string, doc interface{}) droipkg.DroiError {
	s, err := p.get(ctx)
	if err != nil {
		return err
	}

	start := time.Now()
	col := getMongoCollection(s, dbName, collection)

	// Insert document to collection
	dberr := col.Insert(doc)
	err = p.checkDatabaseError(dberr, ctx, s)
	accessLog(ctx, METHOD_CREATE_DOC, fmt.Sprintf("Collection:%s,Stuff:%v", collection, doc), start)
	return err
}

func (p *Pool) Remove(ctx droictx.Context, dbName, collection string, selector interface{}) (err droipkg.DroiError) {
	s, err := p.get(ctx)
	if err != nil {
		return
	}

	start := time.Now()
	col := getMongoCollection(s, dbName, collection)
	dberr := col.Remove(selector)
	err = p.checkDatabaseError(dberr, ctx, s)
	accessLog(ctx, METHOD_DELETE_DOC, fmt.Sprintf("Collection:%s,Selector:%v", collection, selector), start)
	return
}
func (p *Pool) RemoveAll(ctx droictx.Context, dbName, collection string, selector interface{}) (info *mgo.ChangeInfo, err droipkg.DroiError) {
	s, err := p.get(ctx)
	if err != nil {
		return
	}

	start := time.Now()
	col := getMongoCollection(s, dbName, collection)
	info, dberr := col.RemoveAll(selector)
	err = p.checkDatabaseError(dberr, ctx, s)
	accessLog(ctx, METHOD_DELETE_DOCS, fmt.Sprintf("Collection:%s,Selector:%v", collection, selector), start)
	return
}

func (p *Pool) Update(ctx droictx.Context, dbName, collection string, selector interface{}, update interface{}) (err droipkg.DroiError) {
	s, err := p.get(ctx)
	if err != nil {
		return
	}

	start := time.Now()
	col := getMongoCollection(s, dbName, collection)
	dberr := col.Update(selector, update)
	err = p.checkDatabaseError(dberr, ctx, s)
	accessLog(ctx, METHOD_UPDATE_DOC, fmt.Sprintf("Collection:%s,Selector:%v,Stuff:%v", collection, selector, update), start)
	return
}

func (p *Pool) UpdateAll(ctx droictx.Context, dbName, collection string, selector interface{}, update interface{}) (info *mgo.ChangeInfo, err droipkg.DroiError) {
	s, err := p.get(ctx)
	if err != nil {
		return
	}

	start := time.Now()
	col := getMongoCollection(s, dbName, collection)
	info, dberr := col.UpdateAll(selector, update)
	err = p.checkDatabaseError(dberr, ctx, s)

	accessLog(ctx, METHOD_UPDATE_DOCS, fmt.Sprintf("Collection:%s,Selector:%v,Stuff:%v", collection, selector, update), start)
	return
}

func (p *Pool) UpdateId(ctx droictx.Context, dbName, collection string, id interface{}, update interface{}) (err droipkg.DroiError) {
	s, err := p.get(ctx)
	if err != nil {
		return
	}

	start := time.Now()
	col := getMongoCollection(s, dbName, collection)
	dberr := col.UpdateId(id, update)
	err = p.checkDatabaseError(dberr, ctx, s)
	accessLog(ctx, METHOD_UPDATE_DOC, fmt.Sprintf("Collection:%s,Selector:%v,Stuff:%v", collection, id, update), start)
	return
}

func (p *Pool) Upsert(ctx droictx.Context, dbName, collection string, selector interface{}, update interface{}) (info *mgo.ChangeInfo, err droipkg.DroiError) {
	s, err := p.get(ctx)
	if err != nil {
		return
	}

	start := time.Now()
	col := getMongoCollection(s, dbName, collection)
	info, dberr := col.Upsert(selector, update)
	err = p.checkDatabaseError(dberr, ctx, s)
	if err == nil && info.UpsertedId != nil {
		accessLog(ctx, METHOD_CREATE_DOC, fmt.Sprintf("Collection:%s,Selector:%v,Stuff:%v", collection, selector, update), start)
	} else {
		accessLog(ctx, METHOD_UPDATE_DOC, fmt.Sprintf("Collection:%s,Selector:%v,Stuff:%v", collection, selector, update), start)
	}
	return
}

func (p *Pool) BulkInsert(ctx droictx.Context, dbName, collection string, documents []bson.M) (err droipkg.DroiError) {
	s, err := p.get(ctx)
	if err != nil {
		return
	}

	start := time.Now()
	col := getMongoCollection(s, dbName, collection)
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
	err = p.checkDatabaseError(dberr, ctx, s)
	accessLog(ctx, METHOD_BULK_INSERT, collection, start)
	return
}

func (p *Pool) BulkUpsert(ctx droictx.Context, dbName, collection string, selectors, documents []bson.M) (result *mgo.BulkResult, err droipkg.DroiError) {
	s, err := p.get(ctx)
	if err != nil {
		return
	}

	start := time.Now()
	col := getMongoCollection(s, dbName, collection)
	bulk := col.Bulk()
	bulk.Unordered()
	b := len(documents)
	for i := 0; i < b; i++ {
		// NOTE: bson.NewObjectId would fail with goroutine, even there is only one goroutine worker.
		bulk.Upsert(selectors[i], documents[i])
	}
	// Set document _id if not set
	// Insert document to collection
	result, dberr := bulk.Run()
	err = p.checkDatabaseError(dberr, ctx, s)
	accessLog(ctx, METHOD_BULK_UPSERT, collection, start)
	return
}

func (p *Pool) BulkDelete(ctx droictx.Context, dbName, collection string, documents []bson.M) (result *mgo.BulkResult, err droipkg.DroiError) {
	s, err := p.get(ctx)
	if err != nil {
		return
	}

	start := time.Now()
	col := getMongoCollection(s, dbName, collection)
	bulk := col.Bulk()
	bulk.Unordered()
	b := len(documents)
	for i := 0; i < b; i++ {
		// NOTE: bson.NewObjectId would fail with goroutine, even there is only one goroutine worker.
		bulk.Remove(documents[i])
	}
	// Set document _id if not set
	// Insert document to collection
	result, dberr := bulk.Run()
	err = p.checkDatabaseError(dberr, ctx, s)
	accessLog(ctx, METHOD_BULK_DELETE, collection, start)
	return
}

func (p *Pool) GetBulk(ctx droictx.Context, dbName, collection string) (*Bulk, droipkg.DroiError) {
	s, err := p.get(ctx)
	if err != nil {
		return nil, err
	}

	col := getMongoCollection(s, dbName, collection)
	bulk := col.Bulk()
	bulk.Unordered()
	return &Bulk{
		bulk:       bulk,
		collection: collection,
		ctx:        ctx,
		session:    s,
		pool:       p,
	}, nil
}

type Bulk struct {
	bulk       *mgo.Bulk
	pool       *Pool
	ctx        droictx.Context
	session    *Session
	collection string
}

func (b *Bulk) Insert(docs ...interface{}) {
	b.bulk.Insert(docs...)
}
func (b *Bulk) Update(pairs ...interface{}) {
	b.bulk.Update(pairs...)
}
func (b *Bulk) Upsert(pairs ...interface{}) {
	b.bulk.Upsert(pairs...)
}
func (b *Bulk) Remove(selectors ...interface{}) {
	b.bulk.Remove(selectors...)
}
func (b *Bulk) RemoveAll(selectors ...interface{}) {
	b.bulk.RemoveAll(selectors...)
}
func (b *Bulk) Run() (result *mgo.BulkResult, err droipkg.DroiError) {
	start := time.Now()
	result, dberr := b.bulk.Run()
	err = b.pool.checkDatabaseError(dberr, b.ctx, b.session)
	accessLog(b.ctx, METHOD_CREATE_BULK, b.collection, start)
	return
}

func (p *Pool) QueryCount(ctx droictx.Context, dbName, collection string, selector interface{}) (n int, err droipkg.DroiError) {
	s, err := p.get(ctx)
	if err != nil {
		return
	}

	start := time.Now()
	query := getMongoCollection(s, dbName, collection).Find(selector)
	n, dberr := query.Count()
	err = p.checkDatabaseError(dberr, ctx, s)
	logMsg := fmt.Sprintf("Collection:%s,Query:%v",
		collection,
		selector,
	)
	accessLog(ctx, METHOD_READ_DOC, logMsg, start)
	return
}

func (p *Pool) QueryAll(ctx droictx.Context, dbName, collection string, result, selector, fields interface{}, skip, limit int, sort ...string) (err droipkg.DroiError) {
	s, err := p.get(ctx)
	if err != nil {
		return
	}

	start := time.Now()
	query := getMongoCollection(s, dbName, collection).Find(selector)
	if fields != nil {
		query.Select(fields)
	}
	query.Skip(skip)
	query.Limit(limit)
	if len(sort) > 0 {
		query.Sort(sort...)
	}

	dberr := query.All(result)
	err = p.checkDatabaseError(dberr, ctx, s)
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

func (p *Pool) QueryOne(ctx droictx.Context, dbName, collection string, result, selector, fields interface{}, skip, limit int, sort ...string) (err droipkg.DroiError) {
	s, err := p.get(ctx)
	if err != nil {
		return
	}

	start := time.Now()
	query := getMongoCollection(s, dbName, collection).Find(selector)
	if fields != nil {
		query.Select(fields)
	}
	query.Skip(skip)
	query.Limit(limit)
	if len(sort) > 0 {
		query.Sort(sort...)
	}
	dberr := query.One(result)

	err = p.checkDatabaseError(dberr, ctx, s)
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

// FindAndModify can only update one doc
func (p *Pool) FindAndModify(ctx droictx.Context, dbName, collection string, result, selector, update, fields interface{},
	skip, limit int, upsert, returnNew bool, sort ...string) (err droipkg.DroiError) {
	s, err := p.get(ctx)
	if err != nil {
		return
	}
	start := time.Now()
	query := getMongoCollection(s, dbName, collection).Find(selector)
	if fields != nil {
		query.Select(fields)
	}
	query.Skip(skip)
	query.Limit(limit)
	if len(sort) > 0 {
		query.Sort(sort...)
	}

	change := mgo.Change{
		Update:    update,
		ReturnNew: returnNew,
		Upsert:    upsert,
	}
	_, dberr := query.Apply(change, result)

	err = p.checkDatabaseError(dberr, ctx, s)
	logMsg := fmt.Sprintf("Collection:%s,Query:%v,Projection::%v,Sort:%v,Skip:%d,Limit:%d,Change:%v",
		collection,
		selector,
		fields,
		sort,
		skip,
		limit,
		change,
	)
	accessLog(ctx, METHOD_UPDATE_DOC, logMsg, start)
	return
}

func (p *Pool) FindAndRemove(ctx droictx.Context, dbName, collection string, result, selector, fields interface{},
	skip, limit int, sort ...string) (err droipkg.DroiError) {
	s, err := p.get(ctx)
	if err != nil {
		return
	}

	start := time.Now()
	query := getMongoCollection(s, dbName, collection).Find(selector)
	if fields != nil {
		query.Select(fields)
	}
	query.Skip(skip)
	query.Limit(limit)
	if len(sort) > 0 {
		query.Sort(sort...)
	}

	change := mgo.Change{
		Remove: true,
	}
	_, dberr := query.Apply(change, result)

	err = p.checkDatabaseError(dberr, ctx, s)
	logMsg := fmt.Sprintf("Collection:%s,Query:%v,Projection::%v,Sort:%v,Skip:%d,Limit:%d,Change:%v",
		collection,
		selector,
		fields,
		sort,
		skip,
		limit,
		change,
	)
	accessLog(ctx, METHOD_READ_DOC, logMsg, start)
	return
}

func (p *Pool) Indexes(ctx droictx.Context, dbName, collection string) (result []mgo.Index, err droipkg.DroiError) {
	s, err := p.get(ctx)
	if err != nil {
		return
	}
	result, dberr := getMongoCollection(s, dbName, collection).Indexes()
	err = p.checkDatabaseError(dberr, ctx, s)
	return
}

func (p *Pool) CreateIndex(ctx droictx.Context, dbName, collection string, key []string, sparse, unique bool, name string) (err droipkg.DroiError) {
	s, err := p.get(ctx)
	if err != nil {
		return
	}

	index := mgo.Index{
		Key:        key,
		Background: true,
		Sparse:     sparse,
		Unique:     unique,
		Name:       name,
	}
	dberr := getMongoCollection(s, dbName, collection).EnsureIndex(index)
	err = p.checkDatabaseError(dberr, ctx, s)
	if err != nil {
		errLog(ctx, collection+" add index:"+strings.Join(index.Key, ",")+" err:"+err.Error())
	}
	return
}

func (p *Pool) EnsureIndex(ctx droictx.Context, dbName, collection string, index mgo.Index) (err droipkg.DroiError) {
	s, err := p.get(ctx)
	if err != nil {
		return
	}
	dberr := getMongoCollection(s, dbName, collection).EnsureIndex(index)
	err = p.checkDatabaseError(dberr, ctx, s)
	if err != nil {
		errLog(ctx, collection+" add index:"+strings.Join(index.Key, ",")+" err:"+err.Error())
	}
	return
}

func (p *Pool) DropIndex(ctx droictx.Context, dbName, collection string, keys []string) (err droipkg.DroiError) {
	s, err := p.get(ctx)
	if err != nil {
		return
	}

	dberr := getMongoCollection(s, dbName, collection).DropIndex(keys...)
	err = p.checkDatabaseError(dberr, ctx, s)
	return
}

func (p *Pool) DropIndexName(ctx droictx.Context, dbName, collection, name string) (err droipkg.DroiError) {
	s, err := p.get(ctx)
	if err != nil {
		return
	}

	dberr := getMongoCollection(s, dbName, collection).DropIndexName(name)
	err = p.checkDatabaseError(dberr, ctx, s)
	return
}

func (p *Pool) CreateCollection(ctx droictx.Context, dbName, collection string, info *mgo.CollectionInfo) (err droipkg.DroiError) {
	start := time.Now()
	s, err := p.get(ctx)
	if err != nil {
		return
	}

	dberr := getMongoCollection(s, dbName, collection).Create(info)
	err = p.checkDatabaseError(dberr, ctx, s)
	accessLog(ctx, METHOD_CREATE_COL, collection, start)
	return
}

func (p *Pool) DropCollection(ctx droictx.Context, dbName, collection string) (err droipkg.DroiError) {
	start := time.Now()
	s, err := p.get(ctx)
	if err != nil {
		return
	}

	dberr := getMongoCollection(s, dbName, collection).DropCollection()
	err = p.checkDatabaseError(dberr, ctx, s)
	accessLog(ctx, METHOD_DROP_COL, collection, start)
	return
}

func (p *Pool) RenameCollection(ctx droictx.Context, dbName, oldName, newName string) (err droipkg.DroiError) {
	start := time.Now()
	s, err := p.get(ctx)
	if err != nil {
		return
	}
	// result is useless, error has message, so we just throw
	result := bson.M{}
	from := fmt.Sprintf("%s.%s", dbName, oldName)
	to := fmt.Sprintf("%s.%s", dbName, newName)
	dberr := s.Session().Run(bson.D{{"renameCollection", from}, {"to", to}}, result)
	err = p.checkDatabaseError(dberr, ctx, s)
	logMsg := fmt.Sprintf("Rename Collection from:%s to:%s", from, to)
	accessLog(ctx, METHOD_RENAME_COL, logMsg, start)
	return
}

func (p *Pool) Pipe(ctx droictx.Context, dbName, collection string, pipeline, result interface{}) (err droipkg.DroiError) {
	start := time.Now()
	s, err := p.get(ctx)
	if err != nil {
		return
	}
	pipe := getMongoCollection(s, dbName, collection).Pipe(pipeline)
	dberr := pipe.All(result)
	logMsg := fmt.Sprintf("Collection:%s,Pipeline:%v",
		collection,
		pipeline,
	)
	err = p.checkDatabaseError(dberr, ctx, s)
	accessLog(ctx, METHOD_READ_DOC, logMsg, start)
	return
}
