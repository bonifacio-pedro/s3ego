// Package main contains the entry point for running the S3EGO emulator.
//
// This file is intended primarily for development and local testing purposes.
// It initializes the necessary components such as the in-memory database and
// starts the application server on the configured port.
package main

import (
	"github.com/bonifacio-pedro/s3ego/internal/app"
	"github.com/bonifacio-pedro/s3ego/internal/config"
)

// main initializes the database connection, creates the application instance,
// and starts the HTTP server to handle incoming requests for the S3 emulator.
func main() {
	db := config.ConfigDatabase()
	defer db.Close()

	newApp := app.NewApp(db)
	newApp.Run()
}
