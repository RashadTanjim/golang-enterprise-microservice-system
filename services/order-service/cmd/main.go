package main

import (
	"context"
	"enterprise-microservice-system/common/circuitbreaker"
	"enterprise-microservice-system/common/logger"
	"enterprise-microservice-system/common/metrics"
	"enterprise-microservice-system/common/middleware"
	"enterprise-microservice-system/services/order-service/internal/api"
	"enterprise-microservice-system/services/order-service/internal/client"
	"enterprise-microservice-system/services/order-service/internal/config"
	"enterprise-microservice-system/services/order-service/internal/handler"
	"enterprise-microservice-system/services/order-service/internal/model"
	"enterprise-microservice-system/services/order-service/internal/repository"
	"enterprise-microservice-system/services/order-service/internal/service"
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

	log.Info("Starting order service",
		zap.String("port", cfg.Server.Port),
		zap.String("log_level", cfg.Log.Level),
	)

	// Connect to database
	db, err := connectDatabase(cfg, log)
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}

	// Auto-migrate database schema
	if err := db.AutoMigrate(&model.Order{}); err != nil {
		log.Fatal("Failed to migrate database", zap.Error(err))
	}
	log.Info("Database migration completed")

	// Initialize circuit breaker for user service
	cbConfig := circuitbreaker.Config{
		MaxRequests: cfg.CircuitBreaker.MaxRequests,
		Interval:    cfg.CircuitBreaker.Interval,
		Timeout:     cfg.CircuitBreaker.Timeout,
	}
	userServiceCB := circuitbreaker.New("user-service", cbConfig)
	log.Info("Circuit breaker initialized",
		zap.Uint32("max_requests", cbConfig.MaxRequests),
		zap.Duration("interval", cbConfig.Interval),
		zap.Duration("timeout", cbConfig.Timeout),
	)

	// Initialize user service client
	userClient := client.NewUserClient(cfg.UserService.URL, userServiceCB)

	// Initialize dependencies
	orderRepo := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(orderRepo, userClient)
	orderHandler := handler.NewOrderHandler(orderService, log)

	// Initialize metrics
	metricsCollector := metrics.NewMetrics("order_service")

	// Initialize rate limiter
	rateLimiter := middleware.NewRateLimiter(cfg.Server.RateLimit, cfg.Server.RateLimit*2)

	// Start background worker to update circuit breaker metrics
	go updateCircuitBreakerMetrics(userClient, metricsCollector, log)

	// Setup router
	routerSetup := api.NewRouter(orderHandler, log, metricsCollector, rateLimiter)
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
	sqlDB, err := db.DB()
	if err == nil {
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

// updateCircuitBreakerMetrics updates circuit breaker state metrics
func updateCircuitBreakerMetrics(userClient *client.UserClient, metrics *metrics.Metrics, log *logger.Logger) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		state := userClient.GetCircuitBreakerState()
		metrics.SetCircuitState("user-service", state)
	}
}
