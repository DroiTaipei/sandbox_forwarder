package mongo

import (
	"fmt"
	"os"
	"testing"

	dockertest "gopkg.in/ory-am/dockertest.v3"

	"github.com/DroiTaipei/droictx"
	"github.com/DroiTaipei/droipkg"
	"github.com/DroiTaipei/mgo"
	"github.com/DroiTaipei/mgo/bson"
	"github.com/stretchr/testify/assert"
)

var dockerPool *dockertest.Pool
var dockerResource *dockertest.Resource

var testMongoPool *Pool
var testdb = "baas"
var testCol = "_Test"
var testCtx = &droictx.DoneContext{}

func init() {
	droipkg.SetLogger(droipkg.GetDiscardLogger())
}

func BeforeTest() {
	var err error
	dockerPool, err = dockertest.NewPool("")
	if err != nil {
		panic("docker test init fail, error:" + err.Error())
	}

	dockerResource, err = dockerPool.Run("mongo", "3.4", nil)
	if err != nil {
		panic("mongodb 3.4 init fail, error:" + err.Error())
	}

	testPoolConfig := NewDBInfo("testmgo", []string{fmt.Sprintf("localhost:%s", dockerResource.GetPort("27017/tcp"))}, "", "", "", 5, 5, true, false)

	err = dockerPool.Retry(func() error {
		errI := Initialize(testPoolConfig)
		return errI
	})
	if err != nil {
		panic("connect docker fail, error:" + err.Error())
	}

	testMongoPool, err = NewSessionPool(testPoolConfig)
	if err != nil {
		panic("connect docker fail, error:" + err.Error())
	}
}

func AfterTest() {
	if IsAvailable() {
		Close()
	}
	testMongoPool.Close()
	dockerPool.Purge(dockerResource)
}

func TestMain(m *testing.M) {
	BeforeTest()
	retCode := m.Run()
	AfterTest()
	os.Exit(retCode)
}

func TestLog(t *testing.T) {
	SetLogTopic("acc", "std")
	assert.Equal(t, "acc", mongoAccTopic)
	assert.Equal(t, "std", mongoStdTopic)
}

func TestMongoInfo(t *testing.T) {
	assert.True(t, IsAvailable())
	conf := ShowConfig()
	assert.Equal(t, conf["MaxConn"], Len())
	assert.Equal(t, conf["MaxConn"], Cap())
	assert.Equal(t, mgo.Primary, Mode())
}

func TestPing(t *testing.T) {
	err := Ping(testCtx)
	assert.NoError(t, err)
}

func TestGetCollectionOp(t *testing.T) {
	names, err := GetCollectionNames(testCtx, testdb)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(names))

	err = CreateCollection(testCtx, testdb, testCol, &mgo.CollectionInfo{})
	assert.NoError(t, err)

	names, err = GetCollectionNames(testCtx, testdb)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(names))
	assert.Equal(t, testCol, names[0])

	count, err := CollectionCount(testCtx, testdb, testCol)
	assert.NoError(t, err)
	assert.Equal(t, 0, count)

}

func TestRun(t *testing.T) {
	cmd := bson.M{"collstats": testCol}
	result := bson.M{}
	err := DBRun(testCtx, testdb, cmd, result)
	assert.NoError(t, err)
	assert.Equal(t, 0, result["count"].(int))

	result = bson.M{}
	cmd = bson.M{"listCollections": 1}
	err = Run(testCtx, cmd, result)
	assert.NoError(t, err)
}

func TestQuery(t *testing.T) {
	selector := bson.M{}

	// setup doc
	Insert(testCtx, testdb, testCol, bson.M{"test": "1"})

	n, err := QueryCount(testCtx, testdb, testCol, selector)
	assert.NoError(t, err)
	assert.Equal(t, 1, n)

	result := bson.M{}
	err = QueryOne(testCtx, testdb, testCol, &result, selector, bson.M{"test": 1}, 0, 0, "test")
	assert.NoError(t, err)
	assert.Equal(t, "1", result["test"].(string))

	results := []bson.M{}
	err = QueryAll(testCtx, testdb, testCol, &results, selector, bson.M{"test": 1}, 0, 0, "test")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(results))
	assert.Equal(t, "1", results[0]["test"].(string))

	results = []bson.M{}
	err = Pipe(testCtx, testdb, testCol, []bson.M{{"$match": bson.M{}}}, &results)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(results))
	assert.Equal(t, "1", results[0]["test"].(string))

	// teardown
	err = Remove(testCtx, testdb, testCol, selector)
	assert.NoError(t, err)
}

func TestUpdate(t *testing.T) {
	selector := bson.M{"test1": 1}
	doc := bson.M{"test1": 1}
	info, err := Upsert(testCtx, testdb, testCol, selector, doc)
	assert.NoError(t, err)
	assert.Equal(t, 0, info.Matched)
	assert.NotNil(t, info.UpsertedId)

	// query check exists
	result := bson.M{}
	err = QueryOne(testCtx, testdb, testCol, &result, selector, nil, 0, 0)
	assert.NoError(t, err)
	assert.Equal(t, 1, result["test1"].(int))

	info, err = Upsert(testCtx, testdb, testCol, selector, bson.M{"test1": 1, "test2": 0})
	assert.NoError(t, err)
	assert.Equal(t, 1, info.Matched)

	err = Update(testCtx, testdb, testCol, selector, bson.M{"$set": bson.M{"test2": 1}})
	assert.NoError(t, err)

	err = UpdateId(testCtx, testdb, testCol, result["_id"], bson.M{"$set": bson.M{"test2": 2}})
	assert.NoError(t, err)

	info, err = UpdateAll(testCtx, testdb, testCol, selector, bson.M{"$set": bson.M{"test2": 3}})
	assert.NoError(t, err)
	assert.Equal(t, 1, info.Matched)

	// teardown
	info, err = RemoveAll(testCtx, testdb, testCol, bson.M{"test2": 3})
	assert.NoError(t, err)
	assert.Equal(t, 1, info.Removed)
}

func TestFindAndModifyOne(t *testing.T) {
	Insert(testCtx, testdb, testCol, bson.M{"test2": 1})

	oneResult := bson.M{}
	err := FindAndModify(testCtx, testdb, testCol, &oneResult, bson.M{"test2": bson.M{"$eq": 1}}, bson.M{"$inc": bson.M{"test2": 2}}, bson.M{"test2": 1}, 0, 0, true, true, "test2")
	assert.NoError(t, err)

	test2 := oneResult["test2"].(int)
	assert.True(t, int(test2) > 1)
}

func TestFindAndRemoveOne(t *testing.T) {
	oneResult := bson.M{}
	err := FindAndRemove(testCtx, testdb, testCol, &oneResult, bson.M{"test2": bson.M{"$gt": 1}}, bson.M{"test2": 1}, 0, 0, "test2")
	assert.NoError(t, err)
	test2 := oneResult["test2"].(int)
	assert.True(t, int(test2) > 1)
}

func TestBulk(t *testing.T) {
	defer DropCollection(testCtx, testdb, testCol)
	testData := []bson.M{
		bson.M{"test": 1},
		bson.M{"test": 2},
		bson.M{"test": 3},
		bson.M{"test": 4}}
	// insert
	err := BulkInsert(testCtx, testdb, testCol, testData)
	assert.NoError(t, err)

	testSelector := []bson.M{
		bson.M{"test": 1},
		bson.M{"test": 2},
		bson.M{"test": 3},
		bson.M{"test": 4}}
	testData = []bson.M{
		bson.M{"test": 5},
		bson.M{"test": 6},
		bson.M{"test": 7},
		bson.M{"test": 8}}
	// update
	result, err := BulkUpsert(testCtx, testdb, testCol, testSelector, testData)
	assert.NoError(t, err)
	assert.Equal(t, 4, result.Modified)

	// delete
	result, err = BulkDelete(testCtx, testdb, testCol, testData)
	assert.NoError(t, err)
	assert.Equal(t, 4, result.Matched)

	// delete again
	result, err = BulkDelete(testCtx, testdb, testCol, testData)
	assert.NoError(t, err)
	assert.Equal(t, 0, result.Matched)

	bulk, err := GetBulk(testCtx, testdb, testCol)
	assert.NoError(t, err)
	bulk.Insert(bson.M{"test": 1})
	bulk.Upsert(bson.M{"test": 1}, bson.M{"test": 2})
	bulk.Update(bson.M{"test": 2}, bson.M{"test": 3})
	bulk.Remove(bson.M{"test": 3})
	bulk.RemoveAll(bson.M{})

	_, err = bulk.Run()
	assert.NoError(t, err)
}

func TestRename(t *testing.T) {
	// setup create collection
	err := CreateCollection(testCtx, testdb, "testRename", &mgo.CollectionInfo{})
	assert.NoError(t, err)

	err = RenameCollection(testCtx, testdb, "testRename", "testRename2")
	assert.NoError(t, err)

	// fail case
	err = RenameCollection(testCtx, testdb, "testRename", "testRename3")
	assert.Error(t, err)

	// teardown
	DropCollection(testCtx, testdb, "testRename")
	DropCollection(testCtx, testdb, "testRename2")
}

func TestIndex(t *testing.T) {
	// setup create collection
	err := CreateCollection(testCtx, testdb, "testCreateIndex", &mgo.CollectionInfo{})
	assert.NoError(t, err)

	indexes, err := Indexes(testCtx, testdb, "testCreateIndex")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(indexes))

	err = CreateIndex(testCtx, testdb, "testCreateIndex", []string{"testField"}, false, false, "test")
	assert.NoError(t, err)

	indexes, err = Indexes(testCtx, testdb, "testCreateIndex")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(indexes))

	// create index twice with different settings error
	err = CreateIndex(testCtx, testdb, "testCreateIndex", []string{"testField"}, true, false, "test")
	assert.Error(t, err)

	// dup index with same setting would not get error
	sameIndex := mgo.Index{
		Key:        []string{"testField"},
		Background: true,
		Sparse:     false,
		Unique:     false,
		Name:       "test",
	}
	err = EnsureIndex(testCtx, testdb, "testCreateIndex", sameIndex)
	assert.NoError(t, err)

	// dup index with different setting would get error
	sameIndex.Sparse = true
	err = EnsureIndex(testCtx, testdb, "testCreateIndex", sameIndex)
	assert.Error(t, err)

	sameIndex.Key = []string{"test2F"}
	sameIndex.Name = ""
	err = EnsureIndex(testCtx, testdb, "testCreateIndex", sameIndex)
	assert.NoError(t, err)

	err = DropIndex(testCtx, testdb, "testCreateIndex", sameIndex.Key)
	assert.NoError(t, err)

	// if we want use key to drop index, the index name must be default
	err = DropIndex(testCtx, testdb, "testCreateIndex", []string{"testField"})
	assert.Error(t, err)

	err = DropIndexName(testCtx, testdb, "testCreateIndex", "test")
	assert.NoError(t, err)

	indexes, err = Indexes(testCtx, testdb, "testCreateIndex")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(indexes))
	// teardown
	DropCollection(testCtx, testdb, "testCreateIndex")
}

func TestClosedError(t *testing.T) {
	Close()
	err := Ping(testCtx)
	assert.Error(t, err)

	_, err = GetCollectionNames(testCtx, testdb)
	assert.Error(t, err)

	_, err = CollectionCount(testCtx, testdb, testCol)
	assert.Error(t, err)

	cmd := bson.M{"collstats": testCol}
	result := bson.M{}
	err = DBRun(testCtx, testdb, cmd, result)
	assert.Error(t, err)

	err = Run(testCtx, cmd, result)
	assert.Error(t, err)

	doc := bson.M{"123": 1}
	err = Insert(testCtx, testdb, testCol, doc)
	assert.Error(t, err)

	err = Remove(testCtx, testdb, testCol, doc)
	assert.Error(t, err)

	_, err = RemoveAll(testCtx, testdb, testCol, doc)
	assert.Error(t, err)

	err = Update(testCtx, testdb, testCol, doc, doc)
	assert.Error(t, err)

	_, err = UpdateAll(testCtx, testdb, testCol, doc, doc)
	assert.Error(t, err)

	err = UpdateId(testCtx, testdb, testCol, "59cf40762fda3562afa00fab", doc)
	assert.Error(t, err)

	_, err = Upsert(testCtx, testdb, testCol, doc, doc)
	assert.Error(t, err)

	bulkData := []bson.M{bson.M{"test": 1}}
	err = BulkInsert(testCtx, testdb, testCol, bulkData)
	assert.Error(t, err)

	_, err = BulkUpsert(testCtx, testdb, testCol, bulkData, bulkData)
	assert.Error(t, err)

	_, err = BulkDelete(testCtx, testdb, testCol, bulkData)
	assert.Error(t, err)

	_, err = GetBulk(testCtx, testdb, testCol)
	assert.Error(t, err)

	_, err = QueryCount(testCtx, testdb, testCol, doc)
	assert.Error(t, err)

	err = QueryOne(testCtx, testdb, testCol, nil, nil, nil, 0, 0)
	assert.Error(t, err)

	err = QueryAll(testCtx, testdb, testCol, nil, nil, nil, 0, 0)
	assert.Error(t, err)

	err = FindAndModify(testCtx, testdb, testCol, nil, nil, nil, nil, 0, 0, true, true)
	assert.Error(t, err)

	err = FindAndRemove(testCtx, testdb, testCol, nil, nil, nil, 0, 0)
	assert.Error(t, err)

	_, err = Indexes(testCtx, testdb, testCol)
	assert.Error(t, err)

	err = CreateIndex(testCtx, testdb, testCol, nil, false, false, "123")
	assert.Error(t, err)

	err = EnsureIndex(testCtx, testdb, testCol, mgo.Index{})
	assert.Error(t, err)

	err = DropIndexName(testCtx, testdb, testCol, "123")
	assert.Error(t, err)

	err = DropIndex(testCtx, testdb, testCol, nil)
	assert.Error(t, err)

	err = CreateCollection(testCtx, testdb, testCol, nil)
	assert.Error(t, err)

	err = DropCollection(testCtx, testdb, testCol)
	assert.Error(t, err)

	err = RenameCollection(testCtx, testdb, testCol, "new")
	assert.Error(t, err)

	err = Pipe(testCtx, testdb, testCol, nil, nil)
	assert.Error(t, err)
}
