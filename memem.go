package memem

import "sync"

type Cache struct {
	Data map[string]interface{}
	sync.Mutex
}

func NewCache() *Cache {
	m := make(map[string]interface{})
	c := &Cache{
		Data: m,
	}
	return c
}

func (c *Cache) Append(key string, value interface{}) {
	c.Lock()
	c.Data[key] = value
	c.Unlock()
}

func (c *Cache) Get(key string) interface{} {
	return c.Data[key]
}
