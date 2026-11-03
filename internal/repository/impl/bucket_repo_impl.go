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

type bucketRepository struct {
	db *sql.DB
}

// NewBucketRepository creates a new BucketRepository backed by SQLite.
func NewBucketRepository(db *sql.DB) repository.BucketRepository {
	return &bucketRepository{db: db}
}

func (br *bucketRepository) New(bucket *model.Bucket) error {
	res, err := br.db.Exec("INSERT INTO buckets (name) VALUES (?)", bucket.Name)
	if err != nil {
		return fmt.Errorf("failed to insert bucket: %w", err)
	}

	id, err := res.LastInsertId()
	if err == nil {
		bucket.ID = int(id)
	}
	return nil
}

func (br *bucketRepository) Remove(bucketID int) error {
	_, err := br.db.Exec("DELETE FROM files WHERE bucket_id = ?", bucketID)
	if err != nil {
		return fmt.Errorf("failed to remove bucket file records: %w", err)
	}

	_, err = br.db.Exec("DELETE FROM buckets WHERE id = ?", bucketID)
	if err != nil {
		return fmt.Errorf("failed to remove bucket: %w", err)
	}

	return nil
}

func (br *bucketRepository) ExistsByName(bucketName string) (bool, error) {
	var exists bool
	err := br.db.QueryRow("SELECT EXISTS(SELECT 1 FROM buckets WHERE name = ?)", bucketName).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if bucket exists: %w", err)
	}
	return exists, nil
}

func (br *bucketRepository) GetByName(bucketName string) (*model.Bucket, error) {
	row := br.db.QueryRow("SELECT id, name, created_at FROM buckets WHERE name = ?", bucketName)
	var (
		b         model.Bucket
		createdAt string
	)

	if err := row.Scan(&b.ID, &b.Name, &createdAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("bucket not found")
		}
		return nil, fmt.Errorf("error querying bucket: %w", err)
	}

	t, err := time.Parse("2006-01-02 15:04:05", createdAt)
	if err == nil {
		b.CreatedAt = t
	} else {
		b.CreatedAt = time.Now()
	}

	return &b, nil
}

func (br *bucketRepository) ListObjects(bucketID int, prefix string) ([]model.File, error) {
	query := "SELECT id, key, bucket_id, etag, content_type, size, file_path, created_at, last_modified FROM files WHERE bucket_id = ?"
	args := []any{bucketID}

	if prefix != "" {
		query += " AND key LIKE ?"
		args = append(args, prefix+"%")
	}
	query += " ORDER BY key ASC"

	rows, err := br.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query bucket objects: %w", err)
	}
	defer rows.Close()

	files := make([]model.File, 0)
	for rows.Next() {
		var (
			f            model.File
			createdAtStr string
			modifiedStr  string
		)
		err := rows.Scan(
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
			return nil, fmt.Errorf("error scanning object row: %w", err)
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

		files = append(files, f)
	}

	return files, rows.Err()
}

func (br *bucketRepository) ListBuckets() ([]model.Bucket, error) {
	rows, err := br.db.Query("SELECT id, name, created_at FROM buckets ORDER BY name ASC")
	if err != nil {
		return nil, fmt.Errorf("failed to list buckets: %w", err)
	}
	defer rows.Close()

	var buckets []model.Bucket
	for rows.Next() {
		var (
			b            model.Bucket
			createdAtStr string
		)
		if err := rows.Scan(&b.ID, &b.Name, &createdAtStr); err != nil {
			return nil, fmt.Errorf("failed scanning bucket row: %w", err)
		}
		if t, err := time.Parse("2006-01-02 15:04:05", createdAtStr); err == nil {
			b.CreatedAt = t
		}
		buckets = append(buckets, b)
	}
	return buckets, rows.Err()
}
