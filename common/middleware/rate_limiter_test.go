package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRateLimiter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a rate limiter that allows 2 requests per second with burst of 2
	rateLimiter := NewRateLimiter(2, 2)

	// Create a test router
	router := gin.New()
	router.Use(rateLimiter.Middleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Make 3 requests rapidly (should allow 2, reject 1)
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.RemoteAddr = "192.168.1.1:1234" // Same IP for all requests
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if i < 2 {
			// First 2 requests should succeed
			if w.Code != http.StatusOK {
				t.Errorf("Request %d: expected status %d, got %d", i+1, http.StatusOK, w.Code)
			}
		} else {
			// Third request should be rate limited
			if w.Code != http.StatusTooManyRequests {
				t.Errorf("Request %d: expected status %d, got %d", i+1, http.StatusTooManyRequests, w.Code)
			}
		}
	}
}

func TestRateLimiterDifferentIPs(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rateLimiter := NewRateLimiter(1, 1)

	router := gin.New()
	router.Use(rateLimiter.Middleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Requests from different IPs should have separate limiters
	ips := []string{"192.168.1.1:1234", "192.168.1.2:1234", "192.168.1.3:1234"}

	for _, ip := range ips {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.RemoteAddr = ip
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Request from IP %s: expected status %d, got %d", ip, http.StatusOK, w.Code)
		}
	}
}
