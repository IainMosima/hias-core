package middleware

import (
	"github.com/bitbiz/hias-core/shared/auth"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
)

func AuditMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		// Log after request completes
		payload, exists := ctx.Get(AuthPayloadKey)
		if !exists {
			return
		}

		authPayload, ok := payload.(*auth.Payload)
		if !ok {
			return
		}

		statusCode := ctx.Writer.Status()
		if statusCode >= 200 && statusCode < 300 {
			method := ctx.Request.Method
			if method == "POST" || method == "PUT" || method == "PATCH" || method == "DELETE" {
				utils.LogInfo("AUDIT: user=%s role=%s method=%s path=%s status=%d ip=%s",
					authPayload.UserID, authPayload.Role, method, ctx.Request.URL.Path, statusCode, ctx.ClientIP())
			}
		}
	}
}
