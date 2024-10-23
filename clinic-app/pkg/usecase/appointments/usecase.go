package appointments

import (
	"clinic-app/pkg/repository"
	"clinic-app/pkg/usecase"
)

type aptmtUsecaseImpl struct {
	repo repository.AppointmentRepository
}

// NewaptmtUsecase creates a new instance of aptmtUsecaseImpl and returns it as the aptmtUsecase interface
func New(repo repository.AppointmentRepository) usecase.AppointmentUsecase {
	return &aptmtUsecaseImpl{
		repo,
	}
}
