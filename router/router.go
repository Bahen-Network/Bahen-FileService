// router/router.go
package router

import (
	v1 "file-service/api/v1"
	"file-service/objectstorage"
	"github.com/gin-gonic/gin"
)

func SetupRouter(client *objectstorage.BNBClient) *gin.Engine {

	r := gin.Default()
	r.Use(gin.Logger())
	ctrl := v1.NewController(client)

	apiV1 := r.Group("/api/v1")
	{
		apiV1.POST("/upload", ctrl.UploadFolder)
		apiV1.GET("/download", ctrl.DownloadFolder)
		apiV1.GET("/helloWorld", ctrl.HelloWorld)
	}

	return r
}
