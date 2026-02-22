package notifications

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

type EmailAdapter struct {
	client    *ses.Client
	fromEmail string
}

func NewEmailAdapter(client *ses.Client, fromEmail string) *EmailAdapter {
	return &EmailAdapter{
		client:    client,
		fromEmail: fromEmail,
	}
}

func (a *EmailAdapter) Send(ctx context.Context, to, subject, body string) error {
	input := &ses.SendEmailInput{
		Source: aws.String(a.fromEmail),
		Destination: &types.Destination{
			ToAddresses: []string{to},
		},
		Message: &types.Message{
			Subject: &types.Content{
				Data:    aws.String(subject),
				Charset: aws.String("UTF-8"),
			},
			Body: &types.Body{
				Html: &types.Content{
					Data:    aws.String(body),
					Charset: aws.String("UTF-8"),
				},
			},
		},
	}

	_, err := a.client.SendEmail(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to send email via SES: %w", err)
	}

	return nil
}
