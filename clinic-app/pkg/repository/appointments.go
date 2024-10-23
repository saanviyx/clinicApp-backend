package repository

import (
	"clinic-app/pkg/domain/models"
	"clinic-app/pkg/services/factory"
)

// AppointmentRepository defines methods for managing appointments
type AppointmentRepository interface {
	BookAppointment(ftx factory.Service, aptmt models.BookAppointment) error
	GetAppointmentById(ftx factory.Service, appointmentId int) (models.Appointment, error)
	GetPatientHistory(ftx factory.Service, patientId int) ([]models.Appointment, error)
	GetPatientAppointmentHistory(ftx factory.Service) ([]models.Appointment, error)
	CancelAppointment(ftx factory.Service, appointmentId int) error
}
