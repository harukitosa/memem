package memem

import (
	"sync"
	"time"
)

type Cache interface {
	Set(key string, value interface{})
	Get(key string) interface{}
}

type CacheInMemory struct {
	store          map[string]ValueWithTime
	callback       func() interface{}
	cacheClearTime time.Duration
	mx             sync.RWMutex
}

type ValueWithTime struct {
	Value interface{}
	Time  time.Time
}

func NewCache() Cache {
	m := make(map[string]ValueWithTime)
	// cacheClearTime: キャッシュの有効時間、defaultは1時間
	return &CacheInMemory{
		store:          m,
		callback:       nil,
		cacheClearTime: time.Hour,
	}
}

func NewCacheWithCallback(callback func() interface{}) Cache {
	m := make(map[string]ValueWithTime)
	return &CacheInMemory{
		store:          m,
		callback:       callback,
		cacheClearTime: time.Hour,
	}
}

func NewCacheWithClearTime(cleartime time.Duration) Cache {
	m := make(map[string]ValueWithTime)
	return &CacheInMemory{
		store:          m,
		callback:       nil,
		cacheClearTime: cleartime,
	}
}

func NewCacheWithCallbackAndClearTime(callback func() interface{}, cleartime time.Duration) Cache {
	m := make(map[string]ValueWithTime)
	return &CacheInMemory{
		store:          m,
		callback:       callback,
		cacheClearTime: cleartime,
	}
}

func (c *CacheInMemory) Set(key string, value interface{}) {
	c.mx.Lock()
	c.store[key] = ValueWithTime{Time: time.Now(), Value: value}
	c.mx.Unlock()
}

func (c *CacheInMemory) Get(key string) interface{} {
	c.mx.RLock()
	value, ok := c.store[key]
	c.mx.RUnlock()

	diff := time.Now().Sub(value.Time)

	isValidTime := diff <= c.cacheClearTime

	// 値が存在するかつキャッシュ削除期間じゃない
	if ok && isValidTime {
		return value.Value
	}

	// コールバックがあるかつ、値が存在しないか存在しても期限切れの場合コールバックの値を返却する
	if c.callback != nil && (!ok || (ok && !isValidTime)) {
		callbackValue := c.callback()
		c.Set(key, callbackValue)
		return callbackValue
	}

	// それ以外はnilで返却
	return nil
}
