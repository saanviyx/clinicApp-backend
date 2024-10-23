package services

import (
	"clinic-app/pkg/services/factory"
	"database/sql"

	"go.uber.org/zap"
)

// Options holds the configuration for the service setup.
type Options struct {
	Logger *zap.Logger
	DB     *sql.DB
}

// NewService creates a new service instance with the provided dependencies.
func SetupService(opts *Options) error {
	factory.SetUpDependencies(opts.DB)
	return nil
}
