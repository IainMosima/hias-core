package handlers

import (
	"net/http"

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

func (h *ClaimHandler) ListClaims(ctx *gin.Context) {
	pagination := utils.GetPaginationParams(ctx)

	status := ctx.Query("status")
	providerIDStr := ctx.Query("provider")

	if status != "" {
		resp := h.claimSvc.ListClaimsByStatus(ctx.Request.Context(), status, pagination.Page, pagination.PageSize)
		if resp.Error != nil {
			utils.RespondError(ctx, resp.StatusCode, resp.Message)
			return
		}
		countResp := h.claimSvc.GetTotalCount(ctx.Request.Context())
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
		utils.RespondPaginated(ctx, resp.Message, resp.Data, pagination.Page, pagination.PageSize, countResp.Data)
		return
	}

	resp := h.claimSvc.ListClaims(ctx.Request.Context(), pagination.Page, pagination.PageSize)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	countResp := h.claimSvc.GetTotalCount(ctx.Request.Context())
	utils.RespondPaginated(ctx, resp.Message, resp.Data, pagination.Page, pagination.PageSize, countResp.Data)
}

func (h *ClaimHandler) ApproveClaim(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid claim ID")
		return
	}

	resp := h.claimSvc.ApproveClaim(ctx.Request.Context(), id, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

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

func (h *ClaimHandler) ListSLABreached(ctx *gin.Context) {
	pagination := utils.GetPaginationParams(ctx)
	resp := h.claimSvc.ListSLABreached(ctx.Request.Context(), pagination.Page, pagination.PageSize)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

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
