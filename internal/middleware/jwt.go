package middleware

import (
	"net/http"
	"time"

	ginjwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gantoho/go-img-sys/pkg/auth"
	"github.com/gantoho/go-img-sys/pkg/utils"
	"github.com/gin-gonic/gin"
)

// JWTMiddleware returns gin-jwt middleware with custom handlers
func JWTMiddleware(jwtManager *auth.JWTManager) (*ginjwt.GinJWTMiddleware, error) {
	return ginjwt.New(&ginjwt.GinJWTMiddleware{
		Realm:         "api",
		Key:           []byte(jwtManager.GetSecretKey()),
		Timeout:       jwtManager.GetDuration(),
		MaxRefresh:    jwtManager.GetDuration(),
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",

		// Login handler - this would be replaced with your actual login endpoint
		LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			c.JSON(code, gin.H{
				"code":   0,
				"msg":    "success",
				"token":  token,
				"expire": expire.Format("2006-01-02 15:04:05"),
			})
		},

		// Authenticator validates credentials and returns user info
		Authenticator: func(c *gin.Context) (interface{}, error) {
			// This is a placeholder - implement your actual authentication logic
			var req struct {
				Username string `json:"username" binding:"required"`
				Password string `json:"password" binding:"required"`
			}

			if err := c.ShouldBindJSON(&req); err != nil {
				return nil, ginjwt.ErrMissingLoginValues
			}

			// Simple authentication - in production, validate against database
			if req.Username == "admin" && req.Password == "admin123" {
				return &auth.Claims{
					UserID:   "1",
					Username: req.Username,
					Role:     "admin",
				}, nil
			}

			if req.Username == "user" && req.Password == "user123" {
				return &auth.Claims{
					UserID:   "2",
					Username: req.Username,
					Role:     "user",
				}, nil
			}

			return nil, ginjwt.ErrFailedAuthentication
		},

		// PayloadFunc returns the user payload
		PayloadFunc: func(data interface{}) ginjwt.MapClaims {
			if v, ok := data.(*auth.Claims); ok {
				return ginjwt.MapClaims{
					"user_id":  v.UserID,
					"username": v.Username,
					"role":     v.Role,
				}
			}
			return ginjwt.MapClaims{}
		},

		// IdentityKey is the key for the identity in context
		IdentityKey: "identity",

		// IdentityHandler extracts identity from token
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := ginjwt.ExtractClaims(c)
			return &auth.Claims{
				UserID:   claims["user_id"].(string),
				Username: claims["username"].(string),
				Role:     claims["role"].(string),
			}
		},

		// Unauthorized handler
		Unauthorized: func(c *gin.Context, code int, message string) {
			utils.CustomResponse(c, code, message, nil)
		},
	})
}

// Optional middleware that validates JWT but allows requests without it
func OptionalJWTMiddleware(jwtManager *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			token = c.Query("token")
		}

		if token != "" {
			// Remove "Bearer " prefix if present
			if len(token) > 7 && token[:7] == "Bearer " {
				token = token[7:]
			}

			// Try to validate token
			claims, err := jwtManager.ValidateToken(token)
			if err == nil {
				c.Set("user_id", claims.UserID)
				c.Set("username", claims.Username)
				c.Set("role", claims.Role)
				c.Set("authenticated", true)
			}
		}

		c.Next()
	}
}

// RequireRole middleware checks if user has required role
func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists || userRole != role {
			utils.CustomResponse(c, http.StatusForbidden, "insufficient permissions", nil)
			c.Abort()
			return
		}
		c.Next()
	}
}
