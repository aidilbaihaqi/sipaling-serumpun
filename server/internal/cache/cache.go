package cache

import (
	"sync"
	"time"
)

type Item struct {
	Data      []byte
	ExpiresAt time.Time
}

type Cache struct {
	mu   sync.RWMutex
	ttl  time.Duration
	data map[string]Item
}

func New(ttl time.Duration) *Cache {
	return &Cache{ttl: ttl, data: map[string]Item{}}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	it, ok := c.data[key]
	c.mu.RUnlock()
	if !ok || time.Now().After(it.ExpiresAt) {
		return nil, false
	}
	return it.Data, true
}

func (c *Cache) Set(key string, b []byte) {
	c.mu.Lock()
	c.data[key] = Item{Data: b, ExpiresAt: time.Now().Add(c.ttl)}
	c.mu.Unlock()
}
