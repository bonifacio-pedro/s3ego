// Package repository provides database access methods for buckets and files.
package impl

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/bonifacio-pedro/s3ego/internal/model"
	"github.com/bonifacio-pedro/s3ego/internal/repository"
)

// BucketRepository handles operations on buckets and their files in the database.
type bucketRepository struct {
	db *sql.DB
}

// NewBucketRepository creates a new BucketRepository with the given database connection.
func NewBucketRepository(db *sql.DB) repository.BucketRepository {
	return &bucketRepository{db: db}
}

// New inserts a new bucket into the database.
// Returns an error if the insertion fails.
func (br *bucketRepository) New(bucket *model.Bucket) error {
	_, err := br.db.Exec("INSERT INTO buckets (name, url) VALUES (?, ?)", bucket.Name, bucket.Url)
	if err != nil {
		return fmt.Errorf("failed to insert bucket: %w", err)
	}

	return nil
}

// Remove deletes all files associated with a bucket and then removes the bucket itself from the database.
// Returns an error if the deletion fails.
func (br *bucketRepository) Remove(bucketID int) error {
	_, err := br.db.Exec("DELETE FROM files WHERE bucket_id = ?", bucketID)
	if err != nil {
		return fmt.Errorf("failed to remove bucket files: %w", err)
	}

	_, err = br.db.Exec("DELETE FROM buckets WHERE id = ?", bucketID)
	if err != nil {
		return fmt.Errorf("failed to remove bucket: %w", err)
	}

	return nil
}

// ExistsByName checks if a bucket with the given name exists in the database.
// Returns true if it exists, false otherwise, or an error if the query fails.
func (br *bucketRepository) ExistsByName(bucketName string) (bool, error) {
	var exists bool
	err := br.db.QueryRow("SELECT EXISTS(SELECT 1 FROM buckets WHERE name = ?)", bucketName).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if bucket exists: %w", err)
	}

	return exists, nil
}

// GetByName retrieves a bucket by its name.
// Returns a pointer to the Bucket model or an error if the bucket is not found.
func (br *bucketRepository) GetByName(bucketName string) (*model.Bucket, error) {
	row := br.db.QueryRow("SELECT id, name, url FROM buckets WHERE name = ?", bucketName)
	var bucket model.Bucket

	if err := row.Scan(&bucket.ID, &bucket.Name, &bucket.Url); err != nil {
		return nil, errors.New("bucket not found")
	}
	return &bucket, nil
}

// GetFiles retrieves all file keys associated with a given bucket ID.
// Returns a slice of file keys or an error if the query fails or no files are found.
func (br *bucketRepository) GetFiles(bucketID int) ([]string, error) {
	rows, err := br.db.Query("SELECT key FROM files WHERE bucket_id = ?", bucketID)
	if err != nil {
		return nil, errors.New("no files found with that bucket id")
	}
	defer rows.Close()

	keys := make([]string, 0)
	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			return nil, fmt.Errorf("error converting DB row to model in files keys iteration: %w", err)
		}
		keys = append(keys, key)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error in DB rows scanning: %w", err)
	}

	return keys, nil
}

// FileExists checks whether a file with the specified key exists within the given bucket name.
// Returns true if the file exists, false otherwise, or an error if the query fails.
func (br *bucketRepository) FileExists(bucketName string, key string) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM files f
			INNER JOIN buckets b ON f.bucket_id = b.id
			WHERE b.name = ? AND f.key = ?
		)
	`

	err := br.db.QueryRow(query, bucketName, key).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if file exists in bucket: %w", err)
	}
	return exists, nil
}
