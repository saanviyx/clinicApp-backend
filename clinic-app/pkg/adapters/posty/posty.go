package posty

import (
	"database/sql"

	"go.uber.org/zap"
)

// MustInitDatabase initializes a PostgreSQL database connection pool.
//
// Parameters:
//   - dsn: Data Source Name (connection string) for the PostgreSQL database.
//   - maxIdleConns: Maximum number of idle connections in the pool.
//   - maxOpenConns: Maximum number of open connections in the pool.
//   - logger: Logger instance for logging errors and informational messages.
//
// Returns:
//   - A pointer to the sql.DB instance representing the database connection pool.
//   - An error if any occurs during initialization.

func MustInitDatabase(dsn string, maxIdleConns, maxOpenConns int, logger *zap.Logger) (*sql.DB, error) {
	logger.Info("Initializing Database Connections", zap.String("dsn", dsn))

	// Open a new connection to the PostgreSQL database
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Error("Error opening Database Connection", zap.Error(err))
		return nil, err
	}

	// Set the maximum number of idle connections in the pool
	db.SetMaxIdleConns(maxIdleConns)
	// Set the maximum number of open connections in the pool
	db.SetMaxOpenConns(maxOpenConns)

	// Test the database connection by pinging it
	if err := db.Ping(); err != nil {
		logger.Error("Error Ping to Database", zap.Error(err))
		return nil, err
	}

	return db, nil
}
