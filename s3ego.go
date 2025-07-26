package s3ego

import (
	"github.com/bonifacio-pedro/s3ego/internal/app"
	"github.com/bonifacio-pedro/s3ego/internal/config"
	"github.com/bonifacio-pedro/s3ego/internal/service"
)

type S3EGO struct {
	Bucket *service.BucketService
	File   *service.FileService
}

func Start() *S3EGO {
	db := config.ConfigDatabase()
	newApp := app.NewApp(db)

	return &S3EGO{
		Bucket: newApp.BucketService,
		File:   newApp.FileService,
	}
}
