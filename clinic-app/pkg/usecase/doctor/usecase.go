package doctor

import (
	"clinic-app/pkg/repository"
	"clinic-app/pkg/usecase"
)

type doctorsUsecaseImpl struct {
	repo repository.DoctorRepository
}

// NewdoctorsUsecase creates a new instance of doctorsUsecaseImpl and returns it as the doctorsUsecase interface
func New(repo repository.DoctorRepository) usecase.DoctorUsecase {
	return &doctorsUsecaseImpl{
		repo,
	}
}
