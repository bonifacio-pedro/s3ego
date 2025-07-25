package service

import (
	"database/sql"
	"log"
	"s3ego/internal/model"
	"s3ego/internal/repository"
)

type BucketService struct {
	db *sql.DB
}

func NewBucketService(db *sql.DB) *BucketService {
	return &BucketService{db: db}
}

func (bs *BucketService) CreateBucket(name string) (string, error) {
	bucketRepository := repository.NewBucketRepository(bs.db)
	bucketModel := model.CreateBucket(name)

	if err := bucketRepository.CreateBucket(bucketModel); err != nil {
		return "", err
	}

	log.Println("[S3-EMULATOR] CREATED NEW BUCKET: ", bucketModel.Name)

	return bucketModel.Url, nil
}
