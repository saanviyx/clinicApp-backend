package appointments

import (
	"clinic-app/cmd/rest/handler"
	"clinic-app/cmd/rest/middleware"
	"clinic-app/pkg/domain/errors"
	"clinic-app/pkg/domain/models"
	"clinic-app/pkg/services/factory"

	"go.uber.org/zap"
)

func (r *repo) GetPatientHistory(ftx factory.Service, patientId int) ([]models.Appointment, error) {
	// Start a new transaction
	tx, err := ftx.TransactionManager().Begin()
	if err != nil {
		// Log and return error if transaction start fails
		ftx.Logger().Error("Could not begin transaction", zap.Error(err))
		return nil, errors.ErrDatabase
	}
	ftx.Logger().Info("Transaction started for retrieving patient history for doctor")

	var aptmts []models.Appointment

	// Defer a rollback in case of failure
	defer func() {
		if err != nil {
			rollbackErr := ftx.TransactionManager().Rollback(tx)
			if rollbackErr != nil {
				// Log rollback failure
				ftx.Logger().Error("Failed to rollback transaction", zap.Error(rollbackErr))
			}
		}
	}()

	// Execute the query to retrieve patient history by patientId
	rows, err := tx.QueryContext(ftx.Context(), GetPatientAppointmentHistoryQuery, patientId)
	if err != nil {
		// Log and return error if query execution fails
		ftx.Logger().Error("Could not find patient history for doctor", zap.Error(err))
		return nil, errors.ErrNotFound
	}
	defer rows.Close()

	// Scan and collect the appointment records
	for rows.Next() {
		var aptmt models.Appointment
		if err := rows.Scan(
			&aptmt.AppointmentID,
			&aptmt.PatientID,
			&aptmt.DoctorName,
			&aptmt.PatientName,
			&aptmt.StartTime,
			&aptmt.EndTime,
			&aptmt.Status,
		); err != nil {
			return nil, err
		}
		aptmts = append(aptmts, aptmt)
	}

	// Commit the transaction if no errors occurred
	if err := ftx.TransactionManager().Commit(tx); err != nil {
		// Log and return error if commit fails
		ftx.Logger().Error("Could not commit transaction", zap.Error(err))
		return nil, errors.ErrDatabase
	}

	// Log successful retrieval of patient history
	ftx.Logger().Info("Successfully retrieved patient history for doctor",
		zap.Any("Patient History", aptmts),
		zap.Any("Patient ID", patientId),
	)
	middleware.GetTraceParentFromContext(ftx.Context())

	return aptmts, nil
}

func (r *repo) GetPatientAppointmentHistory(ftx factory.Service) ([]models.Appointment, error) {
	// Start a new transaction
	tx, err := ftx.TransactionManager().Begin()
	if err != nil {
		// Log and return error if transaction start fails
		ftx.Logger().Error("Could not begin transaction", zap.Error(err))
		return nil, errors.ErrDatabase
	}
	ftx.Logger().Info("Transaction started for retrieving patient history")

	var aptmts []models.Appointment

	// Get the user ID from the request context
	userId := handler.GetUserId()

	// Defer a rollback in case of failure
	defer func() {
		if err != nil {
			rollbackErr := ftx.TransactionManager().Rollback(tx)
			if rollbackErr != nil {
				// Log rollback failure
				ftx.Logger().Error("Failed to rollback transaction", zap.Error(rollbackErr))
			}
		}
	}()

	// Execute the query to retrieve patient history by userId
	rows, err := tx.QueryContext(ftx.Context(), GetPatientAppointmentHistoryQuery, userId)
	if err != nil {
		// Log and return error if query execution fails
		ftx.Logger().Error("Could not find patient history", zap.Error(err))
		return nil, errors.ErrNotFound
	}
	defer rows.Close()

	// Scan and collect the appointment records
	for rows.Next() {
		var aptmt models.Appointment
		if err := rows.Scan(
			&aptmt.AppointmentID,
			&aptmt.PatientID,
			&aptmt.DoctorName,
			&aptmt.PatientName,
			&aptmt.StartTime,
			&aptmt.EndTime,
			&aptmt.Status,
		); err != nil {
			return nil, err
		}
		aptmts = append(aptmts, aptmt)
	}

	// Commit the transaction if no errors occurred
	if err := ftx.TransactionManager().Commit(tx); err != nil {
		// Log and return error if commit fails
		ftx.Logger().Error("Could not commit transaction", zap.Error(err))
		return nil, errors.ErrDatabase
	}

	// Log successful retrieval of patient history
	ftx.Logger().Info("Successfully retrieved patient history",
		zap.Any("Patient", userId),
	)
	middleware.GetTraceParentFromContext(ftx.Context())

	return aptmts, nil
}
