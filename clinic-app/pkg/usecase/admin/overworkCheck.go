package admin

import (
	"clinic-app/pkg/domain/models"
	"clinic-app/pkg/services/factory"

	"go.uber.org/zap"
)

// MostAppointments retrieves doctors with the most appointments on a specific date.
func (uc *adminUsecaseImpl) MostAppointments(ftx factory.Service, date string) ([]models.DoctorMostAppointments, error) {
	// Call the repository method to get doctors with the most appointments on the specified date
	daptmt, err := uc.repo.GetDoctorsWithMostAppointments(ftx, date)
	if err != nil {
		// Log an error if the repository method fails
		ftx.Logger().Error("Error getting Doctors with Most Appointment", zap.Error(err))
		return daptmt, err
	}

	// Return the list of doctors with the most appointments and a nil error
	return daptmt, nil
}

// OverSixHours retrieves doctors who have worked over six hours on a specific date.
func (uc *adminUsecaseImpl) OverSixHours(ftx factory.Service, date string) ([]models.DoctorOverTime, error) {
	// Call the repository method to get doctors who have worked over six hours on the specified date
	daptmt, err := uc.repo.GetDoctorsWithOverSixHours(ftx, date)
	if err != nil {
		// Log an error if the repository method fails
		ftx.Logger().Error("Error getting Doctors with Over 6 Hours", zap.Error(err))
		return daptmt, err
	}

	// Return the list of doctors who have worked over six hours and a nil error
	return daptmt, nil
}
