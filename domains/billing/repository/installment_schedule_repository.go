package repository

import (
	"context"

	"github.com/bitbiz/hias-core/domains/billing/entity"
	"github.com/google/uuid"
)

type InstallmentScheduleRepository interface {
	Create(ctx context.Context, schedule *entity.InstallmentSchedule) (*entity.InstallmentSchedule, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.InstallmentSchedule, error)
	GetByPolicy(ctx context.Context, policyID uuid.UUID) ([]*entity.InstallmentSchedule, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*entity.InstallmentSchedule, error)
}
