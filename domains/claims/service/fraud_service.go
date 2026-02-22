package service

import (
	"context"
	"github.com/bitbiz/hias-core/domains/claims/entity"
	"github.com/google/uuid"
)

type FraudService interface {
	CheckDuplicate(ctx context.Context, claimNumber string, claimID uuid.UUID) (bool, error)
	CheckFrequency(ctx context.Context, memberID, providerID uuid.UUID, procedureCode string, claimID uuid.UUID) (bool, error)
	CheckAmountThreshold(ctx context.Context, providerID uuid.UUID, procedureCode string, amount int64) (bool, error)
	FlagClaim(ctx context.Context, flag *entity.FraudFlag) error
}
