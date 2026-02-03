package main

import (
	"enterprise-microservice-system/common/logger"
	"enterprise-microservice-system/services/user-service/internal/config"
	"enterprise-microservice-system/services/user-service/migrations"
	"fmt"
	"os"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	log, err := logger.New(cfg.Log.Level)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	db, err := gorm.Open(postgres.Open(cfg.Database.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to access database connection", zap.Error(err))
	}
	defer sqlDB.Close()

	if err := migrations.Run(sqlDB); err != nil {
		log.Fatal("Failed to run migrations", zap.Error(err))
	}

	log.Info("User service migrations complete")
}
