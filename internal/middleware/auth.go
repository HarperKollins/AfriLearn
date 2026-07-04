package middleware

import (
	"database/sql"
	"net/http"
	"sync"
	"time"

	"github.com/afrilearn/curriculum-api/internal/database"
	"github.com/afrilearn/curriculum-api/internal/models"
	"github.com/gin-gonic/gin"
)

// CachedKey stores in-memory metadata for fast key lookup
type CachedKey struct {
	DeveloperName string
	Email         string
	Tier          string
	IsActive      bool
}

var (
	keyCache   = make(map[string]CachedKey)
	cacheMutex sync.RWMutex
)

// UpdateKeyCache adds/updates a key in the in-memory cache
func UpdateKeyCache(apiKey, devName, email, tier string, isActive bool) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	keyCache[apiKey] = CachedKey{
		DeveloperName: devName,
		Email:         email,
		Tier:          tier,
		IsActive:      isActive,
	}
}

// ── Phase 3: Batched request counter ────────────────────────────────────────
// A buffered channel replaces the old unbounded `go func()` per request.
// A single background goroutine drains the channel in batches, writing to DB
// at most once per second — protecting against thundering-herd under load.

const counterBufSize = 4096

var counterCh = make(chan string, counterBufSize)

// InitRateLimiter starts the background goroutine that drains counterCh.
// Call this once at server startup.
func InitRateLimiter() {
	go func() {
		batch := make(map[string]int)
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case key := <-counterCh:
				batch[key]++
			case <-ticker.C:
				if len(batch) == 0 {
					continue
				}
				// Write the whole batch in a single transaction
				tx, err := database.DB.Begin()
				if err != nil {
					batch = make(map[string]int) // drop on error, metric loss acceptable
					continue
				}
				for k, count := range batch {
					_, _ = tx.Exec(
						`UPDATE api_keys SET requests_count = requests_count + $1, updated_at = NOW() WHERE api_key = $2`,
						count, k,
					)
				}
				_ = tx.Commit()
				batch = make(map[string]int)
			}
		}
	}()
}

// APIKeyAuth middleware validates X-API-Key header or api_key query param.
func APIKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			apiKey = c.Query("api_key")
		}

		// No key — allow public/free tier access
		if apiKey == "" {
			c.Header("X-RateLimit-Tier", "public")
			c.Header("X-RateLimit-Limit", "60")
			c.Next()
			return
		}

		// Check fast in-memory cache first
		cacheMutex.RLock()
		cached, found := keyCache[apiKey]
		cacheMutex.RUnlock()

		if !found {
			// Cache miss — query PostgreSQL
			var devName, email, tier string
			var isActive bool
			err := database.DB.QueryRow(`
				SELECT developer_name, email, tier, is_active
				FROM api_keys WHERE api_key = $1
			`, apiKey).Scan(&devName, &email, &tier, &isActive)

			if err == sql.ErrNoRows {
				c.JSON(http.StatusUnauthorized, models.APIResponse{
					Success: false,
					Message: "Invalid API key provided",
				})
				c.Abort()
				return
			}
			if err != nil {
				c.JSON(http.StatusInternalServerError, models.APIResponse{
					Success: false,
					Message: "Failed to authenticate API key",
				})
				c.Abort()
				return
			}

			cached = CachedKey{
				DeveloperName: devName,
				Email:         email,
				Tier:          tier,
				IsActive:      isActive,
			}

			cacheMutex.Lock()
			keyCache[apiKey] = cached
			cacheMutex.Unlock()
		}

		if !cached.IsActive {
			c.JSON(http.StatusForbidden, models.APIResponse{
				Success: false,
				Message: "API key is deactivated",
			})
			c.Abort()
			return
		}

		// Phase 3: non-blocking enqueue into the batch counter channel.
		// If the channel is full (4096 items queued) we simply drop this increment
		// rather than block the request path — counter accuracy is best-effort.
		select {
		case counterCh <- apiKey:
		default:
			// Channel full — skip increment, keep serving traffic
		}

		// Set rate limit headers
		limit := "1000"
		if cached.Tier == "pro" {
			limit = "50000"
		} else if cached.Tier == "enterprise" {
			limit = "unlimited"
		}

		c.Header("X-RateLimit-Tier", cached.Tier)
		c.Header("X-RateLimit-Limit", limit)
		c.Set("developer_email", cached.Email)
		c.Set("developer_tier", cached.Tier)

		c.Next()
	}
}

