package s3ego

import (
	"github.com/bonifacio-pedro/s3ego/internal/app"
	"github.com/bonifacio-pedro/s3ego/internal/config"
	"github.com/bonifacio-pedro/s3ego/internal/domain"
)

type S3EGO struct {
	Bucket *domain.BucketService
	File   *domain.FileService
}

func Start() *S3EGO {
	db := config.ConfigDatabase()
	newApp := app.NewApp(db)

	return &S3EGO{
		Bucket: newApp.BucketService,
		File:   newApp.FileService,
	}
}
