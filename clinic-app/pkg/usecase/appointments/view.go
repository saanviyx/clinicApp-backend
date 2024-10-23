package appointments

import (
	"clinic-app/pkg/domain/models"
	"clinic-app/pkg/services/factory"

	"go.uber.org/zap"
)

// ViewAppointment retrieves details of a specific appointment by its ID.
func (uc *aptmtUsecaseImpl) ViewAppointment(ftx factory.Service, aptmtId int) (models.Appointment, error) {
	// Call the repository method to get the appointment details by ID
	aptmt, err := uc.repo.GetAppointmentById(ftx, aptmtId)
	if err != nil {
		// Log an error if fetching the appointment details fails
		ftx.Logger().Error("Error getting Appointment", zap.Error(err))
		return aptmt, err
	}

	// Return the retrieved appointment details if successful
	return aptmt, nil
}
