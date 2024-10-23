package repository

import (
	"clinic-app/pkg/domain/models"
	"clinic-app/pkg/services/factory"
)

// DoctorRepository defines methods for interacting with doctor data.
type DoctorRepository interface {
	GetAllDoctors(ftx factory.Service) ([]models.Doctor, error)
	GetDoctorById(ftx factory.Service, doctorId int) (models.Doctor, error)
	GetDoctorSlots(ftx factory.Service, doctorId int, isDoctor bool) ([]interface{}, error)
}
