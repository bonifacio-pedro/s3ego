package config

import (
	"database/sql"
	"fmt"
	"log"
	_ "modernc.org/sqlite"
)

func ConfigDatabase() *sql.DB {
	db, err := sql.Open("sqlite", "file::memory:?mode=memory&cache=shared")
	if err != nil {
		log.Fatal(fmt.Sprintf("[S3EGO] %s"), err)
	}
	log.Println("[S3EGO] Started Database")

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS buckets (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT UNIQUE,
			url TEXT UNIQUE
		);
	`)
	if err != nil {
		log.Fatal(fmt.Sprintf("[S3EGO] %s"), err)
	}
	log.Println("[S3EGO] Buckets started")

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
		log.Fatal(fmt.Sprintf("[S3EGO] %s"), err)
	}
	log.Println("[S3EGO] Files table started")

	return db
}
