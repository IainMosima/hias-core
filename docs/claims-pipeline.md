# Claims Processing Pipeline

## Overview

When a claim is submitted via `POST /api/v1/claims`, it goes through an automated pipeline that validates, runs fraud checks, adjudicates, evaluates escalation rules, and publishes events — all in a single synchronous request. The response includes the final status, calculated amounts, adjudication decision, and any fraud flags.

## Pipeline Flow

```
Submit Claim
    |
    v
[1] Create claim record (status: RECEIVED)
    |
    v
[2] Create line items, calculate total_amount
    |
    v
[3] VALIDATE --- fail ---> REJECTED (with reasons)
    |
   pass (status: VALIDATED)
    |
    v
[4] FRAUD CHECKS (run BEFORE adjudication)
    |   - Duplicate check
    |   - Frequency check
    |   - Amount threshold
    |   - Expired contract
    |   - Suspended provider
    |   - Rate card overcharge
    |   - Repeat visit
    |
    v
[5] Check for CRITICAL fraud flags
    |   - If CRITICAL/HIGH fraud -> override to MANUAL_REVIEW after adjudication
    |
    v
[6] ADJUDICATE (applies configurable rules from DB)
    |
    +--> REJECT ---------> REJECTED (with rule results)
    |
    +--> MANUAL_REVIEW ---> MANUAL_REVIEW (flagged for human)
    |
    +--> APPROVE ---------> ADJUDICATED (amounts calculated)
    |
    v
[7] Override: if critical fraud + APPROVE -> force MANUAL_REVIEW
    |
    v
[8] Store adjudication decision
    |
    v
[9] Update claim amounts + status
    |
    v
[10] Evaluate ESCALATION RULES (configurable from DB)
    |   - AMOUNT_EXCEEDS -> escalate to senior reviewer
    |   - FRAUD_FLAG -> escalate if HIGH/CRITICAL severity
    |   - MANUAL_REVIEW -> escalate per rule config
    |   Result: may set status to ESCALATED
    |
    v
[11] Publish ClaimSubmittedEvent to queue (async)
    |
    v
[12] If preauth_id provided, mark PreAuth as CLAIMED
    |
    v
[13] Record timeline entry (claim_status_history)
    |
    v
[14] Return final claim with all data
```

## Stage Details

### Stage 1-2: Claim Creation

- Generates a claim number (`CLM-YYYY-NNNNNN`)
- Calculates `total_amount` from line items (`unit_price * quantity`)
- Stores the claim as `RECEIVED` and creates all line items

**Files:** `services/claims/claim_service_impl.go` (SubmitClaim)

### Stage 3: Validation

The validator performs **pre-adjudication checks** to catch obviously invalid claims early. If any check fails, the claim is rejected immediately without running adjudication.

| Check | Rule | Failure Reason |
|-------|------|----------------|
| Policy exists and is ACTIVE | `policy_active` | "Policy is {status}, must be ACTIVE" |
| Member exists and belongs to policy | `member_active` | "Member not found" or "Member does not belong to this policy" |
| Provider exists and is ACTIVE | `provider_active` | "Provider is {status}, must be ACTIVE" |
| At least one line item | `line_items` | "Claim must have at least one line item" |
| Total amount > 0 | `amount_positive` | "Claim total amount must be positive" |

If validation fails, the claim is rejected with all error reasons joined by `;`.

**Files:** `services/claims/validator_service_impl.go`

### Stage 4: Fraud Checks (BEFORE Adjudication)

Fraud checks run **before** adjudication so that critical fraud flags can influence the adjudication outcome. Each check creates a `FraudFlag` record with a severity level.

| Check | What It Does | Flag Type | Severity |
|-------|-------------|-----------|----------|
| Duplicate | Same claim number, different ID | `DUPLICATE` | HIGH |
| Frequency | Same member+provider+procedure repeated | `FREQUENCY` | MEDIUM |
| Amount Threshold | Total > 500,000 KES (50,000,000 cents) | `AMOUNT_THRESHOLD` | HIGH |
| Expired Contract | Provider contract has expired | `EXPIRED_CONTRACT` | HIGH |
| Suspended Provider | Provider is suspended | `SUSPENDED_PROVIDER` | CRITICAL |
| Rate Card Overcharge | Unit price exceeds provider's agreed rate | `RATE_CARD_OVERCHARGE` | MEDIUM |
| Repeat Visit | Same member+provider within short window | `REPEAT_VISIT` | LOW |

**Files:** `services/claims/fraud_service_impl.go`, `services/claims/claim_service_impl.go` (Step 3)

### Stage 5: Critical Fraud Override

After fraud checks, the system checks if any flag has `HIGH` or `CRITICAL` severity. If so, even if adjudication approves the claim, the decision is overridden to `MANUAL_REVIEW`.

### Stage 6: Adjudication

The adjudicator runs a series of **business rules** and produces an `AdjudicationResult` with a decision (`APPROVE`, `REJECT`, or `MANUAL_REVIEW`), calculated amounts, and detailed rule results.

Rules are evaluated in order. A `FAIL` on any critical rule short-circuits with a `REJECT`.

#### Hardcoded Rules (always run)

| # | Rule | Check | On Fail |
|---|------|-------|---------|
| 1 | Policy Active | Policy status = ACTIVE | REJECT |
| 2 | Member Enrolled | Member exists in DB | REJECT |
| 3 | Provider Active | Provider status = ACTIVE | REJECT |
| 4 | Benefits Exist | Plan has at least one benefit | REJECT |
| 5 | Waiting Period | `service_date >= member.created_at + waiting_period_days` | REJECT |
| 6 | Exclusion Check | Claim diagnosis codes vs plan exclusion ICD codes | REJECT |
| 7 | Annual Limits | Calculate remaining limit, apply co-pay/deductible | Partial approval |
| 8 | Pre-Auth Validation | If preauth_id provided, verify it's APPROVED and not expired | REJECT |
| 9 | Duplicate Check | Via fraud service | MANUAL_REVIEW |

#### Configurable Rules (from adjudication_rules table)

After hardcoded rules pass, the system fetches active `AdjudicationRule` records from the DB and evaluates them:

| Rule Type | What It Does | On Fail |
|-----------|-------------|---------|
| `AMOUNT_THRESHOLD` | Rejects if amount exceeds configured threshold | REJECT |
| `FREQUENCY_LIMIT` | Rejects if member claims exceed monthly limit | REJECT |
| `AUTO_APPROVE` | Auto-approves if amount is below configured threshold | FLAG (skip manual) |

**CRUD for rules:** `GET/POST/PUT/DELETE /api/v1/adjudication-rules`

**Files:** `services/claims/adjudicator_service_impl.go`

#### Amount Calculation

```
payable_amount       = min(total_amount, remaining_annual_limit) - co_pay - deductible
co_pay_amount        = calculated from benefit co-pay rules (percentage or fixed)
deductible_amount    = calculated from benefit deductible rules
member_responsibility = co_pay + deductible + any excess above annual limit
```

### Stage 8-9: Store Decision & Update Claim

The adjudication result is:
1. Stored as an `AdjudicationDecision` record (with JSON rule results)
2. Mapped to a claim status:
   - `APPROVE` -> `ADJUDICATED`
   - `REJECT` -> `REJECTED` (with rejection reason from failed rules)
   - `MANUAL_REVIEW` -> `MANUAL_REVIEW`
3. Claim amounts updated: `approved_amount`, `co_pay_amount`, `member_responsibility`

### Stage 10: Escalation Rules

After status is determined, configurable escalation rules from the `escalation_rules` table are evaluated:

| Condition | Check | Action |
|-----------|-------|--------|
| `AMOUNT_EXCEEDS` | Claim amount > rule threshold | Set status to `ESCALATED`, assign to escalation role |
| `FRAUD_FLAG` | Any fraud flag with HIGH or CRITICAL severity | Set status to `ESCALATED` |
| `MANUAL_REVIEW` | Claim is in MANUAL_REVIEW | Set status to `ESCALATED` |

**CRUD for rules:** `GET/POST/PUT/DELETE /api/v1/escalation-rules`

**Files:** `services/claims/claim_service_impl.go` (evaluateEscalationRules)

### Stage 11: Event Publishing

A `ClaimSubmittedEvent` is published to the `claim-processing` SQS queue (async, non-blocking):
```json
{
  "claim_id": "uuid",
  "claim_number": "CLM-2026-000001",
  "policy_id": "uuid",
  "member_id": "uuid",
  "total_amount": 1000000
}
```

### Stage 13: Timeline Recording

Every status change is recorded in the `claim_status_history` table for the claim timeline:
```json
{
  "from_status": "RECEIVED",
  "to_status": "ADJUDICATED",
  "action": "SUBMIT",
  "performed_by": "uuid",
  "created_at": "2026-03-06T10:00:00Z"
}
```

Retrieve via `GET /api/v1/claims/:id/timeline`.

### Stage 14: Response

The claim is re-fetched from DB to include all updated amounts and status. The response includes:
```json
{
  "status": "success",
  "message": "Claim submitted successfully",
  "data": {
    "id": "...",
    "claim_number": "CLM-2026-000001",
    "status": "ADJUDICATED",
    "total_amount": 1000000,
    "approved_amount": 800000,
    "co_pay_amount": 100000,
    "member_responsibility": 200000,
    "line_items": [...],
    "decision": { "decision": "APPROVE", "payable_amount": 800000, "rule_results": [...] },
    "fraud_flags": [...]
  }
}
```

When fetched via `GET /claims/:id`, additional nested data is included:
- `line_items[]` - individual procedures with approved amounts
- `decision` - full adjudication decision with rule results
- `fraud_flags[]` - any fraud flags raised (with severity and resolution status)

## Post-Pipeline Actions

### Vet Claim (`PUT /claims/:id/vet`)
- Requires claim status: `ADJUDICATED` or `MANUAL_REVIEW` or `ESCALATED`
- Claims officer reviews and optionally adjusts amounts
- Sets status to `VETTED`

### Approve Claim (`PUT /claims/:id/approve`)
- Requires claim status: `ADJUDICATED`, `MANUAL_REVIEW`, `VETTED`, or `ESCALATED`
- Uses amounts from the stored adjudication decision (or vetted amounts)
- Sets status to `APPROVED`
- Publishes `ClaimApprovedEvent` to queue

### Reject Claim (`PUT /claims/:id/reject`)
- Requires claim status: `RECEIVED`, `VALIDATED`, `ADJUDICATED`, `MANUAL_REVIEW`, or `ESCALATED`
- Cannot reject `APPROVED` or `PAID` claims
- Requires a `reason` in the request body

### Ready for Payment (`PUT /claims/:id/ready-for-payment`)
- Requires claim status: `APPROVED`
- Sets status to `READY_FOR_PAYMENT`

### Mark Paid (`PUT /claims/:id/mark-paid`)
- Requires claim status: `READY_FOR_PAYMENT` or `APPROVED`
- Sets status to `PAID`

### Mark Part Paid (`PUT /claims/:id/mark-part-paid`)
- Requires claim status: `READY_FOR_PAYMENT` or `APPROVED`
- Sets status to `PART_PAID`

## Claim Status Lifecycle

```
RECEIVED -> VALIDATED -> ADJUDICATED -> VETTED -> APPROVED -> READY_FOR_PAYMENT -> PAID
                |              |           |                                     -> PART_PAID
                |              |           +-> REJECTED
                |              +--> MANUAL_REVIEW -> VETTED -> APPROVED
                |              |                  -> REJECTED
                |              +--> ESCALATED -> VETTED -> APPROVED
                |              |              -> REJECTED
                |              +--> REJECTED
                +--> REJECTED
```

| Status | Meaning |
|--------|---------|
| `RECEIVED` | Claim created, not yet processed |
| `VALIDATED` | Passed pre-adjudication validation |
| `ADJUDICATED` | Adjudication complete, auto-approved by rules, awaiting human review |
| `MANUAL_REVIEW` | Flagged for human review (fraud, duplicate, etc.) |
| `ESCALATED` | Escalated to senior reviewer per escalation rules |
| `VETTED` | Claims officer has reviewed and verified amounts |
| `PARTIALLY_VETTED` | Partial vetting complete |
| `APPROVED` | Human-approved, ready for payment |
| `REJECTED` | Rejected at any stage (validation, adjudication, or manual) |
| `READY_FOR_PAYMENT` | Approved and queued for payment processing |
| `PAID` | Payment processed (set by billing/remittance) |
| `PART_PAID` | Partial payment processed |

## Bulk Operations

### Bulk Submit (`POST /claims/bulk`)
Submit multiple claims in one request. Each goes through the full pipeline independently.

### CSV Import (`POST /claims/import-csv`)
Import claims from a CSV file. Admin/Claims Officer role required.

### SLA Monitoring (`GET /claims/sla-breached`)
Returns claims that have exceeded their SLA deadline. The `ClaimSLATask` scheduler checks every 15 minutes and sends notifications.

## Key Files Reference

| File | Purpose |
|------|---------|
| `services/claims/claim_service_impl.go` | Orchestrates the full pipeline in SubmitClaim |
| `services/claims/validator_service_impl.go` | Pre-adjudication validation rules |
| `services/claims/adjudicator_service_impl.go` | Business rules engine (eligibility, coverage, limits, configurable rules) |
| `services/claims/fraud_service_impl.go` | Fraud check implementations (7 check types) |
| `domains/claims/entity/adjudication_decision.go` | AdjudicationDecision + AdjudicationResult + RuleResult types |
| `domains/claims/entity/fraud_flag.go` | FraudFlag entity |
| `domains/claims/entity/claim_status_history.go` | Timeline entry entity |
| `domains/claims/repository/claim_repository.go` | ClaimRepository interface |
| `domains/claims/repository/adjudication_repository.go` | AdjudicationRepository interface |
| `domains/claims/repository/adjudication_rule_repository.go` | Configurable adjudication rules |
| `domains/claims/repository/escalation_rule_repository.go` | Configurable escalation rules |
| `infrastructures/scheduler/tasks/claim_sla_task.go` | SLA breach monitoring task |

## Exclusion Configuration

Exclusions are configured per plan via the API:

```
POST /api/v1/plans/:id/exclusions   - Create exclusion
GET  /api/v1/plans/:id/exclusions   - List exclusions for plan
PUT  /api/v1/exclusions/:id         - Update exclusion
DELETE /api/v1/exclusions/:id       - Delete exclusion
```

Request body:
```json
{
  "description": "Cosmetic procedures",
  "type": "cosmetic",
  "icd_codes": ["Z41.1", "Z41.8", "Z41.9"]
}
```

Types: `pre_existing`, `cosmetic`, `experimental`
