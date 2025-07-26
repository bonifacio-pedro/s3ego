package repository

import (
	"database/sql"
	"errors"
	"github.com/bonifacio-pedro/s3ego/internal/model"
	"log"
)

type BucketRepository struct {
	db *sql.DB
}

func NewBucketRepository(db *sql.DB) *BucketRepository {
	return &BucketRepository{db: db}
}

func (br *BucketRepository) CreateBucket(bucket *model.Bucket) error {
	if searchBucket, _ := br.GetBucketByName(bucket.Name); searchBucket != nil {
		return errors.New("bucket with this name already exists")
	}

	_, err := br.db.Exec("INSERT INTO buckets (name, url) VALUES (?, ?)", bucket.Name, bucket.Url)
	if err != nil {
		return err
	}
	log.Println("[S3EGO] Bucket created:", bucket.Name)

	return nil
}

func (br *BucketRepository) GetBucketByUrl(url string) (*model.Bucket, error) {
	row := br.db.QueryRow("SELECT id, name, url FROM buckets WHERE url = ?", url)
	var bucket model.Bucket

	if err := row.Scan(&bucket.ID, &bucket.Name, &bucket.Url); err != nil {
		return nil, err
	}
	return &bucket, nil
}

func (br *BucketRepository) GetBucketByName(bucketName string) (*model.Bucket, error) {
	row := br.db.QueryRow("SELECT id, name, url FROM buckets WHERE name = ?", bucketName)
	var bucket model.Bucket

	if err := row.Scan(&bucket.ID, &bucket.Name, &bucket.Url); err != nil {
		return nil, err
	}
	return &bucket, nil
}

func (br *BucketRepository) GetFiles(bucketID int) ([]string, error) {
	rows, err := br.db.Query("SELECT key FROM files WHERE bucket_id = ?", bucketID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	keys := make([]string, 0)
	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return keys, nil
}
