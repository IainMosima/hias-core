package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bitbiz/hias-core/domains/policy/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/policy/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
)

type underwritingRepository struct {
	store db.Store
}

func NewUnderwritingRepository(store db.Store) domainRepo.UnderwritingRepository {
	return &underwritingRepository{store: store}
}

func (r *underwritingRepository) Create(ctx context.Context, a *entity.UnderwritingAssessment) (*entity.UnderwritingAssessment, error) {
	medDecl := []byte("{}")
	if a.MedicalDeclarations != nil {
		medDecl = a.MedicalDeclarations
	}
	dbA, err := r.store.CreateUnderwritingAssessment(ctx, db.CreateUnderwritingAssessmentParams{
		PolicyID:            a.PolicyID,
		MemberID:            uuidToPgtype(a.MemberID),
		Status:              a.Status,
		Questionnaire:       a.Questionnaire,
		MedicalDeclarations: medDecl,
		CreatedBy:           a.CreatedBy,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create underwriting assessment: %w", err)
	}
	return sqlcUnderwritingToDomain(dbA, ""), nil
}

func (r *underwritingRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.UnderwritingAssessment, error) {
	row, err := r.store.GetUnderwritingByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get underwriting assessment: %w", err)
	}
	return sqlcUnderwritingToDomain(db.UnderwritingAssessment{
		ID: row.ID, PolicyID: row.PolicyID, MemberID: row.MemberID,
		Status: row.Status, Questionnaire: row.Questionnaire,
		MedicalDeclarations: row.MedicalDeclarations, RiskScore: row.RiskScore,
		RiskFlags: row.RiskFlags, DecisionReason: row.DecisionReason,
		AssessedBy: row.AssessedBy, AssessedAt: row.AssessedAt,
		CreatedBy: row.CreatedBy, CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt,
	}, row.AssessedByName), nil
}

func (r *underwritingRepository) GetByPolicyID(ctx context.Context, policyID uuid.UUID) ([]*entity.UnderwritingAssessment, error) {
	dbAs, err := r.store.ListUnderwritingByPolicy(ctx, policyID)
	if err != nil {
		return nil, fmt.Errorf("failed to list underwriting by policy: %w", err)
	}
	result := make([]*entity.UnderwritingAssessment, len(dbAs))
	for i, row := range dbAs {
		result[i] = sqlcUnderwritingToDomain(db.UnderwritingAssessment{
			ID: row.ID, PolicyID: row.PolicyID, MemberID: row.MemberID,
			Status: row.Status, Questionnaire: row.Questionnaire,
			MedicalDeclarations: row.MedicalDeclarations, RiskScore: row.RiskScore,
			RiskFlags: row.RiskFlags, DecisionReason: row.DecisionReason,
			AssessedBy: row.AssessedBy, AssessedAt: row.AssessedAt,
			CreatedBy: row.CreatedBy, CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt,
		}, row.AssessedByName)
	}
	return result, nil
}

func (r *underwritingRepository) GetByMemberID(ctx context.Context, memberID uuid.UUID) (*entity.UnderwritingAssessment, error) {
	row, err := r.store.GetUnderwritingByMember(ctx, uuidToPgtype(memberID))
	if err != nil {
		return nil, fmt.Errorf("failed to get underwriting by member: %w", err)
	}
	return sqlcUnderwritingToDomain(db.UnderwritingAssessment{
		ID: row.ID, PolicyID: row.PolicyID, MemberID: row.MemberID,
		Status: row.Status, Questionnaire: row.Questionnaire,
		MedicalDeclarations: row.MedicalDeclarations, RiskScore: row.RiskScore,
		RiskFlags: row.RiskFlags, DecisionReason: row.DecisionReason,
		AssessedBy: row.AssessedBy, AssessedAt: row.AssessedAt,
		CreatedBy: row.CreatedBy, CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt,
	}, row.AssessedByName), nil
}

func (r *underwritingRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*entity.UnderwritingAssessment, error) {
	dbA, err := r.store.UpdateUnderwritingStatus(ctx, db.UpdateUnderwritingStatusParams{
		ID:     id,
		Status: status,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update underwriting status: %w", err)
	}
	return sqlcUnderwritingToDomain(dbA, ""), nil
}

func (r *underwritingRepository) Update(ctx context.Context, a *entity.UnderwritingAssessment) (*entity.UnderwritingAssessment, error) {
	var riskFlags []byte
	if a.RiskFlags != nil {
		riskFlags = a.RiskFlags
	}
	dbA, err := r.store.UpdateUnderwriting(ctx, db.UpdateUnderwritingParams{
		ID:             a.ID,
		Status:         stringToPgtypeText(a.Status),
		RiskScore:      intToPgtypeInt4(a.RiskScore),
		RiskFlags:      riskFlags,
		DecisionReason: stringToPgtypeText(a.DecisionReason),
		AssessedBy:     uuidToPgtype(a.AssessedBy),
		AssessedAt:     timePtrToPgtypeTimestamptz(a.AssessedAt),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update underwriting: %w", err)
	}
	return sqlcUnderwritingToDomain(dbA, ""), nil
}

func sqlcUnderwritingToDomain(a db.UnderwritingAssessment, assessedByName string) *entity.UnderwritingAssessment {
	var riskScore int
	if a.RiskScore.Valid {
		riskScore = int(a.RiskScore.Int32)
	}
	var riskFlags json.RawMessage
	if a.RiskFlags != nil {
		riskFlags = a.RiskFlags
	} else {
		riskFlags = json.RawMessage("[]")
	}
	var medDecl json.RawMessage
	if a.MedicalDeclarations != nil {
		medDecl = a.MedicalDeclarations
	} else {
		medDecl = json.RawMessage("{}")
	}
	return &entity.UnderwritingAssessment{
		ID:                  a.ID,
		PolicyID:            a.PolicyID,
		MemberID:            pgtypeToUUID(a.MemberID),
		Status:              a.Status,
		Questionnaire:       a.Questionnaire,
		MedicalDeclarations: medDecl,
		RiskScore:           riskScore,
		RiskFlags:           riskFlags,
		DecisionReason:      a.DecisionReason.String,
		AssessedBy:          pgtypeToUUID(a.AssessedBy),
		AssessedByName:      assessedByName,
		AssessedAt:          pgtypeTimestamptzToTimePtr(a.AssessedAt),
		CreatedBy:           a.CreatedBy,
		CreatedAt:           a.CreatedAt,
		UpdatedAt:           a.UpdatedAt,
	}
}
