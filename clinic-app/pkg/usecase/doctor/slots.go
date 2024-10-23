package doctor

import (
	"clinic-app/pkg/services/factory"

	"go.uber.org/zap"
)

// Slots retrieves available time slots for a specific doctor.
// If `check` is true, retrieves slots by doctor, otherwise retrieves slots by patient.
func (uc *doctorsUsecaseImpl) Slots(ftx factory.Service, doctorId int, check bool) ([]interface{}, error) {
	// Call the repository method to get doctor slots
	slots, err := uc.repo.GetDoctorSlots(ftx, doctorId, check)
	if err != nil {
		// Log an error if getting doctor slots fails
		ftx.Logger().Error("Error getting Doctor Slot", zap.Error(err))
		return slots, err
	}
	// Return the available slots if successful
	return slots, nil
}
