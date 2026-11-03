// Package impl provides concrete implementations of repositories.
package impl

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/bonifacio-pedro/s3ego/internal/model"
	"github.com/bonifacio-pedro/s3ego/internal/repository"
)

type fileRepository struct {
	db *sql.DB
}

// NewFileRepository creates a new FileRepository backed by SQLite.
func NewFileRepository(db *sql.DB) repository.FileRepository {
	return &fileRepository{db: db}
}

func (fr *fileRepository) Save(file *model.File) error {
	query := `
		INSERT INTO files (
			key, bucket_id, etag, content_type, size, file_path, created_at, last_modified
		) VALUES (?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		ON CONFLICT(bucket_id, key) DO UPDATE SET
			etag = excluded.etag,
			content_type = excluded.content_type,
			size = excluded.size,
			file_path = excluded.file_path,
			last_modified = CURRENT_TIMESTAMP
	`
	res, err := fr.db.Exec(
		query,
		file.Key,
		file.BucketID,
		file.ETag,
		file.ContentType,
		file.Size,
		file.FilePath,
	)
	if err != nil {
		return fmt.Errorf("error inserting/updating file metadata: %w", err)
	}

	id, err := res.LastInsertId()
	if err == nil && id > 0 {
		file.ID = int(id)
	}
	return nil
}

func (fr *fileRepository) Get(bucketID int, key string) (*model.File, error) {
	row := fr.db.QueryRow(`
		SELECT id, key, bucket_id, etag, content_type, size, file_path, created_at, last_modified
		FROM files
		WHERE bucket_id = ? AND key = ?
	`, bucketID, key)

	var (
		f            model.File
		createdAtStr string
		modifiedStr  string
	)

	err := row.Scan(
		&f.ID,
		&f.Key,
		&f.BucketID,
		&f.ETag,
		&f.ContentType,
		&f.Size,
		&f.FilePath,
		&createdAtStr,
		&modifiedStr,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("NoSuchKey")
		}
		return nil, fmt.Errorf("error querying file metadata: %w", err)
	}

	if t, err := time.Parse("2006-01-02 15:04:05", modifiedStr); err == nil {
		f.LastModified = t
	} else {
		f.LastModified = time.Now()
	}

	if t, err := time.Parse("2006-01-02 15:04:05", createdAtStr); err == nil {
		f.CreatedAt = t
	} else {
		f.CreatedAt = time.Now()
	}

	return &f, nil
}

func (fr *fileRepository) Delete(bucketID int, key string) (*model.File, error) {
	f, err := fr.Get(bucketID, key)
	if err != nil {
		return nil, err
	}

	_, err = fr.db.Exec("DELETE FROM files WHERE bucket_id = ? AND key = ?", bucketID, key)
	if err != nil {
		return nil, fmt.Errorf("error deleting file metadata: %w", err)
	}

	return f, nil
}
