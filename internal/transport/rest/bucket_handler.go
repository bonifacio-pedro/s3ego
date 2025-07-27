// Package rest provides HTTP handlers for bucket and file related operations.
package rest

import (
	"github.com/bonifacio-pedro/s3ego/internal/domain"
	"github.com/gin-gonic/gin"
	"net/http"
)

// BucketHandler handles HTTP requests related to bucket operations.
type BucketHandler struct {
	service *domain.BucketService
}

// NewBucketHandler creates a new BucketHandler with the given BucketService.
func NewBucketHandler(service *domain.BucketService) *BucketHandler {
	return &BucketHandler{service: service}
}

// Create handles POST requests to create a new bucket.
// It expects a bucket name as a URL parameter "name".
// Returns HTTP 201 Created with the bucket URL on success,
// or HTTP 400 Bad Request if an error occurs.
func (bh *BucketHandler) Create(c *gin.Context) {
	bucketName := c.Param("name")

	buckerUrl, err := bh.service.New(bucketName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "bucket created successfully",
		"url":     buckerUrl,
	})
}

// FindAllFiles handles GET requests to list all files in a bucket.
// It expects the bucket name as a URL parameter "bucket".
// Returns HTTP 200 OK with the list of file keys on success,
// or HTTP 400 Bad Request if an error occurs.
func (bh *BucketHandler) FindAllFiles(c *gin.Context) {
	bucketName := c.Param("bucket")

	files, err := bh.service.FindAllFiles(bucketName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Default S3 Headers
	c.Header("x-amz-bucket-region", "us-east-1") // Default region
	c.Header("Content-Type", "application/json")

	c.JSON(http.StatusOK, gin.H{
		"Name":        bucketName,
		"Contents":    files,
		"MaxKeys":     1000,
		"IsTruncated": false,
	})
}

// Delete handles DELETE requests to remove a bucket.
// It expects the bucket name as a URL parameter "bucket".
// Returns HTTP 204 No Content on successful deletion,
// or HTTP 400 Bad Request if an error occurs.
func (bh *BucketHandler) Remove(c *gin.Context) {
	bucketName := c.Param("bucket")

	err := bh.service.Remove(bucketName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
