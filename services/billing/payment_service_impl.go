package billing

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/bitbiz/hias-core/domains/billing/entity"
	billingRepo "github.com/bitbiz/hias-core/domains/billing/repository"
	billingSchema "github.com/bitbiz/hias-core/domains/billing/schema"
	"github.com/bitbiz/hias-core/domains/billing/service"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/bitbiz/hias-core/infrastructures/queue"
	"github.com/bitbiz/hias-core/shared"
	"github.com/google/uuid"
)

type paymentServiceImpl struct {
	paymentRepo  billingRepo.PaymentRepository
	invoiceRepo  billingRepo.InvoiceRepository
	queueManager queue.QueueManager
}

func NewPaymentService(
	paymentRepo billingRepo.PaymentRepository,
	invoiceRepo billingRepo.InvoiceRepository,
	queueManager queue.QueueManager,
) service.PaymentService {
	return &paymentServiceImpl{
		paymentRepo:  paymentRepo,
		invoiceRepo:  invoiceRepo,
		queueManager: queueManager,
	}
}

func (s *paymentServiceImpl) InitiatePayment(ctx context.Context, req billingSchema.InitiatePaymentRequest, createdBy uuid.UUID) *schema.ServiceResponse[billingSchema.PaymentResponse] {
	// Double receipting detection: check for duplicate reference
	if req.ReferenceNumber != "" {
		existing, err := s.paymentRepo.GetByReference(ctx, req.ReferenceNumber)
		if err == nil && existing != nil {
			return schema.NewServiceErrorResponse[billingSchema.PaymentResponse](
				http.StatusConflict,
				"Duplicate payment: reference number already exists",
				nil,
			)
		}
	}

	payment := &entity.Payment{
		Type:            string(shared.PaymentTypePremium),
		Amount:          req.Amount,
		Currency:        string(shared.CurrencyKES),
		Method:          req.Method,
		ReferenceNumber: req.ReferenceNumber,
		Status:          string(shared.PaymentStatusInitiated),
		MaxRetries:      shared.MaxPaymentRetries,
		CreatedBy:       createdBy,
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
	payload, ok := data.(map[string]interface{})
	if !ok {
		return schema.NewServiceErrorResponse[string](http.StatusBadRequest, "Invalid webhook payload", nil)
	}

	resultCode, _ := payload["ResultCode"].(float64)
	accountRef, _ := payload["AccountReference"].(string)

	if accountRef == "" {
		return schema.NewServiceErrorResponse[string](http.StatusBadRequest, "Missing AccountReference in webhook", nil)
	}

	payment, err := s.paymentRepo.GetByReference(ctx, accountRef)
	if err != nil {
		return schema.NewServiceErrorResponse[string](http.StatusNotFound, "Payment not found for reference: "+accountRef, err)
	}

	gatewayJSON, _ := json.Marshal(payload)

	if int(resultCode) != 0 {
		s.paymentRepo.UpdateStatus(ctx, payment.ID, string(shared.PaymentStatusFailed))

		go func() {
			if s.queueManager != nil {
				event := map[string]interface{}{
					"event":      "PaymentFailed",
					"payment_id": payment.ID.String(),
					"reference":  accountRef,
					"reason":     fmt.Sprintf("M-Pesa ResultCode: %d", int(resultCode)),
				}
				eventJSON, _ := json.Marshal(event)
				if err := s.queueManager.Publish(context.Background(), queue.TopicPaymentEvents, eventJSON); err != nil {
					log.Printf("Failed to publish PaymentFailedEvent: %v", err)
				}
			}
		}()

		return schema.NewServiceResponse("payment_failed", http.StatusOK, fmt.Sprintf("Payment failed with ResultCode %d", int(resultCode)))
	}

	_, confirmErr := s.paymentRepo.Confirm(ctx, payment.ID, gatewayJSON)
	if confirmErr != nil {
		return schema.NewServiceErrorResponse[string](http.StatusInternalServerError, "Failed to confirm payment", confirmErr)
	}

	if payment.InvoiceID != uuid.Nil {
		s.invoiceRepo.UpdateStatus(ctx, payment.InvoiceID, string(shared.InvoiceStatusPaid))
	}

	go func() {
		if s.queueManager != nil {
			event := map[string]interface{}{
				"event":      "PaymentConfirmed",
				"payment_id": payment.ID.String(),
				"reference":  accountRef,
				"amount":     payment.Amount,
			}
			eventJSON, _ := json.Marshal(event)
			if err := s.queueManager.Publish(context.Background(), queue.TopicPaymentEvents, eventJSON); err != nil {
				log.Printf("Failed to publish PaymentConfirmedEvent: %v", err)
			}
		}
	}()

	return schema.NewServiceResponse("payment_confirmed", http.StatusOK, "Payment confirmed via M-Pesa webhook")
}

func (s *paymentServiceImpl) RetryPayment(ctx context.Context, paymentID uuid.UUID) *schema.ServiceResponse[billingSchema.PaymentResponse] {
	payment, err := s.paymentRepo.IncrementRetry(ctx, paymentID)
	if err != nil {
		return schema.NewServiceErrorResponse[billingSchema.PaymentResponse](http.StatusInternalServerError, "Failed to retry payment", err)
	}
	return schema.NewServiceResponse(billingSchema.ToPaymentResponse(payment), http.StatusOK, "Payment retried")
}

func (s *paymentServiceImpl) ReconcilePayment(ctx context.Context, paymentID uuid.UUID) *schema.ServiceResponse[billingSchema.PaymentResponse] {
	existing, err := s.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		return schema.NewServiceErrorResponse[billingSchema.PaymentResponse](http.StatusNotFound, "Payment not found", err)
	}
	if existing.Status != string(shared.PaymentStatusConfirmed) {
		return schema.NewServiceErrorResponse[billingSchema.PaymentResponse](
			http.StatusBadRequest,
			fmt.Sprintf("Cannot reconcile payment in %s status; must be CONFIRMED", existing.Status),
			nil,
		)
	}

	// Verify amount against invoice if linked
	reconciliationMsg := "Payment reconciled"
	if existing.InvoiceID != uuid.Nil {
		invoice, invoiceErr := s.invoiceRepo.GetByID(ctx, existing.InvoiceID)
		if invoiceErr == nil {
			if existing.Amount < invoice.Amount {
				shortfall := invoice.Amount - existing.Amount
				reconciliationMsg = fmt.Sprintf("Payment reconciled (UNDERPAYMENT: short by %d cents against invoice %s)", shortfall, invoice.InvoiceNumber)
				log.Printf("Reconciliation: underpayment of %d cents on invoice %s", shortfall, invoice.InvoiceNumber)
			} else if existing.Amount > invoice.Amount {
				overage := existing.Amount - invoice.Amount
				reconciliationMsg = fmt.Sprintf("Payment reconciled (OVERPAYMENT: excess of %d cents against invoice %s)", overage, invoice.InvoiceNumber)
				log.Printf("Reconciliation: overpayment of %d cents on invoice %s", overage, invoice.InvoiceNumber)
			}
		}
	}

	payment, err := s.paymentRepo.Reconcile(ctx, paymentID)
	if err != nil {
		return schema.NewServiceErrorResponse[billingSchema.PaymentResponse](http.StatusInternalServerError, "Failed to reconcile payment", err)
	}
	return schema.NewServiceResponse(billingSchema.ToPaymentResponse(payment), http.StatusOK, reconciliationMsg)
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
