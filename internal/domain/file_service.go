package domain

import "github.com/bonifacio-pedro/s3ego/internal/model"

type FileService interface {
	Get(bucketName string, key string) ([]byte, model.File, error)
	Remove(bucketName string, key string) error
	Upload(bucketName string, data []byte, fileName string) (string, string, error)
}