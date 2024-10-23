package admin

import (
	"clinic-app/pkg/domain/models"
	"clinic-app/pkg/services/factory"

	"go.uber.org/zap"
)

// DoctorsAvailability retrieves the availability of all doctors.
func (uc *adminUsecaseImpl) DoctorsAvailability(ftx factory.Service) ([]models.DoctorAvailability, error) {
	// Call the repository method to get all doctors' availability
	davl, err := uc.repo.GetAllDoctorsAvailability(ftx)
	if err != nil {
		// Log an error if the repository method fails
		ftx.Logger().Error("Error getting Doctors Availability", zap.Error(err))
		return davl, err
	}

	// Return the list of doctor availability and a nil error
	return davl, nil
}
