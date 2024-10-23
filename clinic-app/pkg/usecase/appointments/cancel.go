package appointments

import (
	"clinic-app/pkg/services/factory"

	"go.uber.org/zap"
)

// Cancel removes an existing appointment based on its ID.
func (uc *aptmtUsecaseImpl) Cancel(ftx factory.Service, aptmtId int) error {
	// Call the repository method to cancel the appointment
	err := uc.repo.CancelAppointment(ftx, aptmtId)
	if err != nil {
		// Log an error if the cancellation fails
		ftx.Logger().Error("Error cancelling appointment", zap.Error(err))
		return err
	}

	// Return nil if no error occurred
	return nil
}
