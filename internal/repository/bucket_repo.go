package repository

import "github.com/bonifacio-pedro/s3ego/internal/model"

// BucketRepository interface for decoupling code
type BucketRepository interface {
	New(bucket *model.Bucket) error
	Remove(bucketID int) error
	ExistsByName(bucketName string) (bool, error)
	GetByName(bucketName string) (*model.Bucket, error)
	GetFiles(bucketID int) ([]string, error)
	FileExists(bucketName string, key string) (bool, error)
}
