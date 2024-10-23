package middleware

import (
	"clinic-app/pkg/domain/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// AuthMiddleware checks if the user is authenticated and authorized
func AuthMiddleware(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the token from cookies
		tokenString, err := c.Cookie("token")
		if err != nil {
			fmt.Println("Token missing in cookie") // Log missing token
			c.Abort()                              // Abort the request if token is not found
			return
		}

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"}) // Respond with unauthorized if token is empty
			c.Abort()                                                        // Abort the request
			return
		}

		// Parse and validate the JWT token
		claims := &models.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("your_secret_key"), nil // Return the secret key for token validation
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"}) // Respond with unauthorized if token is invalid
			c.Abort()                                                        // Abort the request
			return
		}

		// Check if the user's role is one of the required roles
		roleAuthorized := false
		for _, role := range requiredRoles {
			if claims.Role == role {
				roleAuthorized = true // User role is authorized
				break
			}
		}

		if !roleAuthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to access this resource"}) // Respond with forbidden if user is not authorized
			c.Abort()                                                                                         // Abort the request
			return
		}

		// Set userID and userRole in the context for further use
		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)
		c.Next() // Proceed to the next handler
	}
}
