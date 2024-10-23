package appointments

import (
	"clinic-app/cmd/rest/handler"
	"clinic-app/cmd/rest/middleware"
	"clinic-app/pkg/domain/errors"
	"clinic-app/pkg/domain/models"
	"clinic-app/pkg/services/factory"
	"database/sql"

	"go.uber.org/zap"
)

func (r *repo) GetAppointmentById(ftx factory.Service, aptmtID int) (models.Appointment, error) {
	// Start a new transaction
	tx, err := ftx.TransactionManager().Begin()
	if err != nil {
		// Log the error if the transaction could not be started
		ftx.Logger().Error("Could not begin transaction", zap.Error(err))
		return models.Appointment{}, errors.ErrDatabase
	}
	ftx.Logger().Info("Transaction started for retrieving appointment")

	var aptmt models.Appointment

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

	// Get the user ID from the request context
	var userID = handler.GetUserId()

	// Query to get the appointment details based on appointment ID and user ID
	row := tx.QueryRowContext(ftx.Context(), GetAppointmentByIdQuery, aptmtID, userID)

	// Scan the result into the Appointment struct
	err = row.Scan(
		&aptmt.AppointmentID,
		&aptmt.PatientID,
		&aptmt.PatientName,
		&aptmt.DoctorName,
		&aptmt.StartTime,
		&aptmt.EndTime,
		&aptmt.Status,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			// Return a not found error if no rows are returned
			return models.Appointment{}, errors.ErrNotFound
		}
		// Return any other errors that occurred during the scan
		return models.Appointment{}, err
	}

	// Log successful retrieval of appointment
	ftx.Logger().Info("Successfully retrieved appointment",
		zap.Any("Appointment", aptmt),
	)
	// Optionally, use the traceparent for logging or tracing purposes
	middleware.GetTraceParentFromContext(ftx.Context())

	return aptmt, nil
}
