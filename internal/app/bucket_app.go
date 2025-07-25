package app

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"s3ego/internal/routes"
	"s3ego/internal/service"
)

type App struct {
	Router        *gin.Engine
	BucketService *service.BucketService
	FileService   *service.FileService
}

func NewApp(db *sql.DB) *App {
	r := routes.HandleRequests(db)
	bucketService := service.NewBucketService(db)
	fileService := service.NewFileService(db)

	return &App{
		Router:        r,
		BucketService: bucketService,
		FileService:   fileService,
	}
}

func (a *App) Run() {
	if err := a.Router.Run(":7777"); err != nil {
		log.Panic(err)
	}
}
