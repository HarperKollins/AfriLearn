package cache

import (
	"sync"
	"time"
)

type cacheItem struct {
	value     interface{}
	expiresAt time.Time
}

type MemoryCache struct {
	items      map[string]cacheItem
	mu         sync.RWMutex
	defaultTTL time.Duration
	hits       int64
	misses     int64
}

var (
	GlobalCache *MemoryCache
	once        sync.Once
)

// InitCache initializes the global cache
func InitCache(defaultTTL time.Duration) *MemoryCache {
	once.Do(func() {
		GlobalCache = &MemoryCache{
			items:      make(map[string]cacheItem),
			defaultTTL: defaultTTL,
		}
	})
	return GlobalCache
}

// GetCache returns the global cache instance
func GetCache() *MemoryCache {
	if GlobalCache == nil {
		return InitCache(1 * time.Hour)
	}
	return GlobalCache
}

// Get retrieves an item from cache
func (c *MemoryCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	item, found := c.items[key]
	c.mu.RUnlock()

	if !found {
		c.mu.Lock()
		c.misses++
		c.mu.Unlock()
		return nil, false
	}

	if time.Now().After(item.expiresAt) {
		c.mu.Lock()
		delete(c.items, key)
		c.misses++
		c.mu.Unlock()
		return nil, false
	}

	c.mu.Lock()
	c.hits++
	c.mu.Unlock()
	return item.value, true
}

// Set stores an item in cache
func (c *MemoryCache) Set(key string, value interface{}, ttl time.Duration) {
	if ttl == 0 {
		ttl = c.defaultTTL
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = cacheItem{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}
}

// Delete removes a key from cache
func (c *MemoryCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Clear flushes all cached items
func (c *MemoryCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]cacheItem)
}

// Stats returns current cache metrics
type CacheStats struct {
	TotalItems int     `json:"total_items"`
	Hits       int64   `json:"hits"`
	Misses     int64   `json:"misses"`
	HitRatio   float64 `json:"hit_ratio"`
}

func (c *MemoryCache) Stats() CacheStats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	totalReqs := c.hits + c.misses
	ratio := 0.0
	if totalReqs > 0 {
		ratio = float64(c.hits) / float64(totalReqs) * 100.0
	}

	return CacheStats{
		TotalItems: len(c.items),
		Hits:       c.hits,
		Misses:     c.misses,
		HitRatio:   ratio,
	}
}
