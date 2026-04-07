package handlers

import (
	"bytes"
	"io"
	"net/http"

	policySchema "github.com/bitbiz/hias-core/domains/policy/schema"
	"github.com/bitbiz/hias-core/domains/policy/service"
	"github.com/bitbiz/hias-core/shared"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MemberHandler struct {
	memberSvc service.MemberService
}

func NewMemberHandler(memberSvc service.MemberService) *MemberHandler {
	return &MemberHandler{memberSvc: memberSvc}
}

// CreateMember godoc
// @Summary      Create a new member (standalone)
// @Description  Creates a new member by specifying the policy ID in the request body
// @Tags         Members
// @Accept       json
// @Produce      json
// @Param        request body policySchema.CreateMemberRequest true "Create member request"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/members [post]
func (h *MemberHandler) CreateMember(ctx *gin.Context) {
	var req policySchema.CreateMemberRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	policyID, err := uuid.Parse(req.PolicyID)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid policy ID")
		return
	}

	enrollReq := policySchema.EnrollMemberRequest{
		NationalID:   req.NationalID,
		Name:         req.Name,
		DateOfBirth:  req.DateOfBirth,
		Gender:       req.Gender,
		Relationship: req.Relationship,
		Phone:        req.Phone,
		Email:        req.Email,
		KRAPin:       req.KRAPin,
		County:       req.County,
		City:         req.City,
		Country:      req.Country,
		Address:      req.Address,
	}

	resp := h.memberSvc.EnrollMember(ctx.Request.Context(), policyID, enrollReq)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// EnrollMember godoc
// @Summary      Enroll a new member in a policy
// @Description  Enrolls a new member under the specified policy
// @Tags         Members
// @Accept       json
// @Produce      json
// @Param        id path string true "Policy ID"
// @Param        request body policySchema.EnrollMemberRequest true "Enroll member request"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/policies/{id}/members [post]
func (h *MemberHandler) EnrollMember(ctx *gin.Context) {
	policyID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid policy ID")
		return
	}

	var req policySchema.EnrollMemberRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.memberSvc.EnrollMember(ctx.Request.Context(), policyID, req)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// VerifyMember godoc
// @Summary      Verify a member
// @Description  Marks the specified member as verified
// @Tags         Members
// @Accept       json
// @Produce      json
// @Param        id path string true "Member ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/members/{id}/verify [put]
func (h *MemberHandler) VerifyMember(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid member ID")
		return
	}

	resp := h.memberSvc.VerifyMember(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// GetMemberEligibility godoc
// @Summary      Get member eligibility
// @Description  Retrieves the eligibility status and details for a specific member
// @Tags         Members
// @Accept       json
// @Produce      json
// @Param        id path string true "Member ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/members/{id}/eligibility [get]
func (h *MemberHandler) GetMemberEligibility(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid member ID")
		return
	}

	resp := h.memberSvc.GetMemberEligibility(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ListMembers godoc
// @Summary      List members of a policy
// @Description  Returns all members enrolled under the specified policy
// @Tags         Members
// @Accept       json
// @Produce      json
// @Param        id path string true "Policy ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/policies/{id}/members [get]
func (h *MemberHandler) ListMembers(ctx *gin.Context) {
	policyID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid policy ID")
		return
	}

	resp := h.memberSvc.ListMembers(ctx.Request.Context(), policyID)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}

	role, _ := ctx.Get("role")
	if roleStr, ok := role.(string); ok && roleStr != string(shared.UserRoleAdmin) {
		for i := range resp.Data {
			maskMemberPII(&resp.Data[i])
		}
	}

	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// GetMember godoc
// @Summary      Get a member by ID
// @Description  Retrieves the details of a specific member
// @Tags         Members
// @Accept       json
// @Produce      json
// @Param        id path string true "Member ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/members/{id} [get]
func (h *MemberHandler) GetMember(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid member ID")
		return
	}

	resp := h.memberSvc.GetMember(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}

	role, _ := ctx.Get("role")
	if roleStr, ok := role.(string); ok && roleStr != string(shared.UserRoleAdmin) {
		maskMemberPII(&resp.Data)
	}

	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// UpdateMember godoc
// @Summary      Update a member
// @Description  Updates the details of an existing member
// @Tags         Members
// @Accept       json
// @Produce      json
// @Param        id path string true "Member ID"
// @Param        request body policySchema.UpdateMemberRequest true "Update member request"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/members/{id} [put]
func (h *MemberHandler) UpdateMember(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid member ID")
		return
	}

	var req policySchema.UpdateMemberRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.memberSvc.UpdateMember(ctx.Request.Context(), id, req)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// RemoveMember godoc
// @Summary      Remove a member
// @Description  Removes a member from their policy
// @Tags         Members
// @Accept       json
// @Produce      json
// @Param        id path string true "Member ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/members/{id} [delete]
func (h *MemberHandler) RemoveMember(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid member ID")
		return
	}

	var req policySchema.RemoveMemberRequest
	_ = ctx.ShouldBindJSON(&req) // optional body

	resp := h.memberSvc.RemoveMember(ctx.Request.Context(), id, req.Reason)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// SuspendMember godoc
// @Summary      Suspend a member
// @Description  Suspends the specified member
// @Tags         Members
// @Accept       json
// @Produce      json
// @Param        id path string true "Member ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/members/{id}/suspend [put]
func (h *MemberHandler) SuspendMember(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid member ID")
		return
	}

	resp := h.memberSvc.SuspendMember(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ReactivateMember godoc
// @Summary      Reactivate a member
// @Description  Reactivates a previously suspended member
// @Tags         Members
// @Accept       json
// @Produce      json
// @Param        id path string true "Member ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/members/{id}/reactivate [put]
func (h *MemberHandler) ReactivateMember(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid member ID")
		return
	}

	resp := h.memberSvc.ReactivateMember(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// BulkEnrollMembers godoc
// @Summary      Bulk enroll members
// @Description  Enrolls multiple members under the specified policy in a single request
// @Tags         Members
// @Accept       json
// @Produce      json
// @Param        id path string true "Policy ID"
// @Param        request body policySchema.BulkEnrollRequest true "Bulk enroll request"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/policies/{id}/members/bulk [post]
func (h *MemberHandler) BulkEnrollMembers(ctx *gin.Context) {
	policyID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid policy ID")
		return
	}

	var req policySchema.BulkEnrollRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.memberSvc.BulkEnrollMembers(ctx.Request.Context(), policyID, req.Members)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// BulkRemoveMembers godoc
// @Summary      Bulk remove members
// @Description  Removes multiple members from the specified policy in a single request
// @Tags         Members
// @Accept       json
// @Produce      json
// @Param        id path string true "Policy ID"
// @Param        request body policySchema.BulkRemoveRequest true "Bulk remove request"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/policies/{id}/members/bulk-remove [post]
func (h *MemberHandler) BulkRemoveMembers(ctx *gin.Context) {
	policyID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid policy ID")
		return
	}

	var req policySchema.BulkRemoveRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	memberIDs := make([]uuid.UUID, 0, len(req.MemberIDs))
	for _, idStr := range req.MemberIDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			utils.RespondError(ctx, http.StatusBadRequest, "Invalid member ID: "+idStr)
			return
		}
		memberIDs = append(memberIDs, id)
	}

	resp := h.memberSvc.BulkRemoveMembers(ctx.Request.Context(), policyID, memberIDs, req.Reason)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ImportMembersCSV godoc
// @Summary      Import members from CSV
// @Description  Imports members into the specified policy from an uploaded CSV file
// @Tags         Members
// @Accept       multipart/form-data
// @Produce      json
// @Param        id path string true "Policy ID"
// @Param        file formData file true "CSV file with member data"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/policies/{id}/members/import [post]
func (h *MemberHandler) ImportMembersCSV(ctx *gin.Context) {
	policyID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid policy ID")
		return
	}

	file, _, err := ctx.Request.FormFile("file")
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "CSV file is required")
		return
	}
	defer file.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, file); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Failed to read CSV file")
		return
	}

	resp := h.memberSvc.ImportMembersCSV(ctx.Request.Context(), policyID, buf.Bytes())
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ListAllMembers godoc
// @Summary      List all members across all policies
// @Description  Returns a paginated list of all members with optional search
// @Tags         Members
// @Produce      json
// @Param        search query string false "Search by name, national ID, email, or phone"
// @Param        page query int false "Page number"
// @Param        page_size query int false "Page size"
// @Success      200 {object} map[string]interface{}
// @Failure      500 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/members [get]
func (h *MemberHandler) ListAllMembers(ctx *gin.Context) {
	pagination := utils.GetPaginationParams(ctx)
	search := ctx.Query("search")

	resp := h.memberSvc.ListMembersFiltered(ctx.Request.Context(), search, pagination.Page, pagination.PageSize)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}

	role, _ := ctx.Get("role")
	if roleStr, ok := role.(string); ok && roleStr != string(shared.UserRoleAdmin) {
		for i := range resp.Data {
			maskMemberPII(&resp.Data[i])
		}
	}

	countResp := h.memberSvc.CountMembersFiltered(ctx.Request.Context(), search)
	utils.RespondPaginated(ctx, resp.Message, resp.Data, pagination.Page, pagination.PageSize, countResp.Data)
}

func maskMemberPII(m *policySchema.MemberResponse) {
	m.NationalID = utils.MaskNationalID(m.NationalID)
	m.Email = utils.MaskEmail(m.Email)
	m.Phone = utils.MaskPhone(m.Phone)
}
