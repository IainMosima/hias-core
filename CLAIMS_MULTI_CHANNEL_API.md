# Claims API — Full Frontend Integration Guide

Complete reference for every claims endpoint, the full lifecycle, and how billing/payment ties in.

---

## Claim Lifecycle

```
                    ┌─────────────────────────────────────────────────────────┐
                    │              SUBMISSION CHANNELS                        │
                    │                                                         │
  Staff ──────────►│ POST /api/v1/claims              (PASETO)               │
  Staff ──────────►│ POST /api/v1/claims/drafts/:id/submit (PASETO)          │
  Provider Portal ►│ POST /api/v1/external/claims     (API Key)              │
  CSV Import ─────►│ POST /api/v1/claims/import-csv   (PASETO)               │
                    └──────────────────┬──────────────────────────────────────┘
                                       ▼
                              ┌─────────────────┐
                              │    RECEIVED      │  ← claim created
                              └────────┬────────┘
                                       ▼
                              ┌─────────────────┐
                              │   VALIDATED      │  ← eligibility + coverage checked
                              └────────┬────────┘
                                       ▼
                              ┌─────────────────┐
                              │  ADJUDICATED     │  ← auto-decision (approve/reject/review)
                              └────────┬────────┘
                                       │
                    ┌──────────────────┼──────────────────┐
                    ▼                  ▼                  ▼
           ┌──────────────┐  ┌────────────────┐  ┌──────────────┐
           │   APPROVED   │  │ MANUAL_REVIEW  │  │   REJECTED   │
           └──────┬───────┘  │  / ESCALATED   │  └──────────────┘
                  │          └───────┬────────┘         │
                  │                  │ (approve/reject) │
                  │                  ▼                  │
                  │          back to APPROVED           │
                  │          or REJECTED                │
                  ▼                                     │
           ┌──────────────┐                             │
           │    VETTED     │  ← claims officer reviews  │
           │ PARTIALLY_VETTED                           │
           └──────┬───────┘                             │
                  ▼                                     │
        ┌───────────────────┐                           │
        │ READY_FOR_PAYMENT │  ← finance confirms      │
        └────────┬──────────┘                           │
                 ▼                                      │
        ┌────────────────┐                              │
        │  PAID / PART_PAID │ ← money sent to provider │
        └────────────────┘                              │
                                                        │
        ┌───── Remittance ◄─── groups approved claims ──┘
        │      by provider into payment batches
        ▼
   Payment initiated (MPESA / Bank Transfer)
        │
        ▼
   Provider receives funds
```

---

## Statuses Reference

| Status | Meaning | Who transitions |
|--------|---------|-----------------|
| `RECEIVED` | Claim created, pipeline starting | System |
| `VALIDATED` | Passed eligibility + coverage checks | System |
| `ADJUDICATED` | Auto-decision made | System |
| `APPROVED` | Approved for payment | System or Admin/Manager |
| `REJECTED` | Denied | System or Admin/Manager |
| `MANUAL_REVIEW` | Needs human review (fraud flags, high amount) | System |
| `ESCALATED` | Escalated to a specific role | System |
| `VETTED` | Claims officer reviewed, amounts confirmed | Admin/ClaimsOfficer |
| `PARTIALLY_VETTED` | Partial vetting done | Admin/ClaimsOfficer |
| `READY_FOR_PAYMENT` | Finance approved for payout | Admin/Finance |
| `PAID` | Fully paid to provider | Admin/Finance |
| `PART_PAID` | Partially paid | Admin/Finance |

---

## 1. Submit Claim (Internal Staff)

```
POST /api/v1/claims
Authorization: Bearer <token>
```

**Request:**

```json
{
  "policy_id": "uuid",
  "member_id": "uuid",
  "provider_id": "uuid",
  "preauth_id": "uuid (optional)",
  "diagnosis_codes": ["J06.9", "R50.9"],
  "service_date": "2026-03-10T00:00:00Z",
  "admission_date": "2026-03-09T00:00:00Z",
  "discharge_date": "2026-03-11T00:00:00Z",
  "notes": "Emergency visit",
  "claim_type": "DIRECT",
  "line_items": [
    {
      "procedure_code": "99213",
      "procedure_name": "Office visit - Level 3",
      "diagnosis_code": "J06.9",
      "quantity": 1,
      "unit_price": 250000
    }
  ]
}
```

| Field | Required | Notes |
|-------|----------|-------|
| `policy_id` | Yes | UUID |
| `member_id` | Yes | UUID |
| `provider_id` | Yes | UUID |
| `preauth_id` | No | Links to pre-authorization if one exists |
| `diagnosis_codes` | Yes | Array of ICD-10 codes |
| `service_date` | Yes | RFC3339 |
| `admission_date` | No | For inpatient claims |
| `discharge_date` | No | For inpatient claims |
| `claim_type` | No | `DIRECT` (default), `REIMBURSEMENT`, `CREDIT`, `EXCEPTION` |
| `line_items` | Yes | Min 1. `unit_price` in cents |

**Response** `200` — claim goes through the full pipeline instantly:

```json
{
  "status": "success",
  "message": "Claim submitted successfully",
  "data": {
    "id": "uuid",
    "claim_number": "CLM-2026-000123",
    "policy_id": "uuid",
    "member_id": "uuid",
    "provider_id": "uuid",
    "status": "ADJUDICATED",
    "total_amount": 250000,
    "approved_amount": 250000,
    "co_pay_amount": 0,
    "member_responsibility": 0,
    "diagnosis_codes": ["J06.9", "R50.9"],
    "service_date": "2026-03-10T00:00:00Z",
    "notes": "Emergency visit",
    "claim_type": "DIRECT",
    "claim_source": "INTERNAL",
    "is_draft": false,
    "sla_breach_at": "2026-03-17T00:00:00Z",
    "line_items": [
      {
        "id": "uuid",
        "procedure_code": "99213",
        "procedure_name": "Office visit - Level 3",
        "diagnosis_code": "J06.9",
        "quantity": 1,
        "unit_price": 250000,
        "total_price": 250000,
        "approved_amount": 250000
      }
    ],
    "decision": {
      "decision": "APPROVE",
      "payable_amount": 250000,
      "member_responsibility": 0,
      "deductible_applied": 0,
      "co_pay_applied": 0,
      "sub_limit_applied": 0,
      "benefit_category": "outpatient",
      "reasons": ["..."],
      "rule_results": ["..."],
      "adjudicated_at": "2026-03-12T10:00:01Z"
    },
    "fraud_flags": [],
    "created_at": "2026-03-12T10:00:00Z",
    "updated_at": "2026-03-12T10:00:01Z"
  }
}
```

**What happens behind the scenes (pipeline):**
1. Claim created with status `RECEIVED`
2. Validation: checks policy active, member active, provider valid
3. Fraud detection: frequency, amount threshold, expired contract, suspended provider, rate card overcharge, repeat visit
4. Adjudication: benefit limits, co-pay, deductible, sub-limits, exclusions
5. Decision: `APPROVE` / `REJECT` / `MANUAL_REVIEW`
6. Escalation rules evaluated (high-amount → escalate to Manager)
7. Event published to queue
8. Pre-auth updated to `CLAIMED` if linked
9. Timeline entry recorded
10. Audit log created

---

## 2. Get Claim

```
GET /api/v1/claims/:id
Authorization: Bearer <token>
```

**Response** `200`: Full claim object (same shape as submit response).

---

## 3. List Claims

```
GET /api/v1/claims
Authorization: Bearer <token>
```

**Query params:**

| Param | Type | Notes |
|-------|------|-------|
| `status` | string | Filter by status: `RECEIVED`, `APPROVED`, `PAID`, etc. |
| `provider` | UUID | Filter by provider |
| `date_from` | string | RFC3339 or `YYYY-MM-DD` |
| `date_to` | string | RFC3339 or `YYYY-MM-DD` |
| `search` | string | Search claim number or status |
| `page` | int | Default 1 |
| `page_size` | int | Default 20 |

**Response** `200` (paginated):

```json
{
  "status": "success",
  "message": "Claims retrieved",
  "data": [ { ... claim objects ... } ],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 150
  }
}
```

---

## 4. Claim Review & Approval Flow

### Approve Claim

```
PUT /api/v1/claims/:id/approve
Authorization: Bearer <token>
```

**Roles:** Admin, Manager

No request body. The system checks the approver's approval limit — if the claim amount exceeds their limit, it's rejected.

**Response** `200`: Updated claim with `status: "APPROVED"`.

### Reject Claim

```
PUT /api/v1/claims/:id/reject
Authorization: Bearer <token>
```

**Roles:** Admin, Manager

```json
{
  "decision": "REJECTED",
  "reason": "Duplicate claim for same service date"
}
```

**Response** `200`: Updated claim with `status: "REJECTED"`, `rejection_reason` populated.

### Vet Claim

Claims officer reviews amounts before payment.

```
PUT /api/v1/claims/:id/vet
Authorization: Bearer <token>
```

**Roles:** Admin, ClaimsOfficer

```json
{
  "vetted_amount": 230000,
  "notes": "Reduced — overcharged for lab work"
}
```

**Response** `200`: Claim with `status: "VETTED"` or `"PARTIALLY_VETTED"`, `vetted_amount`, `vetted_by`, `vetted_at` populated.

- If `vetted_amount == approved_amount` → `VETTED`
- If `vetted_amount < approved_amount` → `PARTIALLY_VETTED`
- Type-specific rules apply (Direct, Reimbursement, Exception each have their own vetting logic)

---

## 5. Payment Flow

### Mark Ready for Payment

Finance confirms the claim is cleared for payout.

```
PUT /api/v1/claims/:id/ready-for-payment
Authorization: Bearer <token>
```

**Roles:** Admin, Finance

**Response** `200`: Claim with `status: "READY_FOR_PAYMENT"`.

### Mark Paid

```
PUT /api/v1/claims/:id/mark-paid
Authorization: Bearer <token>
```

**Roles:** Admin, Finance

**Response** `200`: Claim with `status: "PAID"`.

### Mark Part Paid

```
PUT /api/v1/claims/:id/mark-part-paid
Authorization: Bearer <token>
```

**Roles:** Admin, Finance

**Response** `200`: Claim with `status: "PART_PAID"`.

---

## 6. How Payment Actually Happens (Billing Integration)

After claims are approved, the billing side takes over:

### Remittances — Batching Claims for Provider Payment

```
POST   /api/v1/remittances              ← Create remittance (groups APPROVED claims by provider)
GET    /api/v1/remittances              ← List all remittances
GET    /api/v1/remittances/:id          ← Get remittance details
GET    /api/v1/remittances/:id/export   ← Export payment file (bank/MPESA format)
```

**Create Remittance** — bundles all approved claims for a specific provider:

```json
{
  "provider_id": "uuid"
}
```

**Response:**

```json
{
  "data": {
    "id": "uuid",
    "provider_id": "uuid",
    "total_amount": 1500000,
    "claim_count": 5,
    "status": "PENDING",
    "created_at": "..."
  }
}
```

Remittance statuses: `PENDING` → `PROCESSING` → `SENT` → `CONFIRMED` / `FAILED`

### Payments — Actual Money Transfer

```
POST   /api/v1/payments                 ← Initiate payment (MPESA or bank transfer)
GET    /api/v1/payments                 ← List payments
GET    /api/v1/payments/:id             ← Get payment
PUT    /api/v1/payments/:id/retry       ← Retry failed payment
PUT    /api/v1/payments/:id/reconcile   ← Reconcile with bank statement
```

**Initiate Payment:**

```json
{
  "remittance_id": "uuid",
  "payment_method": "MPESA",
  "reference": "PROV-PAY-001"
}
```

Payment statuses: `INITIATED` → `PROCESSING` → `CONFIRMED` / `FAILED` → `RECONCILED`

### MPESA Webhook (Public)

```
POST /api/v1/webhooks/mpesa
```

Receives payment confirmation from MPESA. Auto-reconciles.

### Provider Statements — Reconciliation

```
GET    /api/v1/providers/:id/statements          ← List statements for provider
POST   /api/v1/providers/:id/statements          ← Upload provider statement
GET    /api/v1/provider-statements/:id           ← Get statement details
GET    /api/v1/provider-statements/:id/line-items ← Line items
POST   /api/v1/provider-statements/:id/reconcile ← Reconcile against claims
```

---

## 7. Claim Timeline

Every status change is recorded as a timeline entry.

```
GET /api/v1/claims/:id/timeline
Authorization: Bearer <token>
```

**Response** `200`:

```json
{
  "data": [
    {
      "id": "uuid",
      "action": "SUBMITTED",
      "from_status": "",
      "to_status": "RECEIVED",
      "performed_by": "uuid",
      "performed_by_name": "John Doe",
      "notes": "",
      "created_at": "2026-03-12T10:00:00Z"
    },
    {
      "id": "uuid",
      "action": "ADJUDICATED",
      "from_status": "VALIDATED",
      "to_status": "ADJUDICATED",
      "performed_by": "00000000-0000-0000-0000-000000000000",
      "performed_by_name": "System",
      "notes": "Auto-approved by adjudication engine",
      "created_at": "2026-03-12T10:00:01Z"
    },
    {
      "id": "uuid",
      "action": "APPROVED",
      "from_status": "ADJUDICATED",
      "to_status": "APPROVED",
      "performed_by": "uuid",
      "performed_by_name": "Jane Manager",
      "notes": "",
      "created_at": "2026-03-12T14:30:00Z"
    },
    {
      "id": "uuid",
      "action": "PAID",
      "from_status": "READY_FOR_PAYMENT",
      "to_status": "PAID",
      "performed_by": "uuid",
      "performed_by_name": "Finance Officer",
      "notes": "",
      "created_at": "2026-03-15T09:00:00Z"
    }
  ]
}
```

---

## 8. Claim Documents

### Upload Document

```
POST /api/v1/claims/:id/documents
Authorization: Bearer <token>
```

```json
{
  "file_name": "invoice.pdf",
  "file_type": "application/pdf",
  "file_size": 102400,
  "s3_key": "claims/uuid/invoice.pdf"
}
```

**Response** `200`:

```json
{
  "data": {
    "id": "uuid",
    "claim_id": "uuid",
    "file_name": "invoice.pdf",
    "file_type": "application/pdf",
    "file_size": 102400,
    "s3_key": "claims/uuid/invoice.pdf",
    "uploaded_by": "uuid",
    "created_at": "2026-03-12T10:00:00Z"
  }
}
```

### List Documents

```
GET /api/v1/claims/:id/documents
```

### Delete Document

```
DELETE /api/v1/claim-documents/:id
```

### Generate Decline Letter (for rejected claims)

```
POST /api/v1/claims/:id/decline-letter
```

Generates a PDF decline letter and stores it as a policy document.

---

## 9. SLA Monitoring

```
GET /api/v1/claims/sla-breached?page=1&page_size=20
Authorization: Bearer <token>
```

Returns claims past their SLA deadline that aren't yet resolved. Each claim has `sla_breach_at` — the system monitors this with a scheduled task every 15 minutes.

---

## 10. Bulk Operations

### Bulk Submit

```
POST /api/v1/claims/bulk
Authorization: Bearer <token>
```

```json
{
  "claims": [
    { "policy_id": "...", "member_id": "...", ... },
    { "policy_id": "...", "member_id": "...", ... }
  ]
}
```

**Response** `200`:

```json
{
  "data": {
    "succeeded": 8,
    "failed": 2,
    "claims": [ ... successful claims ... ],
    "errors": [ "Row 3: Member not active", "Row 7: Invalid provider" ]
  }
}
```

### CSV Import

```
POST /api/v1/claims/import-csv
Authorization: Bearer <token>
Content-Type: multipart/form-data
```

**Roles:** Admin, ClaimsOfficer

Form field: `file` (CSV)

Same response shape as bulk submit.

---

## 11. Draft Claims (Save & Complete Later)

All under PASETO auth. Drafts let staff save incomplete claims.

### Create Draft

```
POST /api/v1/claims/drafts
```

All fields optional:

```json
{
  "policy_id": "uuid",
  "member_id": "uuid",
  "provider_id": "uuid",
  "preauth_id": "uuid",
  "diagnosis_codes": ["J06.9"],
  "service_date": "2026-03-10T00:00:00Z",
  "notes": "Waiting for provider details",
  "claim_type": "DIRECT",
  "line_items": [
    { "procedure_code": "99213", "procedure_name": "Visit", "quantity": 1, "unit_price": 250000 }
  ]
}
```

**Response** `201`: Claim object with `is_draft: true`, `claim_number: "DRF-2026-XXXXXX"`.

### List My Drafts

```
GET /api/v1/claims/drafts?page=1&page_size=20
```

Returns only the authenticated user's drafts.

### Update Draft

```
PUT /api/v1/claims/drafts/:id
```

Same body as create. Only works while `is_draft = true`.

### Submit Draft (finalize)

```
POST /api/v1/claims/drafts/:id/submit
```

No body. Requires `policy_id`, `member_id`, `provider_id` to be set. Runs the full pipeline. Draft is deleted and a real claim is created.

**Error** `400` if required fields missing:

```json
{ "status": "error", "message": "Policy ID is required to submit" }
```

### Delete Draft

```
DELETE /api/v1/claims/drafts/:id
```

Only works while `is_draft = true`.

---

## 12. External Claims (Provider Portal / Partner API)

For external systems submitting claims via API key instead of PASETO.

### Authentication

```
X-API-Key: <64-char hex>
X-API-Secret: <64-char hex>
```

Obtained when an admin creates an API partner (see section 13).

### Submit External Claim

```
POST /api/v1/external/claims
```

```json
{
  "idempotency_key": "prov-claim-2026-001",
  "external_claim_id": "HOSP-INV-12345",
  "member_number": "MBR-2026-000001",
  "diagnosis_codes": ["J06.9"],
  "service_date": "2026-03-10T00:00:00Z",
  "admission_date": "2026-03-09T00:00:00Z",
  "discharge_date": "2026-03-11T00:00:00Z",
  "claim_type": "DIRECT",
  "notes": "Emergency admission",
  "line_items": [
    { "procedure_code": "99223", "procedure_name": "Hospital admission", "quantity": 1, "unit_price": 5000000 },
    { "procedure_code": "36415", "procedure_name": "Blood draw", "quantity": 2, "unit_price": 150000 }
  ],
  "metadata": { "hospital_ref": "REF-001" }
}
```

Key differences from internal submit:
- Uses `member_number` (not `member_id`) — system resolves to internal ID
- Uses `idempotency_key` — prevents duplicates, same key returns existing claim
- Provider resolved from the API partner's linked provider
- Policy resolved from the member's active policy
- `claim_source` set to `PROVIDER_PORTAL` or `PARTNER_API`

**Response** `201`:

```json
{
  "data": {
    "claim_id": "uuid",
    "claim_number": "CLM-2026-000124",
    "status": "RECEIVED",
    "received_at": "2026-03-12T10:30:00Z"
  }
}
```

**Idempotent retry** `200` — same `idempotency_key`:

```json
{
  "message": "Claim already submitted (idempotent)",
  "data": { "claim_id": "uuid", "status": "PROCESSING", ... }
}
```

### Get External Claim Status

```
GET /api/v1/external/claims/:id/status
```

Scoped to the partner's provider — can't query other providers' claims.

**Response** `200`:

```json
{
  "data": {
    "claim_id": "uuid",
    "claim_number": "CLM-2026-000124",
    "status": "PROCESSING",
    "total_amount": 5300000,
    "approved_amount": 0,
    "updated_at": "2026-03-12T11:00:00Z"
  }
}
```

### External Status Mapping

External API uses simplified statuses:

| External Status | Internal Statuses |
|----------------|-------------------|
| `RECEIVED` | RECEIVED, VALIDATED |
| `PROCESSING` | ADJUDICATED, VETTED, PARTIALLY_VETTED |
| `UNDER_REVIEW` | MANUAL_REVIEW, ESCALATED |
| `APPROVED` | APPROVED, READY_FOR_PAYMENT |
| `REJECTED` | REJECTED |
| `SETTLED` | PAID, PART_PAID |

---

## 13. API Partner Admin (Admin only)

Manage API partners that access the external claims API.

### Create Partner

```
POST /api/v1/api-partners
Authorization: Bearer <admin-token>
```

```json
{
  "name": "City Hospital Portal",
  "partner_type": "PROVIDER",
  "provider_id": "uuid-of-provider",
  "rate_limit_per_minute": 60,
  "allowed_claim_types": ["DIRECT"],
  "webhook_url": "https://hospital.example.com/webhooks",
  "contact_email": "it@hospital.example.com"
}
```

| Field | Required | Notes |
|-------|----------|-------|
| `partner_type` | Yes | `PROVIDER`, `PARTNER_NETWORK`, `TPA` |
| `provider_id` | No | Links partner to a provider (required for provider portals) |

**Response** `201` — API secret shown only once:

```json
{
  "data": {
    "id": "uuid",
    "name": "City Hospital Portal",
    "api_key": "a1b2c3d4...64chars",
    "api_secret": "e5f6g7h8...64chars",
    "is_active": true,
    ...
  }
}
```

### Other Endpoints

```
GET    /api/v1/api-partners               ← List partners
GET    /api/v1/api-partners/:id           ← Get partner
PUT    /api/v1/api-partners/:id           ← Update partner
PUT    /api/v1/api-partners/:id/deactivate ← Deactivate (blocks API access)
POST   /api/v1/api-partners/:id/regenerate-key ← New API key + secret
```

---

## 14. Claim Response Object (Full Shape)

Every claim endpoint returns this shape:

```json
{
  "id": "uuid",
  "claim_number": "CLM-2026-000123",
  "policy_id": "uuid",
  "member_id": "uuid",
  "provider_id": "uuid",
  "status": "ADJUDICATED",
  "total_amount": 250000,
  "approved_amount": 250000,
  "co_pay_amount": 0,
  "member_responsibility": 0,
  "diagnosis_codes": ["J06.9"],
  "service_date": "2026-03-10T00:00:00Z",
  "notes": "",
  "claim_type": "DIRECT",
  "vetted_amount": null,
  "vetted_by": null,
  "vetted_at": null,
  "sla_breach_at": "2026-03-17T00:00:00Z",
  "rejection_reason": "",
  "claim_source": "INTERNAL",
  "idempotency_key": "",
  "external_claim_id": "",
  "is_draft": false,
  "draft_completed_at": null,
  "line_items": [ ... ],
  "decision": { ... },
  "fraud_flags": [ ... ],
  "created_at": "2026-03-12T10:00:00Z",
  "updated_at": "2026-03-12T10:00:01Z"
}
```

| Field | Type | Notes |
|-------|------|--------|
| `total_amount` | int64 | In cents. KES 2,500 = `250000` |
| `approved_amount` | int64 | What the insurer will pay |
| `co_pay_amount` | int64 | Member's co-pay portion |
| `member_responsibility` | int64 | Total member owes (co-pay + deductible + over-limit) |
| `claim_source` | string | `INTERNAL`, `PROVIDER_PORTAL`, `PARTNER_API`, `CSV_IMPORT` |
| `is_draft` | bool | `true` = draft not yet submitted |
| `claim_type` | string | `DIRECT`, `REIMBURSEMENT`, `CREDIT`, `EXCEPTION` |
| `decision` | object/null | Only present after adjudication |
| `fraud_flags` | array | Empty if no fraud detected |
| `line_items` | array | Only present on single-claim responses |

---

## 15. Money Convention

All monetary amounts are **BIGINT in cents**:

```
KES 5,000.00   →  500000
KES 80.00      →  8000
KES 0.50       →  50
```

Frontend should divide by 100 for display: `amount / 100`.

---

## 16. Role Permissions Summary

| Endpoint | Allowed Roles |
|----------|---------------|
| Submit / List / Get claim | Any authenticated |
| Approve / Reject | Admin, Manager |
| Vet | Admin, ClaimsOfficer |
| Ready for Payment / Mark Paid / Part Paid | Admin, Finance |
| CSV Import | Admin, ClaimsOfficer |
| Draft CRUD | Any authenticated |
| API Partner Admin | Admin only |
| External claims | API Key auth (no role) |

---

## 17. Error Format

```json
{
  "status": "error",
  "message": "Human-readable error description"
}
```

| Code | Meaning |
|------|---------|
| `400` | Validation error, missing fields |
| `401` | Missing or invalid auth |
| `403` | Insufficient permissions, partner deactivated, wrong provider |
| `404` | Not found |
| `500` | Server error |