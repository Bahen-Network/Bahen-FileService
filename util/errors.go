package util

import (
	"github.com/gin-gonic/gin"
	"log"
)

func (e ErrorResponse) Error() string {
	return e.Message
}

type ErrorResponse struct {
	Code          int    `json:"code"`
	ErrorCode     int    `json:"errorCode"`
	Message       string `json:"message"`
	Detail        string `json:"detail,omitempty"`
	OriginalError error  `json:"original_error,omitempty"`
}

func NewErrorResponse(code int, errorCode int, message string, detail string, originalError error) ErrorResponse {
	return ErrorResponse{
		Code:          code,
		ErrorCode:     errorCode,
		Message:       message,
		Detail:        detail,
		OriginalError: originalError,
	}
}

func HandleErr(err error, msg string) bool {
	if err != nil {
		log.Printf("%s - error: %v", msg, err)
		return true
	}
	return false
}

// Error represents a standard application error.
type Error struct {
	Code    int
	Message string
}

// Error makes it compatible with `error` interface
func (e *Error) Error() string {
	return e.Message
}

var (
	// 1000-1099: Parameter errors
	ObjectNameArgumentError = &Error{Code: 1001, Message: "objectName is required in form data"}
	BucketNameArgumentError = &Error{Code: 1002, Message: "bucketName is required in form data"}

	// 1100-1199: File errors
	FormDataFileRetrieveError = &Error{Code: 1101, Message: "Error in retrieving the folder form file"}
	FileOpenError             = &Error{Code: 1102, Message: "Error in opening the uploaded file"}
	FileReadError             = &Error{Code: 1103, Message: "Reading uploaded file failed"}

	// 1200-1299: Encryption errors
	EncryptionError = &Error{Code: 1201, Message: "Encryption failed"}
	DecryptionError = &Error{Code: 1202, Message: "Decryption failed"}

	// 1300-1399: BNB Client errors
	BNBClientUploadError = &Error{Code: 1301, Message: "Failed to upload to BNBClient"}

	// 1400-1499: Zip operation errors
	ZipEntryCreationError = &Error{Code: 1401, Message: "Failed to create new zip entry"}
	ZipWriteError         = &Error{Code: 1402, Message: "Failed to write decrypted file to new zip"}
	ZipFinalizeError      = &Error{Code: 1403, Message: "Failed to finalize new zip file"}
)

func ReportError(ctx *gin.Context, originalError error, err *Error) {
	if err := ctx.Error(NewErrorResponse(500, err.Code, "Internal Error", err.Message, originalError)); err != nil {
		log.Printf("Unexpected error when adding error to context: %v", err)
	}
}
