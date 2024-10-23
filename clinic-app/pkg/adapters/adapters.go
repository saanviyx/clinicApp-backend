package adapters

import (
	"database/sql"
	"log"

	"go.uber.org/zap"

	// Import PostgreSQL driver for database connection
	_ "github.com/lib/pq"

	// Import the package for initializing the database (adjust import path as needed)
	"clinic-app/pkg/adapters/posty"
)

// Options holds the configuration for setting up the adapters.
type Options struct {
	DBConnStr      string // Connection string for PostgreSQL database
	DBMaxIdleConns int    // Maximum number of idle connections in the pool
	DBMaxOpenConns int    // Maximum number of open connections in the pool
}

// Results holds the initialized adapters.
type Results struct {
	DB     *sql.DB     // Database connection
	Logger *zap.Logger // Logger instance
}

// SetupAdapters initializes and returns the database connection and logger.
func SetupAdapters(opts *Options) (*Results, error) {
	res := &Results{}

	// Initialize the logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Error initializing logger: %v", err) // Fatal error if logger initialization fails
	}
	res.Logger = logger

	// Initialize the database connection using the posty package
	res.DB, err = posty.MustInitDatabase(
		opts.DBConnStr,      // Database connection string
		opts.DBMaxIdleConns, // Maximum number of idle connections
		opts.DBMaxOpenConns, // Maximum number of open connections
		logger)              // Logger instance
	if err != nil {
		logger.Error("Error initializing database", zap.Error(err)) // Log error if database initialization fails
		return nil, err
	}

	logger.Info("Adapters set up successfully") // Log success message
	return res, nil
}
