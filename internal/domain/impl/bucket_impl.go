// Package domain contains business logic and services for managing S3EGO buckets and files.
package impl

import (
	"errors"
	"log"

	"github.com/bonifacio-pedro/s3ego/internal/domain"
	"github.com/bonifacio-pedro/s3ego/internal/model"
	"github.com/bonifacio-pedro/s3ego/internal/repository"
)

// BucketService encapsulates business logic related to S3EGO buckets.
// It acts as an intermediary between the handler layer and the repository.
type bucketService struct {
	repository repository.BucketRepository
}

// NewBucketService returns a new instance of BucketService.
//
// It receives a pointer to a BucketRepository which it uses
// to persist and retrieve bucket data.
func NewBucketService(repository repository.BucketRepository) domain.BucketService {
	return &bucketService{repository: repository}
}

// New creates a new bucket with the given name.
// It returns the bucket URL on success, or an error if the bucket already exists
// or if there was a problem creating it in the repository.
func (bs *bucketService) New(name string) (string, error) {
	bucket := model.NewBucket(name)

	exists, err := bs.repository.ExistsByName(bucket.Name)
	if err != nil {
		return "", err
	}

	if exists {
		return "", errors.New("bucket already exists")
	}

	if err := bs.repository.New(&bucket); err != nil {
		return "", err
	}

	log.Println("[S3EGO] CREATED NEW BUCKET:", bucket.Name)
	return bucket.Url, nil
}

// FindAllFiles returns all file keys stored in a given bucket by name.
// It returns a slice of strings or an error if the bucket doesn't exist
// or if there was an issue fetching the files.
func (bs *bucketService) FindAllFiles(bucketName string) (*[]string, error) {
	bucket, err := bs.repository.GetByName(bucketName)
	if err != nil {
		return nil, err
	}

	files, err := bs.repository.GetFiles(bucket.ID)
	if err != nil {
		return nil, err
	}

	log.Println("[S3EGO] LISTED ALL FILES IN A BUCKET:", bucket.Name)
	return &files, err
}

// Remove deletes a bucket by its name.
// It returns an error if the bucket doesn't exist or fails to be deleted.
func (bs *bucketService) Remove(bucketName string) error {
	bucket, err := bs.repository.GetByName(bucketName)
	if err != nil {
		return err
	}

	if err := bs.repository.Remove(bucket.ID); err != nil {
		return err
	}

	log.Println("[S3EGO] BUCKET DELETED:", bucketName)
	return nil
}
