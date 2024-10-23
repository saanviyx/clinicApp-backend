package appointments

import (
	"clinic-app/cmd/rest/handler"
	"clinic-app/cmd/rest/middleware"
	"clinic-app/pkg/domain/errors"
	"clinic-app/pkg/domain/models"
	"clinic-app/pkg/services/factory"

	"go.uber.org/zap"
)

func (r *repo) BookAppointment(ftx factory.Service, aptmt models.BookAppointment) error {
	// Start a new transaction
	tx, err := ftx.TransactionManager().Begin()
	if err != nil {
		// Log and return error if transaction start fails
		ftx.Logger().Error("Could not begin transaction", zap.Error(err))
		return errors.ErrDatabase
	}
	ftx.Logger().Info("Transaction started for Booking Appointment")

	// Retrieve the user ID from the request context
	userId := handler.GetUserId()
	var appointmentID any
	var result string

	// Defer a rollback in case of any errors
	defer func() {
		if err != nil {
			rollbackErr := ftx.TransactionManager().Rollback(tx)
			if rollbackErr != nil {
				// Log rollback failure
				ftx.Logger().Error("Failed to rollback transaction", zap.Error(rollbackErr))
			}
		}
	}()

	// Execute the query to book an appointment
	err = tx.QueryRowContext(ftx.Context(),
		BookAppointmentQuery,
		aptmt.DoctorID,
		userId,
		aptmt.Date,
		aptmt.StartTime,
		aptmt.EndTime,
	).Scan(&appointmentID, &result)

	if err != nil {
		// Log and return error if query execution fails
		ftx.Logger().Error("Error executing query", zap.Error(err))
		return errors.ErrDatabase
	}

	// Commit the transaction if no errors occurred
	if err := ftx.TransactionManager().Commit(tx); err != nil {
		// Log and return error if commit fails
		ftx.Logger().Error("Could not commit transaction", zap.Error(err))
		return errors.ErrDatabase
	}

	// Handle different results from the query
	switch result {
	case "Valid":
		// Log success and return nil if the appointment was booked successfully
		ftx.Logger().Info("Successfully Booked Appointment",
			zap.Any("Appointment", appointmentID),
		)
		return nil

	case "Appointment Exists":
		// Log and return error if the appointment already exists
		ftx.Logger().Info("Appointment already exists", zap.String("result", result))
		return errors.ErrAppointmentExists

	case "Schedule Not found":
		// Log and return error if the schedule could not be found
		ftx.Logger().Info("Schedule not found", zap.String("result", result))
		return errors.ErrNoSchedule

	case "Doctor Overbooked":
		// Log and return error if the doctor is overbooked
		ftx.Logger().Info("Doctor is overbooked", zap.String("result", result))
		return errors.ErrDoctorOverbooked
	}

	// Handle any unexpected results
	middleware.GetTraceParentFromContext(ftx.Context())

	return nil
}
