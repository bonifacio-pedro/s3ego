// Package main contains the entry point for running the S3EGO emulator server.
package main

import (
	"log"

	"github.com/bonifacio-pedro/s3ego/internal/app"
	"github.com/bonifacio-pedro/s3ego/internal/config"
)

func main() {
	cfg := config.LoadConfig()

	application, err := app.NewApp(cfg.DataDir)
	if err != nil {
		log.Fatalf("[S3EGO] Startup error: %s", err)
	}
	defer application.Close()

	if err := application.Run(cfg.Port); err != nil {
		log.Fatalf("[S3EGO] Server runtime error: %s", err)
	}
}
