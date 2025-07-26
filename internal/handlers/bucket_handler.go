package handlers

import (
	"database/sql"
	"github.com/bonifacio-pedro/s3ego/internal/model"
	"github.com/bonifacio-pedro/s3ego/internal/repository"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type BucketHandler struct {
	db *sql.DB
}

func NewBucketHandler(db *sql.DB) *BucketHandler {
	return &BucketHandler{db: db}
}

func (bh *BucketHandler) CreateBucket(c *gin.Context) {
	bucketName := c.Param("name")

	bucketRepository := repository.NewBucketRepository(bh.db)
	bucket := model.CreateBucket(bucketName)

	if err := bucketRepository.CreateBucket(bucket); err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "bucket created successfully",
		"url":     bucket.Url,
	})
}

func (bh *BucketHandler) FindAllFilesInABucket(c *gin.Context) {
	bucketName := c.Param("bucket")

	bucketRepository := repository.NewBucketRepository(bh.db)

	bucket, err := bucketRepository.GetBucketByName(bucketName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bucket not found"})
		return
	}

	files, err := bucketRepository.GetFiles(bucket.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error getting files"})
	}

	c.JSON(http.StatusOK, gin.H{
		"files": files,
	})
}
