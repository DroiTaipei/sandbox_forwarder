package mongo

import (
	"errors"
	"testing"

	me "github.com/DroiTaipei/droipkg/mongo"
	"github.com/DroiTaipei/mgo"
	"github.com/DroiTaipei/mgo/bson"
	"github.com/stretchr/testify/assert"
)

func TestCheckDatabaseError(t *testing.T) {
	s, err := testMongoPool.get(testCtx)
	assert.NoError(t, err)
	// no error case
	err = testMongoPool.checkDatabaseError(nil, testCtx, s)
	assert.NoError(t, err)

	// reconnect
	s, _ = testMongoPool.get(testCtx)
	err = testMongoPool.checkDatabaseError(errors.New("EOF"), testCtx, s)
	assert.Error(t, err)
	assert.Equal(t, me.APIConnectDatabase, err.ErrorCode())

	// not found
	s, _ = testMongoPool.get(testCtx)
	err = testMongoPool.checkDatabaseError(mgo.ErrNotFound, testCtx, s)
	assert.Error(t, err)
	assert.Equal(t, me.NotFound, err.ErrorCode())

	// collection not found
	s, _ = testMongoPool.get(testCtx)
	err = testMongoPool.checkDatabaseError(errors.New(me.MongoMsgNsNotFound), testCtx, s)
	assert.Error(t, err)
	assert.Equal(t, me.CollectionNotFound, err.ErrorCode())

	// collection not found
	s, _ = testMongoPool.get(testCtx)
	err = testMongoPool.checkDatabaseError(errors.New(me.MongoMsgE11000), testCtx, s)
	assert.Error(t, err)
	assert.Equal(t, me.DocumentConflict, err.ErrorCode())

	// timeout
	s, _ = testMongoPool.get(testCtx)
	err = testMongoPool.checkDatabaseError(errors.New(me.MongoMsgReadTCP), testCtx, s)
	assert.Error(t, err)
	assert.Equal(t, me.Timeout, err.ErrorCode())

	// Collection Conflict
	s, _ = testMongoPool.get(testCtx)
	err = testMongoPool.checkDatabaseError(errors.New(me.MongoMsgCollectionConflict), testCtx, s)
	assert.Error(t, err)
	assert.Equal(t, me.CollectionConflict, err.ErrorCode())

	// query need array input
	s, _ = testMongoPool.get(testCtx)
	err = testMongoPool.checkDatabaseError(errors.New(me.MongoMsgArray), testCtx, s)
	assert.Error(t, err)
	assert.Equal(t, me.QueryInputArray, err.ErrorCode())

	// update need array input
	s, _ = testMongoPool.get(testCtx)
	err = testMongoPool.checkDatabaseError(errors.New(me.MongoMsgEachArray), testCtx, s)
	assert.Error(t, err)
	assert.Equal(t, me.UpdateInputArray, err.ErrorCode())

	// increment need number
	s, _ = testMongoPool.get(testCtx)
	err = testMongoPool.checkDatabaseError(errors.New(me.MongoMsgIncrement), testCtx, s)
	assert.Error(t, err)
	assert.Equal(t, me.IncrementNumeric, err.ErrorCode())

	// regex need string
	s, _ = testMongoPool.get(testCtx)
	err = testMongoPool.checkDatabaseError(errors.New(me.MongoMsgRegexString), testCtx, s)
	assert.Error(t, err)
	assert.Equal(t, me.RegexString, err.ErrorCode())

	// dot field name
	s, _ = testMongoPool.get(testCtx)
	err = testMongoPool.checkDatabaseError(errors.New(me.MongoMsgDotField), testCtx, s)
	assert.Error(t, err)
	assert.Equal(t, me.DotField, err.ErrorCode())

	// dot field name
	s, _ = testMongoPool.get(testCtx)
	err = testMongoPool.checkDatabaseError(errors.New(me.MongoMsgwiredTigerIndex), testCtx, s)
	assert.Error(t, err)
	assert.Equal(t, me.StringIndexTooLong, err.ErrorCode())
}

func TestNeedReconnect(t *testing.T) {
	assert.False(t, needReconnect("i/o timeout "))
	assert.False(t, needReconnect("write tcp"))
	assert.False(t, needReconnect("read tcp"))

	assert.True(t, needReconnect("EOF"))
	assert.True(t, needReconnect("End of file"))
	assert.True(t, needReconnect("Closed explicitly"))
	assert.True(t, needReconnect("no reachable servers"))
	assert.True(t, needReconnect("write results unavailable"))
	assert.True(t, needReconnect("could not find host matching read preference"))
	assert.True(t, needReconnect("None of the hosts"))
}

func TestGetMongoCollection(t *testing.T) {
	s, err := testMongoPool.get(testCtx)
	assert.NoError(t, err)
	defer testMongoPool.put(s)

	col := getMongoCollection(s, testdb, testCol)
	assert.Equal(t, testCol, col.Name)
	assert.Equal(t, testdb+"."+testCol, col.FullName)
}

func TestPoolFindAndModifyOne(t *testing.T) {
	// setup value
	testMongoPool.Insert(testCtx, testdb, testCol, bson.M{"test2": 1})

	oneResult := bson.M{}
	err := testMongoPool.FindAndModify(testCtx, testdb, testCol, &oneResult, bson.M{"test2": bson.M{"$eq": 1}}, bson.M{"$inc": bson.M{"test2": 2}}, nil, 0, 0, true, true)
	assert.NoError(t, err)

	test2 := oneResult["test2"].(int)
	assert.True(t, int(test2) > 1)
}

func TestPoolFindAndRemoveOne(t *testing.T) {
	oneResult := bson.M{}
	err := testMongoPool.FindAndRemove(testCtx, testdb, testCol, &oneResult, bson.M{"test2": bson.M{"$gt": 1}}, nil, 0, 0)
	assert.NoError(t, err)
}

func TestPoolBulk(t *testing.T) {
	defer testMongoPool.DropCollection(testCtx, testdb, testCol)
	testData := []bson.M{
		bson.M{"test": 1},
		bson.M{"test": 2},
		bson.M{"test": 3},
		bson.M{"test": 4}}
	// insert
	err := testMongoPool.BulkInsert(testCtx, testdb, testCol, testData)
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
	result, err := testMongoPool.BulkUpsert(testCtx, testdb, testCol, testSelector, testData)
	assert.NoError(t, err)
	assert.Equal(t, 4, result.Modified)

	// delete
	result, err = testMongoPool.BulkDelete(testCtx, testdb, testCol, testData)
	assert.NoError(t, err)
	assert.Equal(t, 4, result.Matched)

	// delete again
	result, err = testMongoPool.BulkDelete(testCtx, testdb, testCol, testData)
	assert.NoError(t, err)
	assert.Equal(t, 0, result.Matched)
}

func TestPoolRename(t *testing.T) {
	// setup create collection
	err := testMongoPool.CreateCollection(testCtx, testdb, "testRename", &mgo.CollectionInfo{})
	assert.NoError(t, err)

	err = testMongoPool.RenameCollection(testCtx, testdb, "testRename", "testRename2")
	assert.NoError(t, err)

	// fail case
	err = testMongoPool.RenameCollection(testCtx, testdb, "testRename", "testRename3")
	assert.Error(t, err)

	// teardown
	testMongoPool.DropCollection(testCtx, testdb, "testRename")
	testMongoPool.DropCollection(testCtx, testdb, "testRename2")
}

func TestPoolCreateIndex(t *testing.T) {
	// setup create collection
	err := testMongoPool.CreateCollection(testCtx, testdb, "testCreateIndex", &mgo.CollectionInfo{})
	assert.NoError(t, err)

	err = testMongoPool.CreateIndex(testCtx, testdb, "testCreateIndex", []string{"--testField"}, false, false, "test")
	assert.NoError(t, err)
	// teardown
	testMongoPool.DropCollection(testCtx, testdb, "testCreateIndex")
}
