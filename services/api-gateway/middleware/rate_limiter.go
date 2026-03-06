package middleware

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
)

const (
	defaultRateLimit  = 120
	defaultWindowSecs = 60
	cleanupInterval   = 5 * time.Minute
)

type clientEntry struct {
	timestamps []time.Time
	mu         sync.Mutex
}

type RateLimiter struct {
	clients    map[string]*clientEntry
	mu         sync.RWMutex
	limit      int
	windowSize time.Duration
}

func NewRateLimiter(limit int, windowSecs int) *RateLimiter {
	if limit <= 0 {
		limit = defaultRateLimit
	}
	if windowSecs <= 0 {
		windowSecs = defaultWindowSecs
	}

	rl := &RateLimiter{
		clients:    make(map[string]*clientEntry),
		limit:      limit,
		windowSize: time.Duration(windowSecs) * time.Second,
	}

	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		rl.mu.Lock()
		for ip, entry := range rl.clients {
			entry.mu.Lock()
			if len(entry.timestamps) == 0 || now.Sub(entry.timestamps[len(entry.timestamps)-1]) > rl.windowSize {
				delete(rl.clients, ip)
			}
			entry.mu.Unlock()
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) allow(ip string) (bool, int, time.Duration) {
	now := time.Now()
	windowStart := now.Add(-rl.windowSize)

	rl.mu.RLock()
	entry, exists := rl.clients[ip]
	rl.mu.RUnlock()

	if !exists {
		rl.mu.Lock()
		entry = &clientEntry{}
		rl.clients[ip] = entry
		rl.mu.Unlock()
	}

	entry.mu.Lock()
	defer entry.mu.Unlock()

	// Remove expired timestamps
	validIdx := 0
	for _, ts := range entry.timestamps {
		if ts.After(windowStart) {
			entry.timestamps[validIdx] = ts
			validIdx++
		}
	}
	entry.timestamps = entry.timestamps[:validIdx]

	remaining := rl.limit - len(entry.timestamps)
	if remaining <= 0 {
		retryAfter := entry.timestamps[0].Add(rl.windowSize).Sub(now)
		return false, 0, retryAfter
	}

	entry.timestamps = append(entry.timestamps, now)
	return true, remaining - 1, 0
}

func RateLimiterMiddleware(limit, windowSecs int) gin.HandlerFunc {
	rl := NewRateLimiter(limit, windowSecs)

	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()
		allowed, remaining, retryAfter := rl.allow(ip)

		ctx.Header("X-RateLimit-Limit", strconv.Itoa(rl.limit))
		ctx.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))

		if !allowed {
			ctx.Header("Retry-After", strconv.Itoa(int(retryAfter.Seconds())+1))
			utils.RespondError(ctx, http.StatusTooManyRequests, "Rate limit exceeded")
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
