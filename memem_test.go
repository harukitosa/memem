package memem

import (
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	c := NewCache[int]()
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
	c := NewCache[int]()
	c.Set("id", 12)
	if c.Get("id") != 12 {
		t.Error("error not match")
	}
	c.Set("hoge", 100)
	if c.Get("hoge") != 100 {
		t.Error("error")
	}
	c.Clear()
	if c.Get("hoge") == 0 {
		t.Error("error")
	}
}

func TestArrayCache(t *testing.T) {
	c := NewCache[[]string]()
	slice := []string{"Golang", "Java"}
	c.Set("key", slice)
	value := c.Get("key")
	if value[0] != "Golang" {
		t.Error("error not match")
	}
}

func TestGetDataIsNoneCache(t *testing.T) {
	c := NewCache[string]()
	value := c.Get("key")
	if value != "" {
		t.Error("not nil")
	}
}

func TestCallbackCache(t *testing.T) {
	c := NewCacheWithCallback(func() interface{} {
		return "callback result"
	})
	value := c.Get("callback result")
	if value == nil {
		t.Error("is not callback value")
	}
	if value != "callback result" {
		t.Error("is not callback value")
	}
}

func TestWithClearTimeNonClear(t *testing.T) {
	c := NewCacheWithClearTime[string](time.Second)
	c.Set("key", "same value")
	value := c.Get("key")
	if value != "same value" {
		t.Error("data is null")
	}
}

func TestWithClearTimeClear(t *testing.T) {
	c := NewCacheWithClearTime[string](time.Second)
	c.Set("key", "same value")
	time.Sleep(2 * time.Second)
	value := c.Get("key")
	if value != "" {
		t.Error("data is not null")
	}
}

type TestCase struct {
	f       func() interface{}
	t       time.Duration
	isClear bool
}

func TestWithClearTimeAndCallback(t *testing.T) {
	f := func() interface{} {
		return "callback value"
	}
	cases := []struct {
		name string
		in   TestCase
		want string
	}{
		{
			"期限切れしていない場合キャッシュが返却される",
			TestCase{
				f,
				time.Second,
				false,
			},
			"value",
		},
		{
			"期限が過ぎている場合callbackの値が返却される",
			TestCase{
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
			c := NewCacheWithCallbackAndClearTime(tt.in.f, tt.in.t)
			c.Set("key", "value")
			if tt.in.isClear {
				time.Sleep(2 * time.Second)
				value, ok := c.Get("key").(string)
				if !ok {
					t.Error("cast miss")
				}
				if value != tt.want {
					t.Error("is not want value")
				}
			} else {
				value, ok := c.Get("key").(string)
				if !ok {
					t.Error("cast miss")
				}
				if value != tt.want {
					t.Error("is not want value")
				}
			}
		})
	}
}
