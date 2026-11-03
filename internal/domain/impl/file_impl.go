// Package impl contains business logic implementations for S3 objects and buckets.
package impl

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"

	"github.com/bonifacio-pedro/s3ego/internal/domain"
	"github.com/bonifacio-pedro/s3ego/internal/model"
	"github.com/bonifacio-pedro/s3ego/internal/repository"
	"github.com/bonifacio-pedro/s3ego/internal/storage"
)

type fileService struct {
	fileRepository   repository.FileRepository
	bucketRepository repository.BucketRepository
	storage          storage.FileStorage
}

// NewFileService creates a new FileService.
func NewFileService(fileRepo repository.FileRepository, bucketRepo repository.BucketRepository, st storage.FileStorage) domain.FileService {
	return &fileService{
		fileRepository:   fileRepo,
		bucketRepository: bucketRepo,
		storage:          st,
	}
}

func sanitizeKey(key string) string {
	clean := strings.TrimPrefix(path.Clean("/"+strings.TrimPrefix(key, "/")), "/")
	return clean
}

func (fs *fileService) Put(bucketName, key string, src io.Reader) (*model.File, error) {
	cleanKey := sanitizeKey(key)
	if cleanKey == "" || cleanKey == "." {
		return nil, errors.New("InvalidObjectName")
	}

	bucket, err := fs.bucketRepository.GetByName(bucketName)
	if err != nil {
		return nil, errors.New("NoSuchBucket")
	}

	storedPath, size, etag, initialBytes, err := fs.storage.Store(bucketName, cleanKey, src)
	if err != nil {
		return nil, fmt.Errorf("failed storing file: %w", err)
	}

	contentType := model.DetectContentType(initialBytes, cleanKey)

	fileModel := &model.File{
		Key:         cleanKey,
		BucketID:    bucket.ID,
		ETag:        etag,
		ContentType: contentType,
		Size:        size,
		FilePath:    storedPath,
	}

	if err := fs.fileRepository.Save(fileModel); err != nil {
		// Clean up stored file if metadata fails
		fs.storage.Delete(storedPath)
		return nil, fmt.Errorf("failed saving file metadata: %w", err)
	}

	log.Printf("[S3EGO] PutObject success: %s/%s (size=%d bytes, etag=%s, mime=%s)", bucketName, cleanKey, size, etag, contentType)
	return fileModel, nil
}

func (fs *fileService) Get(bucketName, key string) (*model.File, *os.File, error) {
	cleanKey := sanitizeKey(key)
	bucket, err := fs.bucketRepository.GetByName(bucketName)
	if err != nil {
		return nil, nil, errors.New("NoSuchBucket")
	}

	fileMeta, err := fs.fileRepository.Get(bucket.ID, cleanKey)
	if err != nil {
		return nil, nil, errors.New("NoSuchKey")
	}

	fileHandle, err := fs.storage.Open(fileMeta.FilePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed opening stored file: %w", err)
	}

	log.Printf("[S3EGO] GetObject: %s/%s (size=%d)", bucketName, cleanKey, fileMeta.Size)
	return fileMeta, fileHandle, nil
}

func (fs *fileService) Head(bucketName, key string) (*model.File, error) {
	cleanKey := sanitizeKey(key)
	bucket, err := fs.bucketRepository.GetByName(bucketName)
	if err != nil {
		return nil, errors.New("NoSuchBucket")
	}

	fileMeta, err := fs.fileRepository.Get(bucket.ID, cleanKey)
	if err != nil {
		return nil, errors.New("NoSuchKey")
	}

	return fileMeta, nil
}

func (fs *fileService) Delete(bucketName, key string) error {
	cleanKey := sanitizeKey(key)
	bucket, err := fs.bucketRepository.GetByName(bucketName)
	if err != nil {
		return errors.New("NoSuchBucket")
	}

	fileMeta, err := fs.fileRepository.Delete(bucket.ID, cleanKey)
	if err != nil {
		// In S3 standard, deleting a nonexistent key returns 204 No Content
		return nil
	}

	if fileMeta != nil && fileMeta.FilePath != "" {
		_ = fs.storage.Delete(fileMeta.FilePath)
	}

	log.Printf("[S3EGO] DeleteObject: %s/%s", bucketName, cleanKey)
	return nil
}
