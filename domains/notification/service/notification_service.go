package service

import (
	"context"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/google/uuid"
)

type NotificationService interface {
	Send(ctx context.Context, userID uuid.UUID, channel, notifType, subject, body string) *schema.ServiceResponse[string]
	SendBulk(ctx context.Context, userIDs []uuid.UUID, channel, notifType, subject, body string) *schema.ServiceResponse[int]
	MarkRead(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[string]
	GetUnreadCount(ctx context.Context, userID uuid.UUID) *schema.ServiceResponse[int64]
	ListByUser(ctx context.Context, userID uuid.UUID, page, pageSize int) *schema.ServiceResponse[interface{}]
	RetryFailed(ctx context.Context) *schema.ServiceResponse[int]
}
