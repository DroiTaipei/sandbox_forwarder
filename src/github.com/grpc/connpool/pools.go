package connpool

import (
	"errors"
	"sync"
	"sync/atomic"
)

type Pools interface {
	Get(string) (*Pool, error)
	GetRoundRobin() (*Pool, error)
	Connect(string, int) error
}

type Poolsi struct {
	sync.RWMutex
	all   map[string]*Pool
	keys  []string
	robin uint64
}

func NewPools() *Poolsi {
	return &Poolsi{
		all:  make(map[string]*Pool),
		keys: make([]string, 0),
	}
}

func (p *Poolsi) Get(addr string) (*Pool, error) {
	p.RLock()
	defer p.RUnlock()

	pool, ok := p.all[addr]
	if !ok {
		return &Pool{}, errors.New("Pool does not exist in map")
	}

	return pool, nil
}

func (p *Poolsi) GetRoundRobin() (*Pool, error) {

	if len(p.keys) == 0 {
		return nil, errors.New("no servers")
	}

	return p.all[p.keys[atomic.AddUint64(&p.robin, 1)%uint64(len(p.keys))]], nil
}

func (p *Poolsi) Connect(addr string, maxCap int) error {
	p.RLock()
	_, has := p.all[addr]
	p.RUnlock()
	if has {
		return nil
	}

	pool, err := NewPool(addr, maxCap)
	if err != nil {
		return err
	}

	p.Lock()
	defer p.Unlock()
	_, has = p.all[addr]
	if has {
		return nil
	}
	p.all[addr] = pool
	p.keys = append(p.keys, addr)

	return nil
}
