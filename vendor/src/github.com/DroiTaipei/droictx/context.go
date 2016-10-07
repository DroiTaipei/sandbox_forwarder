package droictx

import (
	"io"
)

type contextKV struct {
	key   string
	value interface{}
}

type Context []contextKV

func (c *Context) Set(key string, value interface{}) {
	args := *c
	n := len(args)
	for i := 0; i < n; i++ {
		kv := &args[i]
		if kv.key == key {
			kv.value = value
			return
		}
	}

	ca := cap(args)
	if ca > n {
		args = args[:n+1]
		kv := &args[n]
		kv.key = key
		kv.value = value
		*c = args
		return
	}

	kv := contextKV{
		key:   key,
		value: value,
	}
	*c = append(args, kv)
}

func (c *Context) Get(key string) interface{} {
	args := *c
	n := len(args)
	for i := 0; i < n; i++ {
		kv := &args[i]
		if kv.key == key {
			return kv.value
		}
	}
	return nil
}

func (c *Context) GetString(key string) (value string, ok bool) {
	v := c.Get(key)
	if v == nil {
		return
	}
	value, ok = v.(string)
	return
}

func (c *Context) GetInt(key string) (value int, ok bool) {
	v := c.Get(key)
	if v == nil {
		return
	}
	value, ok = v.(int)
	return
}

func (c *Context) GetInt64(key string) (value int64, ok bool) {
	v := c.Get(key)
	if v == nil {
		return
	}
	value, ok = v.(int64)
	return
}

func (c *Context) Map() (ret map[string]interface{}) {
	ret = make(map[string]interface{})
	args := *c
	n := len(args)
	for i := 0; i < n; i++ {
		ret[args[i].key] = args[i].value
	}
	return
}

func (c *Context) Reset() {
	args := *c
	n := len(args)
	for i := 0; i < n; i++ {
		v := args[i].value
		if vc, ok := v.(io.Closer); ok {
			vc.Close()
		}
	}
	*c = (*c)[:0]
}
