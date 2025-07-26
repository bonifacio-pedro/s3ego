package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/bonifacio-pedro/s3ego/internal/model"
)

type BucketRepository struct {
	db *sql.DB
}

func NewBucketRepository(db *sql.DB) *BucketRepository {
	return &BucketRepository{db: db}
}

func (br *BucketRepository) New(bucket *model.Bucket) error {
	_, err := br.db.Exec("INSERT INTO buckets (name, url) VALUES (?, ?)", bucket.Name, bucket.Url)
	if err != nil {
		return fmt.Errorf("failed to insert bucket: %w", err)
	}

	return nil
}

func (br *BucketRepository) GetByUrl(url string) (*model.Bucket, error) {
	row := br.db.QueryRow("SELECT id, name, url FROM buckets WHERE url = ?", url)
	var bucket model.Bucket

	if err := row.Scan(&bucket.ID, &bucket.Name, &bucket.Url); err != nil {
		return nil, fmt.Errorf("error converting DB row to model: %w", err)
	}
	return &bucket, nil
}

func (br *BucketRepository) ExistsByName(bucketName string) (bool, error) {
	var exists bool
	err := br.db.QueryRow("SELECT EXISTS(SELECT 1 FROM buckets WHERE name = ?)", bucketName).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if bucket exists: %w", err)
	}

	return exists, nil
}

func (br *BucketRepository) GetByName(bucketName string) (*model.Bucket, error) {
	row := br.db.QueryRow("SELECT id, name, url FROM buckets WHERE name = ?", bucketName)
	var bucket model.Bucket

	if err := row.Scan(&bucket.ID, &bucket.Name, &bucket.Url); err != nil {
		return nil, errors.New("bucket not found")
	}
	return &bucket, nil
}

func (br *BucketRepository) GetFiles(bucketID int) ([]string, error) {
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
