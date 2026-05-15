package middleware

import (
	"net/http"

	"github.com/bitbiz/hias-core/domains/claims/repository"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

const (
	APIKeyHeader    = "X-API-Key"
	APISecretHeader = "X-API-Secret"
	APIPartnerKey   = "api_partner"
)

func APIKeyAuthMiddleware(partnerRepo repository.APIPartnerRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		apiKey := ctx.GetHeader(APIKeyHeader)
		if apiKey == "" {
			utils.RespondError(ctx, http.StatusUnauthorized, "X-API-Key header is required")
			ctx.Abort()
			return
		}

		apiSecret := ctx.GetHeader(APISecretHeader)
		if apiSecret == "" {
			utils.RespondError(ctx, http.StatusUnauthorized, "X-API-Secret header is required")
			ctx.Abort()
			return
		}

		partner, err := partnerRepo.GetByAPIKey(ctx.Request.Context(), apiKey)
		if err != nil {
			utils.RespondError(ctx, http.StatusUnauthorized, "Invalid API key")
			ctx.Abort()
			return
		}

		if !partner.IsActive {
			utils.RespondError(ctx, http.StatusForbidden, "API partner is deactivated")
			ctx.Abort()
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(partner.APISecretHash), []byte(apiSecret)); err != nil {
			utils.RespondError(ctx, http.StatusUnauthorized, "Invalid API secret")
			ctx.Abort()
			return
		}

		ctx.Set(APIPartnerKey, partner)
		ctx.Next()
	}
}
