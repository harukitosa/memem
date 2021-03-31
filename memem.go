package memem

import "sync"

type Cache interface {
	Set(key string, value interface{})
	Get(key string) interface{}
}

type CacheInMemory struct {
	store map[string]interface{}
	sync.Mutex
}

func NewCache() Cache {
	m := make(map[string]interface{})
	return &CacheInMemory{
		store: m,
	}
}

func (c *CacheInMemory) Set(key string, value interface{}) {
	c.Lock()
	c.store[key] = value
	c.Unlock()
}

func (c *CacheInMemory) Get(key string) interface{} {
	return c.store[key]
}
