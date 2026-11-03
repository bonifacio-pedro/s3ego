// Package model contains the data models used in the application.
package model

import (
	"mime"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

// File represents the metadata of a file stored within a bucket.
type File struct {
	ID           int       `json:"id" db:"id"`
	Key          string    `json:"key" db:"key"`
	BucketID     int       `json:"bucket_id" db:"bucket_id"`
	ETag         string    `json:"etag" db:"etag"`
	ContentType  string    `json:"content_type" db:"content_type"`
	Size         int64     `json:"size" db:"size"`
	FilePath     string    `json:"file_path" db:"file_path"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	LastModified time.Time `json:"last_modified" db:"last_modified"`
}

// DetectContentType determines the MIME type using initial bytes and file extension fallback.
func DetectContentType(initialBytes []byte, filename string) string {
	if len(initialBytes) > 0 {
		ct := http.DetectContentType(initialBytes)
		if ct != "application/octet-stream" {
			return ct
		}
	}

	ext := strings.ToLower(filepath.Ext(filename))
	if ext != "" {
		if mimeType := mime.TypeByExtension(ext); mimeType != "" {
			return mimeType
		}
	}

	fallbackMap := map[string]string{
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

	if val, ok := fallbackMap[ext]; ok {
		return val
	}

	return "application/octet-stream"
}
