package poke

import (
	"sync"
	"time"
)

// Cache struct to hold a map[string]cacheEntry and a mutex for concurrency control
type Cache struct {
	entries  map[string]cacheEntry
	mutex    sync.Mutex
	duration time.Duration
}

// cacheEntry struct to hold cached data
type cacheEntry struct {
	createdAt time.Time
	data      []byte
}

// NewCache initializes and returns a new Cache
func NewCache(duration time.Duration) *Cache {
	cache := &Cache{
		entries:  make(map[string]cacheEntry),
		duration: duration,
	}

	go cache.reapLoop()

	return cache
}

// Get retrieves data from the cache if it exists and is not expired
func (c *Cache) Get(key string) ([]byte, bool) {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	entry, exists := c.entries[key]
	if !exists {
		return nil, false
	}

	if time.Since(entry.createdAt) > c.duration {
		delete(c.entries, key)
		return nil, false
	}

	return entry.data, true
}

// Add adds data to the cache with the current timestamp
func (c *Cache) Add(key string, data []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.entries[key] = cacheEntry{
		createdAt: time.Now(),
		data:      data,
	}
}

// cache.reapLoop() method that is called when the cache is created (by the NewCache function). Each time an interval (the time.Duration passed to NewCache) passes it should remove any entries that are older than the interval. This makes sure that the cache doesn't grow too large over time. For example, if the interval is 5 seconds, and an entry was added 7 seconds ago, that entry should be removed.
func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.duration)
	defer ticker.Stop()

	for range ticker.C {
		c.mutex.Lock()
		for key, entry := range c.entries {
			if time.Since(entry.createdAt) > c.duration {
				delete(c.entries, key)
			}
		}
		c.mutex.Unlock()
	}
}
