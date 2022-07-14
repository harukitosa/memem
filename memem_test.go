package memem

import (
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	c := NewCache()
	c.Set("id", 12)
	if c.Get("id") != 12 {
		t.Error("error not match")
	}
	c.Set("hoge", "doremifaso")
	if c.Get("hoge") != "doremifaso" {
		t.Error("error")
	}
}

func TestArrayCache(t *testing.T) {
	c := NewCache()
	// var arr[2] string = [2]string {"Golang", "Java"}
	slice := []string{"Golang", "Java"}
	c.Set("key", slice)
	value := c.Get("key")
	v, ok := value.([]string)
	if !ok {
		t.Error("error cast")
	}
	if v[0] != "Golang" {
		t.Error("error not match")
	}
}

func TestGetDataIsNoneCache(t *testing.T) {
	c := NewCache()
	value := c.Get("key")
	if value != nil {
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

func TestWithTimer(t *testing.T) {
	c := NewCacheWithCallback(func() interface{} {
		return "callback result"
	})
	c.Set("key", "same value")
	value := c.GetOrClearIfOverTheTimeLimit("key", 2*time.Second)
	if value != "same value" {
		t.Error("data is null")
	}
}

func TestWithTimerClear(t *testing.T) {
	c := NewCache()
	c.Set("key", "same value")
	time.Sleep(2 * time.Second)
	value := c.GetOrClearIfOverTheTimeLimit("key", 1*time.Second)
	if value != nil {
		t.Error("data is not null")
	}
}

func TestWithTimerClearCallback(t *testing.T) {
	c := NewCacheWithCallback(func() interface{} {
		return "callbackvalue"
	})
	c.Set("key", "same value")
	time.Sleep(2 * time.Second)
	value := c.GetOrClearIfOverTheTimeLimit("key", 1*time.Second)
	if value != "callbackvalue" {
		t.Errorf("data is not callbackvalue")
	}
}
