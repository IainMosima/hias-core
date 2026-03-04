package middleware

import (
	"encoding/json"
	"log"

	auditService "github.com/bitbiz/hias-core/domains/audit/service"
	"github.com/bitbiz/hias-core/shared/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AuditMiddleware(auditSvc auditService.AuditService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Only audit mutating requests
		method := ctx.Request.Method
		if method == "GET" || method == "OPTIONS" || method == "HEAD" {
			ctx.Next()
			return
		}

		ctx.Next()

		// After handler executes, log the audit event
		var userID uuid.UUID
		payload, exists := ctx.Get(AuthPayloadKey)
		if exists {
			if authPayload, ok := payload.(*auth.Payload); ok {
				userID, _ = uuid.Parse(authPayload.UserID)
			}
		}

		if auditSvc != nil && userID != uuid.Nil {
			newValue, _ := json.Marshal(map[string]interface{}{
				"method": method,
				"path":   ctx.Request.URL.Path,
				"status": ctx.Writer.Status(),
			})

			resp := auditSvc.LogEvent(
				ctx.Request.Context(),
				userID,
				"API",
				uuid.Nil,
				"API_CALL",
				nil,
				newValue,
				ctx.ClientIP(),
				ctx.Request.UserAgent(),
			)
			if resp.Error != nil {
				log.Printf("Audit middleware: failed to log event: %v", resp.Error)
			}
		}
	}
}
