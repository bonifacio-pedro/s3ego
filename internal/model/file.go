// Package model contains the data models used in the application.
package model

import "fmt"

// File represents a file stored within a bucket in the S3 emulator.
// It contains an ID, unique key, raw data, and the ID of the bucket it belongs to.
type File struct {
	ID       int    `json:"id"`        // Unique identifier of the file in the database
	Key      string `json:"key"`       // Unique key of the file (usually bucketName/filename)
	Data     []byte `json:"data"`      // Raw binary data of the file
	BucketID uint   `json:"bucket_id"` // Foreign key referencing the bucket this file belongs to
}

// NewFile creates a new File instance given the file data, bucket, and file name.
// It generates the Key by combining the bucket name and file name in the format "bucketName/fileName".
func NewFile(data []byte, bucket Bucket, fileName string) File {
	file := new(File)
	file.Data = data
	file.BucketID = uint(bucket.ID)
	file.Key = fmt.Sprintf("%s/%s", bucket.Name, fileName)

	return *file
}
