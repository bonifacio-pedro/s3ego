package model

import (
	"fmt"
)

type Bucket struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Url   string `json:"url"`
	Files []File `json:"files"`
}

func NewBucket(bucketName string) Bucket {
	bucket := new(Bucket)
	bucket.Name = bucketName
	bucket.Url = fmt.Sprintf("s3ego:7777//%s", bucketName)
	bucket.Files = make([]File, 0)

	return *bucket
}
