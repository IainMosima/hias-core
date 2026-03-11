package identity

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/bitbiz/hias-core/domains/identity/entity"
	"github.com/bitbiz/hias-core/domains/identity/repository"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/bitbiz/hias-core/domains/identity/service"
	"github.com/bitbiz/hias-core/shared"
	"github.com/bitbiz/hias-core/shared/auth"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/google/uuid"
)

type authServiceImpl struct {
	userRepo       repository.UserRepository
	roleRepo       repository.RoleRepository
	permissionRepo repository.PermissionRepository
	tokenMaker     auth.TokenMaker
	config         AuthServiceConfig
}

type AuthServiceConfig struct {
	AccessTokenDuration  interface{}
	RefreshTokenDuration interface{}
}

func NewAuthService(
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	permissionRepo repository.PermissionRepository,
	tokenMaker auth.TokenMaker,
	config AuthServiceConfig,
) service.AuthService {
	return &authServiceImpl{
		userRepo:       userRepo,
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
		tokenMaker:     tokenMaker,
		config:         config,
	}
}

func (s *authServiceImpl) Login(ctx context.Context, req schema.LoginRequest) *schema.ServiceResponse[schema.LoginResponse] {
	// Get user from DB
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		log.Printf("[LOGIN] User not found for email=%s: %v", req.Email, err)
		return schema.NewServiceErrorResponse[schema.LoginResponse](http.StatusUnauthorized, "Invalid credentials", err)
	}

	log.Printf("[LOGIN] User found: id=%s email=%s status=%s hash_len=%d", user.ID, user.Email, user.Status, len(user.PasswordHash))

	// Load role name (same pattern as GetUserByID in user_service_impl.go)
	if role, roleErr := s.roleRepo.GetByID(ctx, user.RoleID); roleErr == nil {
		user.RoleName = role.Name
	}

	// Verify password
	if err := utils.CheckPassword(req.Password, user.PasswordHash); err != nil {
		log.Printf("[LOGIN] Password mismatch for email=%s: %v", req.Email, err)
		return schema.NewServiceErrorResponse[schema.LoginResponse](http.StatusUnauthorized, "Invalid credentials", err)
	}

	if user.Status != string(shared.UserStatusActive) {
		return schema.NewServiceErrorResponse[schema.LoginResponse](http.StatusForbidden, "Account is not active", fmt.Errorf("user status: %s", user.Status))
	}

	// Get permissions for role
	permissions := s.getPermissionStrings(ctx, user.RoleID)

	// Mint PASETO token
	accessToken, payload, err := s.tokenMaker.CreateToken(
		user.ID.String(), user.Email, user.RoleName, permissions, s.config.AccessTokenDuration,
	)
	if err != nil {
		return schema.NewServiceErrorResponse[schema.LoginResponse](http.StatusInternalServerError, "Failed to create token", err)
	}

	// Mint refresh token
	refreshToken, _, err := s.tokenMaker.CreateToken(
		user.ID.String(), user.Email, user.RoleName, permissions, s.config.RefreshTokenDuration,
	)
	if err != nil {
		return schema.NewServiceErrorResponse[schema.LoginResponse](http.StatusInternalServerError, "Failed to create refresh token", err)
	}

	response := schema.LoginResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: payload.ExpiredAt,
		RefreshToken:         refreshToken,
		User:                 schema.ToUserResponse(user),
	}

	return schema.NewServiceResponse(response, http.StatusOK, "Login successful")
}

func (s *authServiceImpl) Register(ctx context.Context, req schema.RegisterRequest) *schema.ServiceResponse[schema.RegisterResponse] {
	// Hash password
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return schema.NewServiceErrorResponse[schema.RegisterResponse](http.StatusInternalServerError, "Registration failed", err)
	}
	log.Printf("[REGISTER] Password hashed: len=%d, starts_with=%s", len(passwordHash), passwordHash[:7])

	// Find role
	roleName := req.RoleName
	if roleName == "" {
		roleName = string(shared.UserRoleAdmin)
	}
	role, err := s.roleRepo.GetByName(ctx, roleName)
	if err != nil {
		return schema.NewServiceErrorResponse[schema.RegisterResponse](
			http.StatusBadRequest,
			fmt.Sprintf("Role '%s' not found. Please ensure roles are seeded.", roleName),
			err,
		)
	}

	// Create user in DB
	user := &entity.User{
		Email:        req.Email,
		Name:         req.Name,
		Phone:        req.Phone,
		NationalID:   req.NationalID,
		PasswordHash: passwordHash,
		RoleID:       role.ID,
		RoleName:     roleName,
		Status:       string(shared.UserStatusActive),
	}

	created, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return schema.NewServiceErrorResponse[schema.RegisterResponse](http.StatusInternalServerError, "Failed to create user", err)
	}

	response := schema.RegisterResponse{
		UserID:  created.ID.String(),
		Email:   created.Email,
		Message: "Registration successful.",
	}

	return schema.NewServiceResponse(response, http.StatusCreated, "User registered")
}

func (s *authServiceImpl) RefreshToken(ctx context.Context, req schema.RefreshTokenRequest) *schema.ServiceResponse[schema.LoginResponse] {
	// Verify the refresh token
	payload, err := s.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		return schema.NewServiceErrorResponse[schema.LoginResponse](http.StatusUnauthorized, "Invalid refresh token", err)
	}

	userID, err := uuid.Parse(payload.UserID)
	if err != nil {
		return schema.NewServiceErrorResponse[schema.LoginResponse](http.StatusUnauthorized, "Invalid user ID in token", err)
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return schema.NewServiceErrorResponse[schema.LoginResponse](http.StatusUnauthorized, "User not found", err)
	}

	// Load role name
	if role, roleErr := s.roleRepo.GetByID(ctx, user.RoleID); roleErr == nil {
		user.RoleName = role.Name
	}

	permissions := s.getPermissionStrings(ctx, user.RoleID)

	accessToken, newPayload, err := s.tokenMaker.CreateToken(
		user.ID.String(), user.Email, user.RoleName, permissions, s.config.AccessTokenDuration,
	)
	if err != nil {
		return schema.NewServiceErrorResponse[schema.LoginResponse](http.StatusInternalServerError, "Failed to create token", err)
	}

	response := schema.LoginResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: newPayload.ExpiredAt,
		User:                 schema.ToUserResponse(user),
	}

	return schema.NewServiceResponse(response, http.StatusOK, "Token refreshed")
}

func (s *authServiceImpl) Logout(ctx context.Context, userID string) *schema.ServiceResponse[string] {
	return schema.NewServiceResponse("Logged out", http.StatusOK, "Logout successful")
}

func (s *authServiceImpl) ChangePassword(ctx context.Context, userID string, req schema.ChangePasswordRequest) *schema.ServiceResponse[string] {
	id, err := uuid.Parse(userID)
	if err != nil {
		return schema.NewServiceErrorResponse[string](http.StatusBadRequest, "Invalid user ID", err)
	}

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[string](http.StatusNotFound, "User not found", err)
	}

	if err := utils.CheckPassword(req.CurrentPassword, user.PasswordHash); err != nil {
		return schema.NewServiceErrorResponse[string](http.StatusUnauthorized, "Current password is incorrect", err)
	}

	newHash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return schema.NewServiceErrorResponse[string](http.StatusInternalServerError, "Failed to hash password", err)
	}

	if err := s.userRepo.UpdatePassword(ctx, id, newHash); err != nil {
		return schema.NewServiceErrorResponse[string](http.StatusInternalServerError, "Failed to update password", err)
	}

	return schema.NewServiceResponse("Password changed successfully", http.StatusOK, "Password changed")
}

func (s *authServiceImpl) getPermissionStrings(ctx context.Context, roleID uuid.UUID) []string {
	if roleID == uuid.Nil {
		return []string{}
	}
	perms, err := s.permissionRepo.ListByRole(ctx, roleID)
	if err != nil {
		return []string{}
	}
	strs := make([]string, len(perms))
	for i, p := range perms {
		strs[i] = p.Resource + ":" + p.Action
	}
	return strs
}
