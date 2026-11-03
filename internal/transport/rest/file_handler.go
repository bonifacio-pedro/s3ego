// Package rest provides HTTP handlers for object operations conforming to AWS S3.
package rest

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/bonifacio-pedro/s3ego/internal/domain"
	"github.com/gin-gonic/gin"
)

// FileHandler handles S3 object-level operations.
type FileHandler struct {
	service domain.FileService
}

// NewFileHandler creates a new FileHandler.
func NewFileHandler(service domain.FileService) *FileHandler {
	return &FileHandler{service: service}
}

// Put handles PUT /{bucket}/{key...} uploading the request body directly as an object.
func (fh *FileHandler) Put(c *gin.Context) {
	bucketName := c.Param("bucket")
	key := strings.TrimPrefix(c.Param("key"), "/")

	if key == "" {
		RespondS3Error(c, http.StatusBadRequest, "InvalidObjectName", "Object key cannot be empty.")
		return
	}

	fileMeta, err := fh.service.Put(bucketName, key, c.Request.Body)
	if err != nil {
		if err.Error() == "NoSuchBucket" {
			RespondS3Error(c, http.StatusNotFound, "NoSuchBucket", "The specified bucket does not exist.")
			return
		}
		if err.Error() == "InvalidObjectName" {
			RespondS3Error(c, http.StatusBadRequest, "InvalidObjectName", "The specified object name is invalid.")
			return
		}
		RespondS3Error(c, http.StatusInternalServerError, "InternalError", err.Error())
		return
	}

	c.Header("ETag", fmt.Sprintf(`"%s"`, fileMeta.ETag))
	c.Header("x-amz-storage-class", "STANDARD")
	c.Header("Content-Length", "0")
	c.Status(http.StatusOK)
}

// Get handles GET /{bucket}/{key...} streaming the object data to the client.
func (fh *FileHandler) Get(c *gin.Context) {
	bucketName := c.Param("bucket")
	key := strings.TrimPrefix(c.Param("key"), "/")

	fileMeta, fileHandle, err := fh.service.Get(bucketName, key)
	if err != nil {
		if err.Error() == "NoSuchBucket" {
			RespondS3Error(c, http.StatusNotFound, "NoSuchBucket", "The specified bucket does not exist.")
			return
		}
		if err.Error() == "NoSuchKey" {
			RespondS3Error(c, http.StatusNotFound, "NoSuchKey", "The specified key does not exist.")
			return
		}
		RespondS3Error(c, http.StatusInternalServerError, "InternalError", err.Error())
		return
	}
	defer fileHandle.Close()

	c.Header("ETag", fmt.Sprintf(`"%s"`, fileMeta.ETag))
	c.Header("x-amz-storage-class", "STANDARD")
	c.Header("Accept-Ranges", "bytes")
	c.Header("Content-Type", fileMeta.ContentType)

	http.ServeContent(c.Writer, c.Request, filepath.Base(fileMeta.Key), fileMeta.LastModified, fileHandle)
}

// Head handles HEAD /{bucket}/{key...} returning object metadata without body.
func (fh *FileHandler) Head(c *gin.Context) {
	bucketName := c.Param("bucket")
	key := strings.TrimPrefix(c.Param("key"), "/")

	fileMeta, err := fh.service.Head(bucketName, key)
	if err != nil {
		if err.Error() == "NoSuchBucket" || err.Error() == "NoSuchKey" {
			c.Status(http.StatusNotFound)
			return
		}
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Header("ETag", fmt.Sprintf(`"%s"`, fileMeta.ETag))
	c.Header("Content-Type", fileMeta.ContentType)
	c.Header("Content-Length", fmt.Sprintf("%d", fileMeta.Size))
	c.Header("Last-Modified", fileMeta.LastModified.UTC().Format(http.TimeFormat))
	c.Header("Accept-Ranges", "bytes")
	c.Header("x-amz-storage-class", "STANDARD")
	c.Status(http.StatusOK)
}

// Delete handles DELETE /{bucket}/{key...} removing the object.
func (fh *FileHandler) Delete(c *gin.Context) {
	bucketName := c.Param("bucket")
	key := strings.TrimPrefix(c.Param("key"), "/")

	err := fh.service.Delete(bucketName, key)
	if err != nil {
		if err.Error() == "NoSuchBucket" {
			RespondS3Error(c, http.StatusNotFound, "NoSuchBucket", "The specified bucket does not exist.")
			return
		}
		RespondS3Error(c, http.StatusInternalServerError, "InternalError", err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}
