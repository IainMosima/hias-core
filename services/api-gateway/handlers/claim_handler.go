package handlers

import (
	"io"
	"net/http"
	"time"

	claimsSchema "github.com/bitbiz/hias-core/domains/claims/schema"
	"github.com/bitbiz/hias-core/domains/claims/service"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ClaimHandler struct {
	claimSvc service.ClaimService
}

func NewClaimHandler(claimSvc service.ClaimService) *ClaimHandler {
	return &ClaimHandler{claimSvc: claimSvc}
}

// SubmitClaim godoc
// @Summary      Submit a new claim
// @Description  Submit a new insurance claim for processing
// @Tags         Claims
// @Accept       json
// @Produce      json
// @Param        request body claimsSchema.SubmitClaimRequest true "Claim submission details"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/claims [post]
func (h *ClaimHandler) SubmitClaim(ctx *gin.Context) {
	var req claimsSchema.SubmitClaimRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.claimSvc.SubmitClaim(ctx.Request.Context(), req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// GetClaim godoc
// @Summary      Get a claim by ID
// @Description  Retrieve a single claim by its unique identifier
// @Tags         Claims
// @Produce      json
// @Param        id path string true "Claim ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/claims/{id} [get]
func (h *ClaimHandler) GetClaim(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid claim ID")
		return
	}

	resp := h.claimSvc.GetClaim(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ListClaims godoc
// @Summary      List claims
// @Description  List claims with optional filtering by status, provider, date range, or search query
// @Tags         Claims
// @Produce      json
// @Param        status query string false "Filter by claim status"
// @Param        provider query string false "Filter by provider ID"
// @Param        date_from query string false "Filter from date (RFC3339)"
// @Param        date_to query string false "Filter to date (RFC3339)"
// @Param        search query string false "Search claim number or status"
// @Param        page query int false "Page number"
// @Param        page_size query int false "Page size"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/claims [get]
func (h *ClaimHandler) ListClaims(ctx *gin.Context) {
	pagination := utils.GetPaginationParams(ctx)

	status := ctx.Query("status")
	providerIDStr := ctx.Query("provider")
	dateFromStr := ctx.Query("date_from")
	dateToStr := ctx.Query("date_to")
	search := ctx.Query("search")

	// If date or search filters are present, use filtered query
	if dateFromStr != "" || dateToStr != "" || search != "" {
		var dateFrom, dateTo *time.Time
		if dateFromStr != "" {
			if t, err := time.Parse(time.RFC3339, dateFromStr); err == nil {
				dateFrom = &t
			} else if t, err := time.Parse("2006-01-02", dateFromStr); err == nil {
				dateFrom = &t
			}
		}
		if dateToStr != "" {
			if t, err := time.Parse(time.RFC3339, dateToStr); err == nil {
				dateTo = &t
			} else if t, err := time.Parse("2006-01-02", dateToStr); err == nil {
				endOfDay := t.Add(24*time.Hour - time.Second)
				dateTo = &endOfDay
			}
		}

		resp := h.claimSvc.ListClaimsFiltered(ctx.Request.Context(), status, dateFrom, dateTo, search, pagination.Page, pagination.PageSize)
		if resp.Error != nil {
			utils.RespondError(ctx, resp.StatusCode, resp.Message)
			return
		}
		countResp := h.claimSvc.CountClaimsFiltered(ctx.Request.Context(), status, dateFrom, dateTo, search)
		if countResp.Error != nil {
			utils.RespondError(ctx, countResp.StatusCode, countResp.Message)
			return
		}
		utils.RespondPaginated(ctx, resp.Message, resp.Data, pagination.Page, pagination.PageSize, countResp.Data)
		return
	}

	if status != "" {
		resp := h.claimSvc.ListClaimsByStatus(ctx.Request.Context(), status, pagination.Page, pagination.PageSize)
		if resp.Error != nil {
			utils.RespondError(ctx, resp.StatusCode, resp.Message)
			return
		}
		countResp := h.claimSvc.GetTotalCount(ctx.Request.Context())
		if countResp.Error != nil {
			utils.RespondError(ctx, countResp.StatusCode, countResp.Message)
			return
		}
		utils.RespondPaginated(ctx, resp.Message, resp.Data, pagination.Page, pagination.PageSize, countResp.Data)
		return
	}

	if providerIDStr != "" {
		providerID, err := uuid.Parse(providerIDStr)
		if err != nil {
			utils.RespondError(ctx, http.StatusBadRequest, "Invalid provider ID")
			return
		}
		resp := h.claimSvc.ListClaimsByProvider(ctx.Request.Context(), providerID, pagination.Page, pagination.PageSize)
		if resp.Error != nil {
			utils.RespondError(ctx, resp.StatusCode, resp.Message)
			return
		}
		countResp := h.claimSvc.GetTotalCount(ctx.Request.Context())
		if countResp.Error != nil {
			utils.RespondError(ctx, countResp.StatusCode, countResp.Message)
			return
		}
		utils.RespondPaginated(ctx, resp.Message, resp.Data, pagination.Page, pagination.PageSize, countResp.Data)
		return
	}

	resp := h.claimSvc.ListClaims(ctx.Request.Context(), pagination.Page, pagination.PageSize)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	countResp := h.claimSvc.GetTotalCount(ctx.Request.Context())
	if countResp.Error != nil {
		utils.RespondError(ctx, countResp.StatusCode, countResp.Message)
		return
	}
	utils.RespondPaginated(ctx, resp.Message, resp.Data, pagination.Page, pagination.PageSize, countResp.Data)
}

// GetTimeline godoc
// @Summary      Get claim timeline
// @Description  Returns the status change history (timeline) for a specific claim
// @Tags         Claims
// @Produce      json
// @Param        id path string true "Claim ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/claims/{id}/timeline [get]
func (h *ClaimHandler) GetTimeline(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid claim ID")
		return
	}

	resp := h.claimSvc.GetClaimTimeline(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ApproveClaim godoc
// @Summary      Approve a claim
// @Description  Approve an insurance claim by its ID
// @Tags         Claims
// @Produce      json
// @Param        id path string true "Claim ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/claims/{id}/approve [put]
func (h *ClaimHandler) ApproveClaim(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid claim ID")
		return
	}

	role, _ := ctx.Get("role")
	roleStr, _ := role.(string)

	resp := h.claimSvc.ApproveClaim(ctx.Request.Context(), id, getUserID(ctx), roleStr)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// RejectClaim godoc
// @Summary      Reject a claim
// @Description  Reject an insurance claim by its ID with a reason
// @Tags         Claims
// @Accept       json
// @Produce      json
// @Param        id path string true "Claim ID"
// @Param        request body claimsSchema.ReviewClaimRequest true "Rejection reason"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/claims/{id}/reject [put]
func (h *ClaimHandler) RejectClaim(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid claim ID")
		return
	}

	var req claimsSchema.ReviewClaimRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.claimSvc.RejectClaim(ctx.Request.Context(), id, req.Reason, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// VetClaim godoc
// @Summary      Vet a claim
// @Description  Vet an insurance claim by its ID with vetting details
// @Tags         Claims
// @Accept       json
// @Produce      json
// @Param        id path string true "Claim ID"
// @Param        request body claimsSchema.VetClaimRequest true "Vetting details"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/claims/{id}/vet [put]
func (h *ClaimHandler) VetClaim(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid claim ID")
		return
	}

	var req claimsSchema.VetClaimRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.claimSvc.VetClaim(ctx.Request.Context(), id, req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// MarkReadyForPayment godoc
// @Summary      Mark claim as ready for payment
// @Description  Mark an approved claim as ready for payment processing
// @Tags         Claims
// @Produce      json
// @Param        id path string true "Claim ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/claims/{id}/ready-for-payment [put]
func (h *ClaimHandler) MarkReadyForPayment(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid claim ID")
		return
	}

	resp := h.claimSvc.MarkReadyForPayment(ctx.Request.Context(), id, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// MarkPaid godoc
// @Summary      Mark claim as paid
// @Description  Mark a claim as fully paid
// @Tags         Claims
// @Produce      json
// @Param        id path string true "Claim ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/claims/{id}/mark-paid [put]
func (h *ClaimHandler) MarkPaid(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid claim ID")
		return
	}

	resp := h.claimSvc.MarkPaid(ctx.Request.Context(), id, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// MarkPartPaid godoc
// @Summary      Mark claim as partially paid
// @Description  Mark a claim as partially paid
// @Tags         Claims
// @Produce      json
// @Param        id path string true "Claim ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/claims/{id}/mark-part-paid [put]
func (h *ClaimHandler) MarkPartPaid(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid claim ID")
		return
	}

	resp := h.claimSvc.MarkPartPaid(ctx.Request.Context(), id, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// BulkSubmitClaims godoc
// @Summary      Bulk submit claims
// @Description  Submit multiple insurance claims in a single request
// @Tags         Claims
// @Accept       json
// @Produce      json
// @Param        request body claimsSchema.BulkSubmitClaimsRequest true "Bulk claim submissions"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/claims/bulk [post]
func (h *ClaimHandler) BulkSubmitClaims(ctx *gin.Context) {
	var req claimsSchema.BulkSubmitClaimsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.claimSvc.BulkSubmitClaims(ctx.Request.Context(), req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ListSLABreached godoc
// @Summary      List SLA-breached claims
// @Description  List all claims that have breached their SLA deadlines
// @Tags         Claims
// @Produce      json
// @Param        page query int false "Page number"
// @Param        page_size query int false "Page size"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/claims/sla-breached [get]
func (h *ClaimHandler) ListSLABreached(ctx *gin.Context) {
	pagination := utils.GetPaginationParams(ctx)
	resp := h.claimSvc.ListSLABreached(ctx.Request.Context(), pagination.Page, pagination.PageSize)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// UploadClaimDocument godoc
// @Summary      Upload a claim document
// @Description  Upload a document associated with a specific claim
// @Tags         Claims
// @Accept       json
// @Produce      json
// @Param        id path string true "Claim ID"
// @Param        request body object true "Document details (file_name, file_type, file_size, s3_key)"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/claims/{id}/documents [post]
func (h *ClaimHandler) UploadClaimDocument(ctx *gin.Context) {
	claimID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid claim ID")
		return
	}

	var req struct {
		FileName string `json:"file_name" binding:"required"`
		FileType string `json:"file_type" binding:"required"`
		FileSize int64  `json:"file_size"`
		S3Key    string `json:"s3_key" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.claimSvc.UploadClaimDocument(ctx.Request.Context(), claimID, req.FileName, req.FileType, req.FileSize, req.S3Key, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ListClaimDocuments godoc
// @Summary      List claim documents
// @Description  List all documents associated with a specific claim
// @Tags         Claims
// @Produce      json
// @Param        id path string true "Claim ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/claims/{id}/documents [get]
func (h *ClaimHandler) ListClaimDocuments(ctx *gin.Context) {
	claimID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid claim ID")
		return
	}

	resp := h.claimSvc.ListClaimDocuments(ctx.Request.Context(), claimID)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ImportClaimsCSV godoc
// @Summary      Import claims from CSV
// @Description  Import multiple claims by uploading a CSV file
// @Tags         Claims
// @Accept       multipart/form-data
// @Produce      json
// @Param        file formData file true "CSV file containing claims data"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/claims/import-csv [post]
func (h *ClaimHandler) ImportClaimsCSV(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "File is required")
		return
	}

	f, err := file.Open()
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, "Failed to open file")
		return
	}
	defer f.Close()

	csvData, err := io.ReadAll(f)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, "Failed to read file")
		return
	}

	resp := h.claimSvc.ImportClaimsCSV(ctx.Request.Context(), csvData, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// DeleteClaimDocument godoc
// @Summary      Delete a claim document
// @Description  Delete a document by its ID
// @Tags         Claims
// @Produce      json
// @Param        id path string true "Document ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/claim-documents/{id} [delete]
func (h *ClaimHandler) DeleteClaimDocument(ctx *gin.Context) {
	docID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid document ID")
		return
	}

	resp := h.claimSvc.DeleteClaimDocument(ctx.Request.Context(), docID)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}
