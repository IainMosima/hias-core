package shared

const (
	// Fraud detection
	FraudAmountThresholdCents int64 = 50_000_000 // 500,000 KES

	// Billing defaults
	InvoiceDueDays         = 30
	PreAuthValidityDays    = 30
	MaxPaymentRetries      = 3
	MaxNotificationRetries = 3

	// Benefit defaults
	DefaultMaxAge = 150
	DefaultMinAge = 0

	// Installment counts per frequency
	InstallmentsPerMonth      = 12
	InstallmentsPerQuarter    = 4
	InstallmentsPerSemiAnnual = 2
	InstallmentsPerAnnual     = 1

	// Quotation defaults
	QuotationValidityDays     = 30
	DefaultMaxDiscountPercent = 10
)
