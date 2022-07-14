package memem

import (
	"sync"
)

type Cache interface {
	Set(key string, value interface{})
	Get(key string) interface{}
}

type CacheInMemory struct {
	store    map[string]interface{}
	callback func() interface{}
	sync.Mutex
}

func NewCacheWithCallback(callback func() interface{}) Cache {
	m := make(map[string]interface{})
	return &CacheInMemory{
		store:    m,
		callback: callback,
	}
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
	value := c.store[key]
	// 値が存在しなくてcallbackがある場合はそれを利用する
	if value == nil && c.callback != nil {
		callbackValue := c.callback()
		c.Set(key, callbackValue)
		return callbackValue
	}
	return c.store[key]
}
