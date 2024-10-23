package authentication

import (
	"clinic-app/pkg/repository"
	"clinic-app/pkg/usecase"
)

type authUsecaseImpl struct {
	repo repository.AuthenticationRepository
}

// New creates a new instance of repository with a database connection
func New(repo repository.AuthenticationRepository) usecase.AuthUsecase {
	return &authUsecaseImpl{
		repo,
	}
}
