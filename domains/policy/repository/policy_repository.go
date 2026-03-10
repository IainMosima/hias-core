package repository

import (
	"context"
	"time"

	"github.com/bitbiz/hias-core/domains/policy/entity"
	"github.com/google/uuid"
)

type PolicyRepository interface {
	Create(ctx context.Context, policy *entity.Policy) (*entity.Policy, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Policy, error)
	GetByNumber(ctx context.Context, number string) (*entity.Policy, error)
	List(ctx context.Context, limit, offset int) ([]*entity.Policy, error)
	ListByStatus(ctx context.Context, status string, limit, offset int) ([]*entity.Policy, error)
	ListExpiringSoon(ctx context.Context, days int) ([]*entity.Policy, error)
	Count(ctx context.Context) (int64, error)
	CountByStatus(ctx context.Context, status string) (int64, error)
	GetActivePoliciesForBilling(ctx context.Context) ([]*entity.Policy, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*entity.Policy, error)
	ActivateWithTimestamp(ctx context.Context, id uuid.UUID) (*entity.Policy, error)
	Update(ctx context.Context, policy *entity.Policy) (*entity.Policy, error)
	UpdatePlanAndPremium(ctx context.Context, id uuid.UUID, planID uuid.UUID, premiumAmount int64) (*entity.Policy, error)
	GetLapsedForTermination(ctx context.Context) ([]*entity.Policy, error)
	GetOverdueForLapse(ctx context.Context) ([]*entity.Policy, error)
	ListFiltered(ctx context.Context, dateFrom, dateTo *time.Time, search string, limit, offset int) ([]*entity.Policy, error)
	CountFiltered(ctx context.Context, dateFrom, dateTo *time.Time, search string) (int64, error)
}
