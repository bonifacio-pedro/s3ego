// Package config provides configuration utilities for the S3EGO project.
package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite" // SQLite pure-Go driver
)

// Config holds runtime configuration options for S3EGO.
type Config struct {
	DataDir string
	Port    string
}

// LoadConfig loads configuration from environment variables with sensible defaults.
func LoadConfig() Config {
	dataDir := os.Getenv("S3EGO_DATA_DIR")
	if dataDir == "" {
		dataDir = "./data"
	}

	port := os.Getenv("S3EGO_PORT")
	if port == "" {
		port = "7777"
	}

	return Config{
		DataDir: dataDir,
		Port:    port,
	}
}

// ConfigDatabase initializes an SQLite database persisted in the configured data directory,
// creates the schema for buckets and file metadata, and returns the *sql.DB connection.
func ConfigDatabase(dataDir string) (*sql.DB, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("[S3EGO] failed to create data directory: %w", err)
	}

	dbPath := filepath.Join(dataDir, "metadata.db")
	db, err := sql.Open("sqlite", fmt.Sprintf("file:%s?cache=shared&_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)", filepath.ToSlash(dbPath)))
	if err != nil {
		return nil, fmt.Errorf("[S3EGO] failed to open metadata database: %w", err)
	}

	// Create buckets table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS buckets (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT UNIQUE NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("[S3EGO] failed to create buckets table: %w", err)
	}

	// Create files metadata table (payload is stored on filesystem)
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS files (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			key TEXT NOT NULL,
			bucket_id INTEGER NOT NULL,
			etag TEXT NOT NULL,
			content_type TEXT DEFAULT 'application/octet-stream',
			size INTEGER DEFAULT 0,
			file_path TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			last_modified DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(bucket_id) REFERENCES buckets(id) ON DELETE CASCADE,
			UNIQUE(bucket_id, key)
		);
	`)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("[S3EGO] failed to create files table: %w", err)
	}

	// Create indexes for fast lookup
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_files_bucket_key ON files(bucket_id, key);",
		"CREATE INDEX IF NOT EXISTS idx_files_etag ON files(etag);",
		"CREATE INDEX IF NOT EXISTS idx_buckets_name ON buckets(name);",
		"CREATE INDEX IF NOT EXISTS idx_files_last_modified ON files(last_modified);",
	}

	for _, indexSQL := range indexes {
		if _, err := db.Exec(indexSQL); err != nil {
			log.Printf("[S3EGO] Warning: Failed to create index: %s", err)
		}
	}

	log.Printf("[S3EGO] Metadata database initialized at %s", dbPath)
	return db, nil
}
