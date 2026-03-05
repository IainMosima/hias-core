package schema

import (
	"time"

	"github.com/bitbiz/hias-core/domains/reinsurance/entity"
	"github.com/google/uuid"
)

type TreatyResponse struct {
	ID             uuid.UUID `json:"id"`
	TreatyNumber   string    `json:"treaty_number"`
	Name           string    `json:"name"`
	TreatyType     string    `json:"treaty_type"`
	Status         string    `json:"status"`
	EffectiveDate  time.Time `json:"effective_date"`
	ExpiryDate     time.Time `json:"expiry_date"`
	RetentionLimit int64     `json:"retention_limit"`
	Currency       string    `json:"currency"`
	Notes          string    `json:"notes,omitempty"`
	CreatedBy      uuid.UUID `json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type TreatyDetailResponse struct {
	TreatyResponse
	Participants    []TreatyParticipantResponse `json:"participants"`
	Layers          []TreatyLayerResponse       `json:"layers"`
	ProfitCommRules []ProfitCommissionResponse  `json:"profit_commission_rules"`
}

type TreatyParticipantResponse struct {
	ID              uuid.UUID `json:"id"`
	TreatyID        uuid.UUID `json:"treaty_id"`
	ReinsurerName   string    `json:"reinsurer_name"`
	SharePercentage float64   `json:"share_percentage"`
	CommissionRate  float64   `json:"commission_rate"`
	IsLead          bool      `json:"is_lead"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type TreatyLayerResponse struct {
	ID               uuid.UUID `json:"id"`
	TreatyID         uuid.UUID `json:"treaty_id"`
	LayerNumber      int       `json:"layer_number"`
	AttachmentPoint  int64     `json:"attachment_point"`
	LayerLimit       int64     `json:"layer_limit"`
	DeductibleAmount int64     `json:"deductible_amount"`
	PremiumRate      float64   `json:"premium_rate"`
	AggregateLimit   *int64    `json:"aggregate_limit,omitempty"`
	AggregateUsed    int64     `json:"aggregate_used"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type CessionResponse struct {
	ID               uuid.UUID `json:"id"`
	CessionNumber    string    `json:"cession_number"`
	TreatyID         uuid.UUID `json:"treaty_id"`
	PolicyID         uuid.UUID `json:"policy_id"`
	TreatyLayerID    uuid.UUID `json:"treaty_layer_id,omitempty"`
	CessionType      string    `json:"cession_type"`
	GrossAmount      int64     `json:"gross_amount"`
	CededAmount      int64     `json:"ceded_amount"`
	RetainedAmount   int64     `json:"retained_amount"`
	CommissionAmount int64     `json:"commission_amount"`
	SharePercentage  float64   `json:"share_percentage"`
	Status           string    `json:"status"`
	CreatedBy        uuid.UUID `json:"created_by"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type RecoveryResponse struct {
	ID                uuid.UUID `json:"id"`
	RecoveryNumber    string    `json:"recovery_number"`
	ClaimID           uuid.UUID `json:"claim_id"`
	TreatyID          uuid.UUID `json:"treaty_id"`
	TreatyLayerID     uuid.UUID `json:"treaty_layer_id,omitempty"`
	CessionID         uuid.UUID `json:"cession_id,omitempty"`
	GrossClaimAmount  int64     `json:"gross_claim_amount"`
	RecoverableAmount int64     `json:"recoverable_amount"`
	RecoveredAmount   int64     `json:"recovered_amount"`
	OutstandingAmount int64     `json:"outstanding_amount"`
	Status            string    `json:"status"`
	WorkflowStatus    string    `json:"workflow_status"`
	Notes             string    `json:"notes,omitempty"`
	CreatedBy         uuid.UUID `json:"created_by"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type RecoveryDetailResponse struct {
	RecoveryResponse
	WorkflowEvents []RecoveryWorkflowEventResponse `json:"workflow_events"`
}

type RecoveryWorkflowEventResponse struct {
	ID          uuid.UUID `json:"id"`
	RecoveryID  uuid.UUID `json:"recovery_id"`
	FromStatus  string    `json:"from_status"`
	ToStatus    string    `json:"to_status"`
	EventType   string    `json:"event_type"`
	Notes       string    `json:"notes,omitempty"`
	PerformedBy uuid.UUID `json:"performed_by"`
	CreatedAt   time.Time `json:"created_at"`
}

type BordereauResponse struct {
	ID              uuid.UUID `json:"id"`
	BordereauNumber string    `json:"bordereau_number"`
	TreatyID        uuid.UUID `json:"treaty_id"`
	BordereauType   string    `json:"bordereau_type"`
	PeriodStart     time.Time `json:"period_start"`
	PeriodEnd       time.Time `json:"period_end"`
	TotalGross      int64     `json:"total_gross"`
	TotalCeded      int64     `json:"total_ceded"`
	TotalCommission int64     `json:"total_commission"`
	ItemCount       int       `json:"item_count"`
	Status          string    `json:"status"`
	CreatedBy       uuid.UUID `json:"created_by"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type BordereauDetailResponse struct {
	BordereauResponse
	Items []BordereauItemResponse `json:"items"`
}

type BordereauItemResponse struct {
	ID               uuid.UUID `json:"id"`
	BordereauID      uuid.UUID `json:"bordereau_id"`
	CessionID        uuid.UUID `json:"cession_id,omitempty"`
	RecoveryID       uuid.UUID `json:"recovery_id,omitempty"`
	PolicyNumber     string    `json:"policy_number,omitempty"`
	ClaimNumber      string    `json:"claim_number,omitempty"`
	GrossAmount      int64     `json:"gross_amount"`
	CededAmount      int64     `json:"ceded_amount"`
	CommissionAmount int64     `json:"commission_amount"`
	CreatedAt        time.Time `json:"created_at"`
}

type ReinsurerStatementResponse struct {
	ID               uuid.UUID `json:"id"`
	StatementNumber  string    `json:"statement_number"`
	TreatyID         uuid.UUID `json:"treaty_id"`
	ParticipantID    uuid.UUID `json:"participant_id"`
	PeriodStart      time.Time `json:"period_start"`
	PeriodEnd        time.Time `json:"period_end"`
	PremiumCeded     int64     `json:"premium_ceded"`
	ClaimsRecovered  int64     `json:"claims_recovered"`
	CommissionDue    int64     `json:"commission_due"`
	ProfitCommission int64     `json:"profit_commission"`
	NetBalance       int64     `json:"net_balance"`
	Status           string    `json:"status"`
	CreatedBy        uuid.UUID `json:"created_by"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type ProfitCommissionResponse struct {
	ID                  uuid.UUID  `json:"id"`
	TreatyID            uuid.UUID  `json:"treaty_id"`
	CommissionType      string     `json:"commission_type"`
	LossRatioFrom       float64    `json:"loss_ratio_from"`
	LossRatioTo         float64    `json:"loss_ratio_to"`
	CommissionRate      float64    `json:"commission_rate"`
	CarryForwardYears   int        `json:"carry_forward_years"`
	CarryForwardBalance int64      `json:"carry_forward_balance"`
	PeriodStart         *time.Time `json:"period_start,omitempty"`
	PeriodEnd           *time.Time `json:"period_end,omitempty"`
	CalculatedAmount    int64      `json:"calculated_amount"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

type ProfitCommissionCalculationResponse struct {
	TreatyID         uuid.UUID `json:"treaty_id"`
	PremiumCeded     int64     `json:"premium_ceded"`
	ClaimsRecovered  int64     `json:"claims_recovered"`
	LossRatio        float64   `json:"loss_ratio"`
	NetProfit        int64     `json:"net_profit"`
	CommissionRate   float64   `json:"commission_rate"`
	CommissionAmount int64     `json:"commission_amount"`
	CarryForward     int64     `json:"carry_forward"`
}

type TreatyAlertResponse struct {
	ID             uuid.UUID  `json:"id"`
	TreatyID       uuid.UUID  `json:"treaty_id"`
	TreatyLayerID  uuid.UUID  `json:"treaty_layer_id,omitempty"`
	AlertType      string     `json:"alert_type"`
	Severity       string     `json:"severity"`
	Message        string     `json:"message"`
	ThresholdValue int64      `json:"threshold_value"`
	CurrentValue   int64      `json:"current_value"`
	IsAcknowledged bool       `json:"is_acknowledged"`
	AcknowledgedBy uuid.UUID  `json:"acknowledged_by,omitempty"`
	AcknowledgedAt *time.Time `json:"acknowledged_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}

type AgedRecoveryBucketResponse struct {
	Bucket           string `json:"bucket"`
	Count            int64  `json:"count"`
	TotalOutstanding int64  `json:"total_outstanding"`
}

type ReinsuranceDashboardResponse struct {
	ActiveTreatyCount    int64   `json:"active_treaty_count"`
	TotalCededPremiums   int64   `json:"total_ceded_premiums"`
	TotalRecoverable     int64   `json:"total_recoverable"`
	TotalRecovered       int64   `json:"total_recovered"`
	TotalOutstanding     int64   `json:"total_outstanding"`
	CessionRatio         float64 `json:"cession_ratio"`
	RecoverySuccessRate  float64 `json:"recovery_success_rate"`
	UnacknowledgedAlerts int64   `json:"unacknowledged_alerts"`
}

// Converter functions

func ToTreatyResponse(t *entity.Treaty) TreatyResponse {
	return TreatyResponse{
		ID: t.ID, TreatyNumber: t.TreatyNumber, Name: t.Name,
		TreatyType: t.TreatyType, Status: t.Status,
		EffectiveDate: t.EffectiveDate, ExpiryDate: t.ExpiryDate,
		RetentionLimit: t.RetentionLimit, Currency: t.Currency,
		Notes: t.Notes, CreatedBy: t.CreatedBy,
		CreatedAt: t.CreatedAt, UpdatedAt: t.UpdatedAt,
	}
}

func ToTreatyParticipantResponse(p *entity.TreatyParticipant) TreatyParticipantResponse {
	return TreatyParticipantResponse{
		ID: p.ID, TreatyID: p.TreatyID, ReinsurerName: p.ReinsurerName,
		SharePercentage: p.SharePercentage, CommissionRate: p.CommissionRate,
		IsLead: p.IsLead, CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt,
	}
}

func ToTreatyLayerResponse(l *entity.TreatyLayer) TreatyLayerResponse {
	return TreatyLayerResponse{
		ID: l.ID, TreatyID: l.TreatyID, LayerNumber: l.LayerNumber,
		AttachmentPoint: l.AttachmentPoint, LayerLimit: l.LayerLimit,
		DeductibleAmount: l.DeductibleAmount, PremiumRate: l.PremiumRate,
		AggregateLimit: l.AggregateLimit, AggregateUsed: l.AggregateUsed,
		CreatedAt: l.CreatedAt, UpdatedAt: l.UpdatedAt,
	}
}

func ToCessionResponse(c *entity.Cession) CessionResponse {
	return CessionResponse{
		ID: c.ID, CessionNumber: c.CessionNumber, TreatyID: c.TreatyID,
		PolicyID: c.PolicyID, TreatyLayerID: c.TreatyLayerID,
		CessionType: c.CessionType, GrossAmount: c.GrossAmount,
		CededAmount: c.CededAmount, RetainedAmount: c.RetainedAmount,
		CommissionAmount: c.CommissionAmount, SharePercentage: c.SharePercentage,
		Status: c.Status, CreatedBy: c.CreatedBy,
		CreatedAt: c.CreatedAt, UpdatedAt: c.UpdatedAt,
	}
}

func ToRecoveryResponse(r *entity.ReinsuranceRecovery) RecoveryResponse {
	return RecoveryResponse{
		ID: r.ID, RecoveryNumber: r.RecoveryNumber, ClaimID: r.ClaimID,
		TreatyID: r.TreatyID, TreatyLayerID: r.TreatyLayerID,
		CessionID: r.CessionID, GrossClaimAmount: r.GrossClaimAmount,
		RecoverableAmount: r.RecoverableAmount, RecoveredAmount: r.RecoveredAmount,
		OutstandingAmount: r.OutstandingAmount, Status: r.Status,
		WorkflowStatus: r.WorkflowStatus, Notes: r.Notes,
		CreatedBy: r.CreatedBy, CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt,
	}
}

func ToRecoveryWorkflowEventResponse(e *entity.RecoveryWorkflowEvent) RecoveryWorkflowEventResponse {
	return RecoveryWorkflowEventResponse{
		ID: e.ID, RecoveryID: e.RecoveryID, FromStatus: e.FromStatus,
		ToStatus: e.ToStatus, EventType: e.EventType, Notes: e.Notes,
		PerformedBy: e.PerformedBy, CreatedAt: e.CreatedAt,
	}
}

func ToBordereauResponse(b *entity.Bordereau) BordereauResponse {
	return BordereauResponse{
		ID: b.ID, BordereauNumber: b.BordereauNumber, TreatyID: b.TreatyID,
		BordereauType: b.BordereauType, PeriodStart: b.PeriodStart,
		PeriodEnd: b.PeriodEnd, TotalGross: b.TotalGross,
		TotalCeded: b.TotalCeded, TotalCommission: b.TotalCommission,
		ItemCount: b.ItemCount, Status: b.Status, CreatedBy: b.CreatedBy,
		CreatedAt: b.CreatedAt, UpdatedAt: b.UpdatedAt,
	}
}

func ToBordereauItemResponse(i *entity.BordereauItem) BordereauItemResponse {
	return BordereauItemResponse{
		ID: i.ID, BordereauID: i.BordereauID, CessionID: i.CessionID,
		RecoveryID: i.RecoveryID, PolicyNumber: i.PolicyNumber,
		ClaimNumber: i.ClaimNumber, GrossAmount: i.GrossAmount,
		CededAmount: i.CededAmount, CommissionAmount: i.CommissionAmount,
		CreatedAt: i.CreatedAt,
	}
}

func ToReinsurerStatementResponse(s *entity.ReinsurerStatement) ReinsurerStatementResponse {
	return ReinsurerStatementResponse{
		ID: s.ID, StatementNumber: s.StatementNumber, TreatyID: s.TreatyID,
		ParticipantID: s.ParticipantID, PeriodStart: s.PeriodStart,
		PeriodEnd: s.PeriodEnd, PremiumCeded: s.PremiumCeded,
		ClaimsRecovered: s.ClaimsRecovered, CommissionDue: s.CommissionDue,
		ProfitCommission: s.ProfitCommission, NetBalance: s.NetBalance,
		Status: s.Status, CreatedBy: s.CreatedBy,
		CreatedAt: s.CreatedAt, UpdatedAt: s.UpdatedAt,
	}
}

func ToProfitCommissionResponse(pc *entity.ProfitCommission) ProfitCommissionResponse {
	return ProfitCommissionResponse{
		ID: pc.ID, TreatyID: pc.TreatyID, CommissionType: pc.CommissionType,
		LossRatioFrom: pc.LossRatioFrom, LossRatioTo: pc.LossRatioTo,
		CommissionRate: pc.CommissionRate, CarryForwardYears: pc.CarryForwardYears,
		CarryForwardBalance: pc.CarryForwardBalance, PeriodStart: pc.PeriodStart,
		PeriodEnd: pc.PeriodEnd, CalculatedAmount: pc.CalculatedAmount,
		CreatedAt: pc.CreatedAt, UpdatedAt: pc.UpdatedAt,
	}
}

func ToTreatyAlertResponse(a *entity.TreatyAlert) TreatyAlertResponse {
	return TreatyAlertResponse{
		ID: a.ID, TreatyID: a.TreatyID, TreatyLayerID: a.TreatyLayerID,
		AlertType: a.AlertType, Severity: a.Severity, Message: a.Message,
		ThresholdValue: a.ThresholdValue, CurrentValue: a.CurrentValue,
		IsAcknowledged: a.IsAcknowledged, AcknowledgedBy: a.AcknowledgedBy,
		AcknowledgedAt: a.AcknowledgedAt, CreatedAt: a.CreatedAt,
	}
}
