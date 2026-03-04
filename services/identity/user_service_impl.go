package identity

import (
	"context"
	"fmt"
	"log"
	"net/http"

	auditService "github.com/bitbiz/hias-core/domains/audit/service"
	"github.com/bitbiz/hias-core/domains/identity/entity"
	"github.com/bitbiz/hias-core/domains/identity/repository"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/bitbiz/hias-core/domains/identity/service"
	"github.com/bitbiz/hias-core/shared"
	"github.com/google/uuid"
)

type userServiceImpl struct {
	userRepo repository.UserRepository
	roleRepo repository.RoleRepository
	auditSvc auditService.AuditService
}

func NewUserService(
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	auditSvc auditService.AuditService,
) service.UserService {
	return &userServiceImpl{
		userRepo: userRepo,
		roleRepo: roleRepo,
		auditSvc: auditSvc,
	}
}

func (s *userServiceImpl) CreateUser(ctx context.Context, req schema.CreateUserRequest, createdBy uuid.UUID) *schema.ServiceResponse[schema.UserResponse] {
	// Find role
	role, err := s.roleRepo.GetByName(ctx, req.RoleName)
	if err != nil {
		return schema.NewServiceErrorResponse[schema.UserResponse](http.StatusBadRequest, fmt.Sprintf("Role '%s' not found", req.RoleName), err)
	}

	user := &entity.User{
		Email:      req.Email,
		Name:       req.Name,
		Phone:      req.Phone,
		NationalID: req.NationalID,
		RoleID:     role.ID,
		RoleName:   role.Name,
		Status:     string(shared.UserStatusActive),
		CreatedBy:  createdBy,
	}

	created, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return schema.NewServiceErrorResponse[schema.UserResponse](http.StatusInternalServerError, "Failed to create user", err)
	}
	created.RoleName = role.Name

	s.logAudit(ctx, createdBy, string(shared.AuditEntityTypeUser), created.ID, string(shared.AuditActionCreate))

	return schema.NewServiceResponse(schema.ToUserResponse(created), http.StatusCreated, "User created")
}

func (s *userServiceImpl) GetUserByID(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[schema.UserResponse] {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[schema.UserResponse](http.StatusNotFound, "User not found", err)
	}

	// Load role name
	if role, roleErr := s.roleRepo.GetByID(ctx, user.RoleID); roleErr == nil {
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
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[schema.UserResponse](http.StatusNotFound, "User not found", err)
	}

	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Phone != nil {
		user.Phone = *req.Phone
	}
	if req.NationalID != nil {
		user.NationalID = *req.NationalID
	}

	updated, err := s.userRepo.Update(ctx, user)
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

	if role, roleErr := s.roleRepo.GetByID(ctx, req.RoleID); roleErr == nil {
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
		return schema.NewServiceErrorResponse[int64](http.StatusInternalServerError, "Failed to get count", err)
	}
	return schema.NewServiceResponse(count, http.StatusOK, "Count retrieved")
}

func (s *userServiceImpl) logAudit(ctx context.Context, userID uuid.UUID, entityType string, entityID uuid.UUID, action string) {
	if s.auditSvc != nil {
		resp := s.auditSvc.LogEvent(ctx, userID, entityType, entityID, action, nil, nil, "", "")
		if resp.Error != nil {
			log.Printf("Failed to log audit: %v", resp.Error)
		}
	}
}
