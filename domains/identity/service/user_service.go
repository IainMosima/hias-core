package service

import (
	"context"

	"github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/google/uuid"
)

type UserService interface {
	CreateUser(ctx context.Context, req schema.CreateUserRequest, createdBy uuid.UUID) *schema.ServiceResponse[schema.UserResponse]
	GetUserByID(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[schema.UserResponse]
	ListUsers(ctx context.Context, page, pageSize int) *schema.ServiceResponse[[]schema.UserResponse]
	UpdateUser(ctx context.Context, id uuid.UUID, req schema.UpdateUserRequest) *schema.ServiceResponse[schema.UserResponse]
	AssignRole(ctx context.Context, userID uuid.UUID, req schema.AssignRoleRequest) *schema.ServiceResponse[schema.UserResponse]
	UpdateStatus(ctx context.Context, id uuid.UUID, req schema.UpdateStatusRequest) *schema.ServiceResponse[schema.UserResponse]
	GetTotalCount(ctx context.Context) *schema.ServiceResponse[int64]
}
