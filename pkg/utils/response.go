package utils

import (
	"net/http"
	"sync"
	"time"

	"github.com/gantoho/go-img-sys/pkg/errors"
	"github.com/gin-gonic/gin"
)

const Version = "1.0.0"

type ResponseMetadata struct {
	Version   string `json:"version"`
	Timestamp int64  `json:"timestamp"`
	Duration  int64  `json:"duration_ms"`
}

type Response struct {
	Code     int              `json:"code"`
	Message  string           `json:"message"`
	Data     interface{}      `json:"data,omitempty"`
	Metadata ResponseMetadata `json:"metadata"`
}

var (
	startTimesMu sync.RWMutex
	startTimes   = make(map[*gin.Context]time.Time)
)

func RecordStartTime(ctx *gin.Context) {
	startTimesMu.Lock()
	defer startTimesMu.Unlock()
	startTimes[ctx] = time.Now()
}

func GetDuration(ctx *gin.Context) int64 {
	startTimesMu.Lock()
	defer startTimesMu.Unlock()
	if startTime, ok := startTimes[ctx]; ok {
		delete(startTimes, ctx)
		return time.Since(startTime).Milliseconds()
	}
	return 0
}

func SuccessResponse(ctx *gin.Context, data interface{}) {
	response := Response{
		Code:    http.StatusOK,
		Message: "success",
		Data:    data,
		Metadata: ResponseMetadata{
			Version:   Version,
			Timestamp: time.Now().Unix(),
			Duration:  GetDuration(ctx),
		},
	}
	ctx.JSON(http.StatusOK, response)
}

func ErrorResponse(ctx *gin.Context, err *errors.AppError) {
	response := Response{
		Code:    err.Code,
		Message: err.Message,
		Metadata: ResponseMetadata{
			Version:   Version,
			Timestamp: time.Now().Unix(),
			Duration:  GetDuration(ctx),
		},
	}
	ctx.JSON(err.Code, response)
}

func CustomResponse(ctx *gin.Context, code int, message string, data interface{}) {
	response := Response{
		Code:    code,
		Message: message,
		Data:    data,
		Metadata: ResponseMetadata{
			Version:   Version,
			Timestamp: time.Now().Unix(),
			Duration:  GetDuration(ctx),
		},
	}
	ctx.JSON(code, response)
}
