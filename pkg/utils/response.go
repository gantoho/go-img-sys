package utils

import (
	"net/http"
	"time"

	"github.com/gantoho/go-img-sys/pkg/errors"
	"github.com/gin-gonic/gin"
)

// ResponseMetadata contains response metadata
type ResponseMetadata struct {
	Version   string `json:"version"`
	Timestamp int64  `json:"timestamp"`
	Duration  int64  `json:"duration_ms"`
}

// Response represents a unified API response
type Response struct {
	Code     int              `json:"code"`
	Message  string           `json:"message"`
	Data     interface{}      `json:"data,omitempty"`
	Metadata ResponseMetadata `json:"metadata"`
}

// startTime for tracking request duration
var (
	startTimes = make(map[*gin.Context]time.Time)
)

// RecordStartTime records the request start time
func RecordStartTime(ctx *gin.Context) {
	startTimes[ctx] = time.Now()
}

// GetDuration returns the request duration in milliseconds
func GetDuration(ctx *gin.Context) int64 {
	if startTime, ok := startTimes[ctx]; ok {
		delete(startTimes, ctx) // Clean up
		return time.Since(startTime).Milliseconds()
	}
	return 0
}

// SuccessResponse returns a success response with metadata
func SuccessResponse(ctx *gin.Context, data interface{}) {
	response := Response{
		Code:    http.StatusOK,
		Message: "success",
		Data:    data,
		Metadata: ResponseMetadata{
			Version:   "1.0.0",
			Timestamp: time.Now().Unix(),
			Duration:  GetDuration(ctx),
		},
	}
	ctx.JSON(http.StatusOK, response)
}

// ErrorResponse returns an error response with metadata
func ErrorResponse(ctx *gin.Context, err *errors.AppError) {
	response := Response{
		Code:    err.Code,
		Message: err.Message,
		Metadata: ResponseMetadata{
			Version:   "1.0.0",
			Timestamp: time.Now().Unix(),
			Duration:  GetDuration(ctx),
		},
	}
	ctx.JSON(err.Code, response)
}

// CustomResponse returns a custom response with metadata
func CustomResponse(ctx *gin.Context, code int, message string, data interface{}) {
	response := Response{
		Code:    code,
		Message: message,
		Data:    data,
		Metadata: ResponseMetadata{
			Version:   "1.0.0",
			Timestamp: time.Now().Unix(),
			Duration:  GetDuration(ctx),
		},
	}
	ctx.JSON(code, response)
}
