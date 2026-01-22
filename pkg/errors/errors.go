package errors

import "net/http"

type AppError struct {
	Code    int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

// Common errors
var (
	ErrFileNotFound    = &AppError{Code: http.StatusNotFound, Message: "file not found"}
	ErrFileUploadFail  = &AppError{Code: http.StatusInternalServerError, Message: "file upload failed"}
	ErrInvalidFileType = &AppError{Code: http.StatusBadRequest, Message: "invalid file type"}
	ErrFileTooLarge    = &AppError{Code: http.StatusBadRequest, Message: "file too large"}
	ErrInternalServer  = &AppError{Code: http.StatusInternalServerError, Message: "internal server error"}
	ErrDirectoryFail   = &AppError{Code: http.StatusInternalServerError, Message: "directory operation failed"}
	ErrNoFiles         = &AppError{Code: http.StatusNotFound, Message: "no files found"}
)

func NewError(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

func NewErrorWithCause(code int, message string, err error) *AppError {
	return &AppError{Code: code, Message: message, Err: err}
}
