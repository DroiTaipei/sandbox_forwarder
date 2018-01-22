package connpool

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"log"
	"math"
	"net"
	"os"
	"strconv"
	"testing"
	"time"
)

const (
	serverNum = 3
)

var (
	servers  []*server
	gServers []*grpc.Server
	cleaners []func()
)

func BeforeTest() {
	var sCleanUp, gsCleanUp func()
	servers, sCleanUp = startServers(serverNum, math.MaxUint32)
	cleaners = append(cleaners, sCleanUp)
	gServers, gsCleanUp = startGrpcServers(servers)
	cleaners = append(cleaners, gsCleanUp)
}

func TestPoolUpperLimit(t *testing.T) {
	assert := assert.New(t)
	p, err := NewPool(servers[0].getAddr(), 1)
	assert.Nil(err)
	c1, err := p.Get()
	assert.NotNil(c1)
	assert.Nil(err)
	c2, err := p.Get()
	assert.Nil(c2)
	assert.NotNil(err)
}

func TestPoolPut(t *testing.T) {
	assert := assert.New(t)
	p, err := NewPool(servers[0].getAddr(), 1)
	assert.Nil(err)
	c1, err := p.Get()
	assert.NotNil(c1)
	assert.Nil(err)
	p.Put(c1)
	c2, err := p.Get()
	assert.NotNil(c2)
	assert.Nil(err)
	assert.Equal(c1, c2)
}

func TestPoolClose(t *testing.T) {
	assert := assert.New(t)
	p, err := NewPool(servers[0].getAddr(), 1)
	assert.Nil(err)
	c1, err := p.Get()
	assert.NotNil(c1)
	assert.Nil(err)
	err = p.Close(c1)
	assert.Nil(err)
	c2, err := p.Get()
	assert.NotNil(c2)
	assert.Nil(err)
	assert.NotEqual(c1, c2)
}

// Do somethings after all test cases
func AfterTest() {
	for idx := range cleaners {
		cleaners[idx]()
	}
}

func TestMain(m *testing.M) {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags)

	BeforeTest()
	retCode := m.Run()
	AfterTest()
	os.Exit(retCode)
}

func startServers(numServers int, maxStreams uint32) ([]*server, func()) {
	var servers []*server
	for i := 0; i < numServers; i++ {
		s := newTestServer()
		servers = append(servers, s)
		go s.start(0, maxStreams)
		s.wait(2 * time.Second)
	}
	return servers, func() {
		for i := 0; i < numServers; i++ {
			servers[i].stop()
		}
	}
}
func startGrpcServers(ss []*server) ([]*grpc.Server, func()) {
	numServers := len(ss)
	var gServers []*grpc.Server

	for i := 0; i < numServers; i++ {
		gs := grpc.NewServer()
		go gs.Serve(ss[i].lis)
		gServers = append(gServers, gs)
	}
	return gServers, func() {
		for i := 0; i < numServers; i++ {
			gServers[i].Stop()
		}
	}
}

// Following copy from https://github.com/grpc/grpc-go/blob/master/call_test.go#L117

type server struct {
	lis        net.Listener
	port       string
	startedErr chan error // sent nil or an error after server starts
}

func newTestServer() *server {
	return &server{startedErr: make(chan error, 1)}
}

// start starts server. Other goroutines should block on s.startedErr for further operations.
func (s *server) start(port int, maxStreams uint32) {
	var err error
	if port == 0 {
		s.lis, err = net.Listen("tcp", "localhost:0")
	} else {
		s.lis, err = net.Listen("tcp", "localhost:"+strconv.Itoa(port))
	}
	if err != nil {
		s.startedErr <- fmt.Errorf("failed to listen: %v", err)
		return
	}
	_, p, err := net.SplitHostPort(s.lis.Addr().String())
	if err != nil {
		s.startedErr <- fmt.Errorf("failed to parse listener address: %v", err)
		return
	}
	s.port = p
	s.startedErr <- nil

}

func (s *server) wait(timeout time.Duration) {
	select {
	case err := <-s.startedErr:
		if err != nil {
			log.Fatal(err)
		}
	case <-time.After(timeout):
		log.Fatalf("Timed out after %v waiting for server to be ready", timeout)
	}
}

func (s *server) stop() {
	s.lis.Close()
}

func (s *server) getAddr() string {
	return s.lis.Addr().String()
}
