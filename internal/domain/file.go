package domain

import (
	"fmt"
	"github.com/bonifacio-pedro/s3ego/internal/model"
	"github.com/bonifacio-pedro/s3ego/internal/repository"
	"log"
)

type FileService struct {
	fileRepository   *repository.FileRepository
	bucketRepository *repository.BucketRepository
}

func NewFileService(fileRepository *repository.FileRepository, bucketRepository *repository.BucketRepository) *FileService {
	return &FileService{fileRepository: fileRepository, bucketRepository: bucketRepository}
}

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
