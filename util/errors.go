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
	StatusCode int
	ErrorCode  int
	Message    string
}

// Error makes it compatible with `error` interface
func (e *Error) Error() string {
	return e.Message
}

var (
	// 1000-1099: Parameter errors
	ObjectNameArgumentError = &Error{StatusCode: 400, ErrorCode: 1001, Message: "objectName is required in form data"}
	BucketNameArgumentError = &Error{StatusCode: 400, ErrorCode: 1002, Message: "bucketName is required in form data"}

	// 1100-1199: File errors
	FormDataFileRetrieveError = &Error{StatusCode: 500, ErrorCode: 1101, Message: "Error in retrieving the folder form file"}
	FileOpenError             = &Error{StatusCode: 500, ErrorCode: 1102, Message: "Error in opening the uploaded file"}
	FileReadError             = &Error{StatusCode: 500, ErrorCode: 1103, Message: "Reading uploaded file failed"}

	// 1200-1299: Encryption errors
	EncryptionError = &Error{StatusCode: 500, ErrorCode: 1201, Message: "Encryption failed"}
	DecryptionError = &Error{StatusCode: 500, ErrorCode: 1202, Message: "Decryption failed"}

	// 1300-1399: BNB Client errors
	BNBClientUploadError = &Error{StatusCode: 500, ErrorCode: 1301, Message: "Failed to upload to BNBClient"}

	// 1400-1499: Zip operation errors
	ZipEntryCreationError = &Error{StatusCode: 500, ErrorCode: 1401, Message: "Failed to create new zip entry"}
	ZipWriteError         = &Error{StatusCode: 500, ErrorCode: 1402, Message: "Failed to write decrypted file to new zip"}
	ZipFinalizeError      = &Error{StatusCode: 500, ErrorCode: 1403, Message: "Failed to finalize new zip file"}
)

func ReportError(ctx *gin.Context, originalError error, err *Error) {
	if err := ctx.Error(NewErrorResponse(err.StatusCode, err.ErrorCode, "Internal Error", err.Message, originalError)); err != nil {
		log.Printf("Unexpected error when adding error to context: %v", err)
	}
}
