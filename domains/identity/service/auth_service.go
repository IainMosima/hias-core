package service

import (
	"context"

	"github.com/bitbiz/hias-core/domains/identity/schema"
)

type AuthService interface {
	Login(ctx context.Context, req schema.LoginRequest) *schema.ServiceResponse[schema.LoginResponse]
	Register(ctx context.Context, req schema.RegisterRequest) *schema.ServiceResponse[schema.RegisterResponse]
	RefreshToken(ctx context.Context, req schema.RefreshTokenRequest) *schema.ServiceResponse[schema.LoginResponse]
	Logout(ctx context.Context, userID string) *schema.ServiceResponse[string]
	ForgotPassword(ctx context.Context, req schema.ForgotPasswordRequest) *schema.ServiceResponse[string]
	ResetPassword(ctx context.Context, req schema.ResetPasswordRequest) *schema.ServiceResponse[string]
	ConfirmSignUp(ctx context.Context, req schema.ConfirmSignUpRequest) *schema.ServiceResponse[string]
}
