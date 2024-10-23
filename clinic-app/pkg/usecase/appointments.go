package usecase

import (
	"clinic-app/pkg/domain/models"
	"clinic-app/pkg/services/factory"
)

// AppointmentUsecase defines methods for managing appointments.
type AppointmentUsecase interface {
	Book(ftx factory.Service, aptmt models.BookAppointment) error
	ViewAppointment(ftx factory.Service, appointmentId int) (models.Appointment, error)
	PatientHistoryForDoctor(ftx factory.Service, patientId int) ([]models.Appointment, error)
	PatientHistory(ftx factory.Service) ([]models.Appointment, error)
	Cancel(ftx factory.Service, appointmentId int) error
}
