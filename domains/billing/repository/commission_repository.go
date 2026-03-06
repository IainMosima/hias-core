package repository

import (
	"context"
	"github.com/bitbiz/hias-core/domains/billing/entity"
	"github.com/google/uuid"
)

type CommissionRuleRepository interface {
	Create(ctx context.Context, rule *entity.CommissionRule) (*entity.CommissionRule, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.CommissionRule, error)
	ListByPlan(ctx context.Context, planID uuid.UUID) ([]*entity.CommissionRule, error)
}

type CommissionPaymentRepository interface {
	Create(ctx context.Context, payment *entity.CommissionPayment) (*entity.CommissionPayment, error)
	List(ctx context.Context, limit, offset int) ([]*entity.CommissionPayment, error)
	ListByIntermediary(ctx context.Context, intermediaryID uuid.UUID, limit, offset int) ([]*entity.CommissionPayment, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*entity.CommissionPayment, error)
}
