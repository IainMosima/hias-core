package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware(allowedOrigins string) gin.HandlerFunc {
	origins := strings.Split(allowedOrigins, ",")

	return func(ctx *gin.Context) {
		origin := ctx.Request.Header.Get("Origin")

		for _, allowed := range origins {
			if strings.TrimSpace(allowed) == origin || strings.TrimSpace(allowed) == "*" {
				ctx.Header("Access-Control-Allow-Origin", origin)
				break
			}
		}

		ctx.Header("Access-Control-Allow-Credentials", "true")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		ctx.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return
		}

		ctx.Next()
	}
}
