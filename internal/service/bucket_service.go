package service

import (
	"database/sql"
	"github.com/bonifacio-pedro/s3ego/internal/model"
	"github.com/bonifacio-pedro/s3ego/internal/repository"
	"log"
)

type BucketService struct {
	db *sql.DB
}

func NewBucketService(db *sql.DB) *BucketService {
	return &BucketService{db: db}
}

func (bs *BucketService) CreateBucket(name string) (string, error) {
	bucketRepository := repository.NewBucketRepository(bs.db)
	bucket := model.CreateBucket(name)

	if err := bucketRepository.CreateBucket(bucket); err != nil {
		return "", err
	}

	log.Println("[S3EG0] CREATED NEW BUCKET: ", bucket.Name)

	return bucket.Url, nil
}

func (bs *BucketService) FindAllFilesInABucket(bucketName string) (*[]string, error) {
	bucketRepository := repository.NewBucketRepository(bs.db)

	bucket, err := bucketRepository.GetBucketByName(bucketName)
	if err != nil {
		return nil, err
	}

	files, err := bucketRepository.GetFiles(bucket.ID)
	if err != nil {
		return nil, err
	}

	log.Println("[S3EGO] LISTED ALL FILES IN A BUCKET: ", bucket.Name)

	return &files, err
}
