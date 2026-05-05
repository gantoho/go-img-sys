package cache

import (
	"sync"
	"time"
)

type CacheEntry struct {
	Value     interface{}
	ExpiresAt time.Time
}

type Cache struct {
	mu     sync.RWMutex
	items  map[string]CacheEntry
	stopCh chan struct{}
}

func NewCache() *Cache {
	cache := &Cache{
		items:  make(map[string]CacheEntry),
		stopCh: make(chan struct{}),
	}
	go cache.cleanup()
	return cache
}

func (c *Cache) Close() {
	close(c.stopCh)
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = CacheEntry{
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, exists := c.items[key]
	if !exists {
		return nil, false
	}
	if time.Now().After(entry.ExpiresAt) {
		return nil, false
	}
	return entry.Value, true
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]CacheEntry)
}

func (c *Cache) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-c.stopCh:
			return
		case <-ticker.C:
			c.mu.Lock()
			now := time.Now()
			for key, entry := range c.items {
				if now.After(entry.ExpiresAt) {
					delete(c.items, key)
				}
			}
			c.mu.Unlock()
		}
	}
}

func (c *Cache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}
