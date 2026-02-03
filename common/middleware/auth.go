package middleware

import (
	"enterprise-microservice-system/common/auth"
	apperrors "enterprise-microservice-system/common/errors"
	"enterprise-microservice-system/common/response"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	contextKeyAuthClaims  = "auth_claims"
	contextKeyAuthRoles   = "auth_roles"
	contextKeyAuthSubject = "auth_subject"
)

// AuthMiddleware validates JWT tokens and stores claims in the request context.
func AuthMiddleware(cfg auth.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := extractBearerToken(c.GetHeader("Authorization"))
		if err != nil {
			response.Error(c, apperrors.New(apperrors.ErrCodeUnauthorized, err.Error(), nil))
			c.Abort()
			return
		}

		claims, err := auth.ParseToken(cfg, token)
		if err != nil {
			response.Error(c, apperrors.New(apperrors.ErrCodeUnauthorized, "invalid or expired token", err))
			c.Abort()
			return
		}

		c.Set(contextKeyAuthClaims, claims)
		c.Set(contextKeyAuthRoles, claims.Roles)
		c.Set(contextKeyAuthSubject, claims.Subject)
		c.Next()
	}
}

// RequireRoles enforces role-based authorization for handlers.
func RequireRoles(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(roles) == 0 {
			c.Next()
			return
		}

		rolesValue, exists := c.Get(contextKeyAuthRoles)
		if !exists {
			response.Error(c, apperrors.New(apperrors.ErrCodeUnauthorized, "missing authentication context", nil))
			c.Abort()
			return
		}

		assignedRoles, ok := rolesValue.([]string)
		if !ok || len(assignedRoles) == 0 {
			response.Error(c, apperrors.New(apperrors.ErrCodeUnauthorized, "missing roles", nil))
			c.Abort()
			return
		}

		if !hasAnyRole(assignedRoles, roles) {
			response.Error(c, apperrors.New(apperrors.ErrCodeForbidden, "insufficient permissions", nil))
			c.Abort()
			return
		}

		c.Next()
	}
}

func extractBearerToken(headerValue string) (string, error) {
	if headerValue == "" {
		return "", errors.New("missing authorization header")
	}

	parts := strings.SplitN(headerValue, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", errors.New("invalid authorization header")
	}

	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", errors.New("empty bearer token")
	}

	return token, nil
}

func hasAnyRole(assigned []string, required []string) bool {
	for _, role := range assigned {
		for _, needed := range required {
			if role == needed {
				return true
			}
		}
	}
	return false
}

// GetAuthSubject retrieves the subject from context when needed by handlers.
func GetAuthSubject(c *gin.Context) (string, bool) {
	value, exists := c.Get(contextKeyAuthSubject)
	if !exists {
		return "", false
	}

	subject, ok := value.(string)
	return subject, ok
}

// GetAuthRoles retrieves roles from context when needed by handlers.
func GetAuthRoles(c *gin.Context) ([]string, bool) {
	value, exists := c.Get(contextKeyAuthRoles)
	if !exists {
		return nil, false
	}

	roles, ok := value.([]string)
	return roles, ok
}
