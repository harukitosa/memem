package memem

import (
	"sync"
	"time"
)

type Cache[K comparable, T any] interface {
	Set(key K, value T)
	Get(key K) T
	Clear()
}

type CacheInMemory[K comparable, T any] struct {
	store          map[K]ValueWithTime[T]
	callback       func() T
	cacheClearTime time.Duration
	mx             sync.RWMutex
}

type ValueWithTime[T any] struct {
	Value T
	Time  time.Time
}

func NewCache[K comparable, T any]() Cache[K, T] {
	m := make(map[K]ValueWithTime[T])
	// cacheClearTime: キャッシュの有効時間、defaultはISUCONに合わせて60s
	return &CacheInMemory[K, T]{
		store:          m,
		callback:       nil,
		cacheClearTime: time.Second * 60,
	}
}

func NewCacheWithCallback[K comparable, T any](callback func() T) Cache[K, T] {
	m := make(map[K]ValueWithTime[T])
	return &CacheInMemory[K, T]{
		store:          m,
		callback:       callback,
		cacheClearTime: time.Second * 60,
	}
}

func NewCacheWithClearTime[K comparable, T any](cleartime time.Duration) Cache[K, T] {
	m := make(map[K]ValueWithTime[T])
	return &CacheInMemory[K, T]{
		store:          m,
		callback:       nil,
		cacheClearTime: cleartime,
	}
}

func NewCacheWithCallbackAndClearTime[K comparable, T any](callback func() T, cleartime time.Duration) Cache[K, T] {
	m := make(map[K]ValueWithTime[T])
	return &CacheInMemory[K, T]{
		store:          m,
		callback:       callback,
		cacheClearTime: cleartime,
	}
}

func (c *CacheInMemory[K, T]) Set(key K, value T) {
	c.mx.Lock()
	c.store[key] = ValueWithTime[T]{Time: time.Now(), Value: value}
	c.mx.Unlock()
}

func (c *CacheInMemory[K, T]) Get(key K) T {
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

func (c *CacheInMemory[K, T]) Clear() {
	c.store = make(map[K]ValueWithTime[T])
}
