package middleware

import (
	"enterprise-microservice-system/common/errors"
	"enterprise-microservice-system/common/logger"
	"enterprise-microservice-system/common/response"
	"fmt"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RecoveryMiddleware recovers from panics and logs the error
func RecoveryMiddleware(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic with stack trace
				log.Error("Panic recovered",
					zap.Any("error", err),
					zap.String("stack", string(debug.Stack())),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
				)

				// Return error response
				appErr := errors.NewInternal(
					fmt.Sprintf("internal server error: %v", err),
					nil,
				)
				response.Error(c, appErr)
				c.Abort()
			}
		}()

		c.Next()
	}
}
