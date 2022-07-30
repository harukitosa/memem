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

// CacheSyncStore is the store implemented with sync.Map
// Note: use it only when the number of read is much larger than the number of write.
type CacheSyncStore[T any] struct {
	callback       func() T
	cacheClearTime time.Duration
	store          sync.Map
}

func (c *CacheSyncStore[T]) Set(key string, value T) {
	c.store.Store(key, ValueWithTime[T]{Time: time.Now(), Value: value})
}

func (c *CacheSyncStore[T]) Get(key string) T {
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

	// Loadに失敗 or キャッシュ削除期間 or キャッシュが見つかったけど、ValueWithTime[T]ではない

	// コールバックがある
	if c.callback != nil {
		callbackValue := c.callback()
		c.Set(key, callbackValue)
		return callbackValue
	}

	var res T
	// それ以外はnilで返却
	return res
}

func (c *CacheSyncStore[T]) Clear() {
	c.store = sync.Map{}
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

func NewCache[T any](opt ...Mode) Cache[T] {
	// cacheClearTime: キャッシュの有効時間、defaultはISUCONに合わせて60s

	aggregatedOpt := aggregateOptions(opt...)
	if aggregatedOpt&UseSyncMap != 0 {
		return &CacheSyncStore[T]{
			store:          sync.Map{},
			callback:       nil,
			cacheClearTime: time.Second * 60,
		}
	}

	m := make(map[string]ValueWithTime[T])

	return &CacheInMemory[T]{
		store:          m,
		callback:       nil,
		cacheClearTime: time.Second * 60,
	}
}

func NewCacheWithCallback[T any](callback func() T, opt ...Mode) Cache[T] {
	aggregatedOpt := aggregateOptions(opt...)
	if aggregatedOpt&UseSyncMap != 0 {
		return &CacheSyncStore[T]{
			store:          sync.Map{},
			callback:       callback,
			cacheClearTime: time.Second * 60,
		}
	}

	m := make(map[string]ValueWithTime[T])

	return &CacheInMemory[T]{
		store:          m,
		callback:       callback,
		cacheClearTime: time.Second * 60,
	}
}

func NewCacheWithClearTime[T any](cleartime time.Duration, opt ...Mode) Cache[T] {
	aggregatedOpt := aggregateOptions(opt...)
	if aggregatedOpt&UseSyncMap != 0 {
		return &CacheSyncStore[T]{
			store:          sync.Map{},
			callback:       nil,
			cacheClearTime: cleartime,
		}
	}

	m := make(map[string]ValueWithTime[T])
	return &CacheInMemory[T]{
		store:          m,
		callback:       nil,
		cacheClearTime: cleartime,
	}
}

func NewCacheWithCallbackAndClearTime[T any](callback func() T, cleartime time.Duration, opt ...Mode) Cache[T] {
	aggregatedOpt := aggregateOptions(opt...)
	if aggregatedOpt&UseSyncMap != 0 {
		return &CacheSyncStore[T]{
			store:          sync.Map{},
			callback:       callback,
			cacheClearTime: cleartime,
		}
	}

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
