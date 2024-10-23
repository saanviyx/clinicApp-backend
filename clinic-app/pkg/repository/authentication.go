package repository

import (
	"clinic-app/pkg/domain/models"
	"clinic-app/pkg/services/factory"
)

// AuthenticationRepository defines methods for user authentication and registration.
type AuthenticationRepository interface {
	RegisterUser(ftx factory.Service, user models.User) (int, error)
	LoginUser(ftx factory.Service, username, password string) (models.User, error)
}
