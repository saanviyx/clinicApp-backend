package middleware

import (
	"clinic-app/internal/constants"
	"clinic-app/pkg/services/factory"
	"context"
	"fmt"
	"math/rand/v2"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// TraceMiddleware generates and logs a traceparent header for tracing requests
func TraceMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate a traceparent value for tracing
		traceparent := GenerateTraceParent()

		// Add the traceparent value to the request context
		c.Request = c.Request.WithContext(AddTraceParentToContext(c.Request.Context(), traceparent))

		// Create a factory instance using the traceparent
		ftx, _ := factory.NewFactoryFromTraceParent(traceparent)

		// Store the factory instance in the context
		c.Set("ftx", ftx)

		// Set the traceparent header in the request and response
		c.Request.Header.Set(constants.TraceparentHeader, traceparent)
		c.Header(constants.TraceparentHeader, traceparent)

		// Continue to the next handler in the chain
		c.Next()
	}
}

// GenerateTraceParent generates a simple traceparent header value
func GenerateTraceParent() string {
	// Create a traceparent value with a random trace ID and span ID
	traceparent := fmt.Sprintf("00-%016x-%08x-01", rand.Uint64(), rand.Uint32())
	return traceparent
}

// AddTraceParentToContext adds a traceparent header to the request context.
func AddTraceParentToContext(ctx context.Context, traceparent string) context.Context {
	// Add the traceparent value to the context
	ctx = context.WithValue(ctx, constants.TraceparentHeader, traceparent)
	return ctx
}

// GetTraceParentFromContext retrieves the traceparent header from the request context.
func GetTraceParentFromContext(ctx context.Context) string {
	// Retrieve the traceparent value from the context
	traceparent, _ := ctx.Value(constants.TraceparentHeader).(string)
	return traceparent
}
