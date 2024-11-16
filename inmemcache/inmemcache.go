package inmemcache

import (
	"errors"
	"sync"
)

type Cache interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, bool)
	Delete(key string) error
}

type InMemCache struct {
	cache map[string]interface{}
	mu    sync.Mutex
}

func NewInMemCache() *InMemCache {
	return &InMemCache{
		cache: make(map[string]interface{}),
		mu:    sync.Mutex{},
	}
}

func (c *InMemCache) Set(key string, value interface{}) error {
	if key == "" {
		return errors.New("key cannot be empty")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = value
	return nil
}

func (c *InMemCache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	value, ok := c.cache[key]
	return value, ok
}

func (c *InMemCache) Delete(key string) error {
	if key == "" {
		return errors.New("key cannot be empty")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	_, ok := c.cache[key]
	if !ok {
		return errors.New("key does not exist")
	}

	delete(c.cache, key)
	return nil
}
