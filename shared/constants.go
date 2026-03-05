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

	// Underwriting thresholds
	UnderwritingAutoApproveThreshold = 30 // risk score < 30 → auto-approve
	UnderwritingReferThreshold       = 60 // 30-60 → refer, > 60 → decline

	// Claims SLA
	ClaimSLAHours = 48 // SLA breach threshold in hours

	// Number format prefixes
	CaseNumberPrefix      = "CASE"
	StatementNumberPrefix = "STMT"

	// Reinsurance number format prefixes
	TreatyNumberPrefix             = "TRY"
	CessionNumberPrefix            = "CES"
	RecoveryNumberPrefix           = "REC"
	BordereauNumberPrefix          = "BDX"
	ReinsurerStatementNumberPrefix = "RST"

	// Reinsurance thresholds
	CatastropheThresholdCents int64 = 500_000_000 // 5,000,000 KES
	AggregateWarningPercent         = 80
)
