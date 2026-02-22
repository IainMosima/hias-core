package utils

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func SetAuthCookie(ctx *gin.Context, name, value, domain string, duration time.Duration) {
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     name,
		Value:    value,
		Domain:   domain,
		Path:     "/",
		MaxAge:   int(duration.Seconds()),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})
}

func ClearAuthCookie(ctx *gin.Context, name, domain string) {
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     name,
		Value:    "",
		Domain:   domain,
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})
}
