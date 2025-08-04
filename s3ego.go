// Package s3ego provides a lightweight S3 bucket emulator built in Go.
//
// S3EGO allows you to create buckets, upload files, and retrieve them
// via simple REST endpoints or by using it as a Go module (library).
//
// The emulator is designed primarily for local development and testing,
// storing data in an in-memory SQLite database for quick and easy setup.
//
// This project aims to facilitate unit testing and local development scenarios
// that require an S3-like interface without needing access to actual cloud storage.
//
// Note: The project is still under active development, and some features are yet to be added,
// such as advanced bucket policies, authentication, and multipart uploads.
package s3ego

import (
	"github.com/bonifacio-pedro/s3ego/internal/app"
	"github.com/bonifacio-pedro/s3ego/internal/config"
	"github.com/bonifacio-pedro/s3ego/internal/domain"
)

// S3EGO is the main struct exposing the bucket and file services for the emulator.
type S3EGO struct {
	Bucket domain.BucketService
	File   domain.FileService
}

// Start initializes the emulator by configuring the in-memory database and
// creating the application with all its services.
//
// Returns a pointer to an S3EGO instance that gives access to the bucket and file services.
func Start() *S3EGO {
	db := config.ConfigDatabase()
	newApp := app.NewApp(db)

	return &S3EGO{
		Bucket: newApp.BucketService,
		File:   newApp.FileService,
	}
}
