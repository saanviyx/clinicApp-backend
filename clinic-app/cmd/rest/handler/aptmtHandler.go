package handler

import (
	"clinic-app/pkg/domain/errors"
	"clinic-app/pkg/domain/models"
	"clinic-app/pkg/services/factory"
	"clinic-app/pkg/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AppointmentHandler struct holds the AppointmentUsecase
type AppointmentHandler struct {
	AptmtUsecase usecase.AppointmentUsecase // Use case for appointment operations
}

// NewAppointmentHandler creates a new instance of AppointmentHandler
func NewAppointmentHandler(uc usecase.AppointmentUsecase) *AppointmentHandler {
	return &AppointmentHandler{
		AptmtUsecase: uc, // Assigning use case to handler
	}
}

// Book handles booking an appointment
func (h *AppointmentHandler) Book(c *gin.Context) {
	ftx := c.MustGet("ftx").(factory.Service) // Extract service from context

	var aptmt models.BookAppointment
	if err := c.ShouldBindJSON(&aptmt); err != nil { // Bind JSON input to aptmt model
		ftx.Logger().Error("Invalid input", zap.Error(err))            // Log error if JSON binding fails
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"}) // Return bad request error
		return
	}

	err := h.AptmtUsecase.Book(ftx, aptmt) // Call use case to book appointment

	switch err {
	case errors.ErrDatabase:
		ftx.Logger().Error("Booking failed", zap.Error(err))                     // Log database error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Booking failed"}) // Return internal server error

	case errors.ErrAppointmentExists:
		ftx.Logger().Info("Appointment already exists for this time slot")                             // Log appointment exists error
		c.JSON(http.StatusConflict, gin.H{"message": "Appointment already exists for this time slot"}) // Return conflict error

	case errors.ErrDuration:
		ftx.Logger().Info("Appointment duration is too long")                                                                          // Log invalid duration error
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "Appointment duration is invalid. Minimum - 15 minutes, Maximum - 2 hours"}) // Return not acceptable error

	case errors.ErrNoSchedule:
		ftx.Logger().Info("Schedule not found for Doctor")                                  // Log no schedule error
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "Schedule not found for Doctor"}) // Return not acceptable error

	case errors.ErrDoctorOverbooked:
		ftx.Logger().Info("All Appointments Booked")                                                               // Log overbooked error
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "Cannot Schedule Appointment. All Appointments Booked"}) // Return not acceptable error

	default:
		if err != nil {
			ftx.Logger().Error("Unknown error occurred", zap.Error(err))                        // Log unknown error
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An unknown error occurred"}) // Return internal server error
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "Appointment booked successfully"}) // Return success message
		}
	}
}

// View handles retrieving an appointment by its ID
func (h *AppointmentHandler) View(c *gin.Context) {
	ftx := c.MustGet("ftx").(factory.Service) // Extract service from context

	appointmentID, err := strconv.Atoi(c.Param("id")) // Convert appointment ID from string to integer
	if err != nil {
		ftx.Logger().Error("Invalid appointment ID", zap.Error(err))            // Log invalid ID error
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid appointment ID"}) // Return bad request error
		return
	}

	appointment, err := h.AptmtUsecase.ViewAppointment(ftx, appointmentID) // Call use case to view appointment
	if err == errors.ErrNotFound {
		c.JSON(http.StatusOK, gin.H{"message": "Appointment does not exist"}) // Return appointment not found
		return
	} else if err != nil {
		ftx.Logger().Error("Failed to retrieve appointment", zap.Error(err))                     // Log retrieval error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve appointment"}) // Return internal server error
		return
	}

	c.JSON(http.StatusOK, gin.H{"appointment": appointment}) // Return appointment details
}

// PatientHistoryForDoctor handles retrieving a patient's appointment history for a specific doctor
func (h *AppointmentHandler) PatientHistoryForDoctor(c *gin.Context) {
	ftx := c.MustGet("ftx").(factory.Service) // Extract service from context

	patientIDStr := c.Param("id")                // Get patient ID from URL parameter
	patientID, err := strconv.Atoi(patientIDStr) // Convert patient ID to integer
	if err != nil {
		ftx.Logger().Error("Invalid patient ID", zap.Error(err))            // Log invalid patient ID error
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid patient ID"}) // Return bad request error
		return
	}

	appointments, err := h.AptmtUsecase.PatientHistoryForDoctor(ftx, patientID) // Call use case to retrieve patient history for doctor
	if err != nil {
		ftx.Logger().Error("Failed to retrieve patient history", zap.Error(err))                     // Log retrieval error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve patient history"}) // Return internal server error
		return
	}

	if appointments == nil {
		c.JSON(http.StatusOK, gin.H{"message": "Patient does not exist"}) // Return message if patient doesn't exist
		return
	}

	c.JSON(http.StatusOK, gin.H{"appointments": appointments}) // Return patient history for doctor
}

// PatientHistory handles retrieving a patient's full appointment history
func (h *AppointmentHandler) PatientHistory(c *gin.Context) {
	ftx := c.MustGet("ftx").(factory.Service) // Extract service from context

	appointments, err := h.AptmtUsecase.PatientHistory(ftx) // Call use case to retrieve patient history
	if err != nil {
		ftx.Logger().Error("Failed to retrieve patient history", zap.Error(err))                     // Log retrieval error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve patient history"}) // Return internal server error
		return
	}

	if appointments == nil {
		c.JSON(http.StatusOK, gin.H{"message": "Patient does not exist"}) // Return message if patient doesn't exist
		return
	}

	c.JSON(http.StatusOK, gin.H{"appointments": appointments}) // Return patient history
}

// Cancel handles canceling an appointment
func (h *AppointmentHandler) Cancel(c *gin.Context) {
	ftx := c.MustGet("ftx").(factory.Service) // Extract service from context

	appointmentID, err := strconv.Atoi(c.Param("id")) // Convert appointment ID from string to integer
	if err != nil {
		ftx.Logger().Error("Invalid appointment ID", zap.Error(err))            // Log invalid appointment ID error
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid appointment ID"}) // Return bad request error
		return
	}

	err = h.AptmtUsecase.Cancel(ftx, appointmentID) // Call use case to cancel appointment
	if err != nil {
		ftx.Logger().Error("Cancellation failed", zap.Error(err))                     // Log cancellation error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cancellation failed"}) // Return internal server error
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Appointment canceled successfully"}) // Return success message
}
