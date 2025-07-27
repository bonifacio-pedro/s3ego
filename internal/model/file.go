// Package model contains the data models used in the application.
package model

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

// File represents a file stored within a bucket in the S3 emulator.
// It contains an ID, unique key, raw data, metadata, and the ID of the bucket it belongs to.
type File struct {
	ID           int       `json:"id" db:"id"`                       // Unique identifier of the file in the database
	Key          string    `json:"key" db:"key"`                     // Unique key of the file (usually bucketName/filename)
	Data         []byte    `json:"data" db:"data"`                   // Raw binary data of the file
	BucketID     uint      `json:"bucket_id" db:"bucket_id"`         // Foreign key referencing the bucket this file belongs to
	ETag         string    `json:"etag" db:"etag"`                   // MD5 hash of the file content for integrity
	ContentType  string    `json:"content_type" db:"content_type"`   // MIME type of the file
	Size         int64     `json:"size" db:"size"`                   // Size of the file in bytes
	CreatedAt    time.Time `json:"created_at" db:"created_at"`       // Timestamp when file was created
	LastModified time.Time `json:"last_modified" db:"last_modified"` // Timestamp when file was last modified
}

// NewFile creates a new File instance given the file data, bucket, and file name.
// It generates the Key by combining the bucket name and file name in the format "bucketName/fileName".
// It also calculates metadata like ETag, ContentType, and Size automatically.
func NewFile(data []byte, bucket Bucket, fileName string) File {
	now := time.Now()

	file := File{
		Data:         data,
		BucketID:     uint(bucket.ID),
		Key:          fmt.Sprintf("%s/%s", bucket.Name, fileName),
		ETag:         calculateETag(data),
		ContentType:  detectContentType(data, fileName),
		Size:         int64(len(data)),
		CreatedAt:    now,
		LastModified: now,
	}

	return file
}

// calculateETag calculates the MD5 hash of file data (S3-style ETag)
func calculateETag(data []byte) string {
	return fmt.Sprintf("%x", md5.Sum(data))
}

// detectContentType tries to detect the content type of a file based on its data and extension.
// It uses http.DetectContentType first, and falls back to file extension if the result is generic.
func detectContentType(data []byte, filename string) string {
	contentType := http.DetectContentType(data)
	if contentType != "application/octet-stream" {
		return contentType
	}

	// fallback based on extension
	ext := strings.ToLower(filepath.Ext(filename))
	extToContentType := map[string]string{
		".json": "application/json",
		".xml":  "application/xml",
		".txt":  "text/plain",
		".html": "text/html",
		".htm":  "text/html",
		".css":  "text/css",
		".js":   "application/javascript",
		".csv":  "text/csv",
		".pdf":  "application/pdf",
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".svg":  "image/svg+xml",
		".webp": "image/webp",
		".zip":  "application/zip",
		".gz":   "application/gzip",
	}

	if fallback, ok := extToContentType[ext]; ok {
		return fallback
	}

	return "application/octet-stream"
}
