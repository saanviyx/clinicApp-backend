package handler

import (
	"clinic-app/pkg/domain/models"
	"clinic-app/pkg/services/factory"
	"clinic-app/pkg/usecase"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
)

// Secret key for JWT
var jwtKey = []byte("your_secret_key")

// Variables to store user role and ID globally
var role = ""
var UserID = 0

// AuthHandler struct holds the AuthUsecase to handle authentication
type AuthHandler struct {
	AuthUsecase usecase.AuthUsecase
}

// NewAuthHandler initializes a new AuthHandler with the provided usecase
func NewAuthHandler(uc usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		AuthUsecase: uc,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	ftx := c.MustGet("ftx").(factory.Service) // Get service from context
	var user models.User

	// Bind incoming JSON to user struct
	if err := c.ShouldBindJSON(&user); err != nil {
		ftx.Logger().Error("Invalid input", zap.Error(err))            // Log error
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"}) // Return bad request
		return
	}

	// Call usecase to register the user
	userID, err := h.AuthUsecase.RegisterUser(ftx, user)
	if err != nil {
		ftx.Logger().Error("Registration failed", zap.Error(err))                     // Log registration failure
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Registration failed"}) // Return internal server error
		return
	}

	// Convert userID to string for logging
	userIDStr := strconv.Itoa(userID)

	// Log registration success
	ftx.Logger().Info("Registration successful", zap.String("UserID:", userIDStr))
	c.JSON(http.StatusOK, gin.H{"message": "Registration successful"}) // Return success response
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	ftx := c.MustGet("ftx").(factory.Service) // Get service from context
	var credentials models.Credentials

	// Bind incoming JSON to credentials struct
	if err := c.ShouldBindJSON(&credentials); err != nil {
		ftx.Logger().Error("Invalid input", zap.Error(err))            // Log error
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"}) // Return bad request
		return
	}

	// Call usecase to login user
	user, err := h.AuthUsecase.LoginUser(ftx, credentials)
	if err != nil {
		ftx.Logger().Error("Login failed", zap.Error(err))                     // Log login failure
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"}) // Return unauthorized error
		return
	}

	// Set role and user ID globally
	role = user.Role
	UserID = user.ID

	// Create JWT token with expiration time
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &models.Claims{
		UserID: user.ID,
		Role:   user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(), // Token expiration
		},
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)                         // Sign token with secret key
	c.SetCookie("token", tokenString, 3600, "/", "localhost", false, true) // Set token in cookie

	// Add token to Authorization header
	c.Request.Header.Add("Authorization", tokenString)
	if err != nil {
		ftx.Logger().Error("Failed to create JWT token", zap.Error(err))           // Log token creation failure
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not log in"}) // Return internal server error
		return
	}

	// Return login success response
	c.JSON(http.StatusOK, gin.H{"message": "Login Successful"})
}

// GetRole returns the current role of the user
func GetRole() string {
	return role
}

// GetUserId returns the current user ID
func GetUserId() int {
	return UserID
}
