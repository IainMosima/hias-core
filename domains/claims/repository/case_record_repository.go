package repository

import (
	"context"
	"time"

	"github.com/bitbiz/hias-core/domains/claims/entity"
	"github.com/google/uuid"
)

type CaseRecordRepository interface {
	Create(ctx context.Context, record *entity.CaseRecord) (*entity.CaseRecord, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.CaseRecord, error)
	GetByNumber(ctx context.Context, number string) (*entity.CaseRecord, error)
	GetByPreAuth(ctx context.Context, preauthID uuid.UUID) (*entity.CaseRecord, error)
	ListByPolicy(ctx context.Context, policyID uuid.UUID, limit, offset int) ([]*entity.CaseRecord, error)
	ListByMember(ctx context.Context, memberID uuid.UUID, limit, offset int) ([]*entity.CaseRecord, error)
	ListByProvider(ctx context.Context, providerID uuid.UUID, limit, offset int) ([]*entity.CaseRecord, error)
	ListByStatus(ctx context.Context, status string, limit, offset int) ([]*entity.CaseRecord, error)
	Update(ctx context.Context, record *entity.CaseRecord) (*entity.CaseRecord, error)
	Admit(ctx context.Context, id uuid.UUID, admissionDate time.Time) (*entity.CaseRecord, error)
	Discharge(ctx context.Context, id uuid.UUID, dischargeDate time.Time, actualCost int64) (*entity.CaseRecord, error)
	Close(ctx context.Context, id uuid.UUID) (*entity.CaseRecord, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*entity.CaseRecord, error)
	CountByStatus(ctx context.Context, status string) (int64, error)
}
