package handlers

import (
	"net/http"

	policySchema "github.com/bitbiz/hias-core/domains/policy/schema"
	"github.com/bitbiz/hias-core/domains/policy/service"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MemberHandler struct {
	memberService service.MemberService
}

func NewMemberHandler(memberService service.MemberService) *MemberHandler {
	return &MemberHandler{memberService: memberService}
}

func (h *MemberHandler) EnrollMember(ctx *gin.Context) {
	policyID, err := uuid.Parse(ctx.Param("policyId"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid policy ID")
		return
	}

	var req policySchema.EnrollMemberRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.memberService.EnrollMember(ctx.Request.Context(), policyID, req)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}

	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *MemberHandler) GetMember(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid member ID")
		return
	}

	resp := h.memberService.GetMemberEligibility(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}

	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *MemberHandler) ListMembers(ctx *gin.Context) {
	policyID, err := uuid.Parse(ctx.Param("policyId"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid policy ID")
		return
	}

	resp := h.memberService.ListMembers(ctx.Request.Context(), policyID)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}

	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *MemberHandler) VerifyMember(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid member ID")
		return
	}

	resp := h.memberService.VerifyMember(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}

	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}
