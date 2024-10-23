package migrations

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq" // PostgreSQL driver for database/sql (blank import to initialize it)
	"go.uber.org/zap"
)

// ApplyMigrations applies or rolls back SQL migrations from a directory
func ApplyMigrations(migrationsDir string, rollback bool, DB *sql.DB, Logger *zap.Logger) error {
	Logger.Info("Initialising Migrations", zap.String("Directory:", migrationsDir))
	// Initialize the migrate instance
	driver, err := postgres.WithInstance(DB, &postgres.Config{})
	if err != nil {
		Logger.Error("Could not create postgres driver", zap.Error(err))
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsDir), "postgres", driver)
	if err != nil {
		Logger.Error("Could not initialize migrate instance", zap.Error(err))
		return err
	}

	// Apply or rollback migrations based on the `rollback` flag
	if rollback {
		err := m.Down()
		if err != nil && err != migrate.ErrNoChange {
			Logger.Error("Could not apply rollback migration", zap.Error(err))
			return err
		}
		if err == migrate.ErrNoChange {
			Logger.Info("No migrations to rollback")
		} else {
			Logger.Info("Database rollback (down migrations) completed successfully")
		}
	} else {
		err := m.Up()
		if err != nil && err != migrate.ErrNoChange {
			Logger.Error("Could not apply migration", zap.Error(err))
			return err
		}
		if err == migrate.ErrNoChange {
			Logger.Info("No new migrations to apply")
		} else {
			Logger.Info("Database migrations applied successfully")
		}
	}

	return nil
}
