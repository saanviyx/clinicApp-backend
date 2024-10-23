package usecase

import (
	"clinic-app/pkg/domain/models"
	"clinic-app/pkg/services/factory"
)

// DoctorUsecase defines methods for managing doctors and their slots.
type DoctorUsecase interface {
	AllDoctors(ftx factory.Service) ([]models.Doctor, error)
	DoctorById(ftx factory.Service, doctorId int) (models.Doctor, error)
	Slots(ftx factory.Service, doctorId int, check bool) ([]interface{}, error)
}
