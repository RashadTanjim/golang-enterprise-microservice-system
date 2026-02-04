package api

import (
	"enterprise-microservice-system/common/auth"
	"enterprise-microservice-system/common/logger"
	"enterprise-microservice-system/common/metrics"
	"enterprise-microservice-system/common/middleware"
	auditdocs "enterprise-microservice-system/services/audit-log-service/docs"
	"enterprise-microservice-system/services/audit-log-service/internal/handler"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Router sets up all routes for the audit log service
type Router struct {
	handler     *handler.AuditLogHandler
	logger      *logger.Logger
	metrics     *metrics.Metrics
	rateLimiter *middleware.RateLimiter
	authConfig  auth.Config
}

// NewRouter creates a new router
func NewRouter(
	handler *handler.AuditLogHandler,
	logger *logger.Logger,
	metrics *metrics.Metrics,
	rateLimiter *middleware.RateLimiter,
	authConfig auth.Config,
) *Router {
	return &Router{
		handler:     handler,
		logger:      logger,
		metrics:     metrics,
		rateLimiter: rateLimiter,
		authConfig:  authConfig,
	}
}

// Setup configures all routes
func (r *Router) Setup() *gin.Engine {
	// Set Gin to release mode
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// Global middleware
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.RecoveryMiddleware(r.logger))
	router.Use(middleware.LoggerMiddleware(r.logger))
	router.Use(r.metrics.Middleware())
	router.Use(r.rateLimiter.Middleware())

	// Health check endpoint (no auth required)
	router.GET("/health", r.healthCheck)

	// Metrics endpoint (Prometheus)
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	auditdocs.SwaggerInfo.BasePath = "/api/v1"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 routes
	v1 := router.Group("/api/v1")

	protected := v1.Group("/")
	protected.Use(middleware.AuthMiddleware(r.authConfig))

	auditLogs := protected.Group("/audit-logs")
	{
		auditLogs.POST("", middleware.RequireRoles("admin", "user", "service"), r.handler.CreateAuditLog)
		auditLogs.GET("", middleware.RequireRoles("admin", "user"), r.handler.ListAuditLogs)
		auditLogs.GET("/:id", middleware.RequireRoles("admin", "user"), r.handler.GetAuditLog)
		auditLogs.PUT("/:id", middleware.RequireRoles("admin"), r.handler.UpdateAuditLog)
		auditLogs.DELETE("/:id", middleware.RequireRoles("admin"), r.handler.DeleteAuditLog)
	}

	return router
}

// healthCheck returns service health status
func (r *Router) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "audit-log-service",
	})
}
