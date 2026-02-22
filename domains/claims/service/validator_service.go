package service

import (
	"context"
	"github.com/bitbiz/hias-core/domains/claims/entity"
)

type ValidatorService interface {
	ValidateClaim(ctx context.Context, claim *entity.Claim, lineItems []*entity.ClaimLineItem) (bool, []string, error)
}
