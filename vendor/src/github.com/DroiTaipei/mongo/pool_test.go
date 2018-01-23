package mongo

import (
	"fmt"
	"testing"
	"time"

	"github.com/DroiTaipei/droictx"
	"github.com/DroiTaipei/droipkg"
	me "github.com/DroiTaipei/droipkg/mongo"
	"github.com/DroiTaipei/mgo"
	"github.com/stretchr/testify/assert"
)

func TestNewDBInfo(t *testing.T) {
	dbi := NewDBInfo("testmgo", []string{"127.0.0.1:27017"}, "", "", "", 3, 5, true, false)
	assert.Equal(t, "testmgo", dbi.Name)
	assert.Equal(t, "127.0.0.1:27017", dbi.Addrs[0])
	assert.Equal(t, "", dbi.User)
	assert.Equal(t, "", dbi.Password)
	assert.Equal(t, "", dbi.AuthDatabase)
	assert.Equal(t, 3, dbi.Timeout)
	assert.Equal(t, 5, dbi.MaxConn)
	assert.Equal(t, true, dbi.Direct)
	assert.Equal(t, mgo.Primary, dbi.ReadMode)

	// secondaryPreferred
	dbi = NewDBInfo("testmgo", []string{"127.0.0.1:27017"}, "", "", "", 3, 1, true, true)
	assert.Equal(t, mgo.SecondaryPreferred, dbi.ReadMode)
}

func TestNewSessionPool(t *testing.T) {
	var sPool *Pool

	resource, err := dockerPool.Run("mongo", "3.4", nil)
	assert.NoError(t, err)
	addr := fmt.Sprintf("localhost:%s", resource.GetPort("27017/tcp"))
	testPoolConfig := NewDBInfo("testmgo", []string{addr}, "", "", "", 1, 1, true, false)

	err = dockerPool.Retry(func() error {
		var errMongo error
		sPool, errMongo = NewSessionPool(testPoolConfig)
		if errMongo != nil {
			return errMongo
		}

		return sPool.Ping(testCtx)
	})
	assert.NoError(t, err)

	assert.True(t, sPool.IsAvailable())
	assert.Equal(t, testPoolConfig.MaxConn, sPool.Len())
	assert.Equal(t, testPoolConfig.MaxConn, sPool.Cap())
	assert.Equal(t, testPoolConfig.ReadMode, sPool.Mode())
	assert.Equal(t, testPoolConfig, sPool.Config())

	conf := sPool.ShowConfig()
	assert.Equal(t, testPoolConfig.MaxConn, conf["MaxConn"])
	assert.Equal(t, testPoolConfig.Addrs, conf["Addrs"])
	assert.Equal(t, testPoolConfig.Timeout, conf["Timeout"])
	assert.Equal(t, testPoolConfig.Direct, conf["Direct"])

	conn, err := sPool.get(testCtx)
	assert.NoError(t, err)
	// test addr
	assert.Equal(t, addr, conn.Addr())
	// cannot get since there is only one conn in pool
	timeoutCtx := &droictx.DoneContext{}
	timeoutCtx.SetTimeout(10*time.Millisecond, nil)
	_, err = sPool.get(timeoutCtx)
	assert.Error(t, err)
	assert.Equal(t, me.ErrAPIFullResource, err)

	errTest := droipkg.NewCarrierDroiError(1234567, "unit test")
	timeoutCtx = &droictx.DoneContext{}
	timeoutCtx.SetTimeout(10*time.Millisecond, errTest)
	_, err = sPool.get(timeoutCtx)
	assert.Error(t, err)
	assert.Equal(t, 1234567, errTest.ErrorCode())
	assert.Equal(t, "unit test", errTest.Error())

	sPool.backgroundReconnect(conn)

	// simulate d/c reconnect
	dockerPool.Purge(resource)
	sPool.backgroundReconnect(conn)

	assert.Error(t, <-AlertChannel)
	// cap -1
	assert.Equal(t, testPoolConfig.MaxConn-1, sPool.Cap())

	sPool.Close()
	// connect fail
	sPool, err = NewSessionPool(testPoolConfig)
	assert.Error(t, err)

}

func TestPoolClose(t *testing.T) {
	var testPool *Pool

	resource, err := dockerPool.Run("mongo", "3.4", nil)
	assert.NoError(t, err)

	addr := fmt.Sprintf("localhost:%s", resource.GetPort("27017/tcp"))
	testPoolConfig := NewDBInfo("mgo", []string{addr}, "", "", "", 2, 1, true, false)

	err = dockerPool.Retry(func() error {
		var errMongo error
		testPool, errMongo = NewSessionPool(testPoolConfig)
		if errMongo != nil {
			return errMongo
		}

		return testPool.Ping(testCtx)
	})
	assert.NoError(t, err)

	conn, err := testPool.get(testCtx)
	assert.NoError(t, err)
	go testPool.Close()
	// here, session pool should not close becuase there is an session outside
	assert.Equal(t, testPoolConfig.MaxConn-1, testPool.Len())

	// block close, let close for loop wait
	time.Sleep(2 * time.Second)
	testPool.put(conn)
	time.Sleep(1 * time.Second)
	// closed
	assert.Equal(t, 0, testPool.Len())

	// pool is closed cannot get
	_, err = testPool.get(testCtx)
	assert.Error(t, err)

	// teardown docker
	dockerPool.Purge(resource)
}
