package middleware

import (
	"encoding/json"
	"enterprise-microservice-system/common/auth"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

type errorResponse struct {
	Success bool `json:"success"`
	Error   struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func TestAuthMiddlewareMissingToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := auth.Config{
		Secret:   "test-secret",
		Issuer:   "test-issuer",
		Audience: "test-audience",
		TokenTTL: time.Minute,
	}

	router := gin.New()
	protected := router.Group("/protected")
	protected.Use(AuthMiddleware(cfg))
	protected.GET("", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", recorder.Code)
	}

	var resp errorResponse
	if err := json.NewDecoder(recorder.Body).Decode(&resp); err == nil && resp.Error.Code == "" {
		t.Fatalf("expected error response")
	}
}

func TestAuthMiddlewareInvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := auth.Config{
		Secret:   "test-secret",
		Issuer:   "test-issuer",
		Audience: "test-audience",
		TokenTTL: time.Minute,
	}

	router := gin.New()
	protected := router.Group("/protected")
	protected.Use(AuthMiddleware(cfg))
	protected.GET("", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", recorder.Code)
	}
}

func TestRequireRoles(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := auth.Config{
		Secret:   "test-secret",
		Issuer:   "test-issuer",
		Audience: "test-audience",
		TokenTTL: time.Minute,
	}

	router := gin.New()
	protected := router.Group("/protected")
	protected.Use(AuthMiddleware(cfg))
	protected.Use(RequireRoles("admin"))
	protected.GET("", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	adminToken, err := auth.GenerateToken(cfg, "admin-user", []string{"admin"})
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	userToken, err := auth.GenerateToken(cfg, "user", []string{"user"})
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	cases := []struct {
		name       string
		token      string
		wantStatus int
	}{
		{name: "admin role", token: adminToken, wantStatus: http.StatusOK},
		{name: "missing role", token: userToken, wantStatus: http.StatusForbidden},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			req.Header.Set("Authorization", "Bearer "+tc.token)
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			if recorder.Code != tc.wantStatus {
				t.Fatalf("expected status %d, got %d", tc.wantStatus, recorder.Code)
			}
		})
	}
}
