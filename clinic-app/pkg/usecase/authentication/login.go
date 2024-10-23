package authentication

import (
	"clinic-app/pkg/domain/models"
	"clinic-app/pkg/services/factory"
)

// LoginUser handles user login by verifying credentials.
func (uc *authUsecaseImpl) LoginUser(ftx factory.Service, credentials models.Credentials) (models.User, error) {
	// Attempt to login using the provided username and password
	user, err := uc.repo.LoginUser(ftx, credentials.Username, credentials.Password)
	if err != nil || user.Password != credentials.Password {
		// Log an error if credentials are invalid or the user could not be found
		ftx.Logger().Error("Invalid credentials")
		return user, err
	}

	// Return the user details if credentials are valid
	return user, nil
}
