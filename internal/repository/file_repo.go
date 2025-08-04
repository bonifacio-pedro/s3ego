package repository

import "github.com/bonifacio-pedro/s3ego/internal/model"

// FileRepository interface for decoupling code
type FileRepository interface {
	New(file *model.File) error
	Remove(key string) error
	GetByKey(key string) (*model.File, error)
}
