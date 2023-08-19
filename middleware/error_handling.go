package middleware

import (
	"errors"
	"file-service/util"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func ErrorHandlingMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if len(ctx.Errors) > 0 {
			err := ctx.Errors[0].Err
			log.Printf("ErrorHandlingMiddelware - error: %v", err)
			var customErr util.ErrorResponse

			if errors.As(err, &customErr) {
				ctx.JSON(customErr.Code, customErr)
			} else {
				ctx.JSON(http.StatusInternalServerError, util.NewErrorResponse(http.StatusInternalServerError, 0, "Internal Server Error", err.Error(), err))
			}
		}
	}
}
