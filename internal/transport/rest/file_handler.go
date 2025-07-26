package rest

import (
	"github.com/bonifacio-pedro/s3ego/internal/domain"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

type FileHandler struct {
	service *domain.FileService
}

func NewFileHandler(service *domain.FileService) *FileHandler {
	return &FileHandler{service: service}
}

func (fh *FileHandler) Get(c *gin.Context) {
	bucketName := c.Param("bucket")
	key := strings.TrimPrefix(c.Param("key"), "/")

	file, err := fh.service.Get(bucketName, key)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"file": file})
}

func (fh *FileHandler) New(c *gin.Context) {
	bucketName := c.Param("bucket")

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required, put 'file' in form"})
		return
	}

	fileData, err := getFileData(fileHeader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
		return
	}

	fileKey, err := fh.service.Upload(bucketName, fileData, fileHeader.Filename)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "File uploaded successfully",
		"key":     fileKey,
		"bucket":  bucketName,
	})
}

func getFileData(fileHeader *multipart.FileHeader) ([]byte, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return data, nil
}
