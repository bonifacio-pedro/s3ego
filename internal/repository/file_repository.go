package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/bonifacio-pedro/s3ego/internal/model"
	"log"
)

type FileRepository struct {
	db *sql.DB
}

func NewFileRepository(db *sql.DB) FileRepository {
	return FileRepository{db: db}
}

func (fr *FileRepository) GetFileByKey(key string) (*model.File, error) {
	row := fr.db.QueryRow("SELECT id, key, data, bucket_id FROM files WHERE key = ?", key)
	var f model.File

	if err := row.Scan(&f.ID, &f.Key, &f.Data, &f.BucketID); err != nil {
		return nil, err
	}

	return &f, nil
}

func (fr *FileRepository) CreateFile(file *model.File) error {
	if searchFile, _ := fr.GetFileByKey(file.Key); searchFile != nil {
		return errors.New("file with this name already exists")
	}

	_, err := fr.db.Exec("INSERT INTO files (key, data, bucket_id) VALUES (?, ?, ?)", file.Key, file.Data, file.BucketID)
	if err != nil {
		return err
	}
	log.Println(fmt.Sprintf("[S3-EMULATOR] File (key: %s) uploaded:", file.Key))

	return nil
}
