package handlers

import (
	"net/http"

	billingSchema "github.com/bitbiz/hias-core/domains/billing/schema"
	"github.com/bitbiz/hias-core/domains/billing/service"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PaymentHandler struct {
	paymentSvc service.PaymentService
}

func NewPaymentHandler(paymentSvc service.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentSvc: paymentSvc}
}

// InitiatePayment godoc
// @Summary      Initiate a payment
// @Description  Initiate a new payment transaction
// @Tags         Payments
// @Accept       json
// @Produce      json
// @Param        request body billingSchema.InitiatePaymentRequest true "Payment initiation payload"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/payments [post]
func (h *PaymentHandler) InitiatePayment(ctx *gin.Context) {
	var req billingSchema.InitiatePaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.paymentSvc.InitiatePayment(ctx.Request.Context(), req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ProcessWebhook godoc
// @Summary      Process M-Pesa webhook
// @Description  Receive and process an M-Pesa payment webhook callback
// @Tags         Payments
// @Accept       json
// @Produce      json
// @Param        request body interface{} true "Webhook payload"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Router       /api/v1/webhooks/mpesa [post]
func (h *PaymentHandler) ProcessWebhook(ctx *gin.Context) {
	var data interface{}
	if err := ctx.ShouldBindJSON(&data); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.paymentSvc.ProcessWebhook(ctx.Request.Context(), data)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// GetPayment godoc
// @Summary      Get a payment by ID
// @Description  Retrieve a single payment by its unique identifier
// @Tags         Payments
// @Accept       json
// @Produce      json
// @Param        id path string true "Payment ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/payments/{id} [get]
func (h *PaymentHandler) GetPayment(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid payment ID")
		return
	}

	resp := h.paymentSvc.GetPayment(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ListPayments godoc
// @Summary      List all payments
// @Description  Retrieve a paginated list of payments
// @Tags         Payments
// @Accept       json
// @Produce      json
// @Param        page query int false "Page number"
// @Param        page_size query int false "Page size"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/payments [get]
func (h *PaymentHandler) ListPayments(ctx *gin.Context) {
	pagination := utils.GetPaginationParams(ctx)
	resp := h.paymentSvc.ListPayments(ctx.Request.Context(), pagination.Page, pagination.PageSize)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// RetryPayment godoc
// @Summary      Retry a failed payment
// @Description  Retry a previously failed payment by its unique identifier
// @Tags         Payments
// @Accept       json
// @Produce      json
// @Param        id path string true "Payment ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/payments/{id}/retry [put]
func (h *PaymentHandler) RetryPayment(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid payment ID")
		return
	}

	resp := h.paymentSvc.RetryPayment(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ReconcilePayment godoc
// @Summary      Reconcile a payment
// @Description  Reconcile a payment by its unique identifier
// @Tags         Payments
// @Accept       json
// @Produce      json
// @Param        id path string true "Payment ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/payments/{id}/reconcile [put]
func (h *PaymentHandler) ReconcilePayment(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid payment ID")
		return
	}

	resp := h.paymentSvc.ReconcilePayment(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}
