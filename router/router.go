// router/router.go
package router

import (
	v1 "file-service/api/v1"
	"file-service/middleware"
	"file-service/objectstorage"
	"github.com/gin-gonic/gin"
)

func SetupRouter(client *objectstorage.BNBClient) *gin.Engine {

	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(middleware.ErrorHandlingMiddleware())
	ctrl := v1.NewController(client)

	apiV1 := r.Group("/api/v1")
	{
		// objectName: string
		// bucketName: string
		// folder: file
		apiV1.POST("/objects", ctrl.PostObject)

		// objectName: string
		// bucketName: string
		apiV1.GET("/objects", ctrl.GetObject)

		// bucketName: string
		apiV1.GET("/buckets/objects", ctrl.ListObjects)
		apiV1.GET("/helloWorld", ctrl.HelloWorld)
	}

	return r
}
