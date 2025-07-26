package routes

import (
	"github.com/bonifacio-pedro/s3ego/internal/transport/rest"
	"github.com/gin-gonic/gin"
)

type Router struct {
	rg            *gin.Engine
	bucketHandler *rest.BucketHandler
	fileHandler   *rest.FileHandler
}

func NewRouter(rg *gin.Engine, bucketHandler *rest.BucketHandler, fileHandler *rest.FileHandler) *Router {
	return &Router{rg: rg, bucketHandler: bucketHandler, fileHandler: fileHandler}
}

func (ro *Router) RegisterRoutes() {
	ro.rg.POST("/bucket-emulator/new-bucket/:name", ro.bucketHandler.Create)
	ro.rg.GET("/bucket-emulator/list-files/:bucket", ro.bucketHandler.FindAllFiles)
	ro.rg.DELETE("/bucket-emulator/delete/:bucket", ro.bucketHandler.Delete)
	ro.rg.DELETE("/bucket-emulator/delete/:bucket/*key", ro.bucketHandler.Delete)
	ro.rg.POST("/bucket-emulator/upload-file/:bucket", ro.fileHandler.New)
	ro.rg.GET("/bucket-emulator/get-file/:bucket/*key", ro.fileHandler.Get)
}
