package appointments

import (
	"clinic-app/pkg/domain/models"
	"clinic-app/pkg/services/factory"

	"go.uber.org/zap"
)

// PatientHistoryForDoctor retrieves appointment history for a specific patient.
func (uc *aptmtUsecaseImpl) PatientHistoryForDoctor(ftx factory.Service, patientId int) ([]models.Appointment, error) {
	// Call the repository method to get the patient's appointment history
	paptmt, err := uc.repo.GetPatientHistory(ftx, patientId)
	if err != nil {
		// Log an error if fetching the appointment history fails
		ftx.Logger().Error("Error getting Patient Appointment History for Doctor", zap.Error(err))
		return paptmt, err
	}

	// Return the retrieved appointment history if successful
	return paptmt, nil
}

// PatientHistory retrieves all appointment history for the current patient.
func (uc *aptmtUsecaseImpl) PatientHistory(ftx factory.Service) ([]models.Appointment, error) {
	// Call the repository method to get the appointment history for the current patient
	paptmt, err := uc.repo.GetPatientAppointmentHistory(ftx)
	if err != nil {
		// Log an error if fetching the appointment history fails
		ftx.Logger().Error("Error getting Patient Appointment History", zap.Error(err))
		return paptmt, err
	}

	// Return the retrieved appointment history if successful
	return paptmt, nil
}
