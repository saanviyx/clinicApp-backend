package rest

import (
	"clinic-app/cmd/rest/handler"
	"clinic-app/cmd/rest/middleware"
	"clinic-app/pkg/usecase"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RestHandler defines the interface for setting up routes and configuring the router
type RestHandler interface {
	RegisterRoutes(router *gin.Engine)
	SetupRouter(logger *zap.Logger) *gin.Engine
}

// restHandler implements the RestHandler interface
type restHandler struct {
	authHandler        *handler.AuthHandler
	appointmentHandler *handler.AppointmentHandler
	doctorHandler      *handler.DoctorHandler
	adminHandler       *handler.AdminHandler
}

// NewRestHandler creates a new instance of restHandler with the provided use cases
func NewRestHandler(
	authUc usecase.AuthUsecase,
	aptmtUc usecase.AppointmentUsecase,
	docUc usecase.DoctorUsecase,
	adminUc usecase.AdminUsecase,
) RestHandler {
	return &restHandler{
		authHandler:        handler.NewAuthHandler(authUc),
		appointmentHandler: handler.NewAppointmentHandler(aptmtUc),
		doctorHandler:      handler.NewDoctorHandler(docUc),
		adminHandler:       handler.NewAdminHandler(adminUc),
	}
}

// RegisterRoutes sets up all routes for the application
func (h *restHandler) RegisterRoutes(router *gin.Engine) {
	// Authentication Routes
	authRoutes := router.Group("/")
	{
		authRoutes.POST("/register", h.authHandler.Register) // Register new user
		authRoutes.GET("/login", h.authHandler.Login)        // User login
	}

	// Appointment Routes
	appointmentRoutes := router.Group("/appointment")
	{
		appointmentRoutes.POST("/",
			middleware.AuthMiddleware("patient"), // Apply Authentication Middleware for patient role
			h.appointmentHandler.Book)            // Book an appointment

		appointmentRoutes.GET("/:id",
			middleware.AuthMiddleware("doctor", "patient"), // Apply Authentication Middleware for doctor and patient roles
			h.appointmentHandler.View)                      // View appointment details

		appointmentRoutes.GET("/:id/history",
			middleware.AuthMiddleware("doctor"),          // Apply Authentication Middleware for doctor role
			h.appointmentHandler.PatientHistoryForDoctor) // View patient history for doctor

		appointmentRoutes.GET("/history",
			middleware.AuthMiddleware("patient"), // Apply Authentication Middleware for patient role
			h.appointmentHandler.PatientHistory)  // View patient’s own appointment history

		appointmentRoutes.DELETE("/:id",
			middleware.AuthMiddleware("doctor", "admin"), // Apply Authentication Middleware for doctor and admin roles
			h.appointmentHandler.Cancel)                  // Cancel an appointment
	}

	// Doctor Routes
	doctorRoutes := router.Group("/doctors")
	{
		doctorRoutes.GET("/",
			middleware.AuthMiddleware("patient", "doctor", "admin"), // Apply Authentication Middleware for patient, doctor, and admin roles
			h.doctorHandler.ViewAll)                                 // View all doctors

		doctorRoutes.GET("/:id",
			middleware.AuthMiddleware("patient", "doctor", "admin"), // Apply Authentication Middleware for patient, doctor, and admin roles
			h.doctorHandler.ViewById)                                // View a specific doctor by ID

		doctorRoutes.GET("/:id/slots",
			middleware.AuthMiddleware("patient", "doctor", "admin"), // Apply Authentication Middleware for patient, doctor, and admin roles
			h.doctorHandler.Slots) // View available slots for a doctor
	}

	// Admin Routes
	adminRoutes := router.Group("/")
	{
		adminRoutes.GET("/doctors-availability",
			middleware.AuthMiddleware("admin"), // Apply Authentication Middleware for admin role
			h.adminHandler.Availability)        // View doctor availability

		adminRoutes.GET("/doctors-most-appointments",
			middleware.AuthMiddleware("admin"), // Apply Authentication Middleware for admin role
			h.adminHandler.MostAppointments)    // View doctors with the most appointments

		adminRoutes.GET("/doctors-over-6-hours",
			middleware.AuthMiddleware("admin"), // Apply Authentication Middleware for admin role
			h.adminHandler.OverSixHours)        // View doctors with over 6 hours of appointments
	}
}

// SetupRouter initializes the router and applies middleware
func (h *restHandler) SetupRouter(logger *zap.Logger) *gin.Engine {
	router := gin.Default() // Create a new Gin router
	gin.SetMode(gin.ReleaseMode)
	router.Use(middleware.TraceMiddleware(logger)) // Apply Trace Middleware for logging
	h.RegisterRoutes(router)                       // Register all routes
	return router
}
