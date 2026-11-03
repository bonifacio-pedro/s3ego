// Package model contains the data models used in the application.
package model

import (
	"time"
)

// Bucket represents an S3 bucket in the emulator.
type Bucket struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// NewBucket creates and initializes a new Bucket instance with the given name.
func NewBucket(bucketName string) Bucket {
	return Bucket{
		Name:      bucketName,
		CreatedAt: time.Now(),
	}
}
