package client

import (
	"context"
	"errors"
	geopb "github.com/DroiTaipei/droipkg/geo/protobuf"
	"github.com/DroiTaipei/droipkg/grpc/connpool"
	"google.golang.org/grpc"
	"log"
	"time"
)

type ClientFactory func(*grpc.ClientConn) geopb.GeoClient

type SimplifiedPool interface {
	Get(...string) (*grpc.ClientConn, func(*grpc.ClientConn) error, error)
}

type sPool struct {
	pools connpool.Pools
}

func (sp *sPool) Get(servers ...string) (*grpc.ClientConn, func(*grpc.ClientConn) error, error) {

	if sp.pools == nil {
		return nil, nil, errors.New("pools not initialized")
	}
	if len(servers) > 1 {
		return nil, nil, errors.New("too many servers")
	}
	var p *connpool.Pool
	var err error
	if len(servers) == 0 {
		p, err = sp.pools.GetRoundRobin()
	} else {
		p, err = sp.pools.Get(servers[0])
	}
	if err != nil {
		return nil, nil, err
	}
	conn, err := p.Get()
	if err != nil {
		return nil, nil, err
	}
	return conn, p.Put, nil
}

type Agent struct {
	sp      SimplifiedPool
	factory ClientFactory
}

var poolsFactory func() connpool.Pools

func init() {
	poolsFactory = func() connpool.Pools {
		return connpool.NewPools()
	}
}

func NewAgent(sp SimplifiedPool, cf ClientFactory) *Agent {
	return &Agent{sp, cf}
}

func Initialize(server string, maxConn int) (*Agent, error) {
	pools := poolsFactory()

	err := pools.Connect(server, maxConn)
	if err != nil {
		return nil, errors.New("can't connect to " + server)
	}

	log.Printf("[GEO Client] Sender add server: %v", server)

	return NewAgent(&sPool{pools}, geopb.NewGeoClient), nil
}

func (a *Agent) Echo(headers map[string]string) (map[string]string, error) {

	if headers == nil {
		return nil, errors.New("no headers")
	}

	conn, putback, err := a.sp.Get()
	if err != nil {
		return nil, err
	}

	c := a.factory(conn)
	putback(conn)

	req := geopb.Content{
		Headers: headers,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	r, err := c.Echo(ctx, &req)
	if err != nil {
		return nil, err
	}

	return r.Headers, nil
}

func (a *Agent) GetMaxmindCity(ip, lang string) (*geopb.MaxmindCityInfo, error) {
	conn, putback, err := a.sp.Get()
	if err != nil {
		return nil, err
	}

	c := a.factory(conn)
	putback(conn)

	req := geopb.IP{
		IP:   ip,
		Lang: lang,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return c.GetMaxmindCity(ctx, &req)
}

func (a *Agent) GetIpipCity(ip, lang string) (*geopb.IpipCityInfo, error) {
	conn, putback, err := a.sp.Get()
	if err != nil {
		return nil, err
	}

	c := a.factory(conn)
	putback(conn)

	req := geopb.IP{
		IP:   ip,
		Lang: lang,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return c.GetIpipCity(ctx, &req)
}
