package domain

import (
	"errors"
	"github.com/bonifacio-pedro/s3ego/internal/model"
	"github.com/bonifacio-pedro/s3ego/internal/repository"
	"log"
)

type BucketService struct {
	repository *repository.BucketRepository
}

func NewBucketService(repository *repository.BucketRepository) *BucketService {
	return &BucketService{repository: repository}
}

func (bs *BucketService) New(name string) (string, error) {
	bucket := model.NewBucket(name)

	exists, err := bs.repository.ExistsByName(bucket.Name)
	if err != nil {
		return "", err
	}

	if exists {
		return "", errors.New("Bucket already exists")
	}

	if err := bs.repository.New(&bucket); err != nil {
		return "", err
	}

	log.Println("[S3EG0] CREATED NEW BUCKET: ", bucket.Name)

	return bucket.Url, nil
}

func (bs *BucketService) FindAllFiles(bucketName string) (*[]string, error) {
	bucket, err := bs.repository.GetByName(bucketName)
	if err != nil {
		return nil, err
	}

	files, err := bs.repository.GetFiles(bucket.ID)
	if err != nil {
		return nil, err
	}

	log.Println("[S3EGO] LISTED ALL FILES IN A BUCKET: ", bucket.Name)

	return &files, err
}

func (bs *BucketService) Remove(bucketName string) error {
	bucket, err := bs.repository.GetByName(bucketName)
	if err != nil {
		return err
	}

	if err := bs.repository.Remove(bucket.ID); err != nil {
		return err
	}

	log.Println("[S3EGO] BUCKET DELETED: ", bucketName)

	return nil
}
