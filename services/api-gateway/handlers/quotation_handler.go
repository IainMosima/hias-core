package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	salesSchema "github.com/bitbiz/hias-core/domains/sales/schema"
	"github.com/bitbiz/hias-core/domains/sales/service"
	"github.com/bitbiz/hias-core/services/api-gateway/middleware"
	"github.com/bitbiz/hias-core/shared/auth"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type QuotationHandler struct {
	quotationSvc service.QuotationService
}

func NewQuotationHandler(quotationSvc service.QuotationService) *QuotationHandler {
	return &QuotationHandler{quotationSvc: quotationSvc}
}

func (h *QuotationHandler) CreateQuotation(ctx *gin.Context) {
	var req salesSchema.CreateQuotationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.quotationSvc.CreateQuotation(ctx.Request.Context(), req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *QuotationHandler) GetQuotation(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid quotation ID")
		return
	}

	resp := h.quotationSvc.GetQuotation(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *QuotationHandler) ListQuotations(ctx *gin.Context) {
	pagination := utils.GetPaginationParams(ctx)

	resp := h.quotationSvc.ListQuotations(ctx.Request.Context(), pagination.Page, pagination.PageSize)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	countResp := h.quotationSvc.GetTotalCount(ctx.Request.Context())
	utils.RespondPaginated(ctx, resp.Message, resp.Data, pagination.Page, pagination.PageSize, countResp.Data)
}

func (h *QuotationHandler) ListQuotationsByLead(ctx *gin.Context) {
	leadID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid lead ID")
		return
	}

	pagination := utils.GetPaginationParams(ctx)

	resp := h.quotationSvc.ListQuotationsByLead(ctx.Request.Context(), leadID, pagination.Page, pagination.PageSize)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *QuotationHandler) CreateVersion(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid quotation ID")
		return
	}

	var req salesSchema.CreateQuotationVersionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.quotationSvc.CreateVersion(ctx.Request.Context(), id, req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *QuotationHandler) GetVersion(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid quotation ID")
		return
	}

	versionNumber, err := strconv.Atoi(ctx.Param("version"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid version number")
		return
	}

	resp := h.quotationSvc.GetVersion(ctx.Request.Context(), id, versionNumber)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *QuotationHandler) ListVersions(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid quotation ID")
		return
	}

	resp := h.quotationSvc.ListVersions(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *QuotationHandler) CompareVersions(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid quotation ID")
		return
	}

	versionA, err := strconv.Atoi(ctx.Query("version_a"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid version_a parameter")
		return
	}
	versionB, err := strconv.Atoi(ctx.Query("version_b"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid version_b parameter")
		return
	}

	resp := h.quotationSvc.CompareVersions(ctx.Request.Context(), id, versionA, versionB)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *QuotationHandler) SubmitForApproval(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid quotation ID")
		return
	}

	versionNumber, err := strconv.Atoi(ctx.Param("version"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid version number")
		return
	}

	resp := h.quotationSvc.SubmitForApproval(ctx.Request.Context(), id, versionNumber, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *QuotationHandler) ApproveVersion(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid quotation ID")
		return
	}

	versionNumber, err := strconv.Atoi(ctx.Param("version"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid version number")
		return
	}

	var req salesSchema.ApproveVersionRequest
	_ = ctx.ShouldBindJSON(&req)

	resp := h.quotationSvc.ApproveVersion(ctx.Request.Context(), id, versionNumber, req, getUserID(ctx), getUserRole(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *QuotationHandler) RejectVersion(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid quotation ID")
		return
	}

	versionNumber, err := strconv.Atoi(ctx.Param("version"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid version number")
		return
	}

	var req salesSchema.RejectVersionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.quotationSvc.RejectVersion(ctx.Request.Context(), id, versionNumber, req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *QuotationHandler) IssueQuotation(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid quotation ID")
		return
	}

	resp := h.quotationSvc.IssueQuotation(ctx.Request.Context(), id, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *QuotationHandler) AcceptQuotation(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid quotation ID")
		return
	}

	resp := h.quotationSvc.AcceptQuotation(ctx.Request.Context(), id, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *QuotationHandler) DeclineQuotation(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid quotation ID")
		return
	}

	resp := h.quotationSvc.DeclineQuotation(ctx.Request.Context(), id, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *QuotationHandler) SendToClient(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid quotation ID")
		return
	}

	var req salesSchema.SendQuotationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.quotationSvc.SendToClient(ctx.Request.Context(), id, req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *QuotationHandler) ConvertToPolicy(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid quotation ID")
		return
	}

	var req salesSchema.ConvertToPolicyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.quotationSvc.ConvertToPolicy(ctx.Request.Context(), id, req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *QuotationHandler) UploadDocument(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid quotation ID")
		return
	}

	var meta salesSchema.UploadDocumentMeta
	if err := ctx.ShouldBindJSON(&meta); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// Generate S3 key
	s3Key := fmt.Sprintf("quotations/%s/documents/%s_%s", id.String(), uuid.New().String(), meta.FileName)

	resp := h.quotationSvc.UploadDocument(ctx.Request.Context(), id, meta, s3Key, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *QuotationHandler) ListDocuments(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid quotation ID")
		return
	}

	resp := h.quotationSvc.ListDocuments(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *QuotationHandler) UpdateDocument(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid document ID")
		return
	}

	var req salesSchema.UpdateDocumentMeta
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userRole := getUserRole(ctx)

	resp := h.quotationSvc.UpdateDocument(ctx.Request.Context(), id, req, getUserID(ctx), userRole)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *QuotationHandler) DeleteDocument(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid document ID")
		return
	}

	userRole := getUserRole(ctx)

	resp := h.quotationSvc.DeleteDocument(ctx.Request.Context(), id, getUserID(ctx), userRole)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *QuotationHandler) ExpireQuotations(ctx *gin.Context) {
	resp := h.quotationSvc.ExpireQuotations(ctx.Request.Context())
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func getUserRole(ctx *gin.Context) string {
	payload, exists := ctx.Get(middleware.AuthPayloadKey)
	if !exists {
		return ""
	}
	authPayload, ok := payload.(*auth.Payload)
	if !ok {
		return ""
	}
	return authPayload.Role
}
