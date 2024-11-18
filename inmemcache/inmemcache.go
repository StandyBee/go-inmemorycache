package inmemcache

import (
	"errors"
	"sync"
	"time"
)

type Cache interface {
	Set(key string, value interface{}, cacheTTL time.Duration) error
	Get(key string) (interface{}, error)
	Delete(key string) error
}

type InMemCache struct {
	cache map[string]interface{}
	mu    sync.Mutex
	ttl   map[string]time.Time
}

func NewInMemCache() *InMemCache {
	return &InMemCache{
		cache: make(map[string]interface{}),
		mu:    sync.Mutex{},
		ttl:   make(map[string]time.Time),
	}
}

func (c *InMemCache) Set(key string, value interface{}, cacheTime time.Duration) error {
	if key == "" {
		return errors.New("key cannot be empty")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = value
	c.ttl[key] = time.Now().Add(cacheTime)
	return nil
}

func (c *InMemCache) Get(key string) (interface{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	val, ok := c.cache[key]

	if !ok {
		return nil, errors.New("key not found")
	}

	deadline, ok := c.ttl[key]
	if !ok {
		return nil, errors.New("key not found")
	}
	if time.Now().After(deadline) {
		delete(c.cache, key)
		delete(c.ttl, key)
		return nil, errors.New("key expired")
	}

	return val, nil
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
	delete(c.ttl, key)
	return nil
}

func (c *InMemCache) Cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for key, deadline := range c.ttl {
		if time.Now().After(deadline) {
			delete(c.cache, key)
			delete(c.ttl, key)
		}
	}
}

func (c *InMemCache) StartCleanup(interval time.Duration) {
	go func() {
		for {
			time.Sleep(interval)
			c.Cleanup()
		}
	}()
}
