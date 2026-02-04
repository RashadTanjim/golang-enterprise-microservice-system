package main

import (
	"context"
	"github.com/RashadTanjim/enterprise-microservice-system/common/audit"
	"github.com/RashadTanjim/enterprise-microservice-system/common/auth"
	"github.com/RashadTanjim/enterprise-microservice-system/common/cache"
	"github.com/RashadTanjim/enterprise-microservice-system/common/logger"
	"github.com/RashadTanjim/enterprise-microservice-system/common/metrics"
	"github.com/RashadTanjim/enterprise-microservice-system/common/middleware"
	"enterprise-microservice-system/services/user-service/internal/api"
	"enterprise-microservice-system/services/user-service/internal/config"
	"enterprise-microservice-system/services/user-service/internal/handler"
	"enterprise-microservice-system/services/user-service/internal/repository"
	"enterprise-microservice-system/services/user-service/internal/service"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	log, err := logger.New(cfg.Log.Level)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	log.Info("Starting user service",
		zap.String("port", cfg.Server.Port),
		zap.String("log_level", cfg.Log.Level),
	)

	// Connect to database
	db, err := connectDatabase(cfg, log)
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}

	// Initialize dependencies
	userRepo := repository.NewUserRepository(db)
	userCache, err := cache.NewRedisCache(cache.Config{
		Enabled:    cfg.Redis.Enabled,
		Host:       cfg.Redis.Host,
		Port:       cfg.Redis.Port,
		Password:   cfg.Redis.Password,
		DB:         cfg.Redis.DB,
		DefaultTTL: cfg.Redis.DefaultTTL,
	}, "user-service")
	if err != nil {
		log.Warn("Redis cache disabled", zap.Error(err))
	}
	userService := service.NewUserService(userRepo, userCache)

	auditClient := audit.NewClient(audit.Config{
		Enabled: cfg.AuditLog.Enabled,
		BaseURL: cfg.AuditLog.URL,
		Timeout: cfg.AuditLog.Timeout,
	}, log)

	userHandler := handler.NewUserHandler(userService, auditClient, log)

	authConfig := auth.Config{
		Secret:   cfg.Auth.Secret,
		Issuer:   cfg.Auth.Issuer,
		Audience: cfg.Auth.Audience,
		TokenTTL: cfg.Auth.TokenTTL,
	}
	authHandler := handler.NewAuthHandler(log, auditClient, authConfig, cfg.Auth.ClientID, cfg.Auth.ClientSecret, cfg.Auth.ClientRoles)

	// Initialize metrics
	metricsCollector := metrics.NewMetrics("user_service")

	// Initialize rate limiter
	rateLimiter := middleware.NewRateLimiter(cfg.Server.RateLimit, cfg.Server.RateLimit*2)

	// Setup router
	routerSetup := api.NewRouter(userHandler, authHandler, log, metricsCollector, rateLimiter, authConfig)
	router := routerSetup.Setup()

	// Create HTTP server
	srv := &http.Server{
		Addr:           ":" + cfg.Server.Port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// Start server in a goroutine
	go func() {
		log.Info("Server listening", zap.String("address", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown", zap.Error(err))
	}

	// Close database connection
	if sqlDB != nil {
		sqlDB.Close()
	}

	log.Info("Server exited")
}

// connectDatabase establishes a connection to the PostgreSQL database
func connectDatabase(cfg *config.Config, log *logger.Logger) (*gorm.DB, error) {
	dsn := cfg.Database.DSN()

	// Configure GORM
	gormConfig := &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		PrepareStmt: true,
	}

	// Connect to database with retries
	var db *gorm.DB
	var err error

	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(dsn), gormConfig)
		if err == nil {
			break
		}

		log.Warn("Failed to connect to database, retrying...",
			zap.Int("attempt", i+1),
			zap.Int("max_retries", maxRetries),
			zap.Error(err),
		)
		time.Sleep(time.Second * time.Duration(i+1))
	}

	if err != nil {
		return nil, err
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Info("Database connection established")
	return db, nil
}
