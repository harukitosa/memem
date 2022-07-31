package memem

import (
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	t.Parallel()
	c := NewCache[string, int]()
	c.Set("id", 12)
	if c.Get("id") != 12 {
		t.Error("error not match")
	}
	c.Set("hoge", 100)
	if c.Get("hoge") != 100 {
		t.Error("error")
	}
}

func TestCacheClear(t *testing.T) {
	t.Parallel()
	c := NewCache[string, int]()
	c.Set("id", 12)
	if c.Get("id") != 12 {
		t.Error("error not match")
	}
	c.Set("hoge", 100)
	if c.Get("hoge") != 100 {
		t.Error("error")
	}
	c.Clear()
	if c.Get("hoge") != 0 {
		t.Errorf("error")
	}
}

func TestArrayCache(t *testing.T) {
	t.Parallel()
	c := NewCache[string, []string]()
	slice := []string{"Golang", "Java"}
	c.Set("key", slice)
	value := c.Get("key")
	if value[0] != "Golang" {
		t.Error("error not match")
	}
}

func TestGetDataIsNoneCache(t *testing.T) {
	t.Parallel()
	c := NewCache[string, string]()
	value := c.Get("key")
	if value != "" {
		t.Error("not nil")
	}
}

func TestCallbackCache(t *testing.T) {
	t.Parallel()
	c := NewCacheWithCallback[string](func() string {
		return "callback result"
	})
	value := c.Get("callback result")
	if value != "callback result" {
		t.Error("is not callback value")
	}
}

func TestWithClearTimeNonClear(t *testing.T) {
	t.Parallel()
	c := NewCacheWithClearTime[string, string](time.Second)
	c.Set("key", "same value")
	value := c.Get("key")
	if value != "same value" {
		t.Error("data is null")
	}
}

func TestWithClearTimeClear(t *testing.T) {
	t.Parallel()
	c := NewCacheWithClearTime[string, string](time.Second)
	c.Set("key", "same value")
	time.Sleep(2 * time.Second)
	value := c.Get("key")
	if value != "" {
		t.Error("data is not null")
	}
}

type TestCase[T any] struct {
	f       func() T
	t       time.Duration
	isClear bool
}

func TestWithClearTimeAndCallback(t *testing.T) {
	f := func() string {
		return "callback value"
	}
	cases := []struct {
		name string
		in   TestCase[string]
		want string
	}{
		{
			"期限切れしていない場合キャッシュが返却される",
			TestCase[string]{
				f,
				time.Second,
				false,
			},
			"value",
		},
		{
			"期限が過ぎている場合callbackの値が返却される",
			TestCase[string]{
				f,
				time.Second,
				true,
			},
			"callback value",
		},
	}
	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := NewCacheWithCallbackAndClearTime[string](tt.in.f, tt.in.t)
			c.Set("key", "value")
			if tt.in.isClear {
				time.Sleep(2 * time.Second)
				value := c.Get("key")
				if value != tt.want {
					t.Error("is not want value")
					return
				}
			} else {
				value := c.Get("key")
				if value != tt.want {
					t.Error("is not want value")
					return
				}
			}
		})
	}
}

// test for CacheSyncStore

func Test_CacheSyncStore_Cache(t *testing.T) {
	c := NewCache[string, int](UseSyncMap)
	c.Set("id", 12)
	if c.Get("id") != 12 {
		t.Error("error not match")
	}
	c.Set("hoge", 100)
	if c.Get("hoge") != 100 {
		t.Error("error")
	}
}

func Test_CacheSyncStore_CacheClear(t *testing.T) {
	c := NewCache[string, int](UseSyncMap)
	c.Set("id", 12)
	if c.Get("id") != 12 {
		t.Error("error not match")
	}
	c.Set("hoge", 100)
	if c.Get("hoge") != 100 {
		t.Error("error")
	}
	c.Clear()
	if c.Get("hoge") != 0 {
		t.Error("error")
	}
}

func Test_CacheSyncStore_ArrayCache(t *testing.T) {
	c := NewCache[string, []string](UseSyncMap)
	slice := []string{"Golang", "Java"}
	c.Set("key", slice)
	value := c.Get("key")
	if value[0] != "Golang" {
		t.Error("error not match")
	}
}

func Test_CacheSyncStore_GetDataIsNoneCache(t *testing.T) {
	c := NewCache[string, string](UseSyncMap)
	value := c.Get("key")
	if value != "" {
		t.Error("not nil")
	}
}

func Test_CacheSyncStore_CallbackCache(t *testing.T) {
	c := NewCacheWithCallback[string](func() string {
		return "callback result"
	}, UseSyncMap)
	value := c.Get("callback result")
	if value != "callback result" {
		t.Error("is not callback value")
	}
}

func Test_CacheSyncStore_WithClearTimeNonClear(t *testing.T) {
	c := NewCacheWithClearTime[string, string](time.Second, UseSyncMap)
	c.Set("key", "same value")
	value := c.Get("key")
	if value != "same value" {
		t.Error("data is null")
	}
}

func Test_CacheSyncStore_WithClearTimeClear(t *testing.T) {
	c := NewCacheWithClearTime[string, string](time.Second, UseSyncMap)
	c.Set("key", "same value")
	time.Sleep(2 * time.Second)
	value := c.Get("key")
	if value != "" {
		t.Error("data is not null")
	}
}

func Test_CacheSyncStore_WithClearTimeAndCallback(t *testing.T) {
	f := func() string {
		return "callback value"
	}
	cases := []struct {
		name string
		in   TestCase[string]
		want string
	}{
		{
			"期限切れしていない場合キャッシュが返却される",
			TestCase[string]{
				f,
				time.Second,
				false,
			},
			"value",
		},
		{
			"期限が過ぎている場合callbackの値が返却される",
			TestCase[string]{
				f,
				time.Second,
				true,
			},
			"callback value",
		},
	}
	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := NewCacheWithCallbackAndClearTime[string](tt.in.f, tt.in.t, UseSyncMap)
			c.Set("key", "value")
			if tt.in.isClear {
				time.Sleep(2 * time.Second)
				value := c.Get("key")
				if value != tt.want {
					t.Error("is not want value")
				}
			} else {
				value := c.Get("key")
				if value != tt.want {
					t.Error("is not want value")
				}
			}
		})
	}
}
