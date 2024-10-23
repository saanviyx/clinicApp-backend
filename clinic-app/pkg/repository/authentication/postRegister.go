package authentication

import (
	"clinic-app/cmd/rest/middleware"
	"clinic-app/pkg/domain/errors"
	"clinic-app/pkg/domain/models"
	"clinic-app/pkg/services/factory"

	"go.uber.org/zap"
)

// RegisterUser inserts a new user into the database
func (r *repo) RegisterUser(ftx factory.Service, user models.User) (int, error) {
	var userID int

	// Start a new transaction
	tx, err := ftx.TransactionManager().Begin()
	if err != nil {
		// Log and return error if transaction start fails
		ftx.Logger().Error("Could not begin transaction", zap.Error(err))
		return 0, errors.ErrDatabase
	}
	ftx.Logger().Info("Transaction started for inserting a new user")

	// Execute the insert query within the transaction context
	err = tx.QueryRowContext(ftx.Context(),
		RegisterUserQuery,
		user.Username,
		user.Name,
		user.Email,
		user.Password,
		user.Role).Scan(&userID)

	if err != nil {
		// Log error and rollback transaction if insert fails
		ftx.Logger().Error("Could not register user", zap.Error(err))
		tx.Rollback() // Rollback transaction on error
		return 0, errors.ErrDatabase
	}

	// Commit the transaction if no errors occurred
	if err := ftx.TransactionManager().Commit(tx); err != nil {
		// Log error and return error if commit fails
		ftx.Logger().Error("Could not commit transaction", zap.Error(err))
		return 0, errors.ErrDatabase
	}

	// Log success
	ftx.Logger().Info("Successfully registered a new user",
		zap.Int("UserID", userID),
		zap.String("Username", user.Username),
	)
	middleware.GetTraceParentFromContext(ftx.Context())

	return userID, nil
}
