package migrations

import (
	"database/sql"
	"embed"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed user/*.sql
var userMigrationsFS embed.FS

//go:embed order/*.sql
var orderMigrationsFS embed.FS

// RunUser applies all up migrations for the user service.
func RunUser(db *sql.DB) error {
	return run(db, userMigrationsFS, "user")
}

// RunOrder applies all up migrations for the order service.
func RunOrder(db *sql.DB) error {
	return run(db, orderMigrationsFS, "order")
}

func run(db *sql.DB, migrationsFS embed.FS, dir string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	source, err := iofs.New(migrationsFS, dir)
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance("iofs", source, "postgres", driver)
	if err != nil {
		return err
	}
	defer func() {
		_, _ = m.Close()
	}()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
