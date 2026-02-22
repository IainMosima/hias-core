package service

import "context"

type CognitoService interface {
	SignUp(ctx context.Context, email, password, name, phone string) (string, error)
	SignIn(ctx context.Context, email, password string) (interface{}, error)
	ConfirmSignUp(ctx context.Context, email, code string) error
	AdminSetPassword(ctx context.Context, email, password string) error
	ForgotPassword(ctx context.Context, email string) error
	ConfirmForgotPassword(ctx context.Context, email, code, newPassword string) error
}
