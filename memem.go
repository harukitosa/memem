package memem

import (
	"sync"
	"time"
)

type Cache[T any] interface {
	Set(key string, value T)
	Get(key string) T
	Clear()
}

type CacheInMemory[T any] struct {
	store          map[string]ValueWithTime[T]
	callback       func() T
	cacheClearTime time.Duration
	mx             sync.RWMutex
}

type ValueWithTime[T any] struct {
	Value T
	Time  time.Time
}

func NewCache[T any]() Cache[T] {
	m := make(map[string]ValueWithTime[T])
	// cacheClearTime: キャッシュの有効時間、defaultはISUCONに合わせて60s
	return &CacheInMemory[T]{
		store:          m,
		callback:       nil,
		cacheClearTime: time.Second * 60,
	}
}

func NewCacheWithCallback[T any](callback func() T) Cache[T] {
	m := make(map[string]ValueWithTime[T])
	return &CacheInMemory[T]{
		store:          m,
		callback:       callback,
		cacheClearTime: time.Second * 60,
	}
}

func NewCacheWithClearTime[T any](cleartime time.Duration) Cache[T] {
	m := make(map[string]ValueWithTime[T])
	return &CacheInMemory[T]{
		store:          m,
		callback:       nil,
		cacheClearTime: cleartime,
	}
}

func NewCacheWithCallbackAndClearTime[T any](callback func() T, cleartime time.Duration) Cache[T] {
	m := make(map[string]ValueWithTime[T])
	return &CacheInMemory[T]{
		store:          m,
		callback:       callback,
		cacheClearTime: cleartime,
	}
}

func (c *CacheInMemory[T]) Set(key string, value T) {
	c.mx.Lock()
	c.store[key] = ValueWithTime[T]{Time: time.Now(), Value: value}
	c.mx.Unlock()
}

func (c *CacheInMemory[T]) Get(key string) T {
	c.mx.RLock()
	value, ok := c.store[key]
	c.mx.RUnlock()

	diff := time.Since(value.Time)

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
	var res T
	// それ以外はnilで返却
	return res
}

func (c *CacheInMemory[T]) Clear() {
	c.store = make(map[string]ValueWithTime[T])
}
