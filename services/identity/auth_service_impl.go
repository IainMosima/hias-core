package identity

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bitbiz/hias-core/domains/identity/entity"
	"github.com/bitbiz/hias-core/domains/identity/repository"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/bitbiz/hias-core/domains/identity/service"
	"github.com/bitbiz/hias-core/shared/auth"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/google/uuid"
)

type authServiceImpl struct {
	userRepo       repository.UserRepository
	roleRepo       repository.RoleRepository
	permissionRepo repository.PermissionRepository
	tokenMaker     auth.TokenMaker
	cognitoService service.CognitoService
	accessDuration interface{}
}

func NewAuthService(
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	permissionRepo repository.PermissionRepository,
	tokenMaker auth.TokenMaker,
	cognitoService service.CognitoService,
	accessDuration interface{},
) service.AuthService {
	return &authServiceImpl{
		userRepo:       userRepo,
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
		tokenMaker:     tokenMaker,
		cognitoService: cognitoService,
		accessDuration: accessDuration,
	}
}

func (s *authServiceImpl) Login(ctx context.Context, req schema.LoginRequest) *schema.ServiceResponse[schema.LoginResponse] {
	// Authenticate with Cognito
	_, err := s.cognitoService.SignIn(ctx, req.Email, req.Password)
	if err != nil {
		utils.LogError("Cognito sign in failed: %v", err)
		return schema.NewServiceErrorResponse[schema.LoginResponse](http.StatusUnauthorized, "Invalid credentials", err)
	}

	// Get user from DB
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return schema.NewServiceErrorResponse[schema.LoginResponse](http.StatusNotFound, "User not found", err)
	}

	// Get role
	role, err := s.roleRepo.GetByID(ctx, user.RoleID)
	if err != nil {
		return schema.NewServiceErrorResponse[schema.LoginResponse](http.StatusInternalServerError, "Failed to get role", err)
	}

	// Get permissions
	permissions, err := s.permissionRepo.ListByRole(ctx, user.RoleID)
	if err != nil {
		return schema.NewServiceErrorResponse[schema.LoginResponse](http.StatusInternalServerError, "Failed to get permissions", err)
	}

	permStrings := make([]string, len(permissions))
	for i, p := range permissions {
		permStrings[i] = fmt.Sprintf("%s:%s", p.Resource, p.Action)
	}

	// Create PASETO token
	token, payload, err := s.tokenMaker.CreateToken(user.ID.String(), user.Email, role.Name, permStrings, s.accessDuration)
	if err != nil {
		return schema.NewServiceErrorResponse[schema.LoginResponse](http.StatusInternalServerError, "Failed to create token", err)
	}

	user.RoleName = role.Name

	return schema.NewServiceResponse(schema.LoginResponse{
		AccessToken:          token,
		AccessTokenExpiresAt: payload.ExpiredAt,
		User:                 schema.ToUserResponse(user),
	}, http.StatusOK, "Login successful")
}

func (s *authServiceImpl) Register(ctx context.Context, req schema.RegisterRequest) *schema.ServiceResponse[schema.RegisterResponse] {
	// Default role to Member
	roleName := req.RoleName
	if roleName == "" {
		roleName = "Member"
	}

	role, err := s.roleRepo.GetByName(ctx, roleName)
	if err != nil {
		return schema.NewServiceErrorResponse[schema.RegisterResponse](http.StatusBadRequest, "Invalid role", err)
	}

	// Sign up with Cognito
	cognitoSub, err := s.cognitoService.SignUp(ctx, req.Email, req.Password, req.Name, req.Phone)
	if err != nil {
		return schema.NewServiceErrorResponse[schema.RegisterResponse](http.StatusBadRequest, "Registration failed", err)
	}

	// Create user in DB
	user, err := s.userRepo.Create(ctx, &entity.User{
		CognitoSub: cognitoSub,
		Email:      req.Email,
		Name:       req.Name,
		Phone:      req.Phone,
		NationalID: req.NationalID,
		RoleID:     role.ID,
		Status:     "ACTIVE",
	})
	if err != nil {
		return schema.NewServiceErrorResponse[schema.RegisterResponse](http.StatusInternalServerError, "Failed to create user", err)
	}

	return schema.NewServiceResponse(schema.RegisterResponse{
		UserID:  user.ID.String(),
		Email:   user.Email,
		Message: "Registration successful. Please verify your email.",
	}, http.StatusCreated, "Registration successful")
}

func (s *authServiceImpl) RefreshToken(ctx context.Context, req schema.RefreshTokenRequest) *schema.ServiceResponse[schema.LoginResponse] {
	// Verify the refresh token
	payload, err := s.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		return schema.NewServiceErrorResponse[schema.LoginResponse](http.StatusUnauthorized, "Invalid refresh token", err)
	}

	userID, _ := uuid.Parse(payload.UserID)
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return schema.NewServiceErrorResponse[schema.LoginResponse](http.StatusNotFound, "User not found", err)
	}

	role, err := s.roleRepo.GetByID(ctx, user.RoleID)
	if err != nil {
		return schema.NewServiceErrorResponse[schema.LoginResponse](http.StatusInternalServerError, "Failed to get role", err)
	}

	permissions, err := s.permissionRepo.ListByRole(ctx, user.RoleID)
	if err != nil {
		return schema.NewServiceErrorResponse[schema.LoginResponse](http.StatusInternalServerError, "Failed to get permissions", err)
	}

	permStrings := make([]string, len(permissions))
	for i, p := range permissions {
		permStrings[i] = fmt.Sprintf("%s:%s", p.Resource, p.Action)
	}

	token, tokenPayload, err := s.tokenMaker.CreateToken(user.ID.String(), user.Email, role.Name, permStrings, s.accessDuration)
	if err != nil {
		return schema.NewServiceErrorResponse[schema.LoginResponse](http.StatusInternalServerError, "Failed to create token", err)
	}

	user.RoleName = role.Name

	return schema.NewServiceResponse(schema.LoginResponse{
		AccessToken:          token,
		AccessTokenExpiresAt: tokenPayload.ExpiredAt,
		User:                 schema.ToUserResponse(user),
	}, http.StatusOK, "Token refreshed")
}

func (s *authServiceImpl) Logout(_ context.Context, _ string) *schema.ServiceResponse[string] {
	return schema.NewServiceResponse("Logged out successfully", http.StatusOK, "Logout successful")
}

func (s *authServiceImpl) ForgotPassword(ctx context.Context, req schema.ForgotPasswordRequest) *schema.ServiceResponse[string] {
	err := s.cognitoService.ForgotPassword(ctx, req.Email)
	if err != nil {
		return schema.NewServiceErrorResponse[string](http.StatusBadRequest, "Failed to initiate password reset", err)
	}
	return schema.NewServiceResponse("Password reset code sent", http.StatusOK, "Check your email for the reset code")
}

func (s *authServiceImpl) ResetPassword(ctx context.Context, req schema.ResetPasswordRequest) *schema.ServiceResponse[string] {
	err := s.cognitoService.ConfirmForgotPassword(ctx, req.Email, req.Code, req.NewPassword)
	if err != nil {
		return schema.NewServiceErrorResponse[string](http.StatusBadRequest, "Failed to reset password", err)
	}
	return schema.NewServiceResponse("Password reset successful", http.StatusOK, "Password has been reset")
}

func (s *authServiceImpl) ConfirmSignUp(ctx context.Context, req schema.ConfirmSignUpRequest) *schema.ServiceResponse[string] {
	err := s.cognitoService.ConfirmSignUp(ctx, req.Email, req.Code)
	if err != nil {
		return schema.NewServiceErrorResponse[string](http.StatusBadRequest, "Failed to confirm sign up", err)
	}
	return schema.NewServiceResponse("Email confirmed", http.StatusOK, "Email verification successful")
}
