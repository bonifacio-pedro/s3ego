package service

import (
	"database/sql"
	"fmt"
	"github.com/bonifacio-pedro/s3ego/internal/model"
	"github.com/bonifacio-pedro/s3ego/internal/repository"
	"log"
)

type FileService struct {
	db *sql.DB
}

func NewFileService(db *sql.DB) *FileService {
	return &FileService{db: db}
}

func (fs *FileService) GetFile(bucketUrl string, key string) (*[]byte, error) {
	bucketRepository := repository.NewBucketRepository(fs.db)
	fileRepository := repository.NewFileRepository(fs.db)

	bucket, err := bucketRepository.GetBucketByUrl(bucketUrl)
	if err != nil {
		return nil, err
	}

	file, err := fileRepository.GetFileByKey(key)
	if err != nil {
		return nil, err
	}

	if int(file.BucketID) != bucket.ID {
		return nil, err
	}

	log.Println(fmt.Sprintf("[S3EGO] PULLED NEW FILE: %s/%s", bucket.Name, key))

	return &file.Data, nil
}

func (fs *FileService) UploadFile(bucketUrl string, data *[]byte, fileName string) (string, error) {
	bucketRepository := repository.NewBucketRepository(fs.db)
	fileRepository := repository.NewFileRepository(fs.db)

	bucket, err := bucketRepository.GetBucketByUrl(bucketUrl)
	if err != nil {
		return "", err
	}

	fileModel := model.CreateFile(data, bucket, fileName)
	if err := fileRepository.CreateFile(fileModel); err != nil {
		return "", err
	}

	log.Println(fmt.Sprintf("[S3EGO] RECEIVED NEW FILE: %s/%s", bucket.Name, fileModel.Key))

	return fileModel.Key, nil
}
