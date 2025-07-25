package main

import (
	"s3ego/internal/app"
	"s3ego/internal/config"
)

func main() {
	db := config.ConfigDatabase()
	defer db.Close()

	app := app.NewApp(db)
	app.Run()
}
