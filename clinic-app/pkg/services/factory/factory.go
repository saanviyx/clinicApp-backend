package factory

import (
	"clinic-app/internal/constants"
	"clinic-app/pkg/domain/errors"
	"clinic-app/pkg/services/transaction"
	"context"
	"database/sql"

	"go.uber.org/zap"
)

// Service interface defines the methods available for services.
type Service interface {
	Logger() *zap.Logger
	PSQL() *sql.DB
	TransactionManager() *transaction.TransactionManager
	Context() context.Context
}

// serviceImpl is the concrete implementation of the Service interface.
type ServiceImpl struct {
	logger *zap.Logger
	db     *sql.DB
	trx    *transaction.TransactionManager
	ctx    context.Context
}

var deps *ServiceImpl

// NewFactory creates a new Factory instance.
func NewFactory(db *sql.DB, ctx context.Context) (Service, error) {
	// Initialize the logger
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	// Create a database/transaction manager instance
	trxMgr := transaction.NewTransactionManager(db, logger)

	return &ServiceImpl{
		db:     db,
		logger: logger,
		trx:    trxMgr,
		ctx:    ctx,
	}, nil
}

// NewFactoryFromTraceParent creates a new factory Service from a traceparent string
func NewFactoryFromTraceParent(traceparent string) (Service, error) {
	if deps == nil {
		// Return an error if dependencies have not been set up
		return nil, errors.ErrNotFound
	}
	// Create a transaction context using the factory
	ctx := context.WithValue(context.Background(), constants.TraceparentHeader, traceparent)

	// Return a new Service implementation with logging
	return &ServiceImpl{
		db:     deps.db,
		logger: deps.logger.With(zap.String("traceparent", traceparent)),
		trx:    deps.trx,
		ctx:    ctx,
	}, nil
}

// SetUpServices sets up the services with the provided dependencies.
func SetUpDependencies(db *sql.DB) {
	f, err := NewFactory(db, context.Background())
	if err != nil {
		f.Logger().Error("Error creating Factory instance", zap.Error(err))
	}
	f.Logger().Info("Setting up services...")

	if deps != nil {
		return
	}

	// Initialize the service with the provided options
	deps = &ServiceImpl{
		db:     f.PSQL(),
		logger: f.Logger(),
		trx:    f.TransactionManager(),
		ctx:    f.Context(),
	}

	// Log successful setup
	f.Logger().Info("Services set up successfully")
}

// Logger returns the logger instance.
func (s *ServiceImpl) Logger() *zap.Logger {
	return s.logger
}

// PSQL returns the database connection.
func (s *ServiceImpl) PSQL() *sql.DB {
	return s.db
}

// TransactionManager returns the transaction manager.
func (s *ServiceImpl) TransactionManager() *transaction.TransactionManager {
	return s.trx
}

// Context returns the context with traceparent
func (s *ServiceImpl) Context() context.Context {
	return s.ctx
}
