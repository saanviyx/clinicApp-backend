package handler

import (
	"clinic-app/pkg/services/factory"
	"clinic-app/pkg/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AdminHandler struct holds the AdminUsecase
type AdminHandler struct {
	AdminUsecase usecase.AdminUsecase // Use case for admin operations
}

// NewAdminHandler creates a new instance of AdminHandler
func NewAdminHandler(uc usecase.AdminUsecase) *AdminHandler {
	return &AdminHandler{
		AdminUsecase: uc, // Assigning use case to handler
	}
}

// Availability handles retrieving availability for all doctors
func (h *AdminHandler) Availability(c *gin.Context) {
	ftx := c.MustGet("ftx").(factory.Service) // Extract service from context

	doctorsAvailability, err := h.AdminUsecase.DoctorsAvailability(ftx) // Call use case to get availability
	if err != nil {
		ftx.Logger().Error("Failed to retrieve doctors availability", zap.Error(err))                     // Log error if any
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve doctors availability"}) // Return error response
		return
	}

	c.JSON(http.StatusOK, gin.H{"doctors_availability": doctorsAvailability}) // Return availability response
}

// MostAppointments handles retrieving doctors with the most appointments for a given date
func (h *AdminHandler) MostAppointments(c *gin.Context) {
	ftx := c.MustGet("ftx").(factory.Service)                              // Extract service from context
	date := c.Query("date")                                                // Get date from query parameters
	doctorsAppointments, err := h.AdminUsecase.MostAppointments(ftx, date) // Call use case to get doctors with most appointments
	if err != nil {
		ftx.Logger().Error("Failed to retrieve doctors with most appointments", zap.Error(err))                     // Log error if any
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve doctors with most appointments"}) // Return error response
		return
	}

	c.JSON(http.StatusOK, gin.H{"doctors_appointments": doctorsAppointments}) // Return doctors appointments response
}

// OverSixHours handles retrieving doctors with more than six hours of appointments for a given date
func (h *AdminHandler) OverSixHours(c *gin.Context) {
	ftx := c.MustGet("ftx").(factory.Service) // Extract service from context
	date := c.Query("date")                   // Get date from query parameters

	doctorsAppointments, err := h.AdminUsecase.OverSixHours(ftx, date) // Call use case to get doctors with over six hours of appointments
	if err != nil {
		ftx.Logger().Error("Failed to retrieve doctors with over six hours of appointments", zap.Error(err))                     // Log error if any
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve doctors with over six hours of appointments"}) // Return error response
		return
	}

	if doctorsAppointments == nil {
		c.JSON(http.StatusOK, gin.H{"message": "No Doctors working over 6 hours"}) // Return message if patient doesn't exist
		return
	}

	c.JSON(http.StatusOK, gin.H{"doctors_appointments": doctorsAppointments}) // Return doctors appointments response
}
