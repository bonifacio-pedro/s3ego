package s3emulator

import (
	"database/sql"
	"s3ego/internal/app"
	"s3ego/internal/config"
)

type S3Emulator struct {
	App *app.App
	DB  *sql.DB
}

func Start() *S3Emulator {
	db := config.ConfigDatabase()
	newApp := app.NewApp(db)

	go newApp.Run()

	return &S3Emulator{
		App: newApp,
		DB:  db,
	}
}
