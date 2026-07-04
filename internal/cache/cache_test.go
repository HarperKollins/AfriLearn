package cache

import (
	"testing"
	"time"
)

func TestMemoryCache_SetAndGet(t *testing.T) {
	c := InitCache(1 * time.Hour)
	c.Clear()

	key := "test:key"
	val := "hello_afrilearn"

	c.Set(key, val, 100*time.Millisecond)

	retrieved, found := c.Get(key)
	if !found {
		t.Fatalf("Expected key '%s' to be found in cache", key)
	}
	if retrieved.(string) != val {
		t.Fatalf("Expected value '%s', got '%s'", val, retrieved)
	}
}

func TestMemoryCache_Expiration(t *testing.T) {
	c := InitCache(1 * time.Hour)
	c.Clear()

	key := "test:exp"
	val := "expiring_data"

	c.Set(key, val, 50*time.Millisecond)
	time.Sleep(60 * time.Millisecond)

	_, found := c.Get(key)
	if found {
		t.Fatalf("Expected key '%s' to be expired and not found", key)
	}
}

func TestMemoryCache_DeleteAndClear(t *testing.T) {
	c := InitCache(1 * time.Hour)
	c.Clear()

	c.Set("k1", "v1", 0)
	c.Set("k2", "v2", 0)

	c.Delete("k1")
	_, found := c.Get("k1")
	if found {
		t.Fatalf("Expected 'k1' to be deleted")
	}

	c.Clear()
	stats := c.Stats()
	if stats.TotalItems != 0 {
		t.Fatalf("Expected TotalItems = 0 after Clear(), got %d", stats.TotalItems)
	}
}

func TestMemoryCache_Stats(t *testing.T) {
	c := InitCache(1 * time.Hour)
	c.Clear()

	c.Set("hit_key", "data", 0)
	c.Get("hit_key")  // Hit
	c.Get("miss_key") // Miss

	stats := c.Stats()
	if stats.Hits != 1 {
		t.Errorf("Expected 1 hit, got %d", stats.Hits)
	}
	if stats.Misses != 1 {
		t.Errorf("Expected 1 miss, got %d", stats.Misses)
	}
}
