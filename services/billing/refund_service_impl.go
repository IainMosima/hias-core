package billing

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bitbiz/hias-core/domains/billing/entity"
	billingRepo "github.com/bitbiz/hias-core/domains/billing/repository"
	billingSchema "github.com/bitbiz/hias-core/domains/billing/schema"
	"github.com/bitbiz/hias-core/domains/billing/service"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/bitbiz/hias-core/shared"
	"github.com/google/uuid"
)

type refundServiceImpl struct {
	refundRepo billingRepo.RefundRepository
}

func NewRefundService(refundRepo billingRepo.RefundRepository) service.RefundService {
	return &refundServiceImpl{refundRepo: refundRepo}
}

func (s *refundServiceImpl) RequestRefund(ctx context.Context, req billingSchema.CreateRefundRequest, createdBy uuid.UUID) *schema.ServiceResponse[billingSchema.RefundResponse] {
	policyID, _ := uuid.Parse(req.PolicyID)

	refund := &entity.Refund{
		PolicyID:  policyID,
		Amount:    req.Amount,
		Currency:  string(shared.CurrencyKES),
		Status:    string(shared.RefundStatusPending),
		Reason:    req.Reason,
		CreatedBy: createdBy,
	}

	if req.CreditNoteID != "" {
		creditNoteID, _ := uuid.Parse(req.CreditNoteID)
		refund.CreditNoteID = creditNoteID
	}

	created, err := s.refundRepo.Create(ctx, refund)
	if err != nil {
		return schema.NewServiceErrorResponse[billingSchema.RefundResponse](http.StatusInternalServerError, "Failed to create refund", err)
	}
	return schema.NewServiceResponse(billingSchema.ToRefundResponse(created), http.StatusCreated, "Refund requested")
}

func (s *refundServiceImpl) ApproveRefund(ctx context.Context, id uuid.UUID, approvedBy uuid.UUID) *schema.ServiceResponse[billingSchema.RefundResponse] {
	refund, err := s.refundRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[billingSchema.RefundResponse](http.StatusNotFound, "Refund not found", err)
	}
	if refund.Status != string(shared.RefundStatusPending) {
		return schema.NewServiceErrorResponse[billingSchema.RefundResponse](
			http.StatusBadRequest,
			fmt.Sprintf("Cannot approve refund in %s status", refund.Status),
			nil,
		)
	}

	approved, err := s.refundRepo.Approve(ctx, id, approvedBy)
	if err != nil {
		return schema.NewServiceErrorResponse[billingSchema.RefundResponse](http.StatusInternalServerError, "Failed to approve refund", err)
	}
	return schema.NewServiceResponse(billingSchema.ToRefundResponse(approved), http.StatusOK, "Refund approved")
}

func (s *refundServiceImpl) ProcessRefund(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[billingSchema.RefundResponse] {
	refund, err := s.refundRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[billingSchema.RefundResponse](http.StatusNotFound, "Refund not found", err)
	}
	if refund.Status != string(shared.RefundStatusApproved) {
		return schema.NewServiceErrorResponse[billingSchema.RefundResponse](
			http.StatusBadRequest,
			fmt.Sprintf("Cannot process refund in %s status; must be APPROVED", refund.Status),
			nil,
		)
	}

	processed, err := s.refundRepo.Process(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[billingSchema.RefundResponse](http.StatusInternalServerError, "Failed to process refund", err)
	}
	return schema.NewServiceResponse(billingSchema.ToRefundResponse(processed), http.StatusOK, "Refund processed")
}

func (s *refundServiceImpl) ListRefunds(ctx context.Context, policyID uuid.UUID, page, pageSize int) *schema.ServiceResponse[[]billingSchema.RefundResponse] {
	offset := (page - 1) * pageSize
	refunds, err := s.refundRepo.ListByPolicy(ctx, policyID, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]billingSchema.RefundResponse](http.StatusInternalServerError, "Failed to list refunds", err)
	}
	responses := make([]billingSchema.RefundResponse, len(refunds))
	for i, r := range refunds {
		responses[i] = billingSchema.ToRefundResponse(r)
	}
	return schema.NewServiceResponse(responses, http.StatusOK, "Refunds retrieved")
}
