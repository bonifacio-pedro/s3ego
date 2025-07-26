package repository

import (
	"database/sql"
	"fmt"
	"github.com/bonifacio-pedro/s3ego/internal/model"
	"log"
)

type FileRepository struct {
	db *sql.DB
}

func NewFileRepository(db *sql.DB) *FileRepository {
	return &FileRepository{db: db}
}

func (fr *FileRepository) New(file *model.File) error {
	_, err := fr.db.Exec("INSERT INTO files (key, data, bucket_id) VALUES (?, ?, ?)", file.Key, file.Data, file.BucketID)
	if err != nil {
		return fmt.Errorf("error inserting file DB row into files: %w", err)
	}
	log.Println(fmt.Sprintf("[S3-EMULATOR] File (key: %s) uploaded:", file.Key))

	return nil
}

func (fr *FileRepository) ExistsByKey(key string) (bool, error) {
	var exists bool
	err := fr.db.QueryRow("SELECT EXISTS(SELECT 1 FROM files WHERE key = ?)", key).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if file exists: %w", err)
	}
	return exists, nil
}

func (fr *FileRepository) GetByKey(key string) (*model.File, error) {
	row := fr.db.QueryRow("SELECT id, key, data, bucket_id FROM files WHERE key = ?", key)
	var f model.File

	if err := row.Scan(&f.ID, &f.Key, &f.Data, &f.BucketID); err != nil {
		return nil, fmt.Errorf("error scanning file DB row: %w", err)
	}

	return &f, nil
}
