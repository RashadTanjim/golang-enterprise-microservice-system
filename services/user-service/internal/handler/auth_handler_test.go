package handler

import (
	"bytes"
	"encoding/json"
	"enterprise-microservice-system/common/auth"
	"enterprise-microservice-system/common/logger"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

type tokenResponse struct {
	Success bool `json:"success"`
	Data    struct {
		AccessToken string    `json:"access_token"`
		TokenType   string    `json:"token_type"`
		ExpiresAt   time.Time `json:"expires_at"`
		Roles       []string  `json:"roles"`
	} `json:"data"`
	Error *struct {
		Code string `json:"code"`
	} `json:"error"`
}

func TestIssueTokenSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	log, err := logger.New("info")
	if err != nil {
		t.Fatalf("failed to init logger: %v", err)
	}
	defer log.Sync()

	cfg := auth.Config{
		Secret:   "test-secret",
		Issuer:   "test-issuer",
		Audience: "test-audience",
		TokenTTL: time.Minute,
	}

	h := NewAuthHandler(log, cfg, "admin", "secret", []string{"admin", "user"})

	router := gin.New()
	router.POST("/token", h.IssueToken)

	payload := map[string]interface{}{
		"client_id":     "admin",
		"client_secret": "secret",
		"roles":         []string{"user"},
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/token", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}

	var resp tokenResponse
	if err := json.NewDecoder(recorder.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !resp.Success || resp.Data.AccessToken == "" {
		t.Fatalf("expected token response, got %+v", resp)
	}
}

func TestIssueTokenInvalidCredentials(t *testing.T) {
	gin.SetMode(gin.TestMode)

	log, err := logger.New("info")
	if err != nil {
		t.Fatalf("failed to init logger: %v", err)
	}
	defer log.Sync()

	cfg := auth.Config{
		Secret:   "test-secret",
		Issuer:   "test-issuer",
		Audience: "test-audience",
		TokenTTL: time.Minute,
	}

	h := NewAuthHandler(log, cfg, "admin", "secret", []string{"admin"})

	router := gin.New()
	router.POST("/token", h.IssueToken)

	payload := map[string]string{
		"client_id":     "admin",
		"client_secret": "wrong",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/token", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", recorder.Code)
	}
}

func TestIssueTokenRolesNotAllowed(t *testing.T) {
	gin.SetMode(gin.TestMode)

	log, err := logger.New("info")
	if err != nil {
		t.Fatalf("failed to init logger: %v", err)
	}
	defer log.Sync()

	cfg := auth.Config{
		Secret:   "test-secret",
		Issuer:   "test-issuer",
		Audience: "test-audience",
		TokenTTL: time.Minute,
	}

	h := NewAuthHandler(log, cfg, "admin", "secret", []string{"admin"})

	router := gin.New()
	router.POST("/token", h.IssueToken)

	payload := map[string]interface{}{
		"client_id":     "admin",
		"client_secret": "secret",
		"roles":         []string{"user"},
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/token", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", recorder.Code)
	}
}
