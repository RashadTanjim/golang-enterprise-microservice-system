package auth

import (
	"sync"
	"testing"
	"time"
)

func TestCachedTokenProvider_GetToken(t *testing.T) {
	config := Config{
		Secret:   "test-secret",
		Issuer:   "test-issuer",
		Audience: "test-audience",
		TokenTTL: 1 * time.Hour,
	}
	subject := "test-service"
	roles := []string{"service"}

	provider := NewCachedTokenProvider(config, subject, roles, 5*time.Minute)

	// First call should generate a new token
	token1, err := provider.GetToken()
	if err != nil {
		t.Fatalf("Failed to get token: %v", err)
	}
	if token1 == "" {
		t.Fatal("Expected non-empty token")
	}

	// Second call should return the cached token
	token2, err := provider.GetToken()
	if err != nil {
		t.Fatalf("Failed to get cached token: %v", err)
	}
	if token2 != token1 {
		t.Fatal("Expected cached token to be the same")
	}
}

func TestCachedTokenProvider_TokenExpiry(t *testing.T) {
	config := Config{
		Secret:   "test-secret",
		Issuer:   "test-issuer",
		Audience: "test-audience",
		TokenTTL: 2 * time.Second, // Short TTL for testing
	}
	subject := "test-service"
	roles := []string{"service"}

	// Set renewal threshold to 1 second before expiry
	provider := NewCachedTokenProvider(config, subject, roles, 1*time.Second)

	// Get initial token
	token1, err := provider.GetToken()
	if err != nil {
		t.Fatalf("Failed to get token: %v", err)
	}

	// Wait for renewal threshold to pass
	time.Sleep(1200 * time.Millisecond)

	// Should get a new token now
	token2, err := provider.GetToken()
	if err != nil {
		t.Fatalf("Failed to get renewed token: %v", err)
	}

	// Tokens should be different since we crossed the renewal threshold
	if token2 == token1 {
		t.Fatal("Expected token to be renewed after crossing renewal threshold")
	}
}

func TestCachedTokenProvider_InvalidateToken(t *testing.T) {
	config := Config{
		Secret:   "test-secret",
		Issuer:   "test-issuer",
		Audience: "test-audience",
		TokenTTL: 1 * time.Hour,
	}
	subject := "test-service"
	roles := []string{"service"}

	provider := NewCachedTokenProvider(config, subject, roles, 5*time.Minute)

	// Get initial token
	_, err := provider.GetToken()
	if err != nil {
		t.Fatalf("Failed to get token: %v", err)
	}

	// Invalidate the token
	provider.InvalidateToken()

	// Verify token is invalidated
	provider.mu.RLock()
	if provider.currentToken != "" {
		t.Error("Expected currentToken to be empty after invalidation")
	}
	provider.mu.RUnlock()

	// Next call should generate a new token
	token2, err := provider.GetToken()
	if err != nil {
		t.Fatalf("Failed to get new token after invalidation: %v", err)
	}

	// Token should not be empty
	if token2 == "" {
		t.Fatal("Expected non-empty token after invalidation")
	}
}

func TestCachedTokenProvider_ConcurrentAccess(t *testing.T) {
	config := Config{
		Secret:   "test-secret",
		Issuer:   "test-issuer",
		Audience: "test-audience",
		TokenTTL: 1 * time.Hour,
	}
	subject := "test-service"
	roles := []string{"service"}

	provider := NewCachedTokenProvider(config, subject, roles, 5*time.Minute)

	// Test concurrent access
	var wg sync.WaitGroup
	numGoroutines := 100
	tokens := make([]string, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			token, err := provider.GetToken()
			if err != nil {
				t.Errorf("Failed to get token in goroutine %d: %v", index, err)
				return
			}
			tokens[index] = token
		}(i)
	}

	wg.Wait()

	// All tokens should be the same (cached)
	firstToken := tokens[0]
	for i, token := range tokens {
		if token != firstToken {
			t.Errorf("Token mismatch at index %d: expected %s, got %s", i, firstToken, token)
		}
	}
}

func TestCachedTokenProvider_DefaultRenewalThreshold(t *testing.T) {
	config := Config{
		Secret:   "test-secret",
		Issuer:   "test-issuer",
		Audience: "test-audience",
		TokenTTL: 1 * time.Hour,
	}
	subject := "test-service"
	roles := []string{"service"}

	// Pass 0 for renewal threshold to use default
	provider := NewCachedTokenProvider(config, subject, roles, 0)

	if provider.renewalThreshold != 5*time.Minute {
		t.Fatalf("Expected default renewal threshold of 5 minutes, got %v", provider.renewalThreshold)
	}
}
