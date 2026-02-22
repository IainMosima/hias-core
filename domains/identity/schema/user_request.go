package schema

import "github.com/google/uuid"

type CreateUserRequest struct {
	Email      string `json:"email" binding:"required,email"`
	Name       string `json:"name" binding:"required"`
	Phone      string `json:"phone" binding:"required"`
	NationalID string `json:"national_id"`
	RoleName   string `json:"role_name" binding:"required"`
	Password   string `json:"password" binding:"required,min=8"`
}

type UpdateUserRequest struct {
	Name       *string `json:"name"`
	Phone      *string `json:"phone"`
	NationalID *string `json:"national_id"`
}

type AssignRoleRequest struct {
	RoleID uuid.UUID `json:"role_id" binding:"required"`
}

type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=ACTIVE INACTIVE SUSPENDED"`
}
