package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strings"

	"enterprise-microservice-system/services/migration-service/internal/config"
	"enterprise-microservice-system/services/migration-service/migrations"
	"github.com/RashadTanjim/enterprise-microservice-system/common/logger"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	target := flag.String("target", "all", "migration target: user, order, all")
	flag.Parse()

	cfg := config.Load()

	log, err := logger.New(cfg.LogLevel)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	switch strings.ToLower(strings.TrimSpace(*target)) {
	case "user":
		runUser(cfg, log)
	case "order":
		runOrder(cfg, log)
	case "all":
		runUser(cfg, log)
		runOrder(cfg, log)
	default:
		log.Fatal("Unknown migration target", zap.String("target", *target))
	}
}

func runUser(cfg *config.Config, log *logger.Logger) {
	runMigration("user", cfg.UserDB, migrations.RunUser, log)
}

func runOrder(cfg *config.Config, log *logger.Logger) {
	runMigration("order", cfg.OrderDB, migrations.RunOrder, log)
}

func runMigration(name string, dbConfig config.DatabaseConfig, run func(*sql.DB) error, log *logger.Logger) {
	db, err := gorm.Open(postgres.Open(dbConfig.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database", zap.String("service", name), zap.Error(err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to access database connection", zap.String("service", name), zap.Error(err))
	}
	defer sqlDB.Close()

	if err := run(sqlDB); err != nil {
		log.Fatal("Failed to run migrations", zap.String("service", name), zap.Error(err))
	}

	log.Info("Migrations complete", zap.String("service", name))
}
