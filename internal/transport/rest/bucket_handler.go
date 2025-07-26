package rest

import (
	"github.com/bonifacio-pedro/s3ego/internal/domain"
	"github.com/gin-gonic/gin"
	"net/http"
)

type BucketHandler struct {
	service *domain.BucketService
}

func NewBucketHandler(service *domain.BucketService) *BucketHandler {
	return &BucketHandler{service: service}
}

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

func (bh *BucketHandler) FindAllFiles(c *gin.Context) {
	bucketName := c.Param("bucket")

	files, err := bh.service.FindAllFiles(bucketName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{
		"files": files,
	})
}

func (bh *BucketHandler) Delete(c *gin.Context) {
	bucketName := c.Param("bucket")

	err := bh.service.Remove(bucketName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
