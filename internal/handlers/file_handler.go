package handlers

import (
	"database/sql"
	"github.com/bonifacio-pedro/s3ego/internal/model"
	"github.com/bonifacio-pedro/s3ego/internal/repository"
	fileUtils "github.com/bonifacio-pedro/s3ego/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type FileHandler struct {
	db *sql.DB
}

func NewFileHandler(db *sql.DB) *FileHandler {
	return &FileHandler{db: db}
}

func (fh *FileHandler) GetFile(c *gin.Context) {
	bucketName := c.Param("bucket")
	key := strings.TrimPrefix(c.Param("key"), "/")

	bucketRepository := repository.NewBucketRepository(fh.db)
	fileRepository := repository.NewFileRepository(fh.db)

	bucket, err := bucketRepository.GetBucketByName(bucketName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bucket not found"})
		return
	}

	file, err := fileRepository.GetFileByKey(key)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	if int(file.BucketID) != bucket.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "file does not belong to bucket"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"file": file.Data})
}

func (fh *FileHandler) UploadFile(c *gin.Context) {
	bucketName := c.Param("bucket")
	bucketRepository := repository.NewBucketRepository(fh.db)
	fileRepository := repository.NewFileRepository(fh.db)

	bucket, err := bucketRepository.GetBucketByName(bucketName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bucket not found"})
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required, put 'file' in form"})
		return
	}

	fileData, err := fileUtils.GetFileData(fileHeader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
		return
	}

	fileModel := model.CreateFile(&fileData, bucket, fileHeader.Filename)

	if err := fileRepository.CreateFile(fileModel); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create file"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "file uploaded successfully",
		"key":     fileModel.Key,
		"bucket":  bucket.Name,
	})
}
