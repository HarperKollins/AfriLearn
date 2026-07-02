package middleware

import (
	"database/sql"
	"net/http"
	"sync"

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

// APIKeyAuth middleware validates X-API-Key header or api_key query param
func APIKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			apiKey = c.Query("api_key")
		}

		// If no key provided, allow access under public/free tier with rate limit
		if apiKey == "" {
			c.Header("X-RateLimit-Tier", "public")
			c.Header("X-RateLimit-Limit", "60")
			c.Next()
			return
		}

		// Check fast in-memory cache
		cacheMutex.RLock()
		cached, found := keyCache[apiKey]
		cacheMutex.RUnlock()

		if !found {
			// Query PostgreSQL
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

			// Store in cache
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

		// Increment request counter in background asynchronously
		go func(k string) {
			_, _ = database.DB.Exec(`UPDATE api_keys SET requests_count = requests_count + 1, updated_at = NOW() WHERE api_key = $1`, k)
		}(apiKey)

		// Set rate limit headers according to Tier
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
