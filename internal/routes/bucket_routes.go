package routes

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"s3ego/internal/handlers"
)

func HandleRequests(db *sql.DB) *gin.Engine {
	r := gin.Default()

	bucketHandler := handlers.NewBucketHandler(db)
	fileHandler := handlers.NewFileHandler(db)

	r.POST("/bucket-emulator/new-bucket/:name", bucketHandler.CreateBucket)
	r.POST("/bucket-emulator/upload-file/:bucket", fileHandler.UploadFile)
	r.GET("/bucket-emulator/get-file/:bucket/*key", fileHandler.GetFile)

	return r
}
