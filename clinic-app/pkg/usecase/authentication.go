package usecase

import (
	"clinic-app/pkg/domain/models"
	"clinic-app/pkg/services/factory"
)

// AuthUsecase defines methods for user authentication and registration.
type AuthUsecase interface {
	RegisterUser(ftx factory.Service, user models.User) (int, error)
	LoginUser(ftx factory.Service, credentials models.Credentials) (models.User, error)
}
