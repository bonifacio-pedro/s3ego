package domain

import (
	"errors"
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

func (fs *FileService) Upload(bucketName string, data []byte, fileName string) (string, error) {
	bucket, err := fs.bucketRepository.GetByName(bucketName)
	if err != nil {
		return "", err
	}

	fileModel := model.NewFile(data, *bucket, fileName)

	fileExists, err := fs.fileRepository.ExistsByKey(fileModel.Key)
	if err != nil {
		return "", err
	}

	if fileExists {
		return fileModel.Key, errors.New("file with this key/name already exists")
	}

	if err := fs.fileRepository.New(&fileModel); err != nil {
		return "", err
	}

	log.Println(fmt.Sprintf("[S3EGO] RECEIVED NEW FILE: %s/%s", bucket.Name, fileModel.Key))

	return fileModel.Key, nil
}
