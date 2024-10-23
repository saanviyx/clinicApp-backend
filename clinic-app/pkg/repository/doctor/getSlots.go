package doctor

import (
	"clinic-app/cmd/rest/middleware"
	"clinic-app/pkg/domain/errors"
	"clinic-app/pkg/domain/models"
	"clinic-app/pkg/services/factory"

	"go.uber.org/zap"
)

// GetDoctorSlots retrieves available time slots for a specific doctor
func (r *repo) GetDoctorSlots(ftx factory.Service, doctorId int, check bool) ([]interface{}, error) {
	// Start a new transaction
	tx, err := ftx.TransactionManager().Begin()
	if err != nil {
		ftx.Logger().Error("Could not begin transaction", zap.Error(err))
		return nil, errors.ErrDatabase
	}
	ftx.Logger().Info("Transaction started for getting slots")

	var slots []interface{}

	// Defer a rollback in case of any errors
	defer func() {
		if err != nil {
			rollbackErr := ftx.TransactionManager().Rollback(tx)
			if rollbackErr != nil {
				ftx.Logger().Error("Failed to rollback transaction", zap.Error(rollbackErr))
			}
		}
	}()

	// Determine the query based on the 'check' parameter
	var query string
	if check {
		query = GetSlotsByDoctorQuery
	} else {
		query = GetSlotsByPatientQuery
	}

	// Execute the query to get the time slots
	rows, err := tx.QueryContext(ftx.Context(), query, doctorId)
	if err != nil {
		ftx.Logger().Error("Query failed", zap.Error(err))
		return nil, errors.ErrDatabase
	}
	defer rows.Close()

	// Scan the results into the appropriate slot models
	for rows.Next() {
		if check {
			var slot models.SlotDoc
			if err = rows.Scan(
				&slot.SlotID,
				&slot.AppointmentID,
				&slot.PatientID,
				&slot.PatientName,
				&slot.StartTime,
				&slot.EndTime,
				&slot.IsBooked,
				&slot.Duration,
			); err != nil {
				return nil, err
			}
			slots = append(slots, slot)
		} else {
			var slot models.SlotPat
			if err := rows.Scan(
				&slot.SlotID,
				&slot.AppointmentID,
				&slot.StartTime,
				&slot.EndTime,
				&slot.IsBooked,
				&slot.Duration,
			); err != nil {
				return nil, err
			}
			slots = append(slots, slot)
		}
	}

	// Commit the transaction if no errors occurred
	if err := ftx.TransactionManager().Commit(tx); err != nil {
		ftx.Logger().Error("Could not commit transaction", zap.Error(err))
		return nil, errors.ErrDatabase
	}

	// Log the successfully retrieved slots
	ftx.Logger().Info("Successfully retrieved slots",
		zap.Any("Slots", slots),
	)
	middleware.GetTraceParentFromContext(ftx.Context())

	return slots, nil
}
