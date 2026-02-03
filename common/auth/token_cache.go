package auth

import (
	"sync"
	"time"
)

// CachedTokenProvider provides JWT tokens with caching and automatic renewal.
// It generates a new token only when the current one is expired or about to expire.
type CachedTokenProvider struct {
	mu                sync.RWMutex
	config            Config
	subject           string
	roles             []string
	currentToken      string
	expiresAt         time.Time
	renewalThreshold  time.Duration
}

// NewCachedTokenProvider creates a new cached token provider.
// renewalThreshold determines how long before expiry to renew the token.
// A good default is 5 minutes before expiry.
func NewCachedTokenProvider(config Config, subject string, roles []string, renewalThreshold time.Duration) *CachedTokenProvider {
	if renewalThreshold == 0 {
		// Default to renewing 5 minutes before expiry
		renewalThreshold = 5 * time.Minute
	}

	return &CachedTokenProvider{
		config:           config,
		subject:          subject,
		roles:            roles,
		renewalThreshold: renewalThreshold,
	}
}

// GetToken returns a valid JWT token, generating a new one if necessary.
func (p *CachedTokenProvider) GetToken() (string, error) {
	// First check with read lock
	p.mu.RLock()
	if p.isTokenValid() {
		token := p.currentToken
		p.mu.RUnlock()
		return token, nil
	}
	p.mu.RUnlock()

	// Token is expired or about to expire, acquire write lock to renew
	p.mu.Lock()
	defer p.mu.Unlock()

	// Double-check after acquiring write lock (another goroutine might have renewed it)
	if p.isTokenValid() {
		return p.currentToken, nil
	}

	// Generate new token
	token, err := GenerateToken(p.config, p.subject, p.roles)
	if err != nil {
		return "", err
	}

	// Update cache
	p.currentToken = token
	p.expiresAt = time.Now().UTC().Add(p.config.TokenTTL)

	return token, nil
}

// isTokenValid checks if the current token is still valid.
// Must be called while holding at least a read lock.
func (p *CachedTokenProvider) isTokenValid() bool {
	if p.currentToken == "" {
		return false
	}

	// Renew if we're within the renewal threshold of expiry
	renewalTime := p.expiresAt.Add(-p.renewalThreshold)
	return time.Now().UTC().Before(renewalTime)
}

// InvalidateToken forces the next GetToken call to generate a new token.
// This can be useful for testing or when you know the token has been revoked.
func (p *CachedTokenProvider) InvalidateToken() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.currentToken = ""
	p.expiresAt = time.Time{}
}
