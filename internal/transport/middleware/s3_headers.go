// Package middleware provides Gin middlewares for emulating standard AWS S3 headers and behavior.
package middleware

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// S3HeadersMiddleware injects standard AWS S3 headers and CORS allowances.
func S3HeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := generateRequestID()
		c.Header("Server", "AmazonS3")
		c.Header("x-amz-request-id", reqID)
		c.Header("x-amz-id-2", generateAMZID2())
		c.Header("x-amz-bucket-region", "us-east-1")
		c.Header("Date", time.Now().UTC().Format(http.TimeFormat))

		// Permissive CORS support for browser-based S3 SDK clients
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, HEAD, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header("Access-Control-Expose-Headers", "ETag, x-amz-request-id, x-amz-id-2, Content-Length, Content-Type, Last-Modified")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}

func generateRequestID() string {
	sum := md5.Sum([]byte(uuid.New().String()))
	return fmt.Sprintf("%X", sum[:8])
}

func generateAMZID2() string {
	return uuid.New().String()[:16]
}
