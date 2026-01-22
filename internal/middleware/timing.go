package middleware

import (
	"github.com/gantoho/go-img-sys/pkg/utils"
	"github.com/gin-gonic/gin"
)

// RequestTimingMiddleware records request start time for response timing
func RequestTimingMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		utils.RecordStartTime(ctx)
		ctx.Next()
	}
}
