package doctor

import (
	"clinic-app/pkg/domain/models"
	"clinic-app/pkg/services/factory"

	"go.uber.org/zap"
)

// AllDoctors retrieves a list of all doctors.
func (uc *doctorsUsecaseImpl) AllDoctors(ftx factory.Service) ([]models.Doctor, error) {
	// Call the repository method to get all doctors
	doctor, err := uc.repo.GetAllDoctors(ftx)
	if err != nil {
		// Log an error if getting all doctors fails
		ftx.Logger().Error("Error getting all Doctor Information", zap.Error(err))
		return doctor, err
	}
	// Return the list of doctors if successful
	return doctor, nil
}

// DoctorById retrieves a specific doctor by their ID.
func (uc *doctorsUsecaseImpl) DoctorById(ftx factory.Service, doctorId int) (models.Doctor, error) {
	// Call the repository method to get a doctor by ID
	doctor, err := uc.repo.GetDoctorById(ftx, doctorId)
	if err != nil {
		// Log an error if getting the doctor by ID fails
		ftx.Logger().Error("Error getting Doctor Information by ID", zap.Error(err))
		return doctor, err
	}
	// Return the doctor if successful
	return doctor, nil
}
