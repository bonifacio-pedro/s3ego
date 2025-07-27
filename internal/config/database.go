// Package config provides configuration utilities for the S3EGO project,
// including initialization of the SQLite in-memory database schema.
package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite" // SQLite driver for database/sql
)

// ConfigDatabase initializes an in-memory SQLite database,
// sets up the necessary schema for buckets and files,
// and returns the active *sql.DB connection.
//
// If the database cannot be initialized or any schema fails to create,
// the function will log the error and terminate the application.
func ConfigDatabase() *sql.DB {
	// Open in-memory SQLite DB with shared cache
	db, err := sql.Open("sqlite", "file::memory:?mode=memory&cache=shared")
	if err != nil {
		log.Fatal(fmt.Sprintf("[S3EGO] Failed to open database: %s", err))
	}
	log.Println("[S3EGO] Started Database")

	// Create buckets table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS buckets (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT UNIQUE NOT NULL,
			url TEXT UNIQUE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		log.Fatal(fmt.Sprintf("[S3EGO] Failed to create buckets table: %s", err))
	}
	log.Println("[S3EGO] Buckets table initialized")

	// Create files table with S3 metadata support
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS files (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			key TEXT NOT NULL,
			data BLOB,
			bucket_id INTEGER NOT NULL,
			etag TEXT NOT NULL,
			content_type TEXT DEFAULT 'application/octet-stream',
			size INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			last_modified DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(bucket_id) REFERENCES buckets(id) ON DELETE CASCADE,
			UNIQUE(bucket_id, key)
		);
	`)
	if err != nil {
		log.Fatal(fmt.Sprintf("[S3EGO] Failed to create files table: %s", err))
	}
	log.Println("[S3EGO] Files table initialized")

	// Create indexes for better performance
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_files_bucket_key ON files(bucket_id, key);",
		"CREATE INDEX IF NOT EXISTS idx_files_etag ON files(etag);",
		"CREATE INDEX IF NOT EXISTS idx_buckets_name ON buckets(name);",
		"CREATE INDEX IF NOT EXISTS idx_files_last_modified ON files(last_modified);",
	}

	for _, indexSQL := range indexes {
		_, err = db.Exec(indexSQL)
		if err != nil {
			log.Printf("[S3EGO] Warning: Failed to create index: %s", err)
		}
	}
	log.Println("[S3EGO] Database indexes created")
	
	return db
}
