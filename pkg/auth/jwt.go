package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Claims represents JWT claims structure
type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// JWTManager manages JWT token operations
type JWTManager struct {
	secretKey string
	duration  time.Duration
}

var jwtManager *JWTManager

// InitJWTManager initializes the JWT manager
func InitJWTManager(secretKey string, duration time.Duration) {
	jwtManager = &JWTManager{
		secretKey: secretKey,
		duration:  duration,
	}
}

// GetJWTManager returns the singleton JWT manager
func GetJWTManager() *JWTManager {
	if jwtManager == nil {
		// Default initialization if not initialized
		InitJWTManager("your-secret-key-change-this", 24*time.Hour)
	}
	return jwtManager
}

// GenerateToken generates a JWT token with claims
func (jm *JWTManager) GenerateToken(userID, username, role string) (string, error) {
	claims := &Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jm.duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jm.secretKey))
}

// ValidateToken validates a JWT token and returns claims
func (jm *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jm.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}

// RefreshToken generates a new token with extended expiry
func (jm *JWTManager) RefreshToken(tokenString string) (string, error) {
	claims, err := jm.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	// Create new token with same claims but extended expiry
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(jm.duration))
	claims.IssuedAt = jwt.NewNumericDate(time.Now())

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jm.secretKey))
}

// GetSecretKey returns the secret key (for gin-jwt middleware)
func (jm *JWTManager) GetSecretKey() string {
	return jm.secretKey
}

// GetDuration returns the token duration
func (jm *JWTManager) GetDuration() time.Duration {
	return jm.duration
}
