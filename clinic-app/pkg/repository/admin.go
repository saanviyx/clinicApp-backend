package repository

import (
	"clinic-app/pkg/domain/models"
	"clinic-app/pkg/services/factory"
)

// AdminRepository defines methods for accessing admin-related data.
type AdminRepository interface {
	GetAllDoctorsAvailability(ftx factory.Service) ([]models.DoctorAvailability, error)
	GetDoctorsWithMostAppointments(ftx factory.Service, date string) ([]models.DoctorMostAppointments, error)
	GetDoctorsWithOverSixHours(ftx factory.Service, date string) ([]models.DoctorOverTime, error)
}
