package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/bonifacio-pedro/s3ego/internal/model"
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

	return nil
}

func (fr *FileRepository) Remove(key string) error {
	_, err := fr.db.Exec("DELETE FROM files WHERE key=?", key)
	if err != nil {
		return fmt.Errorf("error deleting file DB row from files: %w", err)
	}

	return nil
}

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
