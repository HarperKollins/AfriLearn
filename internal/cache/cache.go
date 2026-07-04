package cache

import (
	"strings"
	"sync"
	"time"
)

type cacheItem struct {
	value     interface{}
	expiresAt time.Time
	cachedAt  time.Time
}

// prefixStats tracks hit/miss counts per key prefix
type prefixStats struct {
	hits   int64
	misses int64
}

type MemoryCache struct {
	items       map[string]cacheItem
	mu          sync.RWMutex
	defaultTTL  time.Duration
	hits        int64
	misses      int64
	prefixStats map[string]*prefixStats // keyed by prefix (e.g., "curr", "prompt", "embeddings")
}

var (
	GlobalCache *MemoryCache
	once        sync.Once
)

// InitCache initializes the global cache
func InitCache(defaultTTL time.Duration) *MemoryCache {
	once.Do(func() {
		GlobalCache = &MemoryCache{
			items:       make(map[string]cacheItem),
			defaultTTL:  defaultTTL,
			prefixStats: make(map[string]*prefixStats),
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

// keyPrefix extracts the first segment of a colon-separated cache key (e.g., "curr" from "curr:waec:math")
func keyPrefix(key string) string {
	if idx := strings.Index(key, ":"); idx > 0 {
		return key[:idx]
	}
	return key
}

// Get retrieves an item from cache, tracking per-prefix stats
func (c *MemoryCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	item, found := c.items[key]
	c.mu.RUnlock()

	prefix := keyPrefix(key)

	if !found {
		c.mu.Lock()
		c.misses++
		if c.prefixStats[prefix] == nil {
			c.prefixStats[prefix] = &prefixStats{}
		}
		c.prefixStats[prefix].misses++
		c.mu.Unlock()
		return nil, false
	}

	if time.Now().After(item.expiresAt) {
		c.mu.Lock()
		delete(c.items, key)
		c.misses++
		if c.prefixStats[prefix] == nil {
			c.prefixStats[prefix] = &prefixStats{}
		}
		c.prefixStats[prefix].misses++
		c.mu.Unlock()
		return nil, false
	}

	c.mu.Lock()
	c.hits++
	if c.prefixStats[prefix] == nil {
		c.prefixStats[prefix] = &prefixStats{}
	}
	c.prefixStats[prefix].hits++
	c.mu.Unlock()
	return item.value, true
}

// Set stores an item in cache with the given TTL (0 = use defaultTTL)
func (c *MemoryCache) Set(key string, value interface{}, ttl time.Duration) {
	if ttl == 0 {
		ttl = c.defaultTTL
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = cacheItem{
		value:     value,
		expiresAt: time.Now().Add(ttl),
		cachedAt:  time.Now(),
	}
}

// Delete removes a key from cache
func (c *MemoryCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Clear flushes all cached items and resets all counters
func (c *MemoryCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]cacheItem)
	c.hits = 0
	c.misses = 0
	c.prefixStats = make(map[string]*prefixStats)
}

// PrefixStats represents hit/miss data for one cache key prefix
type PrefixStats struct {
	Prefix   string  `json:"prefix"`
	Hits     int64   `json:"hits"`
	Misses   int64   `json:"misses"`
	HitRatio float64 `json:"hit_ratio"`
}

// CacheStats returns current aggregate cache metrics
type CacheStats struct {
	TotalItems  int           `json:"total_items"`
	Hits        int64         `json:"hits"`
	Misses      int64         `json:"misses"`
	HitRatio    float64       `json:"hit_ratio"`
	ByPrefix    []PrefixStats `json:"by_prefix"`
}

// Stats returns current cache metrics including per-prefix breakdown
func (c *MemoryCache) Stats() CacheStats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	totalReqs := c.hits + c.misses
	ratio := 0.0
	if totalReqs > 0 {
		ratio = float64(c.hits) / float64(totalReqs) * 100.0
	}

	// Count live (non-expired) items
	now := time.Now()
	liveCount := 0
	for _, item := range c.items {
		if now.Before(item.expiresAt) {
			liveCount++
		}
	}

	// Build per-prefix breakdown
	byPrefix := make([]PrefixStats, 0, len(c.prefixStats))
	for prefix, ps := range c.prefixStats {
		total := ps.hits + ps.misses
		r := 0.0
		if total > 0 {
			r = float64(ps.hits) / float64(total) * 100.0
		}
		byPrefix = append(byPrefix, PrefixStats{
			Prefix:   prefix,
			Hits:     ps.hits,
			Misses:   ps.misses,
			HitRatio: r,
		})
	}

	return CacheStats{
		TotalItems: liveCount,
		Hits:       c.hits,
		Misses:     c.misses,
		HitRatio:   ratio,
		ByPrefix:   byPrefix,
	}
}
