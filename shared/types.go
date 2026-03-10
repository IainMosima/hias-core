package shared

// PolicyStatus represents the status of a policy
type PolicyStatus string

const (
	PolicyStatusDraft      PolicyStatus = "DRAFT"
	PolicyStatusActive     PolicyStatus = "ACTIVE"
	PolicyStatusLapsed     PolicyStatus = "LAPSED"
	PolicyStatusTerminated PolicyStatus = "TERMINATED"
	PolicyStatusSuspended  PolicyStatus = "SUSPENDED"
)

// MemberStatus represents the status of a policy member
type MemberStatus string

const (
	MemberStatusActive    MemberStatus = "ACTIVE"
	MemberStatusSuspended MemberStatus = "SUSPENDED"
	MemberStatusRemoved   MemberStatus = "REMOVED"
)

// EndorsementType represents the type of policy endorsement
type EndorsementType string

const (
	EndorsementTypeAddMember    EndorsementType = "ADD_MEMBER"
	EndorsementTypeRemoveMember EndorsementType = "REMOVE_MEMBER"
	EndorsementTypeUpdateMember EndorsementType = "UPDATE_MEMBER"
	EndorsementTypePlanChange   EndorsementType = "PLAN_CHANGE"
)

// EndorsementStatus represents the status of an endorsement
type EndorsementStatus string

const (
	EndorsementStatusPending  EndorsementStatus = "PENDING"
	EndorsementStatusApproved EndorsementStatus = "APPROVED"
	EndorsementStatusRejected EndorsementStatus = "REJECTED"
	EndorsementStatusApplied  EndorsementStatus = "APPLIED"
)

// RenewalStatus represents the status of a policy renewal
type RenewalStatus string

const (
	RenewalStatusPending   RenewalStatus = "PENDING"
	RenewalStatusApproved  RenewalStatus = "APPROVED"
	RenewalStatusRejected  RenewalStatus = "REJECTED"
	RenewalStatusCompleted RenewalStatus = "COMPLETED"
	RenewalStatusExpired   RenewalStatus = "EXPIRED"
)

// UnderwritingStatus represents the status of an underwriting assessment
type UnderwritingStatus string

const (
	UnderwritingStatusPending  UnderwritingStatus = "PENDING"
	UnderwritingStatusApproved UnderwritingStatus = "APPROVED"
	UnderwritingStatusDeclined UnderwritingStatus = "DECLINED"
	UnderwritingStatusRefer    UnderwritingStatus = "REFER"
)

// PolicyDocumentType represents the type of generated policy document
type PolicyDocumentType string

const (
	PolicyDocumentTypeWelcomeLetter  PolicyDocumentType = "WELCOME_LETTER"
	PolicyDocumentTypeMemberCard     PolicyDocumentType = "MEMBER_CARD"
	PolicyDocumentTypePolicySchedule PolicyDocumentType = "POLICY_SCHEDULE"
	PolicyDocumentTypeRenewalNotice  PolicyDocumentType = "RENEWAL_NOTICE"
	PolicyDocumentTypeEndorsement    PolicyDocumentType = "ENDORSEMENT"
	PolicyDocumentTypeLOU            PolicyDocumentType = "LOU"
	PolicyDocumentTypeDeclineLetter  PolicyDocumentType = "DECLINE_LETTER"
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
	ClaimStatusReceived        ClaimStatus = "RECEIVED"
	ClaimStatusValidated       ClaimStatus = "VALIDATED"
	ClaimStatusAdjudicated     ClaimStatus = "ADJUDICATED"
	ClaimStatusApproved        ClaimStatus = "APPROVED"
	ClaimStatusRejected        ClaimStatus = "REJECTED"
	ClaimStatusManualReview    ClaimStatus = "MANUAL_REVIEW"
	ClaimStatusPaid            ClaimStatus = "PAID"
	ClaimStatusVetted          ClaimStatus = "VETTED"
	ClaimStatusPartiallyVetted ClaimStatus = "PARTIALLY_VETTED"
	ClaimStatusReadyForPayment ClaimStatus = "READY_FOR_PAYMENT"
	ClaimStatusPartPaid        ClaimStatus = "PART_PAID"
	ClaimStatusEscalated       ClaimStatus = "ESCALATED"
)

// AdjudicationRuleType represents configurable adjudication rule types
type AdjudicationRuleType string

const (
	AdjudicationRuleTypeAmountThreshold AdjudicationRuleType = "AMOUNT_THRESHOLD"
	AdjudicationRuleTypeFrequencyLimit  AdjudicationRuleType = "FREQUENCY_LIMIT"
	AdjudicationRuleTypeBenefitCheck    AdjudicationRuleType = "BENEFIT_CHECK"
	AdjudicationRuleTypeAutoApprove     AdjudicationRuleType = "AUTO_APPROVE"
)

// EscalationConditionType represents escalation rule condition types
type EscalationConditionType string

const (
	EscalationConditionAmountExceeds EscalationConditionType = "AMOUNT_EXCEEDS"
	EscalationConditionFraudFlag     EscalationConditionType = "FRAUD_FLAG"
	EscalationConditionManualReview  EscalationConditionType = "MANUAL_REVIEW"
)

// CommissionPaymentStatus represents commission payment statuses
type CommissionPaymentStatus string

const (
	CommissionPaymentStatusPending   CommissionPaymentStatus = "PENDING"
	CommissionPaymentStatusProcessed CommissionPaymentStatus = "PROCESSED"
	CommissionPaymentStatusFailed    CommissionPaymentStatus = "FAILED"
)

// RefundStatus represents the status of a refund
type RefundStatus string

const (
	RefundStatusPending   RefundStatus = "PENDING"
	RefundStatusApproved  RefundStatus = "APPROVED"
	RefundStatusProcessed RefundStatus = "PROCESSED"
	RefundStatusFailed    RefundStatus = "FAILED"
)

// PremiumLedgerEntryType represents ledger entry types
type PremiumLedgerEntryType string

const (
	PremiumLedgerEntryTypeDebit  PremiumLedgerEntryType = "DEBIT"
	PremiumLedgerEntryTypeCredit PremiumLedgerEntryType = "CREDIT"
)

// ClaimType represents the type of claim
type ClaimType string

const (
	ClaimTypeDirect        ClaimType = "DIRECT"
	ClaimTypeReimbursement ClaimType = "REIMBURSEMENT"
	ClaimTypeCredit        ClaimType = "CREDIT"
	ClaimTypeException     ClaimType = "EXCEPTION"
)

// CaseStatus represents the status of an inpatient case
type CaseStatus string

const (
	CaseStatusScheduled   CaseStatus = "SCHEDULED"
	CaseStatusAdmitted    CaseStatus = "ADMITTED"
	CaseStatusInTreatment CaseStatus = "IN_TREATMENT"
	CaseStatusDischarged  CaseStatus = "DISCHARGED"
	CaseStatusClosed      CaseStatus = "CLOSED"
)

// ProviderTier represents provider tiers
type ProviderTier string

const (
	ProviderTierOne   ProviderTier = "TIER_1"
	ProviderTierTwo   ProviderTier = "TIER_2"
	ProviderTierThree ProviderTier = "TIER_3"
)

// StatementStatus represents the status of a provider statement
type StatementStatus string

const (
	StatementStatusUploaded   StatementStatus = "UPLOADED"
	StatementStatusReconciled StatementStatus = "RECONCILED"
)

// MatchStatus represents statement line item match status
type MatchStatus string

const (
	MatchStatusUnmatched MatchStatus = "UNMATCHED"
	MatchStatusMatched   MatchStatus = "MATCHED"
	MatchStatusDisputed  MatchStatus = "DISPUTED"
)

// PreAuthStatus represents the status of a pre-authorization
type PreAuthStatus string

const (
	PreAuthStatusSubmitted     PreAuthStatus = "SUBMITTED"
	PreAuthStatusUnderReview   PreAuthStatus = "UNDER_REVIEW"
	PreAuthStatusApproved      PreAuthStatus = "APPROVED"
	PreAuthStatusDenied        PreAuthStatus = "DENIED"
	PreAuthStatusInfoRequested PreAuthStatus = "INFO_REQUESTED"
	PreAuthStatusExpired       PreAuthStatus = "EXPIRED"
	PreAuthStatusClaimed       PreAuthStatus = "CLAIMED"
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

// NotificationType represents the type/category of a notification
type NotificationType string

const (
	NotificationTypeQuotation NotificationType = "QUOTATION"
	NotificationTypeApproval  NotificationType = "APPROVAL"
	NotificationTypeClaim     NotificationType = "CLAIM"
	NotificationTypePolicy    NotificationType = "POLICY"
	NotificationTypeDocument  NotificationType = "DOCUMENT"
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
	BenefitCategoryPharmacy   BenefitCategory = "pharmacy"
	BenefitCategorySpecialist BenefitCategory = "specialist"
	BenefitCategoryEmergency  BenefitCategory = "emergency"
	BenefitCategoryChronic    BenefitCategory = "chronic"
	BenefitCategoryWellness   BenefitCategory = "wellness"
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
	UserRoleSalesAgent    UserRole = "SalesAgent"
	UserRoleManager       UserRole = "Manager"
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

// AuditEntityType represents the type of entity being audited
type AuditEntityType string

const (
	AuditEntityTypeClaim             AuditEntityType = "CLAIM"
	AuditEntityTypePolicy            AuditEntityType = "POLICY"
	AuditEntityTypeMember            AuditEntityType = "MEMBER"
	AuditEntityTypePlan              AuditEntityType = "PLAN"
	AuditEntityTypeBenefit           AuditEntityType = "BENEFIT"
	AuditEntityTypeExclusion         AuditEntityType = "EXCLUSION"
	AuditEntityTypePremiumRule       AuditEntityType = "PREMIUM_RULE"
	AuditEntityTypeProviderNetwork   AuditEntityType = "PROVIDER_NETWORK"
	AuditEntityTypeProvider          AuditEntityType = "PROVIDER"
	AuditEntityTypeUser              AuditEntityType = "USER"
	AuditEntityTypeLead              AuditEntityType = "LEAD"
	AuditEntityTypeQuotation         AuditEntityType = "QUOTATION"
	AuditEntityTypeQuotationVersion  AuditEntityType = "QUOTATION_VERSION"
	AuditEntityTypeQuotationDocument AuditEntityType = "QUOTATION_DOCUMENT"
	AuditEntityTypeApprovalLimit     AuditEntityType = "APPROVAL_LIMIT"
	AuditEntityTypeEndorsement       AuditEntityType = "ENDORSEMENT"
	AuditEntityTypeRenewal           AuditEntityType = "RENEWAL"
	AuditEntityTypeUnderwriting      AuditEntityType = "UNDERWRITING"
	AuditEntityTypePolicyDocument    AuditEntityType = "POLICY_DOCUMENT"
	AuditEntityTypeUnderwritingFlag  AuditEntityType = "UNDERWRITING_FLAG"
	AuditEntityTypeUnderwritingRule  AuditEntityType = "UNDERWRITING_RULE"
	AuditEntityTypeCreditNote        AuditEntityType = "CREDIT_NOTE"
	AuditEntityTypeCaseRecord        AuditEntityType = "CASE_RECORD"
	AuditEntityTypeClaimDocument     AuditEntityType = "CLAIM_DOCUMENT"
	AuditEntityTypeProviderStatement AuditEntityType = "PROVIDER_STATEMENT"
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
	FraudFlagDuplicate          FraudFlagType = "DUPLICATE"
	FraudFlagFrequency          FraudFlagType = "FREQUENCY"
	FraudFlagAmountThreshold    FraudFlagType = "AMOUNT_THRESHOLD"
	FraudFlagExpiredContract    FraudFlagType = "EXPIRED_CONTRACT"
	FraudFlagSuspendedProvider  FraudFlagType = "SUSPENDED_PROVIDER"
	FraudFlagRepeatVisit        FraudFlagType = "REPEAT_VISIT"
	FraudFlagRateCardOvercharge FraudFlagType = "RATE_CARD_OVERCHARGE"
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

// PlanStatus represents the status of a plan
type PlanStatus string

const (
	PlanStatusActive   PlanStatus = "ACTIVE"
	PlanStatusInactive PlanStatus = "INACTIVE"
)

// PlanSegment represents the market segment for a plan
type PlanSegment string

const (
	PlanSegmentRetail    PlanSegment = "retail"
	PlanSegmentCorporate PlanSegment = "corporate"
	PlanSegmentSME       PlanSegment = "sme"
)

// SubLimitType represents how sub-limits are applied on a benefit
type SubLimitType string

const (
	SubLimitTypeNone     SubLimitType = "none"
	SubLimitTypePerVisit SubLimitType = "per_visit"
	SubLimitTypePerItem  SubLimitType = "per_item"
)

// WaitingPeriodType represents the type of waiting period
type WaitingPeriodType string

const (
	WaitingPeriodTypeGeneral     WaitingPeriodType = "general"
	WaitingPeriodTypeMaternity   WaitingPeriodType = "maternity"
	WaitingPeriodTypePreExisting WaitingPeriodType = "pre_existing"
	WaitingPeriodTypeChronic     WaitingPeriodType = "chronic"
	WaitingPeriodTypeSurgical    WaitingPeriodType = "surgical"
)

// ExclusionType represents the type of plan exclusion
type ExclusionType string

const (
	ExclusionTypePreExisting  ExclusionType = "pre_existing"
	ExclusionTypeCosmetic     ExclusionType = "cosmetic"
	ExclusionTypeExperimental ExclusionType = "experimental"
)

// PremiumRuleType represents the classification of a premium rule
type PremiumRuleType string

const (
	PremiumRuleTypeBaseRate       PremiumRuleType = "base_rate"
	PremiumRuleTypeAgeBand        PremiumRuleType = "age_band"
	PremiumRuleTypeFamilySize     PremiumRuleType = "family_size"
	PremiumRuleTypeMaternityAddon PremiumRuleType = "maternity_addon"
	PremiumRuleTypeLoading        PremiumRuleType = "loading"
	PremiumRuleTypeDiscount       PremiumRuleType = "discount"
	PremiumRuleTypeCorporateRate  PremiumRuleType = "corporate_rate"
)

// PremiumCalculationType represents how premiums are calculated
type PremiumCalculationType string

const (
	PremiumCalculationTypeFlat      PremiumCalculationType = "flat"
	PremiumCalculationTypePerMember PremiumCalculationType = "per_member"
	PremiumCalculationTypeTiered    PremiumCalculationType = "tiered"
	PremiumCalculationTypePerFamily PremiumCalculationType = "per_family"
)

// DiscountType represents the type of discount applied
type DiscountType string

const (
	DiscountTypePercentage DiscountType = "percentage"
	DiscountTypeFixed      DiscountType = "fixed"
)

// BillingFrequency represents how often billing occurs
type BillingFrequency string

const (
	BillingFrequencyMonthly    BillingFrequency = "monthly"
	BillingFrequencyQuarterly  BillingFrequency = "quarterly"
	BillingFrequencySemiAnnual BillingFrequency = "semi_annual"
	BillingFrequencyAnnual     BillingFrequency = "annual"
)

// InstallmentScheduleStatus represents the status of an installment schedule
type InstallmentScheduleStatus string

const (
	InstallmentScheduleStatusActive    InstallmentScheduleStatus = "ACTIVE"
	InstallmentScheduleStatusCompleted InstallmentScheduleStatus = "COMPLETED"
	InstallmentScheduleStatusCancelled InstallmentScheduleStatus = "CANCELLED"
)

// InstallmentStatus represents the status of an individual installment
type InstallmentStatus string

const (
	InstallmentStatusPending InstallmentStatus = "PENDING"
	InstallmentStatusPaid    InstallmentStatus = "PAID"
	InstallmentStatusOverdue InstallmentStatus = "OVERDUE"
)

// ProviderNetworkStatus represents the status of a provider network entry
type ProviderNetworkStatus string

const (
	ProviderNetworkStatusActive   ProviderNetworkStatus = "ACTIVE"
	ProviderNetworkStatusInactive ProviderNetworkStatus = "INACTIVE"
)

// Currency represents currency codes
type Currency string

const (
	CurrencyKES Currency = "KES"
)

// LeadStatus represents the status of a sales lead
type LeadStatus string

const (
	LeadStatusNew          LeadStatus = "NEW"
	LeadStatusContacted    LeadStatus = "CONTACTED"
	LeadStatusQualified    LeadStatus = "QUALIFIED"
	LeadStatusProposalSent LeadStatus = "PROPOSAL_SENT"
	LeadStatusNegotiation  LeadStatus = "NEGOTIATION"
	LeadStatusWon          LeadStatus = "WON"
	LeadStatusLost         LeadStatus = "LOST"
	LeadStatusDormant      LeadStatus = "DORMANT"
)

// LeadSource represents where a lead came from
type LeadSource string

const (
	LeadSourceDirect   LeadSource = "direct"
	LeadSourceReferral LeadSource = "referral"
	LeadSourceWeb      LeadSource = "web"
	LeadSourceAgent    LeadSource = "agent"
	LeadSourceBroker   LeadSource = "broker"
)

// QuotationStatus represents the lifecycle status of a quotation
type QuotationStatus string

const (
	QuotationStatusDraft           QuotationStatus = "DRAFT"
	QuotationStatusIssued          QuotationStatus = "ISSUED"
	QuotationStatusPendingDecision QuotationStatus = "PENDING_DECISION"
	QuotationStatusAccepted        QuotationStatus = "ACCEPTED"
	QuotationStatusDeclined        QuotationStatus = "DECLINED"
	QuotationStatusExpired         QuotationStatus = "EXPIRED"
	QuotationStatusConverted       QuotationStatus = "CONVERTED"
)

// QuotationType represents the type of quotation
type QuotationType string

const (
	QuotationTypeStandard   QuotationType = "standard"
	QuotationTypeTailorMade QuotationType = "tailor_made"
)

// ApprovalStatus represents the approval status of a quotation version
type ApprovalStatus string

const (
	ApprovalStatusNone     ApprovalStatus = "NONE"
	ApprovalStatusPending  ApprovalStatus = "PENDING"
	ApprovalStatusApproved ApprovalStatus = "APPROVED"
	ApprovalStatusRejected ApprovalStatus = "REJECTED"
)

// LeadActivityType represents the type of activity on a lead
type LeadActivityType string

const (
	LeadActivityTypeCall         LeadActivityType = "call"
	LeadActivityTypeEmail        LeadActivityType = "email"
	LeadActivityTypeMeeting      LeadActivityType = "meeting"
	LeadActivityTypeNote         LeadActivityType = "note"
	LeadActivityTypeFollowUp     LeadActivityType = "follow_up"
	LeadActivityTypeStatusChange LeadActivityType = "status_change"
)

// LoadingType represents the type of loading applied
type LoadingType string

const (
	LoadingTypePercentage LoadingType = "percentage"
	LoadingTypeFixed      LoadingType = "fixed"
)

// UnderwritingFlagType represents types of underwriting flags
type UnderwritingFlagType string

const (
	UnderwritingFlagMaxAge          UnderwritingFlagType = "MAX_AGE"
	UnderwritingFlagMinAge          UnderwritingFlagType = "MIN_AGE"
	UnderwritingFlagDoubleInsurance UnderwritingFlagType = "DOUBLE_INSURANCE"
	UnderwritingFlagPreExisting     UnderwritingFlagType = "PRE_EXISTING_CONDITION"
	UnderwritingFlagBMIThreshold    UnderwritingFlagType = "BMI_THRESHOLD"
	UnderwritingFlagWaitingPeriod   UnderwritingFlagType = "WAITING_PERIOD"
	UnderwritingFlagRenewalSkip     UnderwritingFlagType = "RENEWAL_SKIP"
)

// UnderwritingFlagStatus represents the status of an underwriting flag
type UnderwritingFlagStatus string

const (
	UnderwritingFlagStatusOpen         UnderwritingFlagStatus = "OPEN"
	UnderwritingFlagStatusAcknowledged UnderwritingFlagStatus = "ACKNOWLEDGED"
	UnderwritingFlagStatusResolved     UnderwritingFlagStatus = "RESOLVED"
	UnderwritingFlagStatusOverridden   UnderwritingFlagStatus = "OVERRIDDEN"
)

// UnderwritingRuleType represents types of underwriting rules
type UnderwritingRuleType string

const (
	UnderwritingRuleMaxAge          UnderwritingRuleType = "MAX_AGE"
	UnderwritingRuleMinAge          UnderwritingRuleType = "MIN_AGE"
	UnderwritingRuleDoubleInsurance UnderwritingRuleType = "DOUBLE_INSURANCE"
	UnderwritingRulePreExisting     UnderwritingRuleType = "PRE_EXISTING_CONDITION"
	UnderwritingRuleBMIThreshold    UnderwritingRuleType = "BMI_THRESHOLD"
	UnderwritingRuleWaitingPeriod   UnderwritingRuleType = "WAITING_PERIOD"
)

// CreditNoteStatus represents the status of a credit note
type CreditNoteStatus string

const (
	CreditNoteStatusDraft     CreditNoteStatus = "DRAFT"
	CreditNoteStatusApproved  CreditNoteStatus = "APPROVED"
	CreditNoteStatusApplied   CreditNoteStatus = "APPLIED"
	CreditNoteStatusCancelled CreditNoteStatus = "CANCELLED"
)

// TreatyType represents the type of reinsurance treaty
type TreatyType string

const (
	TreatyTypeQuotaShare TreatyType = "QUOTA_SHARE"
	TreatyTypeXOL        TreatyType = "XOL"
)

// TreatyStatus represents the status of a reinsurance treaty
type TreatyStatus string

const (
	TreatyStatusDraft      TreatyStatus = "DRAFT"
	TreatyStatusActive     TreatyStatus = "ACTIVE"
	TreatyStatusExpired    TreatyStatus = "EXPIRED"
	TreatyStatusTerminated TreatyStatus = "TERMINATED"
)

// CessionType represents the type of reinsurance cession
type CessionType string

const (
	CessionTypePremium CessionType = "PREMIUM"
	CessionTypeClaim   CessionType = "CLAIM"
)

// CessionStatus represents the status of a cession
type CessionStatus string

const (
	CessionStatusPending  CessionStatus = "PENDING"
	CessionStatusBooked   CessionStatus = "BOOKED"
	CessionStatusReversed CessionStatus = "REVERSED"
)

// RecoveryStatus represents the status of a reinsurance recovery
type RecoveryStatus string

const (
	RecoveryStatusNotified      RecoveryStatus = "NOTIFIED"
	RecoveryStatusAcknowledged  RecoveryStatus = "ACKNOWLEDGED"
	RecoveryStatusInfoRequested RecoveryStatus = "INFO_REQUESTED"
	RecoveryStatusApproved      RecoveryStatus = "APPROVED"
	RecoveryStatusPaid          RecoveryStatus = "PAID"
	RecoveryStatusWrittenOff    RecoveryStatus = "WRITTEN_OFF"
)

// RecoveryWorkflowStatus represents workflow stages for recovery processing
type RecoveryWorkflowStatus string

const (
	RecoveryWorkflowNotification   RecoveryWorkflowStatus = "NOTIFICATION"
	RecoveryWorkflowAcknowledgment RecoveryWorkflowStatus = "ACKNOWLEDGMENT"
	RecoveryWorkflowInfoRequest    RecoveryWorkflowStatus = "INFO_REQUEST"
	RecoveryWorkflowApproval       RecoveryWorkflowStatus = "APPROVAL"
	RecoveryWorkflowPayment        RecoveryWorkflowStatus = "PAYMENT"
)

// BordereauType represents the type of bordereau report
type BordereauType string

const (
	BordereauTypePremium BordereauType = "PREMIUM"
	BordereauTypeClaim   BordereauType = "CLAIM"
)

// BordereauStatus represents the status of a bordereau
type BordereauStatus string

const (
	BordereauStatusDraft     BordereauStatus = "DRAFT"
	BordereauStatusFinalized BordereauStatus = "FINALIZED"
	BordereauStatusSent      BordereauStatus = "SENT"
)

// ReinsurerStatementStatus represents the status of a reinsurer statement
type ReinsurerStatementStatus string

const (
	ReinsurerStatementStatusDraft        ReinsurerStatementStatus = "DRAFT"
	ReinsurerStatementStatusIssued       ReinsurerStatementStatus = "ISSUED"
	ReinsurerStatementStatusAcknowledged ReinsurerStatementStatus = "ACKNOWLEDGED"
	ReinsurerStatementStatusSettled      ReinsurerStatementStatus = "SETTLED"
)

// ProfitCommissionType represents the type of profit commission calculation
type ProfitCommissionType string

const (
	ProfitCommissionTypeSlidingScale ProfitCommissionType = "SLIDING_SCALE"
	ProfitCommissionTypeFlat         ProfitCommissionType = "FLAT"
	ProfitCommissionTypeCarryForward ProfitCommissionType = "CARRY_FORWARD"
)

// TreatyAlertType represents the type of treaty alert
type TreatyAlertType string

const (
	TreatyAlertTypeLimitBreach          TreatyAlertType = "LIMIT_BREACH"
	TreatyAlertTypeAggregateWarning     TreatyAlertType = "AGGREGATE_WARNING"
	TreatyAlertTypeCatastropheThreshold TreatyAlertType = "CATASTROPHE_THRESHOLD"
	TreatyAlertTypeExpiryWarning        TreatyAlertType = "EXPIRY_WARNING"
)

// TreatyAlertSeverity represents the severity of a treaty alert
type TreatyAlertSeverity string

const (
	TreatyAlertSeverityLow      TreatyAlertSeverity = "LOW"
	TreatyAlertSeverityMedium   TreatyAlertSeverity = "MEDIUM"
	TreatyAlertSeverityHigh     TreatyAlertSeverity = "HIGH"
	TreatyAlertSeverityCritical TreatyAlertSeverity = "CRITICAL"
)

// AccreditationStatus represents the accreditation status of a provider
type AccreditationStatus string

const (
	AccreditationStatusNone       AccreditationStatus = "NONE"
	AccreditationStatusPending    AccreditationStatus = "PENDING"
	AccreditationStatusAccredited AccreditationStatus = "ACCREDITED"
	AccreditationStatusExpired    AccreditationStatus = "EXPIRED"
	AccreditationStatusRevoked    AccreditationStatus = "REVOKED"
)

// ReportStatus represents the status of a generated report
type ReportStatus string

const (
	ReportStatusGenerating ReportStatus = "GENERATING"
	ReportStatusCompleted  ReportStatus = "COMPLETED"
	ReportStatusFailed     ReportStatus = "FAILED"
	ReportStatusExpired    ReportStatus = "EXPIRED"
)

// ReportFormat represents report export formats
type ReportFormat string

const (
	ReportFormatCSV  ReportFormat = "CSV"
	ReportFormatXLSX ReportFormat = "XLSX"
	ReportFormatPDF  ReportFormat = "PDF"
)

// ReportCategory represents report categories
type ReportCategory string

const (
	ReportCategoryClaims      ReportCategory = "CLAIMS"
	ReportCategoryPremium     ReportCategory = "PREMIUM"
	ReportCategoryMembership  ReportCategory = "MEMBERSHIP"
	ReportCategoryProvider    ReportCategory = "PROVIDER"
	ReportCategoryReinsurance ReportCategory = "REINSURANCE"
	ReportCategoryManagement  ReportCategory = "MANAGEMENT"
)

// ReportType represents types of report definitions
type ReportType string

const (
	ReportTypePreBuilt ReportType = "PRE_BUILT"
	ReportTypeAdHoc    ReportType = "AD_HOC"
)

// Reinsurance audit entity types
const (
	AuditEntityTypeTreaty                AuditEntityType = "TREATY"
	AuditEntityTypeTreatyParticipant     AuditEntityType = "TREATY_PARTICIPANT"
	AuditEntityTypeTreatyLayer           AuditEntityType = "TREATY_LAYER"
	AuditEntityTypeCession               AuditEntityType = "CESSION"
	AuditEntityTypeReinsuranceRecovery   AuditEntityType = "REINSURANCE_RECOVERY"
	AuditEntityTypeRecoveryWorkflowEvent AuditEntityType = "RECOVERY_WORKFLOW_EVENT"
	AuditEntityTypeBordereau             AuditEntityType = "BORDEREAU"
	AuditEntityTypeBordereauItem         AuditEntityType = "BORDEREAU_ITEM"
	AuditEntityTypeReinsurerStatement    AuditEntityType = "REINSURER_STATEMENT"
	AuditEntityTypeProfitCommission      AuditEntityType = "PROFIT_COMMISSION"
	AuditEntityTypeTreatyAlert           AuditEntityType = "TREATY_ALERT"
	AuditEntityTypeReportDefinition      AuditEntityType = "REPORT_DEFINITION"
	AuditEntityTypeReportSchedule        AuditEntityType = "REPORT_SCHEDULE"
	AuditEntityTypeGeneratedReport       AuditEntityType = "GENERATED_REPORT"
)
