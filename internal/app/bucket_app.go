// Package app initializes and runs the main S3EGO application,
// setting up the repository, services, handlers, and routes.
package app

import (
	"database/sql"
	"github.com/bonifacio-pedro/s3ego/internal/domain"
	"github.com/bonifacio-pedro/s3ego/internal/repository"
	"github.com/bonifacio-pedro/s3ego/internal/transport/rest"
	"github.com/bonifacio-pedro/s3ego/internal/transport/routes"
	"github.com/gin-gonic/gin"
	"log"
)

// App represents the main application instance.
// It holds the router and core services (BucketService and FileService).
type App struct {
	Router        *gin.Engine
	BucketService *domain.BucketService
	FileService   *domain.FileService
}

// NewApp initializes the application, wiring together dependencies such as
// repositories, services, handlers, and routes.
// It returns a fully constructed App ready to be run.
func NewApp(db *sql.DB) *App {
	// Set Gin to Release mode (no debug output)
	gin.SetMode(gin.ReleaseMode)

	// Initialize Gin engine
	rg := gin.Default()

	// Repositories
	bucketRepository := repository.NewBucketRepository(db)
	fileRepository := repository.NewFileRepository(db)

	// Services
	bucketService := domain.NewBucketService(bucketRepository)
	fileService := domain.NewFileService(fileRepository, bucketRepository)

	// Handlers (transport layer)
	bucketHandler := rest.NewBucketHandler(bucketService)
	fileHandler := rest.NewFileHandler(fileService)

	// Routes
	router := routes.NewRouter(rg, bucketHandler, fileHandler)
	router.RegisterRoutes()

	return &App{
		Router:        rg,
		BucketService: bucketService,
		FileService:   fileService,
	}
}

// Run starts the HTTP server on port 7777.
// If an error occurs during startup, the application will panic.
func (a *App) Run() {
	if err := a.Router.Run(":7777"); err != nil {
		log.Panic(err)
	}
}
