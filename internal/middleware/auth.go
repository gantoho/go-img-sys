package middleware

import (
	"net/http"

	"github.com/gantoho/go-img-sys/pkg/auth"
	"github.com/gantoho/go-img-sys/pkg/utils"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates API key
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get API key from header
		apiKey := ctx.GetHeader("X-API-Key")
		if apiKey == "" {
			// Try from query parameter
			apiKey = ctx.Query("api_key")
		}

		if apiKey == "" {
			utils.CustomResponse(ctx, http.StatusUnauthorized, "missing API key", nil)
			ctx.Abort()
			return
		}

		// Validate the API key
		keyManager := auth.GetManager()
		if !keyManager.ValidateKey(apiKey) {
			utils.CustomResponse(ctx, http.StatusUnauthorized, "invalid or expired API key", nil)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

// OptionalAuthMiddleware validates API key but allows requests without it
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		apiKey := ctx.GetHeader("X-API-Key")
		if apiKey == "" {
			apiKey = ctx.Query("api_key")
		}

		if apiKey != "" {
			keyManager := auth.GetManager()
			if !keyManager.ValidateKey(apiKey) {
				utils.CustomResponse(ctx, http.StatusUnauthorized, "invalid or expired API key", nil)
				ctx.Abort()
				return
			}
		}

		ctx.Next()
	}
}
