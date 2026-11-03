// Package app initializes and runs the main S3EGO application.
package app

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/bonifacio-pedro/s3ego/internal/config"
	"github.com/bonifacio-pedro/s3ego/internal/domain"
	domainImpl "github.com/bonifacio-pedro/s3ego/internal/domain/impl"
	repoImpl "github.com/bonifacio-pedro/s3ego/internal/repository/impl"
	"github.com/bonifacio-pedro/s3ego/internal/storage"
	"github.com/bonifacio-pedro/s3ego/internal/transport/rest"
	"github.com/bonifacio-pedro/s3ego/internal/transport/routes"
	"github.com/gin-gonic/gin"
)

// App represents the initialized S3EGO server application.
type App struct {
	Router        *gin.Engine
	BucketService domain.BucketService
	FileService   domain.FileService
	Storage       storage.FileStorage
	DB            *sql.DB
	DataDir       string
}

// NewApp initializes repositories, storage, services, handlers, and routes for a given data directory.
func NewApp(dataDir string) (*App, error) {
	gin.SetMode(gin.ReleaseMode)

	db, err := config.ConfigDatabase(dataDir)
	if err != nil {
		return nil, fmt.Errorf("failed initializing database: %w", err)
	}

	st, err := storage.NewDiskStorage(dataDir)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed initializing storage: %w", err)
	}

	bucketRepository := repoImpl.NewBucketRepository(db)
	fileRepository := repoImpl.NewFileRepository(db)

	bucketService := domainImpl.NewBucketService(bucketRepository, st)
	fileService := domainImpl.NewFileService(fileRepository, bucketRepository, st)

	bucketHandler := rest.NewBucketHandler(bucketService)
	fileHandler := rest.NewFileHandler(fileService)

	rg := gin.New()
	rg.Use(gin.Recovery())

	router := routes.NewRouter(rg, bucketHandler, fileHandler)
	router.RegisterRoutes()

	return &App{
		Router:        rg,
		BucketService: bucketService,
		FileService:   fileService,
		Storage:       st,
		DB:            db,
		DataDir:       dataDir,
	}, nil
}

// Run starts the HTTP server on the specified port.
func (a *App) Run(port string) error {
	addr := ":" + port
	log.Printf("[S3EGO] S3 Wire-Compatible Server listening on %s (dataDir=%s)", addr, a.DataDir)
	return a.Router.Run(addr)
}

// Close gracefully closes open database connections.
func (a *App) Close() error {
	if a.DB != nil {
		return a.DB.Close()
	}
	return nil
}
