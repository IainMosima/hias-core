package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const RequestIDKey = "X-Request-ID"

func RequestIDMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestID := ctx.GetHeader(RequestIDKey)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx.Set(RequestIDKey, requestID)
		ctx.Header(RequestIDKey, requestID)
		ctx.Next()
	}
}
