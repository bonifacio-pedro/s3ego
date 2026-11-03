// Package routes provides HTTP route definitions and registration for the S3 emulator.
package routes

import (
	"github.com/bonifacio-pedro/s3ego/internal/transport/middleware"
	"github.com/bonifacio-pedro/s3ego/internal/transport/rest"
	"github.com/gin-gonic/gin"
)

// Router wraps the Gin engine and HTTP handlers for S3 buckets and files.
type Router struct {
	rg            *gin.Engine
	bucketHandler *rest.BucketHandler
	fileHandler   *rest.FileHandler
}

// NewRouter creates a new Router instance.
func NewRouter(rg *gin.Engine, bucketHandler *rest.BucketHandler, fileHandler *rest.FileHandler) *Router {
	return &Router{
		rg:            rg,
		bucketHandler: bucketHandler,
		fileHandler:   fileHandler,
	}
}

// RegisterRoutes configures middlewares and S3 Path-Style routes.
func (ro *Router) RegisterRoutes() {
	ro.rg.Use(middleware.S3HeadersMiddleware())

	// S3 Path-Style Bucket Routes
	ro.rg.PUT("/:bucket", ro.bucketHandler.Create)
	ro.rg.HEAD("/:bucket", ro.bucketHandler.Head)
	ro.rg.GET("/:bucket", ro.bucketHandler.ListObjects)
	ro.rg.DELETE("/:bucket", ro.bucketHandler.Delete)

	// S3 Path-Style Object Routes
	ro.rg.PUT("/:bucket/*key", ro.fileHandler.Put)
	ro.rg.GET("/:bucket/*key", ro.fileHandler.Get)
	ro.rg.HEAD("/:bucket/*key", ro.fileHandler.Head)
	ro.rg.DELETE("/:bucket/*key", ro.fileHandler.Delete)
}
