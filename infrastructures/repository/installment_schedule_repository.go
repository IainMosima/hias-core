package repository

import (
	"context"
	"fmt"

	"github.com/bitbiz/hias-core/domains/billing/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/billing/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
)

type installmentScheduleRepository struct {
	store db.Store
}

func NewInstallmentScheduleRepository(store db.Store) domainRepo.InstallmentScheduleRepository {
	return &installmentScheduleRepository{store: store}
}

func (r *installmentScheduleRepository) Create(ctx context.Context, schedule *entity.InstallmentSchedule) (*entity.InstallmentSchedule, error) {
	dbSchedule, err := r.store.CreateInstallmentSchedule(ctx, db.CreateInstallmentScheduleParams{
		PolicyID:             schedule.PolicyID,
		Frequency:            schedule.Frequency,
		TotalInstallments:    int32(schedule.TotalInstallments),
		AmountPerInstallment: schedule.AmountPerInstallment,
		StartDate:            schedule.StartDate,
		Status:               schedule.Status,
		CreatedBy:            uuidToPgtype(schedule.CreatedBy),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create installment schedule: %w", err)
	}
	return sqlcInstallmentScheduleToDomain(dbSchedule), nil
}

func (r *installmentScheduleRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.InstallmentSchedule, error) {
	dbSchedule, err := r.store.GetInstallmentScheduleByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get installment schedule by ID: %w", err)
	}
	return sqlcInstallmentScheduleToDomain(dbSchedule), nil
}

func (r *installmentScheduleRepository) GetByPolicy(ctx context.Context, policyID uuid.UUID) ([]*entity.InstallmentSchedule, error) {
	dbSchedules, err := r.store.GetInstallmentScheduleByPolicy(ctx, policyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get installment schedules by policy: %w", err)
	}
	schedules := make([]*entity.InstallmentSchedule, len(dbSchedules))
	for i, s := range dbSchedules {
		schedules[i] = sqlcInstallmentScheduleToDomain(s)
	}
	return schedules, nil
}

func (r *installmentScheduleRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*entity.InstallmentSchedule, error) {
	dbSchedule, err := r.store.UpdateInstallmentScheduleStatus(ctx, db.UpdateInstallmentScheduleStatusParams{
		ID:     id,
		Status: status,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update installment schedule status: %w", err)
	}
	return sqlcInstallmentScheduleToDomain(dbSchedule), nil
}

func sqlcInstallmentScheduleToDomain(s db.InstallmentSchedule) *entity.InstallmentSchedule {
	return &entity.InstallmentSchedule{
		ID:                   s.ID,
		PolicyID:             s.PolicyID,
		Frequency:            s.Frequency,
		TotalInstallments:    int(s.TotalInstallments),
		AmountPerInstallment: s.AmountPerInstallment,
		StartDate:            s.StartDate,
		Status:               s.Status,
		CreatedBy:            pgtypeToUUID(s.CreatedBy),
		CreatedAt:            s.CreatedAt,
		UpdatedAt:            s.UpdatedAt,
	}
}
