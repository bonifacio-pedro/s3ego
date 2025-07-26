package repository

import (
	"database/sql"
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
