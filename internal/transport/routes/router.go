// Package routes provides HTTP route definitions and registration for the S3 emulator.
package routes

import (
	"github.com/bonifacio-pedro/s3ego/internal/transport/middleware"
	"github.com/bonifacio-pedro/s3ego/internal/transport/rest"
	"github.com/gin-gonic/gin"
)

// Router wraps the Gin engine and the HTTP handlers for buckets and files.
type Router struct {
	rg            *gin.Engine
	bucketHandler *rest.BucketHandler
	fileHandler   *rest.FileHandler
}

// NewRouter creates a new Router instance with the provided Gin engine and handlers.
//
// Parameters:
//   - rg: the Gin engine instance to register routes on.
//   - bucketHandler: handler responsible for bucket-related endpoints.
//   - fileHandler: handler responsible for file-related endpoints.
//
// Returns a pointer to the newly created Router.
func NewRouter(rg *gin.Engine, bucketHandler *rest.BucketHandler, fileHandler *rest.FileHandler) *Router {
	return &Router{rg: rg, bucketHandler: bucketHandler, fileHandler: fileHandler}
}

// RegisterRoutes configure S3HeadersMiddleware and
// registers all HTTP routes/endpoints for the bucket and file handlers.
//
// It sets up routes for creating buckets, listing files, deleting buckets and files,
// uploading files, and retrieving files from the bucket emulator.
func (ro *Router) RegisterRoutes() {
	ro.rg.Use(middleware.S3HeadersMiddleware())

	ro.rg.POST("/bucket-emulator/new-bucket/:name", ro.bucketHandler.Create)
	ro.rg.GET("/bucket-emulator/list-files/:bucket", ro.bucketHandler.FindAllFiles)
	ro.rg.DELETE("/bucket-emulator/remove-bucket/:bucket", ro.bucketHandler.Remove)
	ro.rg.DELETE("/bucket-emulator/remove-file/:bucket/*key", ro.fileHandler.Remove)
	ro.rg.POST("/bucket-emulator/upload-file/:bucket", ro.fileHandler.New)
	ro.rg.GET("/bucket-emulator/get-file/:bucket/*key", ro.fileHandler.Get)
}
