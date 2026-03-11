package handlers

import (
	"log"
	"net/http"
	"strconv"

	docSchema "github.com/bitbiz/hias-core/domains/document/schema"
	docService "github.com/bitbiz/hias-core/domains/document/service"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	awsSvc "github.com/bitbiz/hias-core/shared/aws"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DocumentHandler struct {
	store  db.Store
	s3Svc  awsSvc.S3Service
	docSvc docService.DocumentService
}

func NewDocumentHandler(store db.Store, s3Svc awsSvc.S3Service, docSvc docService.DocumentService) *DocumentHandler {
	return &DocumentHandler{store: store, s3Svc: s3Svc, docSvc: docSvc}
}

// ListStandaloneDocuments godoc
// @Summary      List all documents across the system
// @Description  Returns a paginated UNION of policy, claim, and quotation documents
// @Tags         Documents
// @Produce      json
// @Param        limit query int false "Limit" default(10)
// @Param        offset query int false "Offset" default(0)
// @Success      200 {object} map[string]interface{}
// @Failure      500 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/documents/standalone [get]
func (h *DocumentHandler) ListStandaloneDocuments(ctx *gin.Context) {
	limit := 10
	offset := 0
	if v := ctx.Query("limit"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	if v := ctx.Query("offset"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	docs, err := h.store.ListStandaloneDocuments(ctx.Request.Context(), db.ListStandaloneDocumentsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		log.Printf("Failed to list standalone documents: %v", err)
		utils.RespondError(ctx, http.StatusInternalServerError, "Failed to list documents")
		return
	}

	total, err := h.store.CountStandaloneDocuments(ctx.Request.Context())
	if err != nil {
		log.Printf("Failed to count standalone documents: %v", err)
		utils.RespondError(ctx, http.StatusInternalServerError, "Failed to count documents")
		return
	}

	responses := make([]docSchema.StandaloneDocumentResponse, len(docs))
	for i, d := range docs {
		responses[i] = docSchema.StandaloneDocumentResponse{
			ID:           d.ID,
			SourceType:   d.SourceType,
			SourceID:     d.SourceID,
			DocumentType: d.DocumentType,
			FileName:     d.FileName,
			FileSize:     d.FileSize,
			S3Key:        d.S3Key,
			CreatedBy:    d.CreatedBy,
			CreatedAt:    d.CreatedAt,
		}
	}

	result := docSchema.StandaloneDocumentListResponse{
		Documents: responses,
		Total:     total,
		Limit:     limit,
		Offset:    offset,
	}

	utils.RespondSuccess(ctx, http.StatusOK, "Documents retrieved", result)
}

// DownloadDocument godoc
// @Summary      Get a presigned download URL for a document
// @Description  Looks up the document across all tables and returns a presigned S3 URL
// @Tags         Documents
// @Produce      json
// @Param        id path string true "Document ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/documents/{id}/download [get]
func (h *DocumentHandler) DownloadDocument(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid document ID")
		return
	}

	doc, err := h.store.FindDocumentS3Key(ctx.Request.Context(), id)
	if err != nil {
		utils.RespondError(ctx, http.StatusNotFound, "Document not found")
		return
	}

	if h.s3Svc == nil {
		utils.RespondError(ctx, http.StatusInternalServerError, "S3 service not configured")
		return
	}

	url, err := h.s3Svc.GetPresignedURL(ctx.Request.Context(), doc.S3Key, 900)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, "Failed to generate download URL")
		return
	}

	utils.RespondSuccess(ctx, http.StatusOK, "Download URL generated", docSchema.DownloadURLResponse{
		DownloadURL: url,
	})
}

// RequestUploadURL godoc
// @Summary      Request a presigned upload URL
// @Description  Creates a document record and returns a presigned S3 PUT URL
// @Tags         Documents
// @Accept       json
// @Produce      json
// @Param        request body docSchema.UploadURLRequest true "Upload URL request"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/documents/upload-url [post]
func (h *DocumentHandler) RequestUploadURL(ctx *gin.Context) {
	var req docSchema.UploadURLRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.docSvc.RequestUploadURL(ctx.Request.Context(), req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// BulkRequestUploadURLs godoc
// @Summary      Request multiple presigned upload URLs
// @Description  Creates multiple document records and returns presigned S3 PUT URLs
// @Tags         Documents
// @Accept       json
// @Produce      json
// @Param        request body docSchema.BulkUploadURLRequest true "Bulk upload URL request"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/documents/bulk-upload-urls [post]
func (h *DocumentHandler) BulkRequestUploadURLs(ctx *gin.Context) {
	var req docSchema.BulkUploadURLRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.docSvc.BulkRequestUploadURLs(ctx.Request.Context(), req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ConfirmUpload godoc
// @Summary      Confirm a document upload
// @Description  Verifies the file exists in S3 and marks the document as ACTIVE
// @Tags         Documents
// @Produce      json
// @Param        id path string true "Document ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/documents/{id}/confirm-upload [post]
func (h *DocumentHandler) ConfirmUpload(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid document ID")
		return
	}

	resp := h.docSvc.ConfirmUpload(ctx.Request.Context(), id, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// GetDocumentDownloadURL godoc
// @Summary      Get download URL for an uploaded document
// @Description  Returns a presigned GET URL for an ACTIVE uploaded document
// @Tags         Documents
// @Produce      json
// @Param        id path string true "Document ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/documents/{id}/download-url [post]
func (h *DocumentHandler) GetDocumentDownloadURL(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid document ID")
		return
	}

	resp := h.docSvc.GetDownloadURL(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// DeleteDocument godoc
// @Summary      Delete an uploaded document
// @Description  Soft deletes the document record and removes the S3 object
// @Tags         Documents
// @Produce      json
// @Param        id path string true "Document ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/documents/{id} [delete]
func (h *DocumentHandler) DeleteUploadedDocument(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid document ID")
		return
	}

	resp := h.docSvc.DeleteDocument(ctx.Request.Context(), id, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ListPolicyUploads godoc
// @Summary      List uploaded documents for a policy
// @Description  Returns all uploaded documents attached to a policy
// @Tags         Documents
// @Produce      json
// @Param        id path string true "Policy ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/policies/{id}/uploads [get]
func (h *DocumentHandler) ListPolicyUploads(ctx *gin.Context) {
	policyID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid policy ID")
		return
	}

	resp := h.docSvc.ListByEntity(ctx.Request.Context(), "policy", policyID)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ListMemberDocuments godoc
// @Summary      List uploaded documents for a member
// @Description  Returns all uploaded documents attached to a member
// @Tags         Documents
// @Produce      json
// @Param        id path string true "Member ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/members/{id}/documents [get]
func (h *DocumentHandler) ListMemberDocuments(ctx *gin.Context) {
	memberID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid member ID")
		return
	}

	resp := h.docSvc.ListByEntity(ctx.Request.Context(), "member", memberID)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}
