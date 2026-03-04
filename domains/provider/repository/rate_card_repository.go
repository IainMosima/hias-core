package repository

import (
	"context"
	"github.com/bitbiz/hias-core/domains/provider/entity"
	"github.com/google/uuid"
)

type RateCardRepository interface {
	Create(ctx context.Context, rateCard *entity.RateCard) (*entity.RateCard, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.RateCard, error)
	GetByProviderAndProcedure(ctx context.Context, providerID uuid.UUID, procedureCode string) (*entity.RateCard, error)
	GetByProviderProcedureAndAge(ctx context.Context, providerID uuid.UUID, procedureCode string, age int) (*entity.RateCard, error)
	ListByProvider(ctx context.Context, providerID uuid.UUID) ([]*entity.RateCard, error)
	Update(ctx context.Context, rateCard *entity.RateCard) (*entity.RateCard, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
