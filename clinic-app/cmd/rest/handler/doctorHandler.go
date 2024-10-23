package handler

import (
	"clinic-app/pkg/services/factory"
	"clinic-app/pkg/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// DoctorHandler struct holds the DoctorUsecase to manage doctor-related operations
type DoctorHandler struct {
	DocUsecase usecase.DoctorUsecase
}

// NewDoctorHandler initializes a new DoctorHandler with the provided usecase
func NewDoctorHandler(uc usecase.DoctorUsecase) *DoctorHandler {
	return &DoctorHandler{
		DocUsecase: uc,
	}
}

// ViewAll handles retrieving all doctors
func (h *DoctorHandler) ViewAll(c *gin.Context) {
	ftx := c.MustGet("ftx").(factory.Service) // Get service from context

	// Call usecase to get all doctors
	doctors, err := h.DocUsecase.AllDoctors(ftx)
	if err != nil {
		ftx.Logger().Error("Failed to retrieve all doctors", zap.Error(err))                     // Log error if retrieval fails
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve all doctors"}) // Return internal server error
		return
	}

	// Return list of doctors
	c.JSON(http.StatusOK, gin.H{"doctors": doctors})
}

// ViewById handles retrieving a doctor by their ID
func (h *DoctorHandler) ViewById(c *gin.Context) {
	ftx := c.MustGet("ftx").(factory.Service) // Get service from context

	// Get doctor ID from URL parameters and convert to integer
	doctorIdStr := c.Param("id")
	doctorId, err := strconv.Atoi(doctorIdStr)
	if err != nil {
		ftx.Logger().Error("Invalid doctor ID", zap.Error(err))            // Log error for invalid ID
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid doctor ID"}) // Return bad request error
		return
	}

	// Call usecase to get doctor by ID
	doctor, err := h.DocUsecase.DoctorById(ftx, doctorId)
	if err != nil {
		ftx.Logger().Error("Failed to retrieve doctor by ID", zap.Error(err))                     // Log error if retrieval fails
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve doctor by ID"}) // Return internal server error
		return
	}

	// Return the doctor details
	c.JSON(http.StatusOK, gin.H{"doctor": doctor})
}

// Slots handles retrieving available slots for a specific doctor
func (h *DoctorHandler) Slots(c *gin.Context) {
	ftx := c.MustGet("ftx").(factory.Service) // Get service from context

	// Get doctor ID from URL parameters and convert to integer
	doctorIdStr := c.Param("id")
	doctorId, err := strconv.Atoi(doctorIdStr)
	if err != nil {
		ftx.Logger().Error("Invalid doctor ID", zap.Error(err))            // Log error for invalid ID
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid doctor ID"}) // Return bad request error
		return
	}

	// Check user role to determine if the user is a doctor or admin
	var role = GetRole()
	check := false
	if role == "doctor" || role == "admin" {
		check = true // Allow specific operations based on the role
	} else {
		check = false
	}

	// Call usecase to get available slots for the doctor
	slots, err := h.DocUsecase.Slots(ftx, doctorId, check)
	if err != nil {
		ftx.Logger().Error("Failed to retrieve slots for doctor", zap.Error(err))                     // Log error if retrieval fails
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve slots for doctor"}) // Return internal server error
		return
	}

	// If no slots are available, return a specific message
	if slots == nil {
		c.JSON(http.StatusOK, gin.H{"message": "No appointments made, all slots available"})
		return
	}

	// Return the available slots
	c.JSON(http.StatusOK, gin.H{"slots": slots})
}
