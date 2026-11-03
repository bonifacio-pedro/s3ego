// Package rest provides HTTP handlers for bucket and file operations complying with AWS S3 wire protocol.
package rest

import (
	"encoding/xml"
	"errors"
	"net/http"

	"github.com/bonifacio-pedro/s3ego/internal/domain"
	"github.com/gin-gonic/gin"
)

// BucketHandler handles S3 bucket-level operations.
type BucketHandler struct {
	service domain.BucketService
}

// NewBucketHandler creates a new BucketHandler.
func NewBucketHandler(service domain.BucketService) *BucketHandler {
	return &BucketHandler{service: service}
}

// Create handles PUT /{bucket} to create a new bucket.
func (bh *BucketHandler) Create(c *gin.Context) {
	bucketName := c.Param("bucket")

	err := bh.service.Create(bucketName)
	if err != nil {
		if errors.Is(err, errors.New("BucketAlreadyOwnedByYou")) || err.Error() == "BucketAlreadyOwnedByYou" {
			RespondS3Error(c, http.StatusConflict, "BucketAlreadyOwnedByYou", "Your previous request to create the named bucket succeeded and you already own it.")
			return
		}
		if errors.Is(err, errors.New("InvalidBucketName")) || err.Error() == "InvalidBucketName" {
			RespondS3Error(c, http.StatusBadRequest, "InvalidBucketName", "The specified bucket is not valid.")
			return
		}
		RespondS3Error(c, http.StatusBadRequest, "InvalidArgument", err.Error())
		return
	}

	c.Header("Location", "/"+bucketName)
	c.Status(http.StatusOK)
}

// Head handles HEAD /{bucket} to check bucket existence.
func (bh *BucketHandler) Head(c *gin.Context) {
	bucketName := c.Param("bucket")

	exists, err := bh.service.Exists(bucketName)
	if err != nil || !exists {
		c.Status(http.StatusNotFound)
		return
	}

	c.Status(http.StatusOK)
}

// ListObjects handles GET /{bucket} to list objects inside a bucket in XML format.
func (bh *BucketHandler) ListObjects(c *gin.Context) {
	bucketName := c.Param("bucket")
	prefix := c.Query("prefix")

	result, err := bh.service.ListObjects(bucketName, prefix)
	if err != nil {
		if err.Error() == "NoSuchBucket" {
			RespondS3Error(c, http.StatusNotFound, "NoSuchBucket", "The specified bucket does not exist.")
			return
		}
		RespondS3Error(c, http.StatusInternalServerError, "InternalError", err.Error())
		return
	}

	xmlBytes, err := xml.MarshalIndent(result, "", "  ")
	if err != nil {
		RespondS3Error(c, http.StatusInternalServerError, "InternalError", "Failed to marshal XML result.")
		return
	}

	c.Data(http.StatusOK, "application/xml; charset=utf-8", append([]byte(xml.Header), xmlBytes...))
}

// Delete handles DELETE /{bucket} to remove a bucket.
func (bh *BucketHandler) Delete(c *gin.Context) {
	bucketName := c.Param("bucket")

	err := bh.service.Delete(bucketName)
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
