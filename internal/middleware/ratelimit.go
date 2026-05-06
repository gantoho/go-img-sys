package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gantoho/go-img-sys/internal/config"
	"github.com/gantoho/go-img-sys/pkg/utils"
	"github.com/gin-gonic/gin"
)

type RateLimiter struct {
	mu              sync.RWMutex
	requestsPerSec  int
	concurrentLimit int
	tokens          map[string]float64
	lastRefill      map[string]time.Time
	concurrent      map[string]int
	stopCh          chan struct{}
}

func NewRateLimiter(requestsPerSec, concurrentLimit int) *RateLimiter {
	rl := &RateLimiter{
		requestsPerSec:  requestsPerSec,
		concurrentLimit: concurrentLimit,
		tokens:          make(map[string]float64),
		lastRefill:      make(map[string]time.Time),
		concurrent:      make(map[string]int),
		stopCh:          make(chan struct{}),
	}
	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) Stop() {
	close(rl.stopCh)
}

func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	if rl.concurrent[ip] >= rl.concurrentLimit {
		return false
	}

	if _, exists := rl.tokens[ip]; !exists {
		rl.tokens[ip] = float64(rl.requestsPerSec)
		rl.lastRefill[ip] = now
		rl.concurrent[ip] = 0
	}

	elapsed := now.Sub(rl.lastRefill[ip]).Seconds()
	rl.tokens[ip] = min(float64(rl.requestsPerSec), rl.tokens[ip]+elapsed*float64(rl.requestsPerSec))
	rl.lastRefill[ip] = now

	if rl.tokens[ip] >= 1 {
		rl.tokens[ip]--
		rl.concurrent[ip]++
		return true
	}

	return false
}

func (rl *RateLimiter) Release(ip string) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if count, exists := rl.concurrent[ip]; exists && count > 0 {
		rl.concurrent[ip]--
	}
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-rl.stopCh:
			return
		case <-ticker.C:
			rl.mu.Lock()
			threshold := time.Now().Add(-30 * time.Minute)
			for ip, last := range rl.lastRefill {
				if last.Before(threshold) {
					delete(rl.tokens, ip)
					delete(rl.lastRefill, ip)
					delete(rl.concurrent, ip)
				}
			}
			rl.mu.Unlock()
		}
	}
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

var rateLimiter *RateLimiter

func RateLimitMiddleware() gin.HandlerFunc {
	cfg := config.GetConfig()
	if rateLimiter == nil {
		rateLimiter = NewRateLimiter(cfg.Rate.RequestsPerSec, cfg.Rate.ConcurrentLimit)
	}
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
