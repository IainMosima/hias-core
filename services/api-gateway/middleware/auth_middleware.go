package middleware

import (
	"net/http"
	"strings"

	"github.com/bitbiz/hias-core/shared/auth"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
)

const (
	AuthorizationHeaderKey  = "authorization"
	AuthorizationTypeBearer = "bearer"
	AuthPayloadKey          = "auth_payload"
)

func AuthMiddleware(tokenMaker auth.TokenMaker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(AuthorizationHeaderKey)

		if len(authorizationHeader) == 0 {
			utils.RespondError(ctx, http.StatusUnauthorized, "Authorization header is required")
			ctx.Abort()
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			utils.RespondError(ctx, http.StatusUnauthorized, "Invalid authorization header format")
			ctx.Abort()
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != AuthorizationTypeBearer {
			utils.RespondError(ctx, http.StatusUnauthorized, "Unsupported authorization type")
			ctx.Abort()
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			utils.RespondError(ctx, http.StatusUnauthorized, "Invalid or expired token")
			ctx.Abort()
			return
		}

		ctx.Set(AuthPayloadKey, payload)
		ctx.Set("user_id", payload.UserID)
		ctx.Set("role", payload.Role)
		ctx.Set("permissions", payload.Permissions)
		ctx.Next()
	}
}
