package appointments

import (
	"clinic-app/pkg/domain/models"
	"clinic-app/pkg/services/factory"

	"go.uber.org/zap"
)

// Book schedules a new appointment.
func (uc *aptmtUsecaseImpl) Book(ftx factory.Service, aptmt models.BookAppointment) error {
	// Call the repository method to book the appointment
	err := uc.repo.BookAppointment(ftx, aptmt)
	if err != nil {
		// Log an error if the booking fails
		ftx.Logger().Error("Error booking appointment", zap.Error(err))
		return err
	}

	// Return nil if no error occurred
	return nil
}
