package service

import (
	"context"
	"github.com/bitbiz/hias-core/domains/claims/entity"
)

type AdjudicatorService interface {
	Adjudicate(ctx context.Context, claim *entity.Claim, lineItems []*entity.ClaimLineItem) (*entity.AdjudicationResult, error)
}
