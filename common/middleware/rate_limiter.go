package middleware

import (
	"enterprise-microservice-system/common/errors"
	"enterprise-microservice-system/common/response"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter implements token bucket rate limiting
type RateLimiter struct {
	limiters sync.Map // map[string]*rate.Limiter
	limit    rate.Limit
	burst    int
}

// NewRateLimiter creates a new rate limiter
// limit is requests per second, burst is the maximum burst size
func NewRateLimiter(requestsPerSecond int, burst int) *RateLimiter {
	return &RateLimiter{
		limit: rate.Limit(requestsPerSecond),
		burst: burst,
	}
}

// getLimiter retrieves or creates a rate limiter for a given key (typically IP address)
func (rl *RateLimiter) getLimiter(key string) *rate.Limiter {
	limiter, exists := rl.limiters.Load(key)
	if !exists {
		limiter = rate.NewLimiter(rl.limit, rl.burst)
		rl.limiters.Store(key, limiter)
	}
	return limiter.(*rate.Limiter)
}

// Middleware returns a Gin middleware for rate limiting
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Use client IP as the key for rate limiting
		clientIP := c.ClientIP()
		limiter := rl.getLimiter(clientIP)

		// Check if request is allowed
		if !limiter.Allow() {
			response.Error(c, errors.NewRateLimit())
			c.Abort()
			return
		}

		c.Next()
	}
}

// CleanupStaleEntries removes limiters that haven't been used recently
// This should be called periodically to prevent memory leaks
func (rl *RateLimiter) CleanupStaleEntries() {
	rl.limiters.Range(func(key, value interface{}) bool {
		limiter := value.(*rate.Limiter)
		// If the limiter's bucket is full, it means it hasn't been used recently
		if limiter.Tokens() == float64(rl.burst) {
			rl.limiters.Delete(key)
		}
		return true
	})
}
