package schema

import (
	"time"

	"github.com/bitbiz/hias-core/domains/identity/entity"
	"github.com/google/uuid"
)

type UserResponse struct {
	ID         uuid.UUID `json:"id"`
	Email      string    `json:"email"`
	Name       string    `json:"name"`
	Phone      string    `json:"phone"`
	NationalID string    `json:"national_id"`
	RoleID     uuid.UUID `json:"role_id"`
	RoleName   string    `json:"role_name,omitempty"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func ToUserResponse(user *entity.User) UserResponse {
	return UserResponse{
		ID:         user.ID,
		Email:      user.Email,
		Name:       user.Name,
		Phone:      user.Phone,
		NationalID: user.NationalID,
		RoleID:     user.RoleID,
		RoleName:   user.RoleName,
		Status:     user.Status,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}
}

func ToUserResponseList(users []*entity.User) []UserResponse {
	responses := make([]UserResponse, len(users))
	for i, user := range users {
		responses[i] = ToUserResponse(user)
	}
	return responses
}
