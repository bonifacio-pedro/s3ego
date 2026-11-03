package repository

import "github.com/bonifacio-pedro/s3ego/internal/model"

// BucketRepository defines metadata persistence operations for S3 buckets.
type BucketRepository interface {
	New(bucket *model.Bucket) error
	Remove(bucketID int) error
	ExistsByName(bucketName string) (bool, error)
	GetByName(bucketName string) (*model.Bucket, error)
	ListObjects(bucketID int, prefix string) ([]model.File, error)
	ListBuckets() ([]model.Bucket, error)
}
