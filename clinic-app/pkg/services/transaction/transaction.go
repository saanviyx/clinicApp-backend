package transaction

import (
	"database/sql"
	"errors"

	"go.uber.org/zap"
)

// TransactionManager manages database transactions.
type TransactionManager struct {
	db     *sql.DB
	logger *zap.Logger
}

// NewTransactionManager creates a new TransactionManager.
func NewTransactionManager(db *sql.DB, logger *zap.Logger) *TransactionManager {
	return &TransactionManager{
		db:     db,
		logger: logger,
	}
}

// Begin starts a new transaction.
func (tm TransactionManager) Begin() (*sql.Tx, error) {
	tm.logger.Info("Starting new transaction")
	tx, err := tm.db.Begin()
	if err != nil {
		tm.logger.Error("Failed to start transaction", zap.Error(err))
		return nil, err
	}
	return tx, nil
}

// Commit commits the transaction.
func (tm *TransactionManager) Commit(tx *sql.Tx) error {
	if tx == nil {
		tm.logger.Error("No transaction to commit")
		return errors.New("no transaction to commit")
	}

	tm.logger.Info("Committing transaction")
	if err := tx.Commit(); err != nil {
		tm.logger.Error("Failed to commit transaction", zap.Error(err))
		return err
	}
	return nil
}

// Rollback rolls back the transaction.
func (tm *TransactionManager) Rollback(tx *sql.Tx) error {
	if tx == nil {
		tm.logger.Error("No transaction to rollback")
		return errors.New("no transaction to rollback")
	}

	tm.logger.Info("Rolling back transaction")
	if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
		tm.logger.Error("Failed to rollback transaction", zap.Error(err))
		return err
	}
	return nil
}
