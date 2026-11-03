package domain

import (
	"io"
	"os"

	"github.com/bonifacio-pedro/s3ego/internal/model"
)

// FileService defines business operations for managing S3 objects.
type FileService interface {
	Put(bucketName, key string, src io.Reader) (*model.File, error)
	Get(bucketName, key string) (*model.File, *os.File, error)
	Head(bucketName, key string) (*model.File, error)
	Delete(bucketName, key string) error
}
