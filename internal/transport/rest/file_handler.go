// Package rest provides HTTP handlers for file and bucket related operations.
package rest

import (
	"crypto/md5"
	"fmt"
	"github.com/bonifacio-pedro/s3ego/internal/domain"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

// FileHandler handles HTTP requests related to file operations.
type FileHandler struct {
	service *domain.FileService
}

// NewFileHandler creates a new FileHandler with the given FileService.
func NewFileHandler(service *domain.FileService) *FileHandler {
	return &FileHandler{service: service}
}

// Get handles GET requests to download a file from a bucket.
// It expects the bucket name as URL parameter "bucket" and the file key as "key".
// Returns HTTP 200 OK with file data on success,
// or HTTP 400 Bad Request if an error occurs.
func (fh *FileHandler) Get(c *gin.Context) {
	bucketName := c.Param("bucket")
	key := strings.TrimPrefix(c.Param("key"), "/")

	fileData, fileModel, err := fh.service.Get(bucketName, key)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// S3 Default Headers
	etag := fmt.Sprintf(`"%x"`, md5.Sum(fileData))
	c.Header("ETag", etag)
	c.Header("Last-Modified", fileModel.LastModified.UTC().Format(time.RFC1123))
	c.Header("Content-Length", fmt.Sprintf("%d", len(fileData)))
	c.Header("Content-Type", fileModel.ContentType)
	c.Header("Accept-Ranges", "bytes")
	c.Header("x-amz-storage-class", "STANDARD")

	c.Data(http.StatusOK, fileModel.ContentType, fileData)
}

// Remove handles DELETE requests to delete a file from a bucket.
// It expects the bucket name as URL parameter "bucket" and the file key as "key".
// Returns HTTP 204 No Content on success,
// or HTTP 400 Bad Request if an error occurs.
func (fh *FileHandler) Remove(c *gin.Context) {
	bucketName := c.Param("bucket")
	key := strings.TrimPrefix(c.Param("key"), "/")

	err := fh.service.Remove(bucketName, key)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

// New handles POST requests to upload a new file to a bucket.
// It expects the bucket name as URL parameter "bucket" and a form file with key "file".
// Returns HTTP 201 Created with the file key and bucket name on success,
// or HTTP 400 Bad Request / 500 Internal Server Error if an error occurs.
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

	fileKey, fileEtag, err := fh.service.Upload(bucketName, fileData, fileHeader.Filename)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// S3 Default Headers
	c.Header("ETag", fileEtag)
	c.Header("x-amz-version-id", "null")
	c.Header("x-amz-storage-class", "STANDARD")

	c.JSON(http.StatusCreated, gin.H{
		"message": "File uploaded successfully",
		"key":     fileKey,
		"bucket":  bucketName,
		"etag":    fileEtag,
	})
}

// getFileData reads all bytes from the uploaded file header.
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
