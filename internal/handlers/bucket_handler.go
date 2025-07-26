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
	name := c.Param("name")

	bucketRepository := repository.NewBucketRepository(bh.db)
	bucketModel := model.CreateBucket(name)

	if err := bucketRepository.CreateBucket(bucketModel); err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "bucket created successfully",
		"url":     bucketModel.Url,
	})
}
