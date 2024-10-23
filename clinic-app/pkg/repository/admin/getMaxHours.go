package admin

import (
	"clinic-app/cmd/rest/middleware"
	"clinic-app/pkg/domain/errors"
	"clinic-app/pkg/domain/models"
	"clinic-app/pkg/services/factory"

	"go.uber.org/zap"
)

func (r *repo) GetDoctorsWithOverSixHours(ftx factory.Service, date string) ([]models.DoctorOverTime, error) {
	// Start a new transaction
	tx, err := ftx.TransactionManager().Begin()
	if err != nil {
		// Log the error if the transaction could not be started
		ftx.Logger().Error("Could not begin transaction", zap.Error(err))
		return nil, errors.ErrDatabase
	}
	ftx.Logger().Info("Transaction started for retrieving all doctor availability")

	var daptmts []models.DoctorOverTime

	// Defer a rollback in case anything fails
	defer func() {
		if err != nil {
			rollbackErr := ftx.TransactionManager().Rollback(tx)
			if rollbackErr != nil {
				// Log error if rollback fails
				ftx.Logger().Error("Failed to rollback transaction", zap.Error(rollbackErr))
			}
		}
	}()

	// Execute the query to retrieve doctors with over six hours of work
	rows, err := tx.QueryContext(ftx.Context(), GetDoctorsWithOverSixHoursQuery, date)
	if err != nil {
		// Log the error if the query failed
		ftx.Logger().Error("Could not find doctor availability", zap.Error(err))
		return nil, errors.ErrNotFound
	}
	defer rows.Close() // Ensure rows are closed after processing

	// Iterate over the result set and scan each row into the model
	for rows.Next() {
		var daptmt models.DoctorOverTime
		if err := rows.Scan(
			&daptmt.DoctorID,
			&daptmt.DoctorName,
			&daptmt.DoctorEmail,
			&daptmt.TotalTime,
		); err != nil {
			return nil, err
		}
		daptmts = append(daptmts, daptmt) // Append each record to the slice
	}

	// Commit the transaction if no errors occurred
	if err := ftx.TransactionManager().Commit(tx); err != nil {
		// Log the error if the transaction could not be committed
		ftx.Logger().Error("Could not commit transaction", zap.Error(err))
		return nil, errors.ErrDatabase
	}

	ftx.Logger().Info("Successfully retrieved doctors' info for most hours")
	// Optionally, use the traceparent for logging or tracing purposes
	middleware.GetTraceParentFromContext(ftx.Context())

	return daptmts, nil
}
