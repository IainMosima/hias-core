package service

import (
	"context"
	"time"

	"github.com/bitbiz/hias-core/domains/claims/entity"
	"github.com/google/uuid"
)

type FraudService interface {
	CheckDuplicate(ctx context.Context, claimNumber string, claimID uuid.UUID) (bool, error)
	CheckFrequency(ctx context.Context, memberID, providerID uuid.UUID, procedureCode string, claimID uuid.UUID) (bool, error)
	CheckAmountThreshold(ctx context.Context, providerID uuid.UUID, procedureCode string, amount int64) (bool, error)
	CheckExpiredContract(ctx context.Context, providerID uuid.UUID, serviceDate time.Time) (bool, error)
	CheckSuspendedProvider(ctx context.Context, providerID uuid.UUID) (bool, error)
	CheckRateCardOvercharge(ctx context.Context, providerID uuid.UUID, procedureCode string, amount int64) (bool, error)
	CheckRepeatVisit(ctx context.Context, memberID, providerID uuid.UUID, procedureCode string, serviceDate time.Time, claimID uuid.UUID) (bool, error)
	FlagClaim(ctx context.Context, flag *entity.FraudFlag) error
}
