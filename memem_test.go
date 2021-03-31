package memem

import "testing"

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
