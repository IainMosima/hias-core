package schema

import (
	"errors"
	"time"
)

type LoginResponse struct {
	AccessToken          string       `json:"access_token"`
	AccessTokenExpiresAt time.Time    `json:"access_token_expires_at"`
	RefreshToken         string       `json:"refresh_token,omitempty"`
	User                 UserResponse `json:"user"`
}

type RegisterResponse struct {
	UserID  string `json:"user_id"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

type ServiceResponse[T any] struct {
	Data       T      `json:"data,omitempty"`
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Error      error  `json:"-"`
}

func NewServiceResponse[T any](data T, statusCode int, message string) *ServiceResponse[T] {
	return &ServiceResponse[T]{
		Data:       data,
		StatusCode: statusCode,
		Message:    message,
	}
}

func NewServiceErrorResponse[T any](statusCode int, message string, err error) *ServiceResponse[T] {
	if err == nil {
		err = errors.New(message)
	}
	return &ServiceResponse[T]{
		StatusCode: statusCode,
		Message:    message,
		Error:      err,
	}
}
