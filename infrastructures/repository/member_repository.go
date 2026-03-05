package repository

import (
	"context"
	"fmt"

	"github.com/bitbiz/hias-core/domains/policy/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/policy/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
)

type memberRepository struct {
	store db.Store
}

func NewMemberRepository(store db.Store) domainRepo.MemberRepository {
	return &memberRepository{store: store}
}

func (r *memberRepository) Create(ctx context.Context, member *entity.Member) (*entity.Member, error) {
	dbMember, err := r.store.CreateMember(ctx, db.CreateMemberParams{
		PolicyID:     member.PolicyID,
		NationalID:   stringToPgtypeText(member.NationalID),
		Name:         member.Name,
		DateOfBirth:  timeToPgtypeDate(member.DateOfBirth),
		Gender:       member.Gender,
		Relationship: member.Relationship,
		MemberNumber: member.MemberNumber,
		Phone:        stringToPgtypeText(member.Phone),
		Email:        stringToPgtypeText(member.Email),
		KraPin:       stringToPgtypeText(member.KRAPin),
		County:       stringToPgtypeText(member.County),
		Address:      stringToPgtypeText(member.Address),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create member: %w", err)
	}
	return sqlcMemberToDomain(dbMember), nil
}

func (r *memberRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Member, error) {
	dbMember, err := r.store.GetMemberByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get member by ID: %w", err)
	}
	return sqlcMemberToDomain(dbMember), nil
}

func (r *memberRepository) GetByNumber(ctx context.Context, number string) (*entity.Member, error) {
	dbMember, err := r.store.GetMemberByNumber(ctx, number)
	if err != nil {
		return nil, fmt.Errorf("failed to get member by number: %w", err)
	}
	return sqlcMemberToDomain(dbMember), nil
}

func (r *memberRepository) GetByNationalID(ctx context.Context, nationalID string) (*entity.Member, error) {
	dbMember, err := r.store.GetMemberByNationalID(ctx, stringToPgtypeText(nationalID))
	if err != nil {
		return nil, fmt.Errorf("failed to get member by national ID: %w", err)
	}
	return sqlcMemberToDomain(dbMember), nil
}

func (r *memberRepository) ListByPolicy(ctx context.Context, policyID uuid.UUID) ([]*entity.Member, error) {
	dbMembers, err := r.store.ListMembersByPolicy(ctx, policyID)
	if err != nil {
		return nil, fmt.Errorf("failed to list members by policy: %w", err)
	}
	members := make([]*entity.Member, len(dbMembers))
	for i, m := range dbMembers {
		members[i] = sqlcMemberToDomain(m)
	}
	return members, nil
}

func (r *memberRepository) CountByPolicy(ctx context.Context, policyID uuid.UUID) (int64, error) {
	count, err := r.store.CountMembersByPolicy(ctx, policyID)
	if err != nil {
		return 0, fmt.Errorf("failed to count members by policy: %w", err)
	}
	return count, nil
}

func (r *memberRepository) Verify(ctx context.Context, id uuid.UUID) (*entity.Member, error) {
	dbMember, err := r.store.VerifyMember(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to verify member: %w", err)
	}
	return sqlcMemberToDomain(dbMember), nil
}

func (r *memberRepository) Update(ctx context.Context, member *entity.Member) (*entity.Member, error) {
	dbMember, err := r.store.UpdateMember(ctx, db.UpdateMemberParams{
		ID:      member.ID,
		Name:    stringToPgtypeText(member.Name),
		Phone:   stringToPgtypeText(member.Phone),
		Email:   stringToPgtypeText(member.Email),
		KraPin:  stringToPgtypeText(member.KRAPin),
		County:  stringToPgtypeText(member.County),
		Address: stringToPgtypeText(member.Address),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update member: %w", err)
	}
	return sqlcMemberToDomain(dbMember), nil
}

func (r *memberRepository) ListActiveByPolicy(ctx context.Context, policyID uuid.UUID) ([]*entity.Member, error) {
	dbMembers, err := r.store.ListActiveMembersByPolicy(ctx, policyID)
	if err != nil {
		return nil, fmt.Errorf("failed to list active members by policy: %w", err)
	}
	members := make([]*entity.Member, len(dbMembers))
	for i, m := range dbMembers {
		members[i] = sqlcMemberToDomain(m)
	}
	return members, nil
}

func (r *memberRepository) CountActiveByPolicy(ctx context.Context, policyID uuid.UUID) (int64, error) {
	count, err := r.store.CountActiveMembersByPolicy(ctx, policyID)
	if err != nil {
		return 0, fmt.Errorf("failed to count active members: %w", err)
	}
	return count, nil
}

func (r *memberRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*entity.Member, error) {
	dbMember, err := r.store.UpdateMemberStatus(ctx, db.UpdateMemberStatusParams{
		ID:     id,
		Status: status,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update member status: %w", err)
	}
	return sqlcMemberToDomain(dbMember), nil
}

func (r *memberRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.store.DeleteMember(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete member: %w", err)
	}
	return nil
}

func sqlcMemberToDomain(m db.Member) *entity.Member {
	return &entity.Member{
		ID:           m.ID,
		PolicyID:     m.PolicyID,
		NationalID:   m.NationalID.String,
		Name:         m.Name,
		DateOfBirth:  pgtypeDateToTime(m.DateOfBirth),
		Gender:       m.Gender,
		Relationship: m.Relationship,
		MemberNumber: m.MemberNumber,
		Phone:        m.Phone.String,
		Email:        m.Email.String,
		KRAPin:       m.KraPin.String,
		County:       m.County.String,
		Address:      m.Address.String,
		Status:       m.Status,
		Verified:     m.Verified,
		VerifiedAt:   pgtypeTimestamptzToTimePtr(m.VerifiedAt),
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}
