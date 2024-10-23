package authentication

import (
	"clinic-app/pkg/domain/models"
	"clinic-app/pkg/services/factory"

	"go.uber.org/zap"
)

// RegisterUser registers a new user and returns the user ID.
func (uc *authUsecaseImpl) RegisterUser(ftx factory.Service, user models.User) (int, error) {
	// Attempt to register the user using the provided user details
	userID, err := uc.repo.RegisterUser(ftx, user)
	if err != nil {
		// Log an error if user registration fails
		ftx.Logger().Error("Failed to register user", zap.Error(err))
		return 0, err
	}

	// Return the user ID if registration is successful
	return userID, nil
}
