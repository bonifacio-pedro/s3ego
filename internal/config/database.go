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
			name TEXT UNIQUE,
			url TEXT UNIQUE
		);
	`)
	if err != nil {
		log.Fatal(fmt.Sprintf("[S3EGO] Failed to create buckets table: %s", err))
	}
	log.Println("[S3EGO] Buckets table initialized")

	// Create files table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS files (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			key TEXT,
			data BLOB,
			bucket_id INTEGER,
			FOREIGN KEY(bucket_id) REFERENCES buckets(id)
		);
	`)
	if err != nil {
		log.Fatal(fmt.Sprintf("[S3EGO] Failed to create files table: %s", err))
	}
	log.Println("[S3EGO] Files table initialized")

	return db
}
