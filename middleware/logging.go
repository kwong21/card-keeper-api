package middleware

import (
	"time"

	logger "card-keeper-api/internal/logging"

	"github.com/gin-gonic/gin"
)

var middlewareLogger = logger.NewLogger()

// LogToFile logs output to a file on disk
func LogToFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		// start time
		startTime := time.Now()
		// Processing request
		c.Next()
		// End time
		endTime := time.Now()
		// execution time
		latencyTime := endTime.Sub(startTime)
		// Request mode
		reqMethod := c.Request.Method
		// Request routing
		reqURI := c.Request.RequestURI
		// Status code
		statusCode := c.Writer.Status()
		// Request IP
		clientIP := c.ClientIP()
		// Log format
		middlewareLogger.Infof("| %3d | %13v | %15s | %s | %s |",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqURI,
		)
	}
}
