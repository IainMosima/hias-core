# Endorsement Workflow — Implementation Guide

## Overview

The endorsement system enables mid-term policy modifications through a controlled approval workflow. Endorsements support four operation types: adding members, removing members, updating member details, and changing plans.

---

## Lifecycle & Status Flow

```
PENDING ──→ APPROVED ──→ APPLIED
   │            │
   ├──→ REJECTED│
   │            │
   └──→ CANCELLED
                └──→ CANCELLED
```

| Status    | Terminal | Description                              |
|-----------|----------|------------------------------------------|
| PENDING   | No       | Created, awaiting approval               |
| APPROVED  | No       | Approved, ready to apply                 |
| APPLIED   | Yes      | Changes executed on policy/members       |
| REJECTED  | Yes      | Denied by approver                       |
| CANCELLED | Yes      | Withdrawn before application             |

---

## API Endpoints

### Nested under policies

| Method | Endpoint                                | Description              |
|--------|----------------------------------------|--------------------------|
| POST   | `/api/v1/policies/:id/endorsements`    | Create endorsement       |
| GET    | `/api/v1/policies/:id/endorsements`    | List policy endorsements |

### Standalone

| Method | Endpoint                              | Description         |
|--------|---------------------------------------|---------------------|
| GET    | `/api/v1/endorsements/:id`            | Get endorsement     |
| PUT    | `/api/v1/endorsements/:id/approve`    | Approve (PENDING)   |
| PUT    | `/api/v1/endorsements/:id/reject`     | Reject (PENDING)    |
| PUT    | `/api/v1/endorsements/:id/apply`      | Apply (APPROVED)    |
| PUT    | `/api/v1/endorsements/:id/cancel`     | Cancel (PENDING/APPROVED) |

---

## Creating an Endorsement

### Request

```
POST /api/v1/policies/:id/endorsements
```

```json
{
  "endorsement_type": "ADD_MEMBER | REMOVE_MEMBER | UPDATE_MEMBER | PLAN_CHANGE",
  "effective_date": "2026-04-01",
  "changes": { ... },
  "reason": "Optional reason text",
  "premium_adjustment": 0
}
```

### Validation Rules (applied on create)

**General (all types):**
- Policy must be ACTIVE
- `effective_date` must be within policy term (`start_date` to `end_date`)
- `effective_date` cannot be more than 90 days in the past

**Per-type payload validation** — see next section.

---

## Changes Payload Per Type

### ADD_MEMBER

Reuses `EnrollMemberRequest`. All standard member enrollment validations apply.

```json
{
  "name": "Jane Doe",
  "date_of_birth": "1990-05-15",
  "gender": "female",
  "relationship": "spouse",
  "national_id": "12345678",
  "phone": "+254700000000",
  "email": "jane@example.com",
  "kra_pin": "A0123456B",
  "county": "Nairobi",
  "address": "P.O. Box 123"
}
```

| Field          | Required | Validation                                      |
|----------------|----------|-------------------------------------------------|
| name           | Yes      | Non-empty string                                |
| date_of_birth  | Yes      | Format: `YYYY-MM-DD`, must be valid date        |
| gender         | Yes      | `male`, `female`, or `other`                    |
| relationship   | Yes      | `principal`, `spouse`, `child`, or `parent`     |
| national_id    | No       |                                                 |
| phone          | No       |                                                 |
| email          | No       |                                                 |
| kra_pin        | No       |                                                 |
| county         | No       |                                                 |
| address        | No       |                                                 |

### REMOVE_MEMBER

```json
{
  "member_id": "uuid",
  "reason": "Optional removal reason"
}
```

| Field     | Required | Validation                                          |
|-----------|----------|-----------------------------------------------------|
| member_id | Yes      | Valid UUID, member exists, belongs to policy, not REMOVED |
| reason    | No       |                                                     |

### UPDATE_MEMBER

```json
{
  "member_id": "uuid",
  "updates": {
    "name": "Updated Name",
    "phone": "+254711111111",
    "email": "new@email.com",
    "kra_pin": "B9876543A",
    "county": "Mombasa",
    "address": "P.O. Box 456"
  }
}
```

| Field     | Required | Validation                                                 |
|-----------|----------|------------------------------------------------------------|
| member_id | Yes      | Valid UUID, member exists, belongs to policy, not REMOVED  |
| updates   | Yes      | At least one field must be non-nil                         |

### PLAN_CHANGE

```json
{
  "new_plan_id": "uuid",
  "reason": "Upgrading to premium plan"
}
```

| Field       | Required | Validation                                         |
|-------------|----------|----------------------------------------------------|
| new_plan_id | Yes      | Valid UUID, plan exists, is ACTIVE, differs from current |
| reason      | No       |                                                    |

---

## Approve / Reject

### Approve

```
PUT /api/v1/endorsements/:id/approve
```

No body required. Sets `approved_by` and `approved_at` from the authenticated user.

### Reject

```
PUT /api/v1/endorsements/:id/reject
```

```json
{
  "reason": "Reason for rejection"
}
```

---

## Apply Endorsement

```
PUT /api/v1/endorsements/:id/apply
```

Only APPROVED endorsements can be applied. The apply step:

1. **Calculates premium adjustment** (if not manually set) — pro-rated for remaining policy days
2. **Executes the change** based on type:
   - ADD_MEMBER: Enrolls member, sets `coverage_start_date = effective_date`
   - REMOVE_MEMBER: Removes member, sets `coverage_end_date = effective_date`
   - UPDATE_MEMBER: Updates member fields
   - PLAN_CHANGE: Changes the policy plan
3. **Applies premium adjustment** to policy (skipped for REMOVE_MEMBER since `RemoveMember` handles it internally)
4. Sets status to APPLIED with `applied_at` timestamp
5. Triggers document generation asynchronously

### Coverage Dates

On apply, endorsements set coverage dates on affected members:

| Endorsement Type | coverage_start_date    | coverage_end_date      |
|------------------|------------------------|------------------------|
| ADD_MEMBER       | = effective_date       | (unchanged)            |
| REMOVE_MEMBER    | (unchanged)            | = effective_date       |
| UPDATE_MEMBER    | (unchanged)            | (unchanged)            |
| PLAN_CHANGE      | (unchanged)            | (unchanged)            |

These dates drive claim eligibility — see "Claims Impact" below.

---

## Cancel Endorsement

```
PUT /api/v1/endorsements/:id/cancel
```

```json
{
  "reason": "No longer needed"
}
```

Allowed from PENDING or APPROVED status only. APPLIED, REJECTED, and CANCELLED endorsements cannot be cancelled.

---

## Claims Impact

The claim validator now enforces member coverage dates:

1. **Member must be ACTIVE** — claims for PENDING, SUSPENDED, or REMOVED members are rejected
2. **Service date must be within coverage period:**
   - If `coverage_start_date` is set and `service_date < coverage_start_date` → rejected
   - If `coverage_end_date` is set and `service_date > coverage_end_date` → rejected

This means:
- A member added via endorsement with `effective_date = March 1` is covered from March 1. Claims with `service_date` before March 1 are rejected.
- A member removed via endorsement with `effective_date = March 15` has coverage until March 15. Claims before that date remain valid; claims after are rejected.

---

## Premium Adjustment

### Auto-calculation

If `premium_adjustment` is not set (or 0) on create, the system calculates it at apply time:

1. Gets current active member count
2. Calculates new count based on endorsement type (+1 for add, -1 for remove)
3. Calls `PremiumRuleService.CalculatePremiumWithMembers()` with new count
4. Computes full delta (new premium - current premium)
5. Pro-rates for remaining days: `adjustment = delta * (remaining_days / total_days)`

### Double Adjustment Fix

Previously, REMOVE_MEMBER endorsements would adjust premium twice:
- Once via `RemoveMember()` (which internally recalculates and creates credit notes)
- Again via the endorsement premium adjustment block

Now the endorsement-level premium adjustment is **skipped for REMOVE_MEMBER** since the downstream service handles it.

---

## Audit Trail

Every endorsement state change is logged via the audit system:

| Action       | Audit Action   | Logged By      |
|--------------|----------------|----------------|
| Create       | CREATE         | requestedBy    |
| Approve      | STATE_CHANGE   | approvedBy     |
| Reject       | STATE_CHANGE   | rejectedBy     |
| Apply        | STATE_CHANGE   | approvedBy     |
| Cancel       | STATE_CHANGE   | cancelledBy    |

---

## Database Migration

**Migration 000030: `member_coverage_dates`**

```sql
ALTER TABLE members ADD COLUMN coverage_start_date TIMESTAMPTZ;
ALTER TABLE members ADD COLUMN coverage_end_date TIMESTAMPTZ;

-- Backfill
UPDATE members SET coverage_start_date = created_at WHERE status IN ('ACTIVE', 'PENDING');
UPDATE members SET coverage_end_date = updated_at WHERE status = 'REMOVED';
```

---

## Files Changed

| File | Change |
|------|--------|
| `shared/types.go` | Added `EndorsementStatusCancelled` |
| `domains/policy/entity/member.go` | Added `CoverageStartDate`, `CoverageEndDate` |
| `domains/policy/repository/member_repository.go` | Added `UpdateCoverageDates` interface method |
| `domains/policy/schema/policy_request.go` | Added `RemoveMemberChanges`, `UpdateMemberChanges`, `CancelEndorsementRequest` |
| `domains/policy/schema/policy_response.go` | Added coverage dates to `MemberResponse` + mapper |
| `domains/policy/service/endorsement_service.go` | Added `CancelEndorsement` to interface |
| `services/policy/endorsement_service_impl.go` | Validation, premium fix, coverage dates, cancel, doc gen |
| `services/claims/validator_service_impl.go` | Member status + coverage date checks |
| `services/api-gateway/handlers/endorsement_handler.go` | Added `CancelEndorsement` handler |
| `services/api-gateway/routes/routes.go` | Added cancel route |
| `services/api-gateway/main.go` | Wired `planRepo` + `policyDocSvc` into endorsement service |
| `infrastructures/db/migration/000030_*` | Coverage date columns migration |
| `infrastructures/db/query/member.sql` | Added `UpdateMemberCoverageDates` query |
| `infrastructures/db/sqlc/*` | Regenerated (via `make sqlc`) |
| `infrastructures/repository/member_repository.go` | Coverage date mapper + `UpdateCoverageDates` impl |

---

## Frontend Integration Notes

Use typed forms per endorsement type — not raw JSON:

| Type          | Form Fields                                                                 |
|---------------|-----------------------------------------------------------------------------|
| ADD_MEMBER    | Name, DOB (date picker), gender (dropdown), relationship (dropdown), national_id, phone, email, kra_pin, county, address |
| REMOVE_MEMBER | Member picker (ACTIVE members on policy) + optional reason text             |
| UPDATE_MEMBER | Member picker → pre-populated editable fields (name, phone, email, etc.)    |
| PLAN_CHANGE   | Plan picker (active plans, exclude current) + optional reason text          |

Additional notes:
- `premium_adjustment`: Display as read-only — backend calculates on apply. Remove from create form.
- `effective_date`: Date picker, default to today
- `reason`: Optional text field for all types
- After apply, re-fetch policy and member data to reflect changes
- Coverage dates now appear in `MemberResponse` as `coverage_start_date` and `coverage_end_date`
