package usecase

import (
	"clinic-app/pkg/domain/models"
	"clinic-app/pkg/services/factory"
)

// AdminUsecase defines the methods for managing administrative use cases.
type AdminUsecase interface {
	DoctorsAvailability(ftx factory.Service) ([]models.DoctorAvailability, error)
	MostAppointments(ftx factory.Service, date string) ([]models.DoctorMostAppointments, error)
	OverSixHours(ftx factory.Service, date string) ([]models.DoctorOverTime, error)
}
