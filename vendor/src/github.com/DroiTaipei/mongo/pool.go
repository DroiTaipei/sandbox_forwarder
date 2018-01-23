package mongo

import (
	"log"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"

	"github.com/DroiTaipei/droictx"
	"github.com/DroiTaipei/droipkg"
	me "github.com/DroiTaipei/droipkg/mongo"
	"github.com/DroiTaipei/mgo"
	conntrack "github.com/mwitkow/go-conntrack"
)

// AlertChannel put error message, wait for outer user (i.e., gobuster) pick and send.
var AlertChannel = make(chan error, 1)

// Default dial timeout value from https://github.com/DroiTaipei/mgo/blob/v2/cluster.go
var syncSocketTimeout = 5 * time.Second

// DBInfo logs the required info for baas mongodb.
type DBInfo struct {
	Name         string
	User         string
	Password     string
	AuthDatabase string
	Addrs        []string
	MaxConn      int
	Timeout      int
	ReadMode     mgo.Mode
	Direct       bool
}

// NewDBInfo returns the DBInfo instance.
func NewDBInfo(name string, addrs []string, user, password, authdbName string,
	timeout, maxConn int, direct, readSecondary bool) *DBInfo {
	readMode := mgo.Primary
	if readSecondary {
		readMode = mgo.SecondaryPreferred
	}
	return &DBInfo{
		MaxConn:      maxConn,
		Name:         name,
		Addrs:        addrs,
		User:         user,
		Password:     password,
		AuthDatabase: authdbName,
		Timeout:      timeout,
		Direct:       direct,
		ReadMode:     readMode,
	}
}

// Session is wrapper for mgo.Session with logging mongos addr
type Session struct {
	addr string
	s    *mgo.Session
}

// Session returns mgo.Session
func (s *Session) Session() *mgo.Session {
	return s.s
}

// Addr returns target mongos addr
func (s *Session) Addr() string {
	return s.addr
}

// Pool is the mgo session pool
type Pool struct {
	cap       int
	config    *DBInfo
	mode      mgo.Mode
	c         chan *Session
	available bool
}

func newSession(dbi *DBInfo, addr []string, mode mgo.Mode) (newSession *mgo.Session, err error) {
	podName, err := os.Hostname()
	if err != nil {
		panic("Get hostname failed")
	}
	conntrackDialer :=
		func(addr *mgo.ServerAddr) (net.Conn, error) {
			conntrackDialer := conntrack.NewDialFunc(
				conntrack.DialWithName(podName),
				conntrack.DialWithTracing(),
				conntrack.DialWithDialer(&net.Dialer{
					Timeout: syncSocketTimeout,
				}),
			)
			return conntrackDialer(addr.TCPAddr().Network(), addr.String())
		}

	dialInfo := mgo.DialInfo{
		Addrs:      addr,
		Direct:     dbi.Direct,
		FailFast:   true,
		Source:     dbi.AuthDatabase,
		Username:   dbi.User,
		Password:   dbi.Password,
		Timeout:    time.Duration(dbi.Timeout) * time.Second,
		PoolLimit:  1,
		DialServer: conntrackDialer,
	}

	for attempts := 1; attempts <= dbi.Timeout; attempts++ {
		newSession, err = mgo.DialWithInfo(&dialInfo)
		if err == nil {
			break
		}
		time.Sleep(time.Duration(attempts) * time.Second)
	}
	if err != nil {
		errStr := "[mongo] no reachable server"
		log.Println(errStr)
		errLog(systemCtx, errStr)
		return
	}
	newSession.SetMode(mode, true)
	newSession.Ping()
	return
}

// NewSessionPool construct connection pool
func NewSessionPool(dbi *DBInfo) (*Pool, error) {
	c := make(chan *Session, dbi.MaxConn)

	addrAllocations := make(map[string]int)
	for _, v := range dbi.Addrs {
		addrAllocations[v] = 0
	}
	// get LiveServers
	rootSession, dialErr := newSession(dbi, dbi.Addrs, mgo.Primary)
	if dialErr != nil {
		errLog(systemCtx, "unable to connect to ["+strings.Join(dbi.Addrs, ",")+"] because:"+dialErr.Error())
		return nil, dialErr
	}
	defer rootSession.Close()
	liveServers := rootSession.LiveServers()
	lengthOfMongos := len(liveServers)
	// shuffle the server lists
	for i := range liveServers {
		j := rand.Intn(i + 1)
		liveServers[i], liveServers[j] = liveServers[j], liveServers[i]
	}
	sessionCount := 0
	for i := 0; i < dbi.MaxConn; i++ {
		addr := liveServers[i%lengthOfMongos]
		newSession, err := newSession(dbi, []string{addr}, dbi.ReadMode)
		if err == nil {
			addrAllocations[addr] = addrAllocations[addr] + 1
			c <- &Session{addr: addr, s: newSession}
			sessionCount++
		}
	}

	for k, v := range addrAllocations {
		log.Println("[mongo] mongos ", k, ": connnections:", v)
	}

	return &Pool{c: c, config: dbi, cap: sessionCount, available: true, mode: dbi.ReadMode}, nil
}

// IsAvailable returns whether Pool availalbe
func (p *Pool) IsAvailable() bool {
	return p.available
}

func (p *Pool) get(ctx droictx.Context) (session *Session, err droipkg.DroiError) {
	if !p.available {
		err = me.ErrMongoPoolClosed
		return
	}

	select {
	case session = <-p.c:
	case <-ctx.Timeout():
		err = ctx.TimeoutErr()
		if err == nil {
			err = me.ErrAPIFullResource
			return
		}
		// wrap
		err = droipkg.NewCarrierDroiError(err.ErrorCode(), err.Error()+": "+me.ErrAPIFullResource.Error())
		return
	}
	return
}

// Len returns current Pool availalbe connections
func (p *Pool) Len() int {
	return len(p.c)
}

// Cap returns Pool capacity
func (p *Pool) Cap() int {
	return p.cap
}

// Mode returns mgo.Mode settings of Pool
func (p *Pool) Mode() mgo.Mode {
	return p.mode
}

// Config returns DBInfo of Pool
func (p *Pool) Config() *DBInfo {
	return p.config
}

func (p *Pool) put(session *Session) {
	p.c <- session
}

// Close gracefull shutdown conns and Pool status
func (p *Pool) Close() {
	// wait all session come back to pool
	p.available = false
	sessionTimeout := p.config.Timeout + 10
	for i := 0; i < sessionTimeout; i++ {
		// len(p.c) should not > than p.cap, but use >= to get along with error
		if len(p.c) >= p.cap {
			break
		}
		time.Sleep(time.Second)
	}

	close(p.c)
	for s := range p.c {
		s.s.Close()
	}
	p.cap = 0
	p.c = nil
}

// ShowConfig returns debug config info
func (p *Pool) ShowConfig() map[string]interface{} {
	config := make(map[string]interface{})
	config["MaxConn"] = p.config.MaxConn
	config["Addrs"] = p.config.Addrs
	config["Timeout"] = p.config.Timeout
	config["Direct"] = p.config.Direct
	return config
}

func (p *Pool) backgroundReconnect(s *Session) {
	s.s.Close()

	retry := 3
	for i := 0; i < retry; i++ {
		// newSession timeout in 210s
		newS, err := newSession(p.config, []string{s.addr}, p.mode)
		if err == nil {
			s.s = newS
			p.put(s)
			return
		}
	}

	// still cannot connet after 10 mins
	errRetryTotalFailed := droipkg.NewError("[mongo] Reconnect failed within 10 mins")
	p.cap--
	select {
	case AlertChannel <- errRetryTotalFailed:
	default:
		// just pass, no spam our alert and non-blocking
	}
}
