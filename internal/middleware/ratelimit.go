package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gantoho/go-img-sys/pkg/utils"
	"github.com/gin-gonic/gin"
)

// RateLimiter implements token bucket algorithm
type RateLimiter struct {
	mu              sync.RWMutex
	requestsPerSec  int
	concurrentLimit int
	tokens          map[string]float64
	lastRefill      map[string]time.Time
	concurrent      map[string]int
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(requestsPerSec, concurrentLimit int) *RateLimiter {
	return &RateLimiter{
		requestsPerSec:  requestsPerSec,
		concurrentLimit: concurrentLimit,
		tokens:          make(map[string]float64),
		lastRefill:      make(map[string]time.Time),
		concurrent:      make(map[string]int),
	}
}

// Allow checks if a request from the given IP is allowed
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	// Check concurrent connections
	if rl.concurrent[ip] >= rl.concurrentLimit {
		return false
	}

	// Initialize if first request from this IP
	if _, exists := rl.tokens[ip]; !exists {
		rl.tokens[ip] = float64(rl.requestsPerSec)
		rl.lastRefill[ip] = now
		rl.concurrent[ip] = 0
	}

	// Refill tokens based on elapsed time
	elapsed := now.Sub(rl.lastRefill[ip]).Seconds()
	rl.tokens[ip] = min(float64(rl.requestsPerSec), rl.tokens[ip]+elapsed*float64(rl.requestsPerSec))
	rl.lastRefill[ip] = now

	// Check if token is available
	if rl.tokens[ip] >= 1 {
		rl.tokens[ip]--
		rl.concurrent[ip]++
		return true
	}

	return false
}

// Release decrements the concurrent counter
func (rl *RateLimiter) Release(ip string) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if count, exists := rl.concurrent[ip]; exists && count > 0 {
		rl.concurrent[ip]--
	}
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

var rateLimiter = NewRateLimiter(100, 10) // 100 requests/sec, 10 concurrent connections per IP

// RateLimitMiddleware applies rate limiting to requests
func RateLimitMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()

		if !rateLimiter.Allow(ip) {
			utils.CustomResponse(ctx, http.StatusTooManyRequests, "rate limit exceeded", nil)
			ctx.Abort()
			return
		}

		ctx.Next()
		rateLimiter.Release(ip)
	}
}
