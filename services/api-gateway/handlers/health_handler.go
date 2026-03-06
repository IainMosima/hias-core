package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Health godoc
// @Summary      Health check
// @Description  Returns the health status of the service
// @Tags         Health
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Router       /health [get]
func (h *HealthHandler) Health(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// Ready godoc
// @Summary      Readiness check
// @Description  Returns the readiness status of the service
// @Tags         Health
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Router       /ready [get]
func (h *HealthHandler) Ready(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "ready",
	})
}
