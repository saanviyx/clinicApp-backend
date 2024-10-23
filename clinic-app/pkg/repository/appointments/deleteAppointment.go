package appointments

import (
	"clinic-app/cmd/rest/middleware"
	"clinic-app/pkg/domain/errors"
	"clinic-app/pkg/services/factory"

	"go.uber.org/zap"
)

func (r *repo) CancelAppointment(ftx factory.Service, appointmentId int) error {
	// Start a new transaction
	tx, err := ftx.TransactionManager().Begin()
	if err != nil {
		// Log the error if the transaction could not be started
		ftx.Logger().Error("Could not begin transaction", zap.Error(err))
		return errors.ErrDatabase
	}
	ftx.Logger().Info("Transaction started for deleting appointment")

	// Defer a rollback in case anything fails
	defer func() {
		if err != nil {
			rollbackErr := ftx.TransactionManager().Rollback(tx)
			if rollbackErr != nil {
				// Log the error if rollback fails
				ftx.Logger().Error("Failed to rollback transaction", zap.Error(rollbackErr))
			}
		}
	}()

	// Execute query to cancel the appointment
	_, err = tx.ExecContext(ftx.Context(), CancelAppointmentQuery, appointmentId)
	if err != nil {
		// Log the error if the appointment could not be found or deleted
		ftx.Logger().Error("Could not find Appointment", zap.Error(err))
		return errors.ErrNotFound
	}

	// Execute query to delete the slot associated with the appointment
	_, err = tx.ExecContext(ftx.Context(), DeleteSlotQuery, appointmentId)
	if err != nil {
		// Log the error if the slot could not be found or deleted
		ftx.Logger().Error("Could not find Appointment", zap.Error(err))
		return errors.ErrNotFound
	}

	// Commit the transaction if no errors occurred
	if err := ftx.TransactionManager().Commit(tx); err != nil {
		// Log the error if the transaction could not be committed
		ftx.Logger().Error("Could not commit transaction", zap.Error(err))
		return errors.ErrDatabase
	}

	ftx.Logger().Info("Successfully cancelled appointment")
	// Optionally, use the traceparent for logging or tracing purposes
	middleware.GetTraceParentFromContext(ftx.Context())

	return nil
}
