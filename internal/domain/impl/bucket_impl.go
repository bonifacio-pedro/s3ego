// Package impl contains the business logic for managing buckets and files.
package impl

import (
	"errors"
	"fmt"
	"log"
	"regexp"

	"github.com/bonifacio-pedro/s3ego/internal/domain"
	"github.com/bonifacio-pedro/s3ego/internal/model"
	"github.com/bonifacio-pedro/s3ego/internal/repository"
	"github.com/bonifacio-pedro/s3ego/internal/storage"
)

var bucketNameRegex = regexp.MustCompile(`^[a-z0-9][a-z0-9.-]{1,61}[a-z0-9]$`)

type bucketService struct {
	repository repository.BucketRepository
	storage    storage.FileStorage
}

// NewBucketService creates a new BucketService.
func NewBucketService(repo repository.BucketRepository, st storage.FileStorage) domain.BucketService {
	return &bucketService{
		repository: repo,
		storage:    st,
	}
}

func (bs *bucketService) validateBucketName(name string) error {
	if len(name) < 3 || len(name) > 63 || !bucketNameRegex.MatchString(name) {
		return errors.New("InvalidBucketName")
	}
	return nil
}

func (bs *bucketService) Create(name string) error {
	if err := bs.validateBucketName(name); err != nil {
		return err
	}

	exists, err := bs.repository.ExistsByName(name)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("BucketAlreadyOwnedByYou")
	}

	bucket := model.NewBucket(name)
	if err := bs.repository.New(&bucket); err != nil {
		return err
	}

	log.Printf("[S3EGO] Created bucket: %s", name)
	return nil
}

func (bs *bucketService) Exists(name string) (bool, error) {
	return bs.repository.ExistsByName(name)
}

func (bs *bucketService) Get(name string) (*model.Bucket, error) {
	b, err := bs.repository.GetByName(name)
	if err != nil {
		return nil, errors.New("NoSuchBucket")
	}
	return b, nil
}

func (bs *bucketService) ListObjects(bucketName, prefix string) (*model.ListBucketResult, error) {
	bucket, err := bs.Get(bucketName)
	if err != nil {
		return nil, err
	}

	files, err := bs.repository.ListObjects(bucket.ID, prefix)
	if err != nil {
		return nil, err
	}

	contents := make([]model.ObjectSummary, 0, len(files))
	for _, f := range files {
		contents = append(contents, model.ObjectSummary{
			Key:          f.Key,
			LastModified: f.LastModified.UTC(),
			ETag:         fmt.Sprintf(`"%s"`, f.ETag),
			Size:         f.Size,
			StorageClass: "STANDARD",
		})
	}

	result := &model.ListBucketResult{
		Xmlns:       "http://s3.amazonaws.com/doc/2006-03-01/",
		Name:        bucketName,
		Prefix:      prefix,
		KeyCount:    len(contents),
		MaxKeys:     1000,
		IsTruncated: false,
		Contents:    contents,
	}

	log.Printf("[S3EGO] Listed %d objects in bucket %s (prefix=%q)", len(contents), bucketName, prefix)
	return result, nil
}

func (bs *bucketService) Delete(name string) error {
	bucket, err := bs.Get(name)
	if err != nil {
		return err
	}

	if err := bs.repository.Remove(bucket.ID); err != nil {
		return err
	}

	if err := bs.storage.DeleteBucketDir(name); err != nil {
		log.Printf("[S3EGO] Warning: failed to delete bucket directory %s: %s", name, err)
	}

	log.Printf("[S3EGO] Deleted bucket: %s", name)
	return nil
}
