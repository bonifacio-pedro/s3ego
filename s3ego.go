// Package s3ego provides an S3 wire-compatible emulator built in Go.
//
// S3EGO implements standard AWS S3 REST endpoints (Path-Style) and can be used
// as a standalone server or imported directly as a Go library for unit and integration testing.
package s3ego

import (
	"log"
	"os"

	"github.com/bonifacio-pedro/s3ego/internal/app"
	"github.com/bonifacio-pedro/s3ego/internal/domain"
	"github.com/bonifacio-pedro/s3ego/internal/storage"
)

// Option defines a functional option for configuring an S3EGO instance.
type Option func(*options)

type options struct {
	dataDir string
	isTemp  bool
}

// WithDataDir configures a custom data directory for persisting files and metadata.
func WithDataDir(dir string) Option {
	return func(o *options) {
		o.dataDir = dir
	}
}

// S3EGO is the main struct exposing services and application components.
type S3EGO struct {
	App     *app.App
	Bucket  domain.BucketService
	File    domain.FileService
	Storage storage.FileStorage
	dataDir string
	isTemp  bool
}

// Start initializes the emulator. If no custom data directory is specified via WithDataDir,
// an isolated temporary directory is created automatically for clean testing.
func Start(opts ...Option) *S3EGO {
	cfg := options{}
	for _, opt := range opts {
		opt(&cfg)
	}

	if cfg.dataDir == "" {
		envDir := os.Getenv("S3EGO_DATA_DIR")
		if envDir != "" {
			cfg.dataDir = envDir
		} else {
			tempDir, err := os.MkdirTemp("", "s3ego-test-*")
			if err != nil {
				log.Fatalf("[S3EGO] Failed to create temp directory: %s", err)
			}
			cfg.dataDir = tempDir
			cfg.isTemp = true
		}
	}

	application, err := app.NewApp(cfg.dataDir)
	if err != nil {
		log.Fatalf("[S3EGO] Failed to initialize application: %s", err)
	}

	return &S3EGO{
		App:     application,
		Bucket:  application.BucketService,
		File:    application.FileService,
		Storage: application.Storage,
		dataDir: cfg.dataDir,
		isTemp:  cfg.isTemp,
	}
}

// Close gracefully closes database connections and cleans up temporary resources.
func (s *S3EGO) Close() {
	if s.App != nil {
		_ = s.App.Close()
	}
	if s.isTemp && s.dataDir != "" {
		_ = os.RemoveAll(s.dataDir)
	}
}
