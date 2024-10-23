package authentication

import (
	"clinic-app/cmd/rest/middleware"
	"clinic-app/pkg/domain/errors"
	"clinic-app/pkg/domain/models"
	"clinic-app/pkg/services/factory"

	"database/sql"

	"go.uber.org/zap"
)

// LoginUser retrieves a user by email and validates the password
func (r *repo) LoginUser(ftx factory.Service, username, password string) (models.User, error) {
	var user models.User

	// Start a new transaction
	tx, err := ftx.TransactionManager().Begin()
	if err != nil {
		// Log and return error if transaction start fails
		ftx.Logger().Error("Could not begin transaction", zap.Error(err))
		return user, errors.ErrDatabase
	}
	ftx.Logger().Info("Transaction started for user login")

	// Retrieve user data based on email and password
	err = tx.QueryRowContext(ftx.Context(), LoginUserQuery, username, password).Scan(&user.ID, &user.Username, &user.Password, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			// Log and return error if the user is not found
			ftx.Logger().Error("User not found", zap.String("User", username))
			tx.Rollback() // Rollback transaction in case of error
			return user, errors.ErrUserNotFound
		}
		// Log and return error for other database retrieval issues
		ftx.Logger().Error("Could not retrieve user", zap.Error(err))
		tx.Rollback() // Rollback transaction in case of error
		return user, errors.ErrDatabase
	}

	// Validate the provided password (no hashing, just comparison)
	if user.Password != password {
		// Log and return error if the password is incorrect
		ftx.Logger().Error("Invalid password")
		tx.Rollback() // Rollback transaction in case of invalid password
		return user, errors.ErrInvalidPassword
	}

	// Commit the transaction if no errors occurred
	if err := ftx.TransactionManager().Commit(tx); err != nil {
		// Log and return error if transaction commit fails
		ftx.Logger().Error("Could not commit transaction", zap.Error(err))
		return user, errors.ErrDatabase
	}

	// Defer a rollback in case of any subsequent errors (redundant here as commit has succeeded)
	defer func() {
		if err != nil {
			rollbackErr := ftx.TransactionManager().Rollback(tx)
			if rollbackErr != nil {
				// Log rollback failure
				ftx.Logger().Error("Failed to rollback transaction", zap.Error(rollbackErr))
			}
		}
	}()

	// Log success and return user information
	ftx.Logger().Info("Successfully logged in user",
		zap.Any("User", user.ID),
	)
	middleware.GetTraceParentFromContext(ftx.Context())
	return user, nil
}
