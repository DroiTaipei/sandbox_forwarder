package connpool

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPoolsRoundRobinError(t *testing.T) {
	assert := assert.New(t)
	ps := NewPools()
	p, err := ps.GetRoundRobin()
	assert.Nil(p)
	assert.NotNil(err)
}

func TestPoolsRoundGetError(t *testing.T) {
	assert := assert.New(t)
	ps := NewPools()
	p, err := ps.Get("Not.Exist.Server")
	assert.NotNil(p)
	assert.NotNil(err)
	assert.Equal("", p.Addr)
}

func TestPoolsGet(t *testing.T) {
	assert := assert.New(t)
	ps := NewPools()
	b := len(servers)
	for i := 0; i < b; i++ {
		err := ps.Connect(servers[i].getAddr(), 1)
		assert.Nil(err)
	}
	p, err := ps.Get(servers[0].getAddr())
	assert.Nil(err)
	assert.Equal(servers[0].getAddr(), p.Addr)
}

func TestPoolsRoundRobin(t *testing.T) {
	assert := assert.New(t)
	ps := NewPools()
	b := len(servers)
	for i := 0; i < b; i++ {
		err := ps.Connect(servers[i].getAddr(), 1)
		assert.Nil(err)
	}
	for i := 0; i < b; i++ {
		p, err := ps.GetRoundRobin()
		assert.Nil(err)
		idx := (i + 1) % len(servers)
		assert.Equal(servers[idx].getAddr(), p.Addr)
	}

}
