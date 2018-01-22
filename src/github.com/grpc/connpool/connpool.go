package connpool

import (
	"errors"
	"google.golang.org/grpc"
	"io"
	"sync/atomic"
)

type Pool struct {
	conns  chan *grpc.ClientConn
	Addr   string
	maxCap int64
	count  int64
}

func NewPool(addr string, maxCap int) (*Pool, error) {
	p := new(Pool)
	p.Addr = addr
	p.conns = make(chan *grpc.ClientConn, maxCap)
	p.maxCap = int64(maxCap)

	conn, err := p.dialNew()
	if err != nil {
		return nil, err
	}
	p.conns <- conn
	return p, nil
}

func (p *Pool) dialNew() (*grpc.ClientConn, error) {

	atomic.AddInt64(&p.count, 1)

	if p.count > p.maxCap {
		atomic.AddInt64(&p.count, -1)
		return nil, errors.New("connection has reached upper limit")
	}
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	// opts = append(opts, grpc.WithKeepaliveParams(keepalive.ClientParameters{Timeout: time.Second * 15}))

	c, err := grpc.Dial(p.Addr, opts...)
	if err != nil {
		atomic.AddInt64(&p.count, -1)
		return nil, err
	}

	return c, nil
}

func (p *Pool) Close(conn io.Closer) error {

	atomic.AddInt64(&p.count, -1)

	return conn.Close()
}

func (p *Pool) Get() (*grpc.ClientConn, error) {

	select {
	case conn := <-p.conns:
		return conn, nil
	default:
		return p.dialNew()
	}
}

func (p *Pool) Put(conn *grpc.ClientConn) error {
	select {
	case p.conns <- conn:
		return nil
	default:
		return p.Close(conn)
	}
}
