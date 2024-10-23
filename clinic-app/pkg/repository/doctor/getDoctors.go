package doctor

import (
	"clinic-app/cmd/rest/middleware"
	"clinic-app/pkg/domain/errors"
	"clinic-app/pkg/domain/models"
	"clinic-app/pkg/services/factory"
	"database/sql"

	"go.uber.org/zap"
)

// GetAllDoctors retrieves a list of all doctors
func (r *repo) GetAllDoctors(ftx factory.Service) ([]models.Doctor, error) {
	// Start a new transaction
	tx, err := ftx.TransactionManager().Begin()
	if err != nil {
		ftx.Logger().Error("Could not begin transaction", zap.Error(err))
		return nil, errors.ErrDatabase
	}
	ftx.Logger().Info("Transaction started for retrieving doctors")

	var doctors []models.Doctor

	// Defer a rollback in case anything fails
	defer func() {
		if err != nil {
			rollbackErr := ftx.TransactionManager().Rollback(tx)
			if rollbackErr != nil {
				ftx.Logger().Error("Failed to rollback transaction", zap.Error(rollbackErr))
			}
		}
	}()

	// Execute the query to get all doctors
	rows, err := tx.QueryContext(ftx.Context(), GetAllDoctorsQuery)
	if err != nil {
		ftx.Logger().Error("Could not find Doctor", zap.Error(err))
		return nil, errors.ErrNotFound
	}
	defer rows.Close()

	// Scan the results into doctor models
	for rows.Next() {
		var doctor models.Doctor
		if err := rows.Scan(
			&doctor.ID,
			&doctor.Name,
			&doctor.Email,
			&doctor.Availability,
		); err != nil {
			ftx.Logger().Error("Error scanning doctor row", zap.Error(err))
			return nil, errors.ErrNotFound
		}
		doctors = append(doctors, doctor)
	}

	// Commit the transaction if no errors occurred
	if err := ftx.TransactionManager().Commit(tx); err != nil {
		ftx.Logger().Error("Could not commit transaction", zap.Error(err))
		return nil, errors.ErrDatabase
	}

	// Log the successfully retrieved doctors
	ftx.Logger().Info("Successfully retrieved all doctors",
		zap.Any("Doctors", doctors),
	)
	middleware.GetTraceParentFromContext(ftx.Context())

	return doctors, nil
}

// GetDoctorById retrieves a specific doctor by ID
func (r *repo) GetDoctorById(ftx factory.Service, doctorId int) (models.Doctor, error) {
	// Start a new transaction
	tx, err := ftx.TransactionManager().Begin()
	if err != nil {
		ftx.Logger().Error("Could not begin transaction", zap.Error(err))
		return models.Doctor{}, errors.ErrDatabase
	}
	ftx.Logger().Info("Transaction started for retrieving doctor by Id")

	// Execute the query to get a doctor by ID
	row := tx.QueryRowContext(ftx.Context(), GetDoctorByIdQuery, doctorId)

	var doctor models.Doctor
	err = row.Scan(
		&doctor.ID,
		&doctor.Name,
		&doctor.Email,
		&doctor.Availability,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Doctor{}, errors.ErrNotFound
		}
		return models.Doctor{}, err
	}

	// Log the successfully retrieved doctor
	ftx.Logger().Info("Successfully retrieved doctor",
		zap.Any("Doctor", doctor),
	)
	middleware.GetTraceParentFromContext(ftx.Context())

	return doctor, nil
}
