package handlers

import (
	"net/http"

	billingSchema "github.com/bitbiz/hias-core/domains/billing/schema"
	"github.com/bitbiz/hias-core/domains/billing/service"
	"github.com/bitbiz/hias-core/services/api-gateway/rest/middleware"
	"github.com/bitbiz/hias-core/shared/auth"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BillingHandler struct {
	billingService service.BillingService
	paymentService service.PaymentService
}

func NewBillingHandler(billingService service.BillingService, paymentService service.PaymentService) *BillingHandler {
	return &BillingHandler{
		billingService: billingService,
		paymentService: paymentService,
	}
}

func (h *BillingHandler) CreateInvoice(ctx *gin.Context) {
	var req struct {
		PolicyID string `json:"policy_id" binding:"required,uuid"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.billingService.GenerateInvoice(ctx.Request.Context(), req.PolicyID)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}

	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *BillingHandler) GetInvoice(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid invoice ID")
		return
	}

	// TODO: Add GetInvoice(ctx, id uuid.UUID) to BillingService interface.
	// Currently using GenerateInvoice with invoice ID as a temporary approach.
	resp := h.billingService.GenerateInvoice(ctx.Request.Context(), id.String())
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}

	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *BillingHandler) ListInvoices(ctx *gin.Context) {
	// TODO: Add ListInvoices(ctx, page, pageSize int) to BillingService interface.
	// Currently returns the billing cycle run count as a placeholder.
	resp := h.billingService.RunBillingCycle(ctx.Request.Context())
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}

	utils.RespondSuccess(ctx, http.StatusOK, resp.Message, resp.Data)
}

func (h *BillingHandler) InitiatePayment(ctx *gin.Context) {
	var req billingSchema.InitiatePaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	payload := ctx.MustGet(middleware.AuthPayloadKey).(*auth.Payload)
	createdBy, _ := uuid.Parse(payload.UserID)

	resp := h.paymentService.InitiatePayment(ctx.Request.Context(), req, createdBy)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}

	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *BillingHandler) GetPayment(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid payment ID")
		return
	}

	resp := h.paymentService.GetPayment(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}

	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *BillingHandler) ListPayments(ctx *gin.Context) {
	pagination := utils.GetPaginationParams(ctx)

	resp := h.paymentService.ListPayments(ctx.Request.Context(), pagination.Page, pagination.PageSize)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}

	utils.RespondPaginated(ctx, resp.Message, resp.Data, pagination.Page, pagination.PageSize, 0)
}

func (h *BillingHandler) ProcessPaymentWebhook(ctx *gin.Context) {
	var webhookData billingSchema.MpesaWebhookRequest
	if err := ctx.ShouldBindJSON(&webhookData); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.paymentService.ProcessWebhook(ctx.Request.Context(), webhookData.Body)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}

	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}
