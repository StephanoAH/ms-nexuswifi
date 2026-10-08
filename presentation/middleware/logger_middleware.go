package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"strings"
	"time"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if the content type is multipart/form-data
		if strings.Contains(c.GetHeader("Content-Type"), "multipart/form-data") {
			// Log request without body if it's multipart (likely contains files)
			log.Printf("Request [multipart]: %s %s from %s", c.Request.Method, c.Request.URL, c.ClientIP())
		} else {
			// Intercept request body for non-multipart requests
			requestBody, err := io.ReadAll(c.Request.Body)
			if err != nil {
				log.Printf("Error reading request body: %v", err)
			}
			// Restore the io.ReadCloser to its original state
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))

			// Log the complete request body
			log.Printf("Request: %s %s %s from %s", c.Request.Method, c.Request.URL, requestBody, c.ClientIP())
		}

		// Set up logging response writer
		blw := &bodyLogWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
			status:         200, // Default if not set
		}
		c.Writer = blw

		// Timing the request
		start := time.Now()
		c.Next()
		latency := time.Since(start)

		// Logging the response details
		log.Printf("Response: %d %s", blw.status, blw.body.String())
		log.Printf("Latency: %v", latency)
	}
}

// bodyLogWriter captures the response body and status code
type bodyLogWriter struct {
	gin.ResponseWriter
	body   *bytes.Buffer
	status int
}

// Write overrides the ResponseWriter Write method to capture response body
func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// WriteHeader overrides the ResponseWriter WriteHeader method to capture status code
func (w *bodyLogWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}
