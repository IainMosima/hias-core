package auth

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

type CognitoService struct {
	client     *cognitoidentityprovider.Client
	clientID   string
	userPoolID string
}

func NewCognitoService(client *cognitoidentityprovider.Client, clientID, userPoolID string) *CognitoService {
	return &CognitoService{
		client:     client,
		clientID:   clientID,
		userPoolID: userPoolID,
	}
}

func (s *CognitoService) SignUp(ctx context.Context, email, password, name, phone string) (string, error) {
	input := &cognitoidentityprovider.SignUpInput{
		ClientId: aws.String(s.clientID),
		Username: aws.String(email),
		Password: aws.String(password),
		UserAttributes: []types.AttributeType{
			{Name: aws.String("email"), Value: aws.String(email)},
			{Name: aws.String("name"), Value: aws.String(name)},
			{Name: aws.String("phone_number"), Value: aws.String(phone)},
		},
	}

	result, err := s.client.SignUp(ctx, input)
	if err != nil {
		return "", fmt.Errorf("cognito sign up failed: %w", err)
	}

	return *result.UserSub, nil
}

func (s *CognitoService) SignIn(ctx context.Context, email, password string) (*cognitoidentityprovider.InitiateAuthOutput, error) {
	input := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: types.AuthFlowTypeUserPasswordAuth,
		ClientId: aws.String(s.clientID),
		AuthParameters: map[string]string{
			"USERNAME": email,
			"PASSWORD": password,
		},
	}

	return s.client.InitiateAuth(ctx, input)
}

func (s *CognitoService) ConfirmSignUp(ctx context.Context, email, code string) error {
	input := &cognitoidentityprovider.ConfirmSignUpInput{
		ClientId:         aws.String(s.clientID),
		Username:         aws.String(email),
		ConfirmationCode: aws.String(code),
	}

	_, err := s.client.ConfirmSignUp(ctx, input)
	return err
}

func (s *CognitoService) AdminSetPassword(ctx context.Context, email, password string) error {
	input := &cognitoidentityprovider.AdminSetUserPasswordInput{
		UserPoolId: aws.String(s.userPoolID),
		Username:   aws.String(email),
		Password:   aws.String(password),
		Permanent:  true,
	}

	_, err := s.client.AdminSetUserPassword(ctx, input)
	return err
}

func (s *CognitoService) ForgotPassword(ctx context.Context, email string) error {
	input := &cognitoidentityprovider.ForgotPasswordInput{
		ClientId: aws.String(s.clientID),
		Username: aws.String(email),
	}

	_, err := s.client.ForgotPassword(ctx, input)
	return err
}

func (s *CognitoService) ConfirmForgotPassword(ctx context.Context, email, code, newPassword string) error {
	input := &cognitoidentityprovider.ConfirmForgotPasswordInput{
		ClientId:         aws.String(s.clientID),
		Username:         aws.String(email),
		ConfirmationCode: aws.String(code),
		Password:         aws.String(newPassword),
	}

	_, err := s.client.ConfirmForgotPassword(ctx, input)
	return err
}
