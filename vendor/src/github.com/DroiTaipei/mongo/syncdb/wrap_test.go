package syncdb

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

var testdb = "baas"
var nulldb = ""
var movedb = "baas1"
var testCol = "_TestSoft"
var nullStr = ""
var testCtx = &droictx.DoneContext{}

var dockerPool *dockertest.Pool
var dockerResource *dockertest.Resource

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
}

func AfterTest() {
	Close()
	dockerPool.Purge(dockerResource)
}

func TestMain(m *testing.M) {
	BeforeTest()
	retCode := m.Run()
	AfterTest()
	os.Exit(retCode)
}

func TestInsertSoftMigration(t *testing.T) {
	doc := bson.M{"test": "for soft migration"}
	err := Insert(testCtx, testdb, movedb, testCol, doc)
	assert.NoError(t, err)

	// check baas and baas1 both have these data
	result1 := bson.M{}
	selector := bson.M{"test": "for soft migration"}
	nullBson := bson.M{}
	err = QueryOne(testCtx, testdb, testCol, &result1, selector, nullBson, 0, 0)
	assert.NoError(t, err)
	data, ok := result1["test"].(string)
	assert.True(t, ok)
	assert.Equal(t, "for soft migration", data)

	result2 := bson.M{}
	err = QueryOne(testCtx, movedb, testCol, &result2, selector, nullBson, 0, 0)
	assert.NoError(t, err)
	data, ok = result2["test"].(string)
	assert.True(t, ok)
	assert.Equal(t, "for soft migration", data)

	// teardown
	err = Remove(testCtx, testdb, movedb, testCol, selector)
	assert.NoError(t, err)

	err = QueryOne(testCtx, testdb, testCol, &result1, selector, nullBson, 0, 0)
	assert.Error(t, err)

	err = QueryOne(testCtx, movedb, testCol, &result2, selector, nullBson, 0, 0)
	assert.Error(t, err)
}

func TestUpsertSoftMigration(t *testing.T) {
	doc := bson.M{"test": "for soft migration"}
	selector := bson.M{"test": "for soft migration"}
	_, err := Upsert(testCtx, testdb, movedb, testCol, selector, doc)
	assert.NoError(t, err)

	// check baas and baas1 both have these data
	result1 := bson.M{}
	nullBson := bson.M{}
	err = QueryOne(testCtx, testdb, testCol, &result1, selector, nullBson, 0, 0)
	assert.NoError(t, err)
	data, ok := result1["test"].(string)
	assert.True(t, ok)
	assert.Equal(t, "for soft migration", data)

	result2 := bson.M{}
	err = QueryOne(testCtx, movedb, testCol, &result2, selector, nullBson, 0, 0)
	assert.NoError(t, err)
	data, ok = result2["test"].(string)
	assert.True(t, ok)
	assert.Equal(t, "for soft migration", data)

	// update
	doc = bson.M{"test": "2"}
	err = Update(testCtx, testdb, movedb, testCol, selector, doc)
	assert.NoError(t, err)

	// check baas and baas1 both have these data
	selector = bson.M{"test": "2"}
	result1 = bson.M{}
	err = QueryOne(testCtx, testdb, testCol, &result1, selector, nullBson, 0, 0)
	assert.NoError(t, err)
	data, ok = result1["test"].(string)
	assert.True(t, ok)
	assert.Equal(t, "2", data)

	result2 = bson.M{}
	err = QueryOne(testCtx, movedb, testCol, &result2, selector, nullBson, 0, 0)
	assert.NoError(t, err)
	data, ok = result2["test"].(string)
	assert.True(t, ok)
	assert.Equal(t, "2", data)

	// teardown
	err = Remove(testCtx, testdb, movedb, testCol, selector)
	assert.NoError(t, err)

	err = QueryOne(testCtx, testdb, testCol, &result1, selector, nullBson, 0, 0)
	assert.Error(t, err)

	err = QueryOne(testCtx, movedb, testCol, &result2, selector, nullBson, 0, 0)
	assert.Error(t, err)
}

func TestSoftMigrationBulk(t *testing.T) {
	defer DropCollection(testCtx, testdb, movedb, testCol)
	testData := []bson.M{
		bson.M{"test": 1},
		bson.M{"test": 2},
		bson.M{"test": 3},
		bson.M{"test": 4}}
	// insert
	err := BulkInsert(testCtx, testdb, movedb, testCol, testData)
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

	// check both db have data

	// update
	result, err := BulkUpsert(testCtx, testdb, movedb, testCol, testSelector, testData)
	assert.NoError(t, err)
	assert.Equal(t, 4, result.Modified)

	// delete
	result, err = BulkDelete(testCtx, testdb, movedb, testCol, testData)
	assert.NoError(t, err)
	assert.Equal(t, 4, result.Matched)

	// delete again
	result, err = BulkDelete(testCtx, testdb, movedb, testCol, testData)
	assert.NoError(t, err)
	assert.Equal(t, 0, result.Matched)
}

func TestFindAndModifyOne(t *testing.T) {
	Insert(testCtx, testdb, nulldb, testCol, bson.M{"test2": 1})

	oneResult := bson.M{}
	err := FindAndModify(testCtx, testdb, nulldb, testCol, &oneResult, bson.M{"test2": bson.M{"$eq": 1}}, bson.M{"$inc": bson.M{"test2": 2}}, nil, 0, 0, true, true)
	assert.NoError(t, err)

	test2 := oneResult["test2"].(int)
	assert.True(t, int(test2) > 1)
}

func TestFindAndRemoveOne(t *testing.T) {
	oneResult := bson.M{}
	err := FindAndRemove(testCtx, testdb, nulldb, testCol, &oneResult, bson.M{"test2": bson.M{"$gt": 1}}, nil, 0, 0)
	assert.NoError(t, err)
}

func TestBulk(t *testing.T) {
	defer DropCollection(testCtx, testdb, nulldb, testCol)
	testData := []bson.M{
		bson.M{"test": 1},
		bson.M{"test": 2},
		bson.M{"test": 3},
		bson.M{"test": 4}}
	// insert
	err := BulkInsert(testCtx, testdb, nulldb, testCol, testData)
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
	result, err := BulkUpsert(testCtx, testdb, nulldb, testCol, testSelector, testData)
	assert.NoError(t, err)
	assert.Equal(t, 4, result.Modified)

	// delete
	result, err = BulkDelete(testCtx, testdb, nulldb, testCol, testData)
	assert.NoError(t, err)
	assert.Equal(t, 4, result.Matched)

	// delete again
	result, err = BulkDelete(testCtx, testdb, nulldb, testCol, testData)
	assert.NoError(t, err)
	assert.Equal(t, 0, result.Matched)
}

func TestRename(t *testing.T) {
	// setup create collection
	err := CreateCollection(testCtx, testdb, nullStr, "testRename", &mgo.CollectionInfo{})
	assert.NoError(t, err)

	err = RenameCollection(testCtx, testdb, nullStr, "testRename", "testRename2")
	assert.NoError(t, err)

	// fail case
	err = RenameCollection(testCtx, testdb, nullStr, "testRename", "testRename3")
	assert.Error(t, err)

	// teardown
	DropCollection(testCtx, testdb, nullStr, "testRename2")
	DropCollection(testCtx, testdb, nullStr, "testRename")
}
