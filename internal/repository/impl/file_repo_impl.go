// Package impl provides concrete implementations of repositories.
package impl

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/bonifacio-pedro/s3ego/internal/model"
	"github.com/bonifacio-pedro/s3ego/internal/repository"
)

// FileRepository handles CRUD operations for files in the database.
type fileRepository struct {
	db *sql.DB
}

// NewFileRepository creates a new FileRepository with the given database connection.
func NewFileRepository(db *sql.DB) repository.FileRepository {
	return &fileRepository{db: db}
}

// New inserts a new file record into the files table.
// Returns an error if the insertion fails.
func (fr *fileRepository) New(file *model.File) error {
	_, err := fr.db.Exec(`
		INSERT INTO files (
			key, data, bucket_id, etag, content_type, size, created_at, last_modified
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		file.Key,
		file.Data,
		file.BucketID,
		file.ETag,
		file.ContentType,
		file.Size,
		file.CreatedAt,
		file.LastModified,
	)
	if err != nil {
		return fmt.Errorf("error inserting file DB row into files: %w", err)
	}

	return nil
}

// Remove deletes a file record from the files table by its key.
// Returns an error if the deletion fails.
func (fr *fileRepository) Remove(key string) error {
	_, err := fr.db.Exec("DELETE FROM files WHERE key=?", key)
	if err != nil {
		return fmt.Errorf("error deleting file DB row from files: %w", err)
	}

	return nil
}

// GetByKey retrieves a file from the database by its key.
// Returns the file model or an error if the file does not exist or scanning fails.
func (fr *fileRepository) GetByKey(key string) (*model.File, error) {
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
