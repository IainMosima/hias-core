# Claims Processing Pipeline

## Overview

When a claim is submitted via `POST /api/v1/claims`, it goes through an automated pipeline that validates, adjudicates, and flags the claim in a single synchronous request. The response includes the final status, calculated amounts, adjudication decision, and any fraud flags.

## Pipeline Flow

```
Submit Claim
    |
    v
[1] Create claim record (status: RECEIVED)
    |
    v
[2] Create line items
    |
    v
[3] VALIDATE --- fail ---> REJECTED (with reasons)
    |
   pass
    |
    v
    (status: VALIDATED)
    |
    v
[4] ADJUDICATE
    |
    +--> REJECT ---------> REJECTED (with rule results)
    |
    +--> MANUAL_REVIEW ---> MANUAL_REVIEW (flagged for human)
    |
    +--> APPROVE ---------> ADJUDICATED (amounts calculated)
    |
    v
[5] Store adjudication decision
    |
    v
[6] Update claim amounts + status
    |
    v
[7] Run fraud checks (frequency + amount threshold)
    |
    v
[8] Return final claim with all data
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

### Stage 4: Adjudication

The adjudicator runs a series of **business rules** and produces an `AdjudicationResult` with a decision (`APPROVE`, `REJECT`, or `MANUAL_REVIEW`), calculated amounts, and detailed rule results.

Rules are evaluated in order. A `FAIL` on any critical rule short-circuits with a `REJECT`.

#### Rule 1: Policy Active (eligibility)
- Same as validation, but part of the adjudicator's own rule results
- Fail -> REJECT

#### Rule 2: Member Enrolled (eligibility)
- Member must exist in the DB
- Fail -> REJECT

#### Rule 3: Provider Active (eligibility)
- Provider must exist and have `status = ACTIVE`
- Fail -> REJECT

#### Rule 4: Benefits Exist (coverage)
- At least one benefit must be configured for the policy's plan
- Fail -> REJECT

#### Rule 5: Waiting Period (eligibility)
- For each benefit with `waiting_period_days > 0`:
  - Calculates `waiting_end = member.created_at + waiting_period_days`
  - If `claim.service_date < waiting_end` -> REJECT
- Uses `member.created_at` as the enrollment date (when the member was added to the policy)
- Example: Benefit has 30-day waiting period. Member enrolled Jan 1. Claim with service date Jan 15 -> REJECTED

**Files:** `services/claims/adjudicator_service_impl.go` (lines 118-135)

#### Rule 6: Exclusion Check (coverage)
- Fetches all exclusions for the policy's plan via `exclusionRepo.ListByPlan`
- Parses the claim's `diagnosis_codes` JSON array
- For each exclusion, parses its `icd_codes` JSON array
- If any claim diagnosis code matches any exclusion ICD code -> REJECT
- Example: Plan excludes cosmetic procedures with ICD codes `["Z41.1", "Z41.8"]`. Claim has diagnosis code `Z41.1` -> REJECTED

**Files:** `services/claims/adjudicator_service_impl.go` (lines 137-157)

#### Rule 7: Annual Limits & Co-Pay Calculation (limits)
- For each benefit in the plan:
  - Queries total approved amount this year for the member + benefit category
  - `remaining = annual_limit - used`
  - If `remaining <= 0` -> skip this benefit (limit exhausted)
  - If `payable_amount > remaining` -> cap at remaining (partial approval)
  - Apply co-pay:
    - `percentage`: `co_pay = payable * co_pay_value / 100`
    - `fixed`: `co_pay = co_pay_value`
  - `payable_amount -= co_pay_amount`
  - Break after first matching benefit

**Amount calculation:**
```
payable_amount      = min(total_amount, remaining_annual_limit) - co_pay
co_pay_amount       = calculated from benefit co-pay rules
member_responsibility = co_pay + any excess above annual limit
```

#### Rule 8: Duplicate Check (fraud)
- Calls `fraudSvc.CheckDuplicate(claimNumber, claimID)`
- Queries fraud_flags table for existing flags with same claim number but different claim ID
- If duplicate found -> `MANUAL_REVIEW` (does NOT auto-reject, sends to human reviewer)

### Stage 5-6: Store Decision & Update Claim

The adjudication result is:
1. Stored as an `AdjudicationDecision` record (with JSON rule results)
2. Mapped to a claim status:
   - `APPROVE` -> claim status `ADJUDICATED`
   - `REJECT` -> claim status `REJECTED` (with rejection reason from failed rules)
   - `MANUAL_REVIEW` -> claim status `MANUAL_REVIEW`
3. Claim amounts updated: `approved_amount`, `co_pay_amount`, `member_responsibility`

### Stage 7: Post-Adjudication Fraud Checks

These run **after** adjudication and create `FraudFlag` records. They do NOT change the claim status — they are informational flags visible when fetching the claim.

#### Frequency Check
- `fraudSvc.CheckFrequency(memberID, providerID, procedureCode, claimID)`
- Queries fraud_flags table for how many times this member+provider+procedure combo has been claimed
- If count > 0 -> creates a `FREQUENCY` flag with `MEDIUM` severity
- Purpose: catch patterns like a member visiting the same provider for the same procedure repeatedly

#### Amount Threshold Check
- `fraudSvc.CheckAmountThreshold(providerID, procedureCode, totalAmount)`
- Simple threshold: flags if `total_amount > 50,000,000` (500,000 KES)
- If exceeded -> creates an `AMOUNT_THRESHOLD` flag with `HIGH` severity
- Purpose: catch abnormally large claims

**Files:** `services/claims/fraud_service_impl.go`, `services/claims/claim_service_impl.go` (lines 167-190)

### Stage 8: Response

The claim is re-fetched from DB to include all updated amounts and status. The response includes:
```json
{
  "id": "...",
  "claim_number": "CLM-2026-000001",
  "status": "ADJUDICATED",
  "total_amount": 1000000,
  "approved_amount": 800000,
  "co_pay_amount": 200000,
  "member_responsibility": 200000,
  ...
}
```

When fetched via `GET /claims/:id`, additional nested data is included:
- `line_items[]` - individual procedures
- `decision` - full adjudication decision with rule results
- `fraud_flags[]` - any fraud flags raised

## Post-Pipeline Actions

### Approve Claim (`PUT /claims/:id/approve`)
- Requires claim status: `ADJUDICATED` or `MANUAL_REVIEW`
- Uses amounts from the stored adjudication decision
- Sets status to `APPROVED`

### Reject Claim (`PUT /claims/:id/reject`)
- Requires claim status: `RECEIVED`, `VALIDATED`, `ADJUDICATED`, or `MANUAL_REVIEW`
- Cannot reject `APPROVED` or `PAID` claims
- Requires a `reason` in the request body

## Claim Status Lifecycle

```
RECEIVED -> VALIDATED -> ADJUDICATED -> APPROVED -> PAID
                |              |
                |              +--> MANUAL_REVIEW -> APPROVED
                |              |                  -> REJECTED
                |              +--> REJECTED
                +--> REJECTED
```

| Status | Meaning |
|--------|---------|
| `RECEIVED` | Claim created, not yet processed |
| `VALIDATED` | Passed pre-adjudication validation |
| `ADJUDICATED` | Adjudication complete, auto-approved by rules, awaiting human approval |
| `MANUAL_REVIEW` | Flagged for human review (e.g., duplicate detected) |
| `APPROVED` | Human-approved, ready for payment |
| `REJECTED` | Rejected at any stage (validation, adjudication, or manual) |
| `PAID` | Payment processed (set by billing/remittance) |

## Key Files Reference

| File | Purpose |
|------|---------|
| `services/claims/claim_service_impl.go` | Orchestrates the full pipeline in SubmitClaim |
| `services/claims/validator_service_impl.go` | Pre-adjudication validation rules |
| `services/claims/adjudicator_service_impl.go` | Business rules engine (eligibility, coverage, limits, fraud) |
| `services/claims/fraud_service_impl.go` | Fraud check implementations (duplicate, frequency, threshold) |
| `domains/claims/entity/adjudication_decision.go` | AdjudicationDecision + AdjudicationResult + RuleResult types |
| `domains/claims/entity/fraud_flag.go` | FraudFlag entity |
| `domains/claims/repository/claim_repository.go` | ClaimRepository interface (UpdateStatus, UpdateAmounts, Reject) |
| `domains/claims/repository/adjudication_repository.go` | AdjudicationRepository interface |

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

## Policy Status Guards

| Action | Required Status | Endpoint |
|--------|----------------|----------|
| Activate | `DRAFT` | `PUT /policies/:id/activate` |
| Lapse | `ACTIVE` | `PUT /policies/:id/lapse` |
| Terminate | `ACTIVE` or `LAPSED` | `PUT /policies/:id/terminate` |
| Reinstate | `LAPSED` | `PUT /policies/:id/reinstate` |