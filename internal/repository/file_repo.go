// Package repository provides database access methods for files.
package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/bonifacio-pedro/s3ego/internal/model"
)

// FileRepository handles CRUD operations for files in the database.
type FileRepository struct {
	db *sql.DB
}

// NewFileRepository creates a new FileRepository with the given database connection.
func NewFileRepository(db *sql.DB) *FileRepository {
	return &FileRepository{db: db}
}

// New inserts a new file record into the files table.
// Returns an error if the insertion fails.
func (fr *FileRepository) New(file *model.File) error {
	_, err := fr.db.Exec("INSERT INTO files (key, data, bucket_id) VALUES (?, ?, ?)", file.Key, file.Data, file.BucketID)
	if err != nil {
		return fmt.Errorf("error inserting file DB row into files: %w", err)
	}

	return nil
}

// Remove deletes a file record from the files table by its key.
// Returns an error if the deletion fails.
func (fr *FileRepository) Remove(key string) error {
	_, err := fr.db.Exec("DELETE FROM files WHERE key=?", key)
	if err != nil {
		return fmt.Errorf("error deleting file DB row from files: %w", err)
	}

	return nil
}

// GetByKey retrieves a file from the database by its key.
// Returns the file model or an error if the file does not exist or scanning fails.
func (fr *FileRepository) GetByKey(key string) (*model.File, error) {
	row := fr.db.QueryRow("SELECT id, key, data, bucket_id FROM files WHERE key = ?", key)
	var f model.File

	if err := row.Scan(&f.ID, &f.Key, &f.Data, &f.BucketID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("file does not exist")
		}
		return nil, fmt.Errorf("error scanning file DB row: %w", err)
	}

	return &f, nil
}
