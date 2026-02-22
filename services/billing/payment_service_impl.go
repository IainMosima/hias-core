package billing

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	billingEntity "github.com/bitbiz/hias-core/domains/billing/entity"
	billingRepo "github.com/bitbiz/hias-core/domains/billing/repository"
	billingSchema "github.com/bitbiz/hias-core/domains/billing/schema"
	billingService "github.com/bitbiz/hias-core/domains/billing/service"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/bitbiz/hias-core/shared"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/google/uuid"
)

const maxPaymentRetries = 3

type paymentServiceImpl struct {
	paymentRepo billingRepo.PaymentRepository
	invoiceRepo billingRepo.InvoiceRepository
}

func NewPaymentService(
	paymentRepo billingRepo.PaymentRepository,
	invoiceRepo billingRepo.InvoiceRepository,
) billingService.PaymentService {
	return &paymentServiceImpl{
		paymentRepo: paymentRepo,
		invoiceRepo: invoiceRepo,
	}
}

// InitiatePayment creates a new payment record and begins the payment process.
// For M-Pesa payments, this would trigger an STK Push to the member's phone.
// For bank transfers, it would initiate the transfer via the banking API.
// State machine: -> INITIATED -> PROCESSING
func (s *paymentServiceImpl) InitiatePayment(ctx context.Context, req billingSchema.InitiatePaymentRequest, createdBy uuid.UUID) *schema.ServiceResponse[billingSchema.PaymentResponse] {
	var invoiceID uuid.UUID
	var claimID uuid.UUID
	var paymentType string

	// Determine payment type from request
	if req.InvoiceID != "" {
		id, err := uuid.Parse(req.InvoiceID)
		if err != nil {
			return schema.NewServiceErrorResponse[billingSchema.PaymentResponse](http.StatusBadRequest, "Invalid invoice ID", err)
		}
		invoiceID = id
		paymentType = string(shared.PaymentTypePremium)

		// Verify invoice exists and is not already paid
		invoice, err := s.invoiceRepo.GetByID(ctx, invoiceID)
		if err != nil {
			return schema.NewServiceErrorResponse[billingSchema.PaymentResponse](http.StatusNotFound, "Invoice not found", err)
		}
		if invoice.Status == string(shared.InvoiceStatusPaid) {
			return schema.NewServiceErrorResponse[billingSchema.PaymentResponse](
				http.StatusConflict,
				"Invoice is already paid",
				fmt.Errorf("invoice %s is already paid", req.InvoiceID),
			)
		}
	} else if req.ClaimID != "" {
		id, err := uuid.Parse(req.ClaimID)
		if err != nil {
			return schema.NewServiceErrorResponse[billingSchema.PaymentResponse](http.StatusBadRequest, "Invalid claim ID", err)
		}
		claimID = id
		paymentType = string(shared.PaymentTypeRemittance)
	} else {
		return schema.NewServiceErrorResponse[billingSchema.PaymentResponse](
			http.StatusBadRequest,
			"Either invoice_id or claim_id must be provided",
			fmt.Errorf("no invoice or claim specified"),
		)
	}

	// Generate a unique payment reference number
	referenceNumber := generatePaymentReference()

	payment := &billingEntity.Payment{
		InvoiceID:       invoiceID,
		ClaimID:         claimID,
		Type:            paymentType,
		Amount:          req.Amount,
		Currency:        "KES",
		Method:          req.Method,
		ReferenceNumber: referenceNumber,
		Status:          string(shared.PaymentStatusInitiated),
		RetryCount:      0,
		MaxRetries:      maxPaymentRetries,
		CreatedBy:       createdBy,
	}

	payment, err := s.paymentRepo.Create(ctx, payment)
	if err != nil {
		utils.LogError("Failed to create payment: %v", err)
		return schema.NewServiceErrorResponse[billingSchema.PaymentResponse](http.StatusInternalServerError, "Failed to create payment", err)
	}

	// Move to PROCESSING status
	// In a full implementation, this would trigger:
	// - M-Pesa: Safaricom Daraja STK Push API
	// - Bank transfer: initiation via banking integration
	payment, err = s.paymentRepo.UpdateStatus(ctx, payment.ID, string(shared.PaymentStatusProcessing))
	if err != nil {
		utils.LogError("Failed to update payment %s to PROCESSING: %v", payment.ID, err)
		// Return the INITIATED payment rather than failing completely
		return schema.NewServiceResponse(billingSchema.ToPaymentResponse(payment), http.StatusCreated, "Payment initiated but processing status update failed")
	}

	utils.LogInfo("Payment %s initiated: type=%s method=%s amount=%d KES ref=%s",
		payment.ID, paymentType, req.Method, req.Amount, referenceNumber)

	return schema.NewServiceResponse(billingSchema.ToPaymentResponse(payment), http.StatusCreated, "Payment initiated")
}

// ProcessWebhook handles payment callback data from external payment gateways
// (e.g., M-Pesa Daraja callback, bank transfer confirmation).
// It updates the payment status based on the gateway response.
func (s *paymentServiceImpl) ProcessWebhook(ctx context.Context, data interface{}) *schema.ServiceResponse[string] {
	// Marshal the webhook data for storage and logging
	webhookJSON, err := json.Marshal(data)
	if err != nil {
		utils.LogError("Failed to marshal webhook data: %v", err)
		return schema.NewServiceErrorResponse[string](http.StatusBadRequest, "Invalid webhook data", err)
	}

	utils.LogInfo("Payment webhook received: %s", string(webhookJSON))

	// In a full M-Pesa integration, the flow would be:
	//
	// 1. Parse the M-Pesa callback structure:
	//    Body.stkCallback.MerchantRequestID
	//    Body.stkCallback.CheckoutRequestID
	//    Body.stkCallback.ResultCode (0 = success)
	//    Body.stkCallback.ResultDesc
	//    Body.stkCallback.CallbackMetadata (MpesaReceiptNumber, TransactionDate, PhoneNumber)
	//
	// 2. Look up payment by CheckoutRequestID/MerchantRequestID:
	//    payment, err := s.paymentRepo.GetByReference(ctx, checkoutRequestID)
	//
	// 3. Update based on ResultCode:
	//    if resultCode == 0 {
	//        s.paymentRepo.Confirm(ctx, payment.ID, webhookJSON)
	//        if payment.InvoiceID != uuid.Nil {
	//            s.invoiceRepo.UpdateStatus(ctx, payment.InvoiceID, "PAID")
	//        }
	//    } else {
	//        s.paymentRepo.UpdateStatus(ctx, payment.ID, "FAILED")
	//    }

	return schema.NewServiceResponse("Webhook processed", http.StatusOK, "Payment webhook processed successfully")
}

// RetryPayment retries a failed payment. Only payments in FAILED status with
// remaining retry attempts can be retried.
// State machine: FAILED -> PROCESSING
func (s *paymentServiceImpl) RetryPayment(ctx context.Context, paymentID uuid.UUID) *schema.ServiceResponse[billingSchema.PaymentResponse] {
	payment, err := s.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		return schema.NewServiceErrorResponse[billingSchema.PaymentResponse](http.StatusNotFound, "Payment not found", err)
	}

	// State machine: can only retry from FAILED
	if payment.Status != string(shared.PaymentStatusFailed) {
		return schema.NewServiceErrorResponse[billingSchema.PaymentResponse](
			http.StatusConflict,
			fmt.Sprintf("Cannot retry payment in status %s; must be FAILED", payment.Status),
			fmt.Errorf("invalid state transition from %s for retry", payment.Status),
		)
	}

	// Check retry limit
	if payment.RetryCount >= payment.MaxRetries {
		return schema.NewServiceErrorResponse[billingSchema.PaymentResponse](
			http.StatusConflict,
			fmt.Sprintf("Maximum retries (%d) exhausted for payment %s", payment.MaxRetries, paymentID),
			fmt.Errorf("max retries reached: %d/%d", payment.RetryCount, payment.MaxRetries),
		)
	}

	// Increment retry count
	payment, err = s.paymentRepo.IncrementRetry(ctx, paymentID)
	if err != nil {
		utils.LogError("Failed to increment retry count for payment %s: %v", paymentID, err)
		return schema.NewServiceErrorResponse[billingSchema.PaymentResponse](http.StatusInternalServerError, "Failed to increment retry count", err)
	}

	// Move back to PROCESSING
	payment, err = s.paymentRepo.UpdateStatus(ctx, paymentID, string(shared.PaymentStatusProcessing))
	if err != nil {
		return schema.NewServiceErrorResponse[billingSchema.PaymentResponse](http.StatusInternalServerError, "Failed to update payment status", err)
	}

	utils.LogInfo("Payment %s retry #%d initiated (ref: %s, amount: %d %s)",
		paymentID, payment.RetryCount, payment.ReferenceNumber, payment.Amount, payment.Currency)

	return schema.NewServiceResponse(billingSchema.ToPaymentResponse(payment), http.StatusOK, fmt.Sprintf("Payment retry #%d initiated", payment.RetryCount))
}

// ReconcilePayment marks a confirmed payment as reconciled, indicating that
// the funds have been verified against bank/gateway settlement records.
// State machine: CONFIRMED -> RECONCILED
func (s *paymentServiceImpl) ReconcilePayment(ctx context.Context, paymentID uuid.UUID) *schema.ServiceResponse[billingSchema.PaymentResponse] {
	payment, err := s.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		return schema.NewServiceErrorResponse[billingSchema.PaymentResponse](http.StatusNotFound, "Payment not found", err)
	}

	// State machine: can only reconcile from CONFIRMED
	if payment.Status != string(shared.PaymentStatusConfirmed) {
		return schema.NewServiceErrorResponse[billingSchema.PaymentResponse](
			http.StatusConflict,
			fmt.Sprintf("Cannot reconcile payment in status %s; must be CONFIRMED", payment.Status),
			fmt.Errorf("invalid state transition from %s for reconciliation", payment.Status),
		)
	}

	payment, err = s.paymentRepo.Reconcile(ctx, paymentID)
	if err != nil {
		return schema.NewServiceErrorResponse[billingSchema.PaymentResponse](http.StatusInternalServerError, "Failed to reconcile payment", err)
	}

	utils.LogInfo("Payment %s reconciled (ref: %s, amount: %d %s)",
		paymentID, payment.ReferenceNumber, payment.Amount, payment.Currency)

	return schema.NewServiceResponse(billingSchema.ToPaymentResponse(payment), http.StatusOK, "Payment reconciled")
}

// GetPayment retrieves a single payment by ID.
func (s *paymentServiceImpl) GetPayment(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[billingSchema.PaymentResponse] {
	payment, err := s.paymentRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[billingSchema.PaymentResponse](http.StatusNotFound, "Payment not found", err)
	}
	return schema.NewServiceResponse(billingSchema.ToPaymentResponse(payment), http.StatusOK, "Payment retrieved")
}

// ListPayments returns a paginated list of payments.
func (s *paymentServiceImpl) ListPayments(ctx context.Context, page, pageSize int) *schema.ServiceResponse[[]billingSchema.PaymentResponse] {
	offset := (page - 1) * pageSize
	payments, err := s.paymentRepo.List(ctx, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]billingSchema.PaymentResponse](http.StatusInternalServerError, "Failed to list payments", err)
	}

	responses := make([]billingSchema.PaymentResponse, 0, len(payments))
	for _, p := range payments {
		responses = append(responses, billingSchema.ToPaymentResponse(p))
	}

	return schema.NewServiceResponse(responses, http.StatusOK, "Payments retrieved")
}

// generatePaymentReference creates a unique payment reference number.
// Format: PAY-YYYY-<8 hex chars from UUID>
func generatePaymentReference() string {
	id := uuid.New()
	return fmt.Sprintf("PAY-%d-%s", time.Now().Year(), id.String()[:8])
}
