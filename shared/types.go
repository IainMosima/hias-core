package shared

// PolicyStatus represents the status of a policy
type PolicyStatus string

const (
	PolicyStatusDraft      PolicyStatus = "DRAFT"
	PolicyStatusActive     PolicyStatus = "ACTIVE"
	PolicyStatusLapsed     PolicyStatus = "LAPSED"
	PolicyStatusTerminated PolicyStatus = "TERMINATED"
)

// ProviderStatus represents the status of a provider
type ProviderStatus string

const (
	ProviderStatusPending       ProviderStatus = "PENDING"
	ProviderStatusCredentialing ProviderStatus = "CREDENTIALING"
	ProviderStatusActive        ProviderStatus = "ACTIVE"
	ProviderStatusSuspended     ProviderStatus = "SUSPENDED"
	ProviderStatusTerminated    ProviderStatus = "TERMINATED"
)

// ClaimStatus represents the status of a claim
type ClaimStatus string

const (
	ClaimStatusReceived     ClaimStatus = "RECEIVED"
	ClaimStatusValidated    ClaimStatus = "VALIDATED"
	ClaimStatusAdjudicated  ClaimStatus = "ADJUDICATED"
	ClaimStatusApproved     ClaimStatus = "APPROVED"
	ClaimStatusRejected     ClaimStatus = "REJECTED"
	ClaimStatusManualReview ClaimStatus = "MANUAL_REVIEW"
	ClaimStatusPaid         ClaimStatus = "PAID"
)

// PreAuthStatus represents the status of a pre-authorization
type PreAuthStatus string

const (
	PreAuthStatusSubmitted    PreAuthStatus = "SUBMITTED"
	PreAuthStatusUnderReview  PreAuthStatus = "UNDER_REVIEW"
	PreAuthStatusApproved     PreAuthStatus = "APPROVED"
	PreAuthStatusDenied       PreAuthStatus = "DENIED"
	PreAuthStatusInfoRequested PreAuthStatus = "INFO_REQUESTED"
	PreAuthStatusExpired      PreAuthStatus = "EXPIRED"
	PreAuthStatusClaimed      PreAuthStatus = "CLAIMED"
)

// PaymentStatus represents the status of a payment
type PaymentStatus string

const (
	PaymentStatusInitiated  PaymentStatus = "INITIATED"
	PaymentStatusProcessing PaymentStatus = "PROCESSING"
	PaymentStatusConfirmed  PaymentStatus = "CONFIRMED"
	PaymentStatusFailed     PaymentStatus = "FAILED"
	PaymentStatusReconciled PaymentStatus = "RECONCILED"
	PaymentStatusCancelled  PaymentStatus = "CANCELLED"
)

// PaymentMethod represents how a payment is made
type PaymentMethod string

const (
	PaymentMethodMpesa        PaymentMethod = "MPESA"
	PaymentMethodBankTransfer PaymentMethod = "BANK_TRANSFER"
)

// PaymentType represents the type of payment
type PaymentType string

const (
	PaymentTypePremium    PaymentType = "PREMIUM"
	PaymentTypeRemittance PaymentType = "REMITTANCE"
)

// NotificationChannel represents notification delivery channels
type NotificationChannel string

const (
	NotificationChannelSMS   NotificationChannel = "SMS"
	NotificationChannelEmail NotificationChannel = "EMAIL"
	NotificationChannelInApp NotificationChannel = "IN_APP"
	NotificationChannelPush  NotificationChannel = "PUSH"
)

// NotificationStatus represents the status of a notification
type NotificationStatus string

const (
	NotificationStatusPending   NotificationStatus = "PENDING"
	NotificationStatusSent      NotificationStatus = "SENT"
	NotificationStatusDelivered NotificationStatus = "DELIVERED"
	NotificationStatusFailed    NotificationStatus = "FAILED"
	NotificationStatusRead      NotificationStatus = "READ"
)

// BenefitCategory represents benefit categories
type BenefitCategory string

const (
	BenefitCategoryOutpatient BenefitCategory = "outpatient"
	BenefitCategoryInpatient  BenefitCategory = "inpatient"
	BenefitCategoryDental     BenefitCategory = "dental"
	BenefitCategoryOptical    BenefitCategory = "optical"
	BenefitCategoryMaternity  BenefitCategory = "maternity"
)

// CoPayType represents co-pay calculation types
type CoPayType string

const (
	CoPayTypePercentage CoPayType = "percentage"
	CoPayTypeFixed      CoPayType = "fixed"
)

// UserStatus represents the status of a user
type UserStatus string

const (
	UserStatusActive    UserStatus = "ACTIVE"
	UserStatusInactive  UserStatus = "INACTIVE"
	UserStatusSuspended UserStatus = "SUSPENDED"
)

// UserRole names
type UserRole string

const (
	UserRoleAdmin         UserRole = "Admin"
	UserRoleUnderwriter   UserRole = "Underwriter"
	UserRoleClaimsOfficer UserRole = "ClaimsOfficer"
	UserRoleFinance       UserRole = "Finance"
	UserRoleProvider      UserRole = "Provider"
	UserRoleMember        UserRole = "Member"
)

// ProviderType represents the type of healthcare provider
type ProviderType string

const (
	ProviderTypeHospital ProviderType = "hospital"
	ProviderTypeClinic   ProviderType = "clinic"
	ProviderTypePharmacy ProviderType = "pharmacy"
	ProviderTypeLab      ProviderType = "lab"
)

// PlanType represents insurance plan types
type PlanType string

const (
	PlanTypeIndividual PlanType = "individual"
	PlanTypeGroup      PlanType = "group"
)

// MemberRelationship represents member relationship to policyholder
type MemberRelationship string

const (
	MemberRelationshipPrincipal MemberRelationship = "principal"
	MemberRelationshipSpouse    MemberRelationship = "spouse"
	MemberRelationshipChild     MemberRelationship = "child"
	MemberRelationshipParent    MemberRelationship = "parent"
)

// AuditAction represents audit event actions
type AuditAction string

const (
	AuditActionCreate      AuditAction = "CREATE"
	AuditActionUpdate      AuditAction = "UPDATE"
	AuditActionDelete      AuditAction = "DELETE"
	AuditActionStateChange AuditAction = "STATE_CHANGE"
)

// AdjudicationDecision represents the adjudication engine decision
type AdjudicationDecision string

const (
	AdjudicationDecisionApprove      AdjudicationDecision = "APPROVE"
	AdjudicationDecisionReject       AdjudicationDecision = "REJECT"
	AdjudicationDecisionManualReview AdjudicationDecision = "MANUAL_REVIEW"
)

// RuleCategory represents adjudication rule categories
type RuleCategory string

const (
	RuleCategoryEligibility RuleCategory = "eligibility"
	RuleCategoryCoverage    RuleCategory = "coverage"
	RuleCategoryLimits      RuleCategory = "limits"
	RuleCategoryFraud       RuleCategory = "fraud"
)

// RuleResult represents the result of an adjudication rule
type RuleResultStatus string

const (
	RuleResultPass RuleResultStatus = "PASS"
	RuleResultFail RuleResultStatus = "FAIL"
	RuleResultFlag RuleResultStatus = "FLAG"
)

// FraudFlagType represents types of fraud flags
type FraudFlagType string

const (
	FraudFlagDuplicate       FraudFlagType = "DUPLICATE"
	FraudFlagFrequency       FraudFlagType = "FREQUENCY"
	FraudFlagAmountThreshold FraudFlagType = "AMOUNT_THRESHOLD"
)

// FraudSeverity represents severity levels
type FraudSeverity string

const (
	FraudSeverityLow      FraudSeverity = "LOW"
	FraudSeverityMedium   FraudSeverity = "MEDIUM"
	FraudSeverityHigh     FraudSeverity = "HIGH"
	FraudSeverityCritical FraudSeverity = "CRITICAL"
)

// InvoiceStatus represents the status of an invoice
type InvoiceStatus string

const (
	InvoiceStatusPending   InvoiceStatus = "PENDING"
	InvoiceStatusPaid      InvoiceStatus = "PAID"
	InvoiceStatusOverdue   InvoiceStatus = "OVERDUE"
	InvoiceStatusCancelled InvoiceStatus = "CANCELLED"
)

// RemittanceStatus represents the status of a remittance
type RemittanceStatus string

const (
	RemittanceStatusPending    RemittanceStatus = "PENDING"
	RemittanceStatusProcessing RemittanceStatus = "PROCESSING"
	RemittanceStatusSent       RemittanceStatus = "SENT"
	RemittanceStatusConfirmed  RemittanceStatus = "CONFIRMED"
	RemittanceStatusFailed     RemittanceStatus = "FAILED"
)

// ContractStatus represents the status of a provider contract
type ContractStatus string

const (
	ContractStatusActive     ContractStatus = "ACTIVE"
	ContractStatusExpired    ContractStatus = "EXPIRED"
	ContractStatusTerminated ContractStatus = "TERMINATED"
)

// Gender represents gender values
type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
	GenderOther  Gender = "other"
)
