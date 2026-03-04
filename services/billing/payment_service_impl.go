package billing

import (
	"context"
	"net/http"

	"github.com/bitbiz/hias-core/domains/billing/entity"
	billingRepo "github.com/bitbiz/hias-core/domains/billing/repository"
	billingSchema "github.com/bitbiz/hias-core/domains/billing/schema"
	"github.com/bitbiz/hias-core/domains/billing/service"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/bitbiz/hias-core/shared"
	"github.com/google/uuid"
)

type paymentServiceImpl struct {
	paymentRepo billingRepo.PaymentRepository
	invoiceRepo billingRepo.InvoiceRepository
}

func NewPaymentService(
	paymentRepo billingRepo.PaymentRepository,
	invoiceRepo billingRepo.InvoiceRepository,
) service.PaymentService {
	return &paymentServiceImpl{
		paymentRepo: paymentRepo,
		invoiceRepo: invoiceRepo,
	}
}

func (s *paymentServiceImpl) InitiatePayment(ctx context.Context, req billingSchema.InitiatePaymentRequest, createdBy uuid.UUID) *schema.ServiceResponse[billingSchema.PaymentResponse] {
	payment := &entity.Payment{
		Type:       string(shared.PaymentTypePremium),
		Amount:     req.Amount,
		Currency:   string(shared.CurrencyKES),
		Method:     req.Method,
		Status:     string(shared.PaymentStatusInitiated),
		MaxRetries: shared.MaxPaymentRetries,
		CreatedBy:  createdBy,
	}

	if req.InvoiceID != "" {
		invoiceID, _ := uuid.Parse(req.InvoiceID)
		payment.InvoiceID = invoiceID
	}

	created, err := s.paymentRepo.Create(ctx, payment)
	if err != nil {
		return schema.NewServiceErrorResponse[billingSchema.PaymentResponse](http.StatusInternalServerError, "Failed to initiate payment", err)
	}

	return schema.NewServiceResponse(billingSchema.ToPaymentResponse(created), http.StatusCreated, "Payment initiated")
}

func (s *paymentServiceImpl) ProcessWebhook(ctx context.Context, data interface{}) *schema.ServiceResponse[string] {
	return schema.NewServiceResponse("processed", http.StatusOK, "Webhook processed")
}

func (s *paymentServiceImpl) RetryPayment(ctx context.Context, paymentID uuid.UUID) *schema.ServiceResponse[billingSchema.PaymentResponse] {
	payment, err := s.paymentRepo.IncrementRetry(ctx, paymentID)
	if err != nil {
		return schema.NewServiceErrorResponse[billingSchema.PaymentResponse](http.StatusInternalServerError, "Failed to retry payment", err)
	}
	return schema.NewServiceResponse(billingSchema.ToPaymentResponse(payment), http.StatusOK, "Payment retried")
}

func (s *paymentServiceImpl) ReconcilePayment(ctx context.Context, paymentID uuid.UUID) *schema.ServiceResponse[billingSchema.PaymentResponse] {
	payment, err := s.paymentRepo.Reconcile(ctx, paymentID)
	if err != nil {
		return schema.NewServiceErrorResponse[billingSchema.PaymentResponse](http.StatusInternalServerError, "Failed to reconcile payment", err)
	}
	return schema.NewServiceResponse(billingSchema.ToPaymentResponse(payment), http.StatusOK, "Payment reconciled")
}

func (s *paymentServiceImpl) GetPayment(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[billingSchema.PaymentResponse] {
	payment, err := s.paymentRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[billingSchema.PaymentResponse](http.StatusNotFound, "Payment not found", err)
	}
	return schema.NewServiceResponse(billingSchema.ToPaymentResponse(payment), http.StatusOK, "Payment retrieved")
}

func (s *paymentServiceImpl) ListPayments(ctx context.Context, page, pageSize int) *schema.ServiceResponse[[]billingSchema.PaymentResponse] {
	offset := (page - 1) * pageSize
	payments, err := s.paymentRepo.List(ctx, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]billingSchema.PaymentResponse](http.StatusInternalServerError, "Failed to list payments", err)
	}

	responses := make([]billingSchema.PaymentResponse, len(payments))
	for i, p := range payments {
		responses[i] = billingSchema.ToPaymentResponse(p)
	}

	return schema.NewServiceResponse(responses, http.StatusOK, "Payments retrieved")
}
