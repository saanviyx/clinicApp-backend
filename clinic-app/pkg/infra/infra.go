package infra

import (
	"clinic-app/pkg/infra/migrations"
	"database/sql"

	"go.uber.org/zap"
)

type Options struct {
	DB            *sql.DB
	Logger        *zap.Logger
	MigrationPath string
}

type Infrastructure struct {
	DB     *sql.DB
	Logger *zap.Logger
}

// NewInfrastructure creates a new instance of Infrastructure with database and logger
func NewInfrastructure(opts *Options) (*Infrastructure, error) {
	// Create Infrastructure instance
	infra := &Infrastructure{
		DB:     opts.DB,
		Logger: opts.Logger,
	}

	// Apply migrations
	if err := migrations.ApplyMigrations(opts.MigrationPath, false, infra.DB, infra.Logger); err != nil {
		opts.Logger.Error("Error Applying Migrations", zap.Error(err))
		return nil, err
	}

	opts.Logger.Info("Infrastructure set up Successfully")
	return infra, nil
}
