package repository

import "github.com/bonifacio-pedro/s3ego/internal/model"

// FileRepository defines metadata persistence operations for files/objects.
type FileRepository interface {
	Save(file *model.File) error
	Get(bucketID int, key string) (*model.File, error)
	Delete(bucketID int, key string) (*model.File, error)
}
