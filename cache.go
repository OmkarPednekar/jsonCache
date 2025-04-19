package jsonCache

import (
	"time"
)

type Cache struct {
	store map[string]CacheElement
}

type CacheElement struct {
	value  []byte
	expiry time.Time
}

func New() *Cache {
	//inialize a new cache engine
	cache := &Cache{
		store: make(map[string]CacheElement),
	}
	return cache
}

func (c *Cache) Set(key string, value []byte, ttl int) {
	expiration := time.Now().Add(time.Duration(ttl) * time.Millisecond)
	var toStore = &CacheElement{
		value:  value,
		expiry: expiration,
	}
	c.store[key] = *toStore
}

func (c *Cache) Get(key string) []byte {
	entry, found := c.store[key]
	if !found {
		return nil
	}
	if time.Now().After(entry.expiry) {
		delete(c.store, key)
		return nil
	}
	return entry.value
}
