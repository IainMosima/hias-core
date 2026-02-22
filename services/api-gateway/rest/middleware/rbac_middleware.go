package middleware

import (
	"net/http"

	"github.com/bitbiz/hias-core/shared/auth"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
)

func RequireRole(roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		payload, exists := ctx.Get(AuthPayloadKey)
		if !exists {
			utils.RespondError(ctx, http.StatusUnauthorized, "Authentication required")
			ctx.Abort()
			return
		}

		authPayload, ok := payload.(*auth.Payload)
		if !ok {
			utils.RespondError(ctx, http.StatusUnauthorized, "Invalid auth payload")
			ctx.Abort()
			return
		}

		for _, role := range roles {
			if authPayload.Role == role {
				ctx.Next()
				return
			}
		}

		utils.RespondError(ctx, http.StatusForbidden, "Insufficient permissions")
		ctx.Abort()
	}
}

func RequirePermission(resource, action string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		payload, exists := ctx.Get(AuthPayloadKey)
		if !exists {
			utils.RespondError(ctx, http.StatusUnauthorized, "Authentication required")
			ctx.Abort()
			return
		}

		authPayload, ok := payload.(*auth.Payload)
		if !ok {
			utils.RespondError(ctx, http.StatusUnauthorized, "Invalid auth payload")
			ctx.Abort()
			return
		}

		required := resource + ":" + action
		for _, perm := range authPayload.Permissions {
			if perm == required {
				ctx.Next()
				return
			}
		}

		utils.RespondError(ctx, http.StatusForbidden, "Insufficient permissions")
		ctx.Abort()
	}
}
