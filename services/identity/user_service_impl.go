package identity

import (
	"context"
	"net/http"

	"github.com/bitbiz/hias-core/domains/identity/entity"
	"github.com/bitbiz/hias-core/domains/identity/repository"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/bitbiz/hias-core/domains/identity/service"
	"github.com/bitbiz/hias-core/shared/auth"
	"github.com/google/uuid"
)

type userServiceImpl struct {
	userRepo       repository.UserRepository
	roleRepo       repository.RoleRepository
	permissionRepo repository.PermissionRepository
	cognitoService service.CognitoService
	tokenMaker     auth.TokenMaker
}

func NewUserService(
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	permissionRepo repository.PermissionRepository,
	cognitoService service.CognitoService,
	tokenMaker auth.TokenMaker,
) service.UserService {
	return &userServiceImpl{
		userRepo:       userRepo,
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
		cognitoService: cognitoService,
		tokenMaker:     tokenMaker,
	}
}

func (s *userServiceImpl) CreateUser(ctx context.Context, req schema.CreateUserRequest, createdBy uuid.UUID) *schema.ServiceResponse[schema.UserResponse] {
	role, err := s.roleRepo.GetByName(ctx, req.RoleName)
	if err != nil {
		return schema.NewServiceErrorResponse[schema.UserResponse](http.StatusBadRequest, "Invalid role", err)
	}

	cognitoSub, err := s.cognitoService.SignUp(ctx, req.Email, req.Password, req.Name, req.Phone)
	if err != nil {
		return schema.NewServiceErrorResponse[schema.UserResponse](http.StatusBadRequest, "Failed to create Cognito user", err)
	}

	// Auto-confirm the user
	_ = s.cognitoService.AdminSetPassword(ctx, req.Email, req.Password)

	user, err := s.userRepo.Create(ctx, &entity.User{
		CognitoSub: cognitoSub,
		Email:      req.Email,
		Name:       req.Name,
		Phone:      req.Phone,
		NationalID: req.NationalID,
		RoleID:     role.ID,
		Status:     "ACTIVE",
		CreatedBy:  createdBy,
	})
	if err != nil {
		return schema.NewServiceErrorResponse[schema.UserResponse](http.StatusInternalServerError, "Failed to create user", err)
	}

	user.RoleName = role.Name
	return schema.NewServiceResponse(schema.ToUserResponse(user), http.StatusCreated, "User created successfully")
}

func (s *userServiceImpl) GetUserByID(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[schema.UserResponse] {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[schema.UserResponse](http.StatusNotFound, "User not found", err)
	}

	role, _ := s.roleRepo.GetByID(ctx, user.RoleID)
	if role != nil {
		user.RoleName = role.Name
	}

	return schema.NewServiceResponse(schema.ToUserResponse(user), http.StatusOK, "User retrieved")
}

func (s *userServiceImpl) ListUsers(ctx context.Context, page, pageSize int) *schema.ServiceResponse[[]schema.UserResponse] {
	offset := (page - 1) * pageSize
	users, err := s.userRepo.List(ctx, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]schema.UserResponse](http.StatusInternalServerError, "Failed to list users", err)
	}
	return schema.NewServiceResponse(schema.ToUserResponseList(users), http.StatusOK, "Users retrieved")
}

func (s *userServiceImpl) UpdateUser(ctx context.Context, id uuid.UUID, req schema.UpdateUserRequest) *schema.ServiceResponse[schema.UserResponse] {
	existing, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[schema.UserResponse](http.StatusNotFound, "User not found", err)
	}

	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Phone != nil {
		existing.Phone = *req.Phone
	}
	if req.NationalID != nil {
		existing.NationalID = *req.NationalID
	}

	updated, err := s.userRepo.Update(ctx, existing)
	if err != nil {
		return schema.NewServiceErrorResponse[schema.UserResponse](http.StatusInternalServerError, "Failed to update user", err)
	}

	return schema.NewServiceResponse(schema.ToUserResponse(updated), http.StatusOK, "User updated")
}

func (s *userServiceImpl) AssignRole(ctx context.Context, userID uuid.UUID, req schema.AssignRoleRequest) *schema.ServiceResponse[schema.UserResponse] {
	updated, err := s.userRepo.UpdateRole(ctx, userID, req.RoleID)
	if err != nil {
		return schema.NewServiceErrorResponse[schema.UserResponse](http.StatusInternalServerError, "Failed to assign role", err)
	}

	role, _ := s.roleRepo.GetByID(ctx, req.RoleID)
	if role != nil {
		updated.RoleName = role.Name
	}

	return schema.NewServiceResponse(schema.ToUserResponse(updated), http.StatusOK, "Role assigned")
}

func (s *userServiceImpl) UpdateStatus(ctx context.Context, id uuid.UUID, req schema.UpdateStatusRequest) *schema.ServiceResponse[schema.UserResponse] {
	updated, err := s.userRepo.UpdateStatus(ctx, id, req.Status)
	if err != nil {
		return schema.NewServiceErrorResponse[schema.UserResponse](http.StatusInternalServerError, "Failed to update status", err)
	}
	return schema.NewServiceResponse(schema.ToUserResponse(updated), http.StatusOK, "Status updated")
}

func (s *userServiceImpl) GetTotalCount(ctx context.Context) *schema.ServiceResponse[int64] {
	count, err := s.userRepo.Count(ctx)
	if err != nil {
		return schema.NewServiceErrorResponse[int64](http.StatusInternalServerError, "Failed to count users", err)
	}
	return schema.NewServiceResponse(count, http.StatusOK, "Count retrieved")
}
