package middleware

import (
	"log"
	"net/http"

	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				utils.RespondError(ctx, http.StatusInternalServerError, "Internal server error")
				ctx.Abort()
			}
		}()
		ctx.Next()
	}
}

func RequestIDMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestID := ctx.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = ctx.GetHeader("X-Trace-ID")
		}
		if requestID == "" {
			requestID = generateRequestID()
		}
		ctx.Set("request_id", requestID)
		ctx.Header("X-Request-ID", requestID)
		ctx.Next()
	}
}

func generateRequestID() string {
	return "req-" + randomHex(8)
}

func randomHex(n int) string {
	const hex = "0123456789abcdef"
	b := make([]byte, n)
	for i := range b {
		b[i] = hex[i%len(hex)]
	}
	return string(b)
}
