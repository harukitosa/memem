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

// CacheSyncStore is the store implemented with sync.Map
// Note: use it only when the number of read is much larger than the number of write.
type CacheSyncStore[K comparable, T any] struct {
	callback       func() T
	cacheClearTime time.Duration
	store          sync.Map
}

func (c *CacheSyncStore[K, T]) Set(key K, value T) {
	c.store.Store(key, ValueWithTime[T]{Time: time.Now(), Value: value})
}

func (c *CacheSyncStore[K, T]) Get(key K) T {
	v, ok := c.store.Load(key)
	if ok {
		value, ok := v.(ValueWithTime[T])
		if ok {
			diff := time.Since(value.Time)

			isValidTime := diff <= c.cacheClearTime
			// 値が存在するかつキャッシュ削除期間じゃない
			if ok && isValidTime {
				return value.Value
			}

			// 値が存在するけど、キャッシュ削除期間なので、callbackに流す

		} else {
			// mememがバグっていない限り、ここには来ないはず
			// callbackに流すと、バグに気が付けないので、panicで知らせる
			panic("shouldn't reach here")
		}
	}

	// Loadに失敗 or キャッシュ削除期間 -> callbackを試みる

	if c.callback != nil {
		callbackValue := c.callback()
		c.Set(key, callbackValue)
		return callbackValue
	}

	var res T
	// それ以外はnilで返却
	return res
}

func (c *CacheSyncStore[K, T]) Clear() {
	c.store = sync.Map{}
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

type Mode uint

const (
	// Use sync.Map as backend
	// Note: use it only when the number of read is much larger than the number of write.
	UseSyncMap Mode = 1 << iota
)

func aggregateOptions(opt ...Mode) Mode {
	var ret Mode
	for _, o := range opt {
		ret |= o
	}
	return ret
}

func NewCache[K comparable, T any](opt ...Mode) Cache[K, T] {
	// cacheClearTime: キャッシュの有効時間、defaultはISUCONに合わせて60s

	aggregatedOpt := aggregateOptions(opt...)
	if aggregatedOpt&UseSyncMap != 0 {
		return &CacheSyncStore[K, T]{
			store:          sync.Map{},
			callback:       nil,
			cacheClearTime: time.Second * 60,
		}
	}

	m := make(map[K]ValueWithTime[T])

	return &CacheInMemory[K, T]{
		store:          m,
		callback:       nil,
		cacheClearTime: time.Second * 60,
	}
}

func NewCacheWithCallback[K comparable, T any](callback func() T, opt ...Mode) Cache[K, T] {
	aggregatedOpt := aggregateOptions(opt...)
	if aggregatedOpt&UseSyncMap != 0 {
		return &CacheSyncStore[K, T]{
			store:          sync.Map{},
			callback:       callback,
			cacheClearTime: time.Second * 60,
		}
	}

	m := make(map[K]ValueWithTime[T])

	return &CacheInMemory[K, T]{
		store:          m,
		callback:       callback,
		cacheClearTime: time.Second * 60,
	}
}

func NewCacheWithClearTime[K comparable, T any](cleartime time.Duration, opt ...Mode) Cache[K, T] {
	aggregatedOpt := aggregateOptions(opt...)
	if aggregatedOpt&UseSyncMap != 0 {
		return &CacheSyncStore[K, T]{
			store:          sync.Map{},
			callback:       nil,
			cacheClearTime: cleartime,
		}
	}

	m := make(map[K]ValueWithTime[T])
	return &CacheInMemory[K, T]{
		store:          m,
		callback:       nil,
		cacheClearTime: cleartime,
	}
}

func NewCacheWithCallbackAndClearTime[K comparable, T any](callback func() T, cleartime time.Duration, opt ...Mode) Cache[K, T] {
	aggregatedOpt := aggregateOptions(opt...)
	if aggregatedOpt&UseSyncMap != 0 {
		return &CacheSyncStore[K, T]{
			store:          sync.Map{},
			callback:       callback,
			cacheClearTime: cleartime,
		}
	}

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

	if ok && isValidTime {
		// 値が存在するかつキャッシュ削除期間じゃない
		return value.Value
	}

	// 値が存在しない or キャッシュ削除期間　-> callbackを試みる

	if c.callback != nil {
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
