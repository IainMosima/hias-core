package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
)

type rateLimiter struct {
	requests map[string][]time.Time
	mu       sync.Mutex
	limit    int
	window   time.Duration
}

func newRateLimiter(limit int, window time.Duration) *rateLimiter {
	return &rateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

func (rl *rateLimiter) allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)

	// Clean old entries
	valid := make([]time.Time, 0)
	for _, t := range rl.requests[key] {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}

	if len(valid) >= rl.limit {
		rl.requests[key] = valid
		return false
	}

	rl.requests[key] = append(valid, now)
	return true
}

func RateLimitMiddleware(requestsPerMinute int) gin.HandlerFunc {
	limiter := newRateLimiter(requestsPerMinute, time.Minute)

	return func(ctx *gin.Context) {
		key := ctx.ClientIP()
		if !limiter.allow(key) {
			utils.RespondError(ctx, http.StatusTooManyRequests, "Rate limit exceeded")
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
