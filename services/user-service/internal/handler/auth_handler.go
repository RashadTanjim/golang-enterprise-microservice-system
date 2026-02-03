package handler

import (
	"enterprise-microservice-system/common/auth"
	"enterprise-microservice-system/common/errors"
	"enterprise-microservice-system/common/logger"
	"enterprise-microservice-system/common/response"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AuthHandler handles authentication token issuance.
type AuthHandler struct {
	logger       *logger.Logger
	authConfig   auth.Config
	clientID     string
	clientSecret string
	clientRoles  []string
}

// TokenRequest represents the token request payload.
type TokenRequest struct {
	ClientID     string   `json:"client_id" binding:"required"`
	ClientSecret string   `json:"client_secret" binding:"required"`
	Roles        []string `json:"roles"`
}

// TokenResponse represents the token response payload.
type TokenResponse struct {
	AccessToken string    `json:"access_token"`
	TokenType   string    `json:"token_type"`
	ExpiresAt   time.Time `json:"expires_at"`
	Roles       []string  `json:"roles"`
}

// NewAuthHandler creates a new auth handler.
func NewAuthHandler(log *logger.Logger, authConfig auth.Config, clientID, clientSecret string, clientRoles []string) *AuthHandler {
	return &AuthHandler{
		logger:       log,
		authConfig:   authConfig,
		clientID:     clientID,
		clientSecret: clientSecret,
		clientRoles:  clientRoles,
	}
}

// IssueToken validates client credentials and issues a JWT.
// @Summary Issue a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param token body TokenRequest true "Token request"
// @Success 200 {object} response.Response{data=TokenResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Router /auth/token [post]
func (h *AuthHandler) IssueToken(c *gin.Context) {
	var req TokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid token request", zap.Error(err))
		response.Error(c, err)
		return
	}

	if req.ClientID != h.clientID || req.ClientSecret != h.clientSecret {
		response.Error(c, errors.New(errors.ErrCodeUnauthorized, "invalid client credentials", nil))
		return
	}

	roles := h.clientRoles
	if len(req.Roles) > 0 {
		if !rolesAllowed(req.Roles, h.clientRoles) {
			response.Error(c, errors.New(errors.ErrCodeForbidden, "requested roles are not allowed", nil))
			return
		}
		roles = req.Roles
	}

	token, err := auth.GenerateToken(h.authConfig, req.ClientID, roles)
	if err != nil {
		h.logger.Error("Failed to generate token", zap.Error(err))
		response.Error(c, errors.New(errors.ErrCodeInternal, "failed to generate token", err))
		return
	}

	expiresAt := time.Now().UTC().Add(h.authConfig.TokenTTL)
	response.Success(c, TokenResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresAt:   expiresAt,
		Roles:       roles,
	})
}

func rolesAllowed(requested []string, allowed []string) bool {
	allowedSet := make(map[string]struct{}, len(allowed))
	for _, role := range allowed {
		allowedSet[role] = struct{}{}
	}

	for _, role := range requested {
		if _, ok := allowedSet[role]; !ok {
			return false
		}
	}

	return true
}
