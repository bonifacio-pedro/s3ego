package domain

import "github.com/bonifacio-pedro/s3ego/internal/model"

// BucketService defines business operations for managing S3 buckets.
type BucketService interface {
	Create(name string) error
	Exists(name string) (bool, error)
	Get(name string) (*model.Bucket, error)
	ListObjects(bucketName, prefix string) (*model.ListBucketResult, error)
	Delete(name string) error
}
