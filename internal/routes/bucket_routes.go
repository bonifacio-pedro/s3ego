package routes

import (
	"database/sql"
	"github.com/bonifacio-pedro/s3ego/internal/handlers"
	"github.com/gin-gonic/gin"
)

func HandleRequests(db *sql.DB) *gin.Engine {
	r := gin.Default()

	bucketHandler := handlers.NewBucketHandler(db)
	fileHandler := handlers.NewFileHandler(db)

	r.POST("/bucket-emulator/new-bucket/:name", bucketHandler.CreateBucket)
	r.POST("/bucket-emulator/upload-file/:bucket", fileHandler.UploadFile)
	r.GET("/bucket-emulator/list-files/:bucket", bucketHandler.FindAllFilesInABucket)

	r.GET("/bucket-emulator/get-file/:bucket/*key", fileHandler.GetFile)

	return r
}
