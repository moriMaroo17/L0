package main

import (
	"fmt"
	"sync"
)

type Cache struct {
	locker      sync.RWMutex
	memoryCache map[string]*Data
}

func NewCache() Cache {
	return Cache{memoryCache: make(map[string]*Data)}
}

func (c *Cache) Put(key string, value Data) {
	c.locker.Lock()
	defer c.locker.Unlock()
	c.memoryCache[key] = &value
}

func (c *Cache) Get(key string) (Data, error) {
	c.locker.RLock()
	defer c.locker.RUnlock()
	data, ok := c.memoryCache[key]
	if !ok {
		return Data{}, fmt.Errorf("no such cache key %s", key)
	}
	return *data, nil
}

func (c *Cache) CheckEmpty() bool {
	return len(c.memoryCache) == 0
}

func (c *Cache) Restore(p DBExecutor, restoreCh <-chan Data) {
	go p.Backup()
	for data := range restoreCh {
		go c.Put(data.Payment.Transaction, data)
	}
}
