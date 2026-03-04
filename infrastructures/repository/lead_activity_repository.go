package repository

import (
	"context"
	"fmt"

	"github.com/bitbiz/hias-core/domains/sales/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/sales/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
)

type leadActivityRepository struct {
	store db.Store
}

func NewLeadActivityRepository(store db.Store) domainRepo.LeadActivityRepository {
	return &leadActivityRepository{store: store}
}

func (r *leadActivityRepository) Create(ctx context.Context, activity *entity.LeadActivity) (*entity.LeadActivity, error) {
	dbActivity, err := r.store.CreateLeadActivity(ctx, db.CreateLeadActivityParams{
		LeadID:       activity.LeadID,
		ActivityType: activity.ActivityType,
		Description:  stringToPgtypeText(activity.Description),
		ScheduledAt:  timePtrToPgtypeTimestamptz(activity.ScheduledAt),
		CompletedAt:  timePtrToPgtypeTimestamptz(activity.CompletedAt),
		CreatedBy:    uuidToPgtype(activity.CreatedBy),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create lead activity: %w", err)
	}
	return sqlcLeadActivityToDomain(dbActivity), nil
}

func (r *leadActivityRepository) ListByLead(ctx context.Context, leadID uuid.UUID) ([]*entity.LeadActivity, error) {
	dbActivities, err := r.store.ListLeadActivitiesByLead(ctx, leadID)
	if err != nil {
		return nil, fmt.Errorf("failed to list lead activities: %w", err)
	}
	activities := make([]*entity.LeadActivity, len(dbActivities))
	for i, a := range dbActivities {
		activities[i] = sqlcLeadActivityToDomain(a)
	}
	return activities, nil
}

func sqlcLeadActivityToDomain(a db.LeadActivity) *entity.LeadActivity {
	return &entity.LeadActivity{
		ID:           a.ID,
		LeadID:       a.LeadID,
		ActivityType: a.ActivityType,
		Description:  a.Description.String,
		ScheduledAt:  pgtypeTimestamptzToTimePtr(a.ScheduledAt),
		CompletedAt:  pgtypeTimestamptzToTimePtr(a.CompletedAt),
		CreatedBy:    pgtypeToUUID(a.CreatedBy),
		CreatedAt:    a.CreatedAt,
	}
}
