package repository

import (
	"context"

	"github.com/bitbiz/hias-core/domains/claims/entity"
	"github.com/google/uuid"
)

type EscalationRuleRepository interface {
	Create(ctx context.Context, rule *entity.EscalationRule) (*entity.EscalationRule, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.EscalationRule, error)
	List(ctx context.Context, limit, offset int) ([]*entity.EscalationRule, error)
	ListActive(ctx context.Context) ([]*entity.EscalationRule, error)
	Update(ctx context.Context, rule *entity.EscalationRule) (*entity.EscalationRule, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
