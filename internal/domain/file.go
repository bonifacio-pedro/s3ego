// Package domain contains the business logic for managing buckets and files.
package domain

import (
	"fmt"
	"github.com/bonifacio-pedro/s3ego/internal/model"
	"github.com/bonifacio-pedro/s3ego/internal/repository"
	"log"
)

// FileService provides methods to manage files within buckets.
// It communicates with FileRepository and BucketRepository to perform CRUD operations.
type FileService struct {
	fileRepository   *repository.FileRepository
	bucketRepository *repository.BucketRepository
}

// NewFileService creates a new FileService with the provided file and bucket repositories.
func NewFileService(fileRepository *repository.FileRepository, bucketRepository *repository.BucketRepository) *FileService {
	return &FileService{fileRepository: fileRepository, bucketRepository: bucketRepository}
}

// Get retrieves the file data by bucket name and file key.
// Returns the file data bytes or an error if the bucket or file doesn't exist,
// or if the file does not belong to the specified bucket.
func (fs *FileService) Get(bucketName string, key string) ([]byte, error) {
	bucket, err := fs.bucketRepository.GetByName(bucketName)
	if err != nil {
		return nil, err
	}

	file, err := fs.fileRepository.GetByKey(key)
	if err != nil {
		return nil, err
	}

	if int(file.BucketID) != bucket.ID {
		return nil, err
	}

	log.Println(fmt.Sprintf("[S3EGO] PULLED NEW FILE: %s/%s", bucket.Name, key))
	return file.Data, nil
}

// Remove deletes a file specified by bucket name and key.
// Returns an error if the bucket or file does not exist,
// or if the file does not belong to the specified bucket.
func (fs *FileService) Remove(bucketName string, key string) error {
	bucket, err := fs.bucketRepository.GetByName(bucketName)
	if err != nil {
		return err
	}

	file, err := fs.fileRepository.GetByKey(key)
	if err != nil {
		return err
	}

	if file.BucketID != uint(bucket.ID) {
		return fmt.Errorf("this file is not in %s bucket", bucket.Name)
	}

	err = fs.fileRepository.Remove(key)
	if err != nil {
		return err
	}

	log.Println(fmt.Sprintf("[S3EGO] FILE REMOVED: %s/%s", bucket.Name, key))
	return nil
}

// Upload stores a new file in the specified bucket.
// It returns the key of the stored file or an error if the bucket does not exist,
// if the file already exists in the bucket, or if there was a failure during insertion.
func (fs *FileService) Upload(bucketName string, data []byte, fileName string) (string, error) {
	bucket, err := fs.bucketRepository.GetByName(bucketName)
	if err != nil {
		return "", err
	}

	fileModel := model.NewFile(data, *bucket, fileName)

	fileExists, err := fs.bucketRepository.FileExists(bucketName, fileModel.Key)
	if err != nil {
		return "", err
	}

	if fileExists {
		return fileModel.Key, fmt.Errorf("file %s already exists in %s bucket", fileModel.Key, bucketName)
	}

	if err := fs.fileRepository.New(&fileModel); err != nil {
		return "", err
	}

	log.Println(fmt.Sprintf("[S3EGO] RECEIVED NEW FILE: %s/%s", bucket.Name, fileModel.Key))
	return fileModel.Key, nil
}
