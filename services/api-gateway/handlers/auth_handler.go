package handlers

import (
	"net/http"

	"github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/bitbiz/hias-core/domains/identity/service"
	"github.com/bitbiz/hias-core/services/api-gateway/middleware"
	"github.com/bitbiz/hias-core/shared/auth"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authSvc service.AuthService
	userSvc service.UserService
}

func NewAuthHandler(authSvc service.AuthService, userSvc service.UserService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc, userSvc: userSvc}
}

// Login godoc
// @Summary      User login
// @Description  Authenticate a user and return access tokens
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body schema.LoginRequest true "Login credentials"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Router       /api/v1/auth/login [post]
func (h *AuthHandler) Login(ctx *gin.Context) {
	var req schema.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.authSvc.Login(ctx.Request.Context(), req)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}

	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// Register godoc
// @Summary      Register a new user
// @Description  Create a new user account
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body schema.RegisterRequest true "Registration details"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Router       /api/v1/auth/register [post]
func (h *AuthHandler) Register(ctx *gin.Context) {
	var req schema.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.authSvc.Register(ctx.Request.Context(), req)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}

	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// RefreshToken godoc
// @Summary      Refresh access token
// @Description  Refresh an expired access token using a refresh token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body schema.RefreshTokenRequest true "Refresh token"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Router       /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(ctx *gin.Context) {
	var req schema.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.authSvc.RefreshToken(ctx.Request.Context(), req)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}

	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// Logout godoc
// @Summary      User logout
// @Description  Invalidate the current user session
// @Tags         Auth
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Security     BearerAuth
// @Router       /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(ctx *gin.Context) {
	payload := ctx.MustGet(middleware.AuthPayloadKey).(*auth.Payload)

	resp := h.authSvc.Logout(ctx.Request.Context(), payload.UserID)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}

	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// GetProfile godoc
// @Summary      Get current user profile
// @Description  Returns the profile of the currently authenticated user
// @Tags         Profile
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Failure      401 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/profile [get]
func (h *AuthHandler) GetProfile(ctx *gin.Context) {
	userID := getUserID(ctx)

	resp := h.userSvc.GetUserByID(ctx.Request.Context(), userID)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// UpdateProfile godoc
// @Summary      Update current user profile
// @Description  Updates the name and/or phone of the currently authenticated user
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Param        request body schema.UpdateProfileRequest true "Profile update details"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/profile [put]
func (h *AuthHandler) UpdateProfile(ctx *gin.Context) {
	userID := getUserID(ctx)

	var req schema.UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	updateReq := schema.UpdateUserRequest{
		Name:  req.Name,
		Phone: req.Phone,
	}

	resp := h.userSvc.UpdateUser(ctx.Request.Context(), userID, updateReq)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ChangePassword godoc
// @Summary      Change user password
// @Description  Changes the password of the currently authenticated user
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body schema.ChangePasswordRequest true "Password change details"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/auth/change-password [put]
func (h *AuthHandler) ChangePassword(ctx *gin.Context) {
	userID := getUserID(ctx)

	var req schema.ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.authSvc.ChangePassword(ctx.Request.Context(), userID.String(), req)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}
