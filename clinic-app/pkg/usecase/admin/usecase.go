package admin

import (
	"clinic-app/pkg/repository"
	"clinic-app/pkg/usecase"
)

type adminUsecaseImpl struct {
	repo repository.AdminRepository
}

// NewadminUsecase creates a new instance of adminUsecaseImpl and returns it as the adminUsecase interface
func New(repo repository.AdminRepository) usecase.AdminUsecase {
	return &adminUsecaseImpl{
		repo,
	}
}
