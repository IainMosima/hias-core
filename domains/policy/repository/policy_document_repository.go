package repository

import (
	"context"

	"github.com/bitbiz/hias-core/domains/policy/entity"
	"github.com/google/uuid"
)

type PolicyDocumentRepository interface {
	Create(ctx context.Context, doc *entity.PolicyDocument) (*entity.PolicyDocument, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.PolicyDocument, error)
	ListByPolicy(ctx context.Context, policyID uuid.UUID) ([]*entity.PolicyDocument, error)
	ListByMember(ctx context.Context, memberID uuid.UUID) ([]*entity.PolicyDocument, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
