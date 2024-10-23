package errors

import (
	"fmt"
	"net/http"
)

// AppError represents an application-level error with an associated HTTP status code
type ClinicAppError struct {
	Code    int    // HTTP status code
	Message string // Error message
}

// Error implements the error interface for AppError
func (e *ClinicAppError) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s", e.Code, e.Message)
}

// NewAppError creates a new AppError with a status code and message
func NewClinicAppError(code int, message string) *ClinicAppError {
	return &ClinicAppError{
		Code:    code,
		Message: message,
	}
}

var (
	ErrNotFound          = NewClinicAppError(http.StatusNotFound, "Resource not found")
	ErrBadRequest        = NewClinicAppError(http.StatusBadRequest, "Bad request")
	ErrDatabase          = NewClinicAppError(http.StatusInternalServerError, "Database error")
	ErrUserNotFound      = NewClinicAppError(http.StatusNotFound, "User not found")
	ErrInvalidPassword   = NewClinicAppError(http.StatusUnauthorized, "Invalid password")
	ErrUnauthorized      = NewClinicAppError(http.StatusUnauthorized, "Unauthorized access")
	ErrDuration          = NewClinicAppError(http.StatusNotAcceptable, "Invalid Appointment Duration")
	ErrAppointmentExists = NewClinicAppError(http.StatusBadRequest, "Appointment already exists for this time")
	ErrNoSchedule        = NewClinicAppError(http.StatusNotFound, "No schedule found for the doctor")
	ErrDoctorOverbooked  = NewClinicAppError(http.StatusNotAcceptable, "Doctor is overbooked")
)
