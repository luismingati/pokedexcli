package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	cacheEntry map[string]cacheEntry
	mu         sync.Mutex
	interval   time.Duration
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		cacheEntry: make(map[string]cacheEntry),
		interval:   interval,
	}
	go c.cleanUp()
	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cacheEntry[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, found := c.cacheEntry[key]

	if !found {
		return nil, false
	}

	return entry.val, true
}

func (c *Cache) cleanUp() {
	for {
		time.Sleep(c.interval)
		c.mu.Lock()
		now := time.Now()
		for key, value := range c.cacheEntry {
			if now.Sub(value.createdAt) > c.interval {
				delete(c.cacheEntry, key)
			}
		}
		c.mu.Unlock()
	}
}

func (c *Cache) List() {
	for key := range c.cacheEntry {
		println(key)
	}
}
