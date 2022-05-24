package main

import "fmt"

type Cache struct {
	memoryCache map[string]*Data
}

func NewCache() Cache {
	return Cache{memoryCache: make(map[string]*Data)}
}

func (c *Cache) Put(key string, value Data) {
	c.memoryCache[key] = &value
}

func (c *Cache) Get(key string) (Data, error) {
	data, ok := c.memoryCache[key]
	if !ok {
		return Data{}, fmt.Errorf("no such cache key %s", key)
	}
	return *data, nil
}
