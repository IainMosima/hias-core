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

// CreateQuotation godoc
// @Summary      Create a new quotation
// @Description  Create a new insurance quotation with the provided details
// @Tags         Quotations
// @Accept       json
// @Produce      json
// @Param        request body salesSchema.CreateQuotationRequest true "Quotation creation request"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/quotations [post]
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

// GetQuotation godoc
// @Summary      Get a quotation by ID
// @Description  Retrieve a specific quotation by its unique identifier
// @Tags         Quotations
// @Accept       json
// @Produce      json
// @Param        id path string true "Quotation ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/quotations/{id} [get]
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

// ListQuotations godoc
// @Summary      List quotations
// @Description  Retrieve a paginated list of all quotations
// @Tags         Quotations
// @Accept       json
// @Produce      json
// @Param        page query int false "Page number"
// @Param        page_size query int false "Page size"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/quotations [get]
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

// ListQuotationsByLead godoc
// @Summary      List quotations by lead
// @Description  Retrieve a paginated list of quotations associated with a specific lead
// @Tags         Quotations
// @Accept       json
// @Produce      json
// @Param        id path string true "Lead ID"
// @Param        page query int false "Page number"
// @Param        page_size query int false "Page size"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/leads/{id}/quotations [get]
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

// CreateVersion godoc
// @Summary      Create a new quotation version
// @Description  Create a new version for an existing quotation
// @Tags         Quotations
// @Accept       json
// @Produce      json
// @Param        id path string true "Quotation ID"
// @Param        request body salesSchema.CreateQuotationVersionRequest true "Version creation request"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/quotations/{id}/versions [post]
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

// GetVersion godoc
// @Summary      Get a specific quotation version
// @Description  Retrieve a specific version of a quotation by quotation ID and version number
// @Tags         Quotations
// @Accept       json
// @Produce      json
// @Param        id path string true "Quotation ID"
// @Param        version path int true "Version number"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/quotations/{id}/versions/{version} [get]
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

// ListVersions godoc
// @Summary      List quotation versions
// @Description  Retrieve all versions of a specific quotation
// @Tags         Quotations
// @Accept       json
// @Produce      json
// @Param        id path string true "Quotation ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/quotations/{id}/versions [get]
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

// CompareVersions godoc
// @Summary      Compare quotation versions
// @Description  Compare two versions of a quotation side by side
// @Tags         Quotations
// @Accept       json
// @Produce      json
// @Param        id path string true "Quotation ID"
// @Param        version_a query int true "First version number"
// @Param        version_b query int true "Second version number"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/quotations/{id}/versions/compare [get]
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

// SubmitForApproval godoc
// @Summary      Submit quotation version for approval
// @Description  Submit a specific version of a quotation for approval review
// @Tags         Quotations
// @Accept       json
// @Produce      json
// @Param        id path string true "Quotation ID"
// @Param        version path int true "Version number"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/quotations/{id}/versions/{version}/submit-approval [put]
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

// ApproveVersion godoc
// @Summary      Approve a quotation version
// @Description  Approve a specific version of a quotation
// @Tags         Quotations
// @Accept       json
// @Produce      json
// @Param        id path string true "Quotation ID"
// @Param        version path int true "Version number"
// @Param        request body salesSchema.ApproveVersionRequest false "Approval details"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/quotations/{id}/versions/{version}/approve [put]
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

// RejectVersion godoc
// @Summary      Reject a quotation version
// @Description  Reject a specific version of a quotation with a reason
// @Tags         Quotations
// @Accept       json
// @Produce      json
// @Param        id path string true "Quotation ID"
// @Param        version path int true "Version number"
// @Param        request body salesSchema.RejectVersionRequest true "Rejection details"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/quotations/{id}/versions/{version}/reject [put]
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

// IssueQuotation godoc
// @Summary      Issue a quotation
// @Description  Issue an approved quotation, making it available to send to the client
// @Tags         Quotations
// @Accept       json
// @Produce      json
// @Param        id path string true "Quotation ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/quotations/{id}/issue [put]
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

// AcceptQuotation godoc
// @Summary      Accept a quotation
// @Description  Accept an issued quotation on behalf of the client
// @Tags         Quotations
// @Accept       json
// @Produce      json
// @Param        id path string true "Quotation ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/quotations/{id}/accept [put]
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

// DeclineQuotation godoc
// @Summary      Decline a quotation
// @Description  Decline an issued quotation on behalf of the client
// @Tags         Quotations
// @Accept       json
// @Produce      json
// @Param        id path string true "Quotation ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/quotations/{id}/decline [put]
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

// SendToClient godoc
// @Summary      Send quotation to client
// @Description  Send an issued quotation to the client via the specified delivery method
// @Tags         Quotations
// @Accept       json
// @Produce      json
// @Param        id path string true "Quotation ID"
// @Param        request body salesSchema.SendQuotationRequest true "Send request details"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/quotations/{id}/send [put]
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

// ConvertToPolicy godoc
// @Summary      Convert quotation to policy
// @Description  Convert an accepted quotation into a new insurance policy
// @Tags         Quotations
// @Accept       json
// @Produce      json
// @Param        id path string true "Quotation ID"
// @Param        request body salesSchema.ConvertToPolicyRequest true "Conversion request details"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/quotations/{id}/convert [post]
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

// UploadDocument godoc
// @Summary      Upload a quotation document
// @Description  Upload a document and attach it to a specific quotation
// @Tags         Quotations
// @Accept       json
// @Produce      json
// @Param        id path string true "Quotation ID"
// @Param        request body salesSchema.UploadDocumentMeta true "Document metadata"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/quotations/{id}/documents [post]
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

// ListDocuments godoc
// @Summary      List quotation documents
// @Description  Retrieve all documents attached to a specific quotation
// @Tags         Quotations
// @Accept       json
// @Produce      json
// @Param        id path string true "Quotation ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/quotations/{id}/documents [get]
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

// UpdateDocument godoc
// @Summary      Update a quotation document
// @Description  Update metadata of a quotation document by its ID
// @Tags         Quotations
// @Accept       json
// @Produce      json
// @Param        id path string true "Document ID"
// @Param        request body salesSchema.UpdateDocumentMeta true "Updated document metadata"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/quotation-documents/{id} [put]
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

// DeleteDocument godoc
// @Summary      Delete a quotation document
// @Description  Delete a quotation document by its ID
// @Tags         Quotations
// @Accept       json
// @Produce      json
// @Param        id path string true "Document ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/quotation-documents/{id} [delete]
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

// ExpireQuotations godoc
// @Summary      Expire quotations
// @Description  Trigger expiration of all quotations that have passed their validity date
// @Tags         Quotations
// @Accept       json
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/quotations/expire [post]
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
