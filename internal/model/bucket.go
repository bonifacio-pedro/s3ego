// Package model contains the data models used in the application.
package model

import (
	"fmt"
)

// Bucket represents an S3 bucket in the emulator.
// It holds a unique identifier, name, URL, and associated files.
type Bucket struct {
	ID    int    `json:"id"`    // Unique identifier of the bucket in the database
	Name  string `json:"name"`  // Name of the bucket
	Url   string `json:"url"`   // Base URL of the bucket
	Files []File `json:"files"` // List of files contained in the bucket
}

// NewBucket creates and initializes a new Bucket instance with the given name.
// The bucket URL is generated in the format "s3ego:7777//<bucketName>".
// It initializes the Files slice as empty.
func NewBucket(bucketName string) Bucket {
	bucket := new(Bucket)
	bucket.Name = bucketName
	bucket.Url = fmt.Sprintf("s3ego:7777//%s", bucketName)
	bucket.Files = make([]File, 0)

	return *bucket
}
