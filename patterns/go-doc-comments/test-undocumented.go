package cache

import (
	"errors"
	"sync"
	"time"
)

type Cache struct {
	mu    sync.RWMutex
	items map[string]*Item
	ttl   time.Duration
}

type Item struct {
	Value     interface{}
	ExpiresAt time.Time
}

func NewCache(ttl time.Duration) *Cache {
	return &Cache{
		items: make(map[string]*Item),
		ttl:   ttl,
	}
}

func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = &Item{
		Value:     value,
		ExpiresAt: time.Now().Add(c.ttl),
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return nil, false
	}

	if time.Now().After(item.ExpiresAt) {
		return nil, false
	}

	return item.Value, true
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}

func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]*Item)
}

func (c *Cache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.items)
}

func (c *Cache) Keys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := make([]string, 0, len(c.items))
	for k := range c.items {
		keys = append(keys, k)
	}
	return keys
}

func (c *Cache) Cleanup() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	deleted := 0

	for key, item := range c.items {
		if now.After(item.ExpiresAt) {
			delete(c.items, key)
			deleted++
		}
	}

	return deleted
}

const (
	DefaultTTL = 5 * time.Minute
	MaxTTL     = 24 * time.Hour
)

var ErrKeyNotFound = errors.New("key not found")
var ErrExpired = errors.New("item expired")
