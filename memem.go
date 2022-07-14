package memem

import (
	"sync"
	"time"
)

type Cache interface {
	Set(key string, value interface{})
	Get(key string) interface{}
	GetOrClearIfOverTheTimeLimit(key string, clearTime time.Duration) interface{}
}

type CacheInMemory struct {
	store    map[string]ValueWithTime
	callback func() interface{}
	sync.Mutex
}

type ValueWithTime struct {
	Value interface{}
	Time  time.Time
}

func NewCacheWithCallback(callback func() interface{}) Cache {
	m := make(map[string]ValueWithTime)
	return &CacheInMemory{
		store:    m,
		callback: callback,
	}
}

func NewCache() Cache {
	m := make(map[string]ValueWithTime)
	return &CacheInMemory{
		store: m,
	}
}

func (c *CacheInMemory) Set(key string, value interface{}) {
	c.Lock()
	c.store[key] = ValueWithTime{Time: time.Now(), Value: value}
	c.Unlock()
}

func (c *CacheInMemory) Get(key string) interface{} {
	value, ok := c.store[key]
	// 値が存在しなくてcallbackがある場合はそれを利用する
	if !ok && c.callback != nil {
		callbackValue := c.callback()
		c.Set(key, callbackValue)
		return callbackValue
	}
	return value.Value
}

func (c *CacheInMemory) GetOrClearIfOverTheTimeLimit(key string, clearTime time.Duration) interface{} {
	value, ok := c.store[key]
	if !ok && c.callback != nil {
		callbackValue := c.callback()
		c.Set(key, callbackValue)
		return callbackValue
	}
	now := time.Now()
	diff := now.Sub(value.Time)
	if diff <= clearTime {
		return value.Value
	}

	if c.callback != nil {
		callbackValue := c.callback()
		c.Set(key, callbackValue)
		return callbackValue
	}
	return nil
}
