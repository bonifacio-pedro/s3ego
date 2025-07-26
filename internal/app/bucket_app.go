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

type App struct {
	Router        *gin.Engine
	BucketService *domain.BucketService
	FileService   *domain.FileService
}

func NewApp(db *sql.DB) *App {
	gin.SetMode(gin.ReleaseMode)
	rg := gin.Default()

	// Repositories
	bucketRepository := repository.NewBucketRepository(db)
	fileRepository := repository.NewFileRepository(db)

	// Services
	bucketService := domain.NewBucketService(bucketRepository)
	fileService := domain.NewFileService(fileRepository, bucketRepository)

	// Transport layer
	bucketHandler := rest.NewBucketHandler(bucketService)
	fileHandler := rest.NewFileHandler(fileService)

	// Routing
	router := routes.NewRouter(rg, bucketHandler, fileHandler)
	router.RegisterRoutes()

	return &App{
		Router:        rg,
		BucketService: bucketService,
		FileService:   fileService,
	}
}

func (a *App) Run() {
	if err := a.Router.Run(":7777"); err != nil {
		log.Panic(err)
	}
}
