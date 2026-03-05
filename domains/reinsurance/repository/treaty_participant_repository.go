package repository

import (
	"context"

	"github.com/bitbiz/hias-core/domains/reinsurance/entity"
	"github.com/google/uuid"
)

type TreatyParticipantRepository interface {
	Create(ctx context.Context, participant *entity.TreatyParticipant) (*entity.TreatyParticipant, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.TreatyParticipant, error)
	ListByTreaty(ctx context.Context, treatyID uuid.UUID) ([]*entity.TreatyParticipant, error)
	Update(ctx context.Context, participant *entity.TreatyParticipant) (*entity.TreatyParticipant, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetTotalShareByTreaty(ctx context.Context, treatyID uuid.UUID) (float64, error)
}
