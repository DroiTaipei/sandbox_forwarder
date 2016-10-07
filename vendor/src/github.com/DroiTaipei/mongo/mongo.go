package mongo

import (
	"log"
	"sync"
	"time"

	"strings"

	"github.com/DroiTaipei/droipkg"
	"github.com/DroiTaipei/mgo"
)

const (
	MAX_CONNECT_RETRY = 20
)

var mgoDatabaseName string
var mgoDefaultPK = "_id"
var mgoReadSessionPool chan *RootSession
var mgoWriteSessionPool chan *RootSession
var mgoShardSession *mgo.Session
var mgoConfig *DBInfo
var reconnectFlag = false
var reconnectWG sync.WaitGroup
var mgoLiveServers []string

// DBInfo logs the required info for baas mongodb. Common dbapi don't need Shard.
type DBInfo struct {
	MaxConn       int
	MaxIdle       int
	Name          string
	Addrs         []string
	User          string
	Password      string
	Database      string
	AuthDatabase  string
	Timeout       int
	ShardUser     string
	ShardPassword string
	Direct        bool
}

// NewDBInfo returns the DBInfo instance.
func NewDBInfo(name string, addrs []string, user, password, dbName, authdbName string,
	timeout, maxConn, maxIdle int, shardUser, shardPassword string, direct bool) *DBInfo {
	return &DBInfo{
		MaxConn:       maxConn,
		MaxIdle:       maxIdle,
		Name:          name,
		Addrs:         addrs,
		User:          user,
		Password:      password,
		Database:      dbName,
		AuthDatabase:  authdbName,
		Timeout:       timeout,
		ShardUser:     shardUser,
		ShardPassword: shardPassword,
		Direct:        direct,
	}
}

// 4 log the entry which session connecting to, used by debug/error log.
type RootSession struct {
	Session *mgo.Session
	Mongos  string
}

// getWriteSess retrieves the write session which setMode =primary
func getWriteSess() (*mgo.Session, string) {
	if reconnectFlag {
		reconnectWG.Wait()
	}
	conn := <-mgoWriteSessionPool
	newSession := conn.Session.Clone()
	addr := conn.Mongos
	mgoWriteSessionPool <- conn

	return newSession, addr
}

// getReadSess retrieves the read session which setMode = SecondaryPreferred
func getReadSess() (*mgo.Session, string) {
	if reconnectFlag {
		reconnectWG.Wait()
	}
	conn := <-mgoReadSessionPool
	newSession := conn.Session.Clone()
	addr := conn.Mongos
	mgoReadSessionPool <- conn

	return newSession, addr
}

func getShardSess() *mgo.Session {
	if reconnectFlag {
		reconnectWG.Wait()
	}

	newSession := mgoShardSession.Clone()

	return newSession
}

// Database common operation function
func getMongoCollection(s *mgo.Session, colName string) *mgo.Collection {
	return s.DB(mgoDatabaseName).C(colName)
}

// Reconnect does the reconnection process only one process in one period of time.
func Reconnect() (err error) {
	if reconnectFlag {
		reconnectWG.Wait()
	} else {
		// only one can do this thing
		reconnectFlag = true
		reconnectWG.Add(1)
		timeout := mgoConfig.Timeout
		maxAttempts := timeout / 3
		// reconnect
		for attempts := 1; attempts <= maxAttempts; attempts++ {
			err = connectMongo()
			if err == nil {
				break
			}
			time.Sleep(time.Second)
		}
		reconnectWG.Done()
		// let reqs come in the nearest time period pass this action
		go func() {
			time.Sleep(time.Duration(timeout) * time.Second)
			reconnectFlag = false
		}()
	}
	return err
}

func SetLogger(l droipkg.Logger) {
	droipkg.SetLogger(l)
}

// NewSession new and return a session and it's connecting mognos address.
func NewSession(dbi *DBInfo, addr []string, mode mgo.Mode) (newSession *mgo.Session, mongosAddr string, err error) {
	dialInfo := mgo.DialInfo{
		Addrs:     addr,
		Direct:    dbi.Direct,
		FailFast:  true,
		Source:    dbi.AuthDatabase,
		Username:  dbi.User,
		Password:  dbi.Password,
		Timeout:   time.Duration(dbi.Timeout) * time.Second,
		PoolLimit: 1,
	}
	maxAttempts := 20

	for attempts := 1; attempts <= maxAttempts; attempts++ {
		newSession, err = mgo.DialWithInfo(&dialInfo)
		if err == nil {
			break
		}
		errLog(systemCtx, "no reachable server")
		time.Sleep(time.Duration(attempts) * time.Second)
	}
	if err != nil {
		return
	}
	newSession.SetMode(mode, true)
	newSession.Ping()
	mongosAddr = newSession.Server().Addr
	// fmt.Println("connect mongos:", newSession.Server().Addr)
	return
}

// newShardSession new and return the session of sharding user.
func newShardSession(dbi *DBInfo) (sharderSession *mgo.Session, err error) {
	shardDialInfo := mgo.DialInfo{
		Addrs:     dbi.Addrs,
		Direct:    dbi.Direct,
		FailFast:  true,
		Source:    dbi.AuthDatabase,
		Username:  dbi.ShardUser,
		Password:  dbi.ShardPassword,
		Timeout:   time.Duration(dbi.Timeout) * time.Second,
		PoolLimit: 1,
	}
	maxAttempts := 10
	for attempts := 1; attempts <= maxAttempts; attempts++ {
		sharderSession, err = mgo.DialWithInfo(&shardDialInfo)
		if err == nil {
			break
		}
		errLog(systemCtx, "shard session no reachable server")
		time.Sleep(time.Duration(attempts) * time.Second)
	}
	return
}

// connectMongo init or refresh the connection in Mgo
func connectMongo() (err error) {
	dbi := mgoConfig
	readPool := make(chan *RootSession, dbi.MaxConn)
	writePool := make(chan *RootSession, dbi.MaxConn)

	var shardSession *mgo.Session
	var newSession *mgo.Session
	var mongosAddr string
	addrAllocations := make(map[string]int)
	for _, v := range dbi.Addrs {
		addrAllocations[v] = 0
	}
	// get LiveServers
	rootSession, _, dialErr := NewSession(dbi, dbi.Addrs,mgo.Primary)
	if dialErr != nil {
		errLog(systemCtx, "unable to connect to ["+strings.Join(dbi.Addrs, ",")+"] because:"+dialErr.Error())
		err = dialErr
		return
	}
	defer rootSession.Close()
	liveServers := rootSession.LiveServers()
	droipkg.GetLogger().Debug("live:", liveServers)
	// read
	lengthOfMongos := len(liveServers)
	for i := 0; i < dbi.MaxConn; i++ {
		newSession, mongosAddr, err = NewSession(dbi, []string{liveServers[i%lengthOfMongos]}, mgo.SecondaryPreferred)
		if err == nil {
			addrAllocations[mongosAddr] = addrAllocations[mongosAddr] + 1
			newConn := &RootSession{Session: newSession, Mongos: mongosAddr}
			readPool <- newConn
		}
	}
	if mgoReadSessionPool != nil {
		time.Sleep(time.Second)
		sessionPoolClose(mgoReadSessionPool)
	}
	mgoReadSessionPool = readPool

	// write
	for i := 0; i < dbi.MaxConn; i++ {
		newSession, mongosAddr, err = NewSession(dbi, []string{liveServers[i%lengthOfMongos]}, mgo.Primary)
		if err == nil {
			addrAllocations[mongosAddr] = addrAllocations[mongosAddr] + 1
			newConn := &RootSession{Session: newSession, Mongos: mongosAddr}
			writePool <- newConn
		}
	}
	if mgoWriteSessionPool != nil {
		time.Sleep(time.Second)
		sessionPoolClose(mgoWriteSessionPool)
	}
	mgoWriteSessionPool = writePool

	for k, v := range addrAllocations {
		droipkg.GetLogger().Debug("mongos ", k, ": connnections:", v)
		log.Println("mongos ", k, ": connnections:", v)
	}
	addrAllocations = nil
	mgoLiveServers = nil
	mgoLiveServers = liveServers

	// short cut for not using shardSession
	if len(dbi.ShardUser) < 1 {
		return
	}
	shardSession, err = newShardSession(dbi)
	if err != nil {
		errLog(systemCtx, "unable to connect to ["+strings.Join(dbi.Addrs, ",")+"] because:"+err.Error())
		return
	}
	if mgoShardSession != nil {
		shardClose()
	}
	mgoShardSession = shardSession

	return
}

func sessionPoolClose(mgoSessionPool chan *RootSession) {
	for i := 0; i < len(mgoSessionPool); i++ {
		con := <-mgoSessionPool
		con.Session.Close()
	}
	close(mgoSessionPool)
}

func shardClose() {
	mgoShardSession.Close()
}

func Close() {
	sessionPoolClose(mgoWriteSessionPool)
	sessionPoolClose(mgoReadSessionPool)
	shardClose()
	mgoDatabaseName = ""
	mgoDefaultPK = "_id"
	mgoConfig = nil
	reconnectFlag = false
	mgoLiveServers = []string{}
}

// Initialize init mongo instance
func Initialize(dbi *DBInfo, PKField string, l droipkg.Logger) (err error) {
	SetLogger(l)
	mgoConfig = dbi
	mgoDatabaseName = dbi.Database
	mgoDefaultPK = PKField
	err = connectMongo()
	return
}

func CheckMongos() error {
	rootSession, _, err := NewSession(mgoConfig, mgoConfig.Addrs, mgo.Primary)
	if err != nil {
		errLog(systemCtx, "unable to connect to ["+strings.Join(mgoConfig.Addrs, ",")+"] because:"+err.Error())
		return err
	}
	defer rootSession.Close()
	liveServers := rootSession.LiveServers()
	oldLiveServers := mgoLiveServers
	if len(liveServers) > len(oldLiveServers) {
		errLog(systemCtx, "API try rebalance mongos")
		Reconnect()
	}
	return nil
}

func ShowConfig() map[string]interface{} {
	config := make(map[string]interface{})
	config["MaxConn"] = mgoConfig.MaxConn
	config["MaxIdle"] = mgoConfig.MaxIdle
	config["Addrs"] = mgoConfig.Addrs
	config["Database"] = mgoConfig.Database
	config["Timeout"] = mgoConfig.Timeout
	config["Direct"] = mgoConfig.Direct
	return config
}
