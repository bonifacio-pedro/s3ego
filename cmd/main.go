package main

import (
	"github.com/bonifacio-pedro/s3ego/internal/app"
	"github.com/bonifacio-pedro/s3ego/internal/config"
)

func main() {
	db := config.ConfigDatabase()
	defer db.Close()

	newApp := app.NewApp(db)
	newApp.Run()
}
