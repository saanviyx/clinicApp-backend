package admin

import (
	"clinic-app/cmd/rest/middleware"
	"clinic-app/pkg/domain/errors"
	"clinic-app/pkg/domain/models"
	"clinic-app/pkg/services/factory"

	"go.uber.org/zap"
)

func (r *repo) GetAllDoctorsAvailability(ftx factory.Service) ([]models.DoctorAvailability, error) {
	// Start a new transaction
	tx, err := ftx.TransactionManager().Begin()
	if err != nil {
		ftx.Logger().Error("Could not begin transaction", zap.Error(err))
		return nil, errors.ErrDatabase
	}
	ftx.Logger().Info("Transaction started for retrieving doctors' availability")

	var avls []models.DoctorAvailability

	// Defer a rollback in case anything fails during the transaction
	defer func() {
		if err != nil {
			rollbackErr := ftx.TransactionManager().Rollback(tx)
			if rollbackErr != nil {
				ftx.Logger().Error("Failed to rollback transaction", zap.Error(rollbackErr))
			}
		}
	}()

	// Execute the query to retrieve doctor availability data
	rows, err := tx.QueryContext(ftx.Context(), GetAllDoctorsAvailabilityQuery)
	if err != nil {
		ftx.Logger().Error("Could not retrieve doctor availability", zap.Error(err))
		return nil, errors.ErrNotFound
	}
	defer rows.Close() // Ensure rows are closed after processing

	// Iterate over the result set and scan each row into the model
	for rows.Next() {
		var avl models.DoctorAvailability
		if err := rows.Scan(
			&avl.DoctorID,
			&avl.DoctorName,
			&avl.DoctorEmail,
			&avl.Date,
			&avl.TotalAppointments,
			&avl.TotalTime,
			&avl.Availability,
		); err != nil {
			return nil, err
		}
		avls = append(avls, avl) // Append each availability record to the slice
	}

	// Commit the transaction if no errors occurred
	if err := ftx.TransactionManager().Commit(tx); err != nil {
		ftx.Logger().Error("Could not commit transaction", zap.Error(err))
		return nil, errors.ErrDatabase
	}

	ftx.Logger().Info("Successfully retrieved doctor availability",
		zap.Any("Doctors Availability", avls),
	)

	// Optionally, you might use the traceparent here for logging or tracing
	middleware.GetTraceParentFromContext(ftx.Context())

	return avls, nil
}
