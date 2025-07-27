// Package middleware provides Gin middlewares for the S3EGO project.
package middleware

import (
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"time"
)

// S3HeadersMiddleware is a Gin middleware that adds typical S3-like headers
// to every HTTP response. It emulates the behavior of Amazon S3 by injecting
// headers such as x-amz-request-id, x-amz-id-2, Date, and Server.
//
// This middleware is useful for simulating real AWS S3 responses in local
// development and testing environments.
func S3HeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Server", "s3ego/1.0")
		c.Header("x-amz-request-id", generateRequestID())
		c.Header("x-amz-id-2", generateAMZID2())
		c.Header("Date", time.Now().UTC().Format(time.RFC1123))

		c.Next()
	}
}

// generateRequestID returns a pseudo-random request ID using an MD5 hash
// of a newly generated UUID. This mimics the x-amz-request-id header
// returned by AWS S3.
func generateRequestID() string {
	return fmt.Sprintf("%x", md5.Sum([]byte(uuid.New().String())))
}

// generateAMZID2 returns a pseudo-random string used to simulate the
// x-amz-id-2 header returned by AWS S3. It returns the first 16 characters
// of a UUID.
func generateAMZID2() string {
	return uuid.New().String()[:16]
}
