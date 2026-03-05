# HIAS Core — System Architecture

## 1. Tech Stack

| Layer | Technology | Purpose |
|-------|-----------|---------|
| Language | Go 1.24+ | Backend runtime |
| HTTP Framework | Gin | REST API routing, middleware |
| RPC | gRPC | Inter-service communication |
| Database | PostgreSQL | Primary data store |
| Query Builder | SQLC | Type-safe SQL → Go code generation |
| Migrations | golang-migrate | Schema versioning |
| Cache | Redis | Session cache, rate limiting |
| Auth | PASETO v4 | Stateless token authentication |
| Identity | AWS Cognito | User pool (optional integration) |
| Queue | Watermill + AWS SQS | Async event processing |
| Scheduler | robfig/cron v3 | Background task scheduling |
| Docs | Swaggo | OpenAPI/Swagger generation |

---

## 2. DDD + Clean Architecture

The codebase follows Domain-Driven Design with Clean Architecture. Dependencies flow inward — outer layers depend on inner layers, never the reverse.

```
┌─────────────────────────────────────────────────────┐
│                   API Gateway                        │
│  (REST handlers, gRPC server, middleware, routes)    │
├─────────────────────────────────────────────────────┤
│                   Service Layer                      │
│  (Business logic implementations, orchestration)     │
├─────────────────────────────────────────────────────┤
│               Infrastructure Layer                   │
│  (DB repos, cache, queue, workers, scheduler)        │
├─────────────────────────────────────────────────────┤
│                   Domain Layer                       │
│  (Entities, repository interfaces, service           │
│   interfaces, DTOs/schemas — INTERFACES ONLY)        │
└─────────────────────────────────────────────────────┘
```

### Layer Responsibilities

**Domain Layer** (`domains/`)
- Entity definitions (structs)
- Repository interfaces (data access contracts)
- Service interfaces (business logic contracts)
- Request/Response DTOs (`schema/` sub-packages)
- No implementations — only interfaces and types

**Infrastructure Layer** (`infrastructures/`)
- PostgreSQL repository implementations (SQLC-generated)
- Redis cache implementations
- SQS queue producers/consumers
- Scheduler task implementations
- External service clients

**Service Layer** (`services/`)
- Business logic implementations
- Orchestration across multiple repositories
- Validation, calculation, and state machine logic
- Audit logging side effects

**API Gateway** (`services/api-gateway/`)
- Gin HTTP handlers (REST)
- gRPC server definitions
- Middleware (auth, RBAC, CORS, logging)
- Route registration
- Request parsing and response formatting

### Directory Structure

```
hias-core/
├── domains/
│   ├── identity/          # Users, roles, auth
│   │   └── schema/        # DTOs
│   ├── product/           # Plans, benefits, exclusions, premium rules
│   │   └── schema/
│   ├── policy/            # Policies, members, endorsements, renewals, underwriting
│   │   └── schema/
│   ├── claims/            # Claims, adjudication, cases
│   │   └── schema/
│   ├── billing/           # Invoices, payments, installments, credit notes, statements
│   │   └── schema/
│   ├── provider/          # Providers, contracts, rate cards
│   │   └── schema/
│   ├── sales/             # Leads, quotations, approval limits
│   │   └── schema/
│   ├── preauth/           # Pre-authorizations
│   │   └── schema/
│   ├── notification/      # Notifications
│   │   └── schema/
│   ├── audit/             # Audit events
│   │   └── schema/
│   ├── analytics/         # Dashboard, KPIs
│   │   └── schema/
│   └── reinsurance/       # Treaties, cessions, recoveries, bordereaux, statements, alerts
│       └── schema/
├── infrastructures/
│   ├── db/                # SQLC generated code, migrations
│   ├── cache/             # Redis implementations
│   ├── queue/             # SQS/Watermill implementations
│   ├── scheduler/         # Cron task implementations
│   └── external/          # External service clients
├── services/
│   ├── api-gateway/
│   │   ├── handlers/      # HTTP handlers per domain
│   │   ├── middleware/     # Auth, RBAC, CORS
│   │   └── routes/        # Route registration
│   ├── identity/          # Auth, user service implementations
│   ├── product/           # Plan, benefit, premium rule implementations
│   ├── policy/            # Policy, member, endorsement implementations
│   ├── claims/            # Claim, adjudicator, validator implementations
│   ├── billing/           # Invoice, payment, credit note implementations
│   ├── provider/          # Provider, contract implementations
│   ├── sales/             # Lead, quotation implementations
│   ├── preauth/           # Pre-auth implementations
│   ├── notification/      # Notification implementations
│   ├── audit/             # Audit implementations
│   ├── analytics/         # Analytics implementations
│   └── reinsurance/       # Treaty, cession, recovery implementations
└── shared/
    ├── types.go           # All status enums and typed constants
    ├── constants.go       # System-wide constants and thresholds
    ├── auth/              # PASETO token maker, payload
    └── utils/             # Response helpers, number generators
```

---

## 3. Domain Map

### 12 Domains and Their Entities

```
┌─────────────┐     ┌──────────────┐     ┌──────────────┐
│  Identity    │     │   Product    │     │   Provider   │
│  ─────────   │     │  ──────────  │     │  ──────────  │
│  User        │     │  Plan        │     │  Provider    │
│  Role        │     │  Benefit     │     │  Contract    │
│              │     │  SubBenefit  │     │  RateCard    │
│              │     │  Exclusion   │     │              │
│              │     │  PremiumRule │     │              │
│              │     │  UWRule      │     │              │
│              │     │  ProvNetwork │     │              │
└──────┬───────┘     └──────┬───────┘     └──────┬───────┘
       │                    │                     │
       │    ┌───────────────┼─────────────────────┘
       │    │               │
       ▼    ▼               ▼
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│    Sales     │     │    Policy    │     │    Claims    │
│  ──────────  │     │  ──────────  │     │  ──────────  │
│  Lead        │────▶│  Policy      │────▶│  Claim       │
│  LeadActivity│     │  Member      │     │  LineItem    │
│  Quotation   │     │  Endorsement │     │  Adjudication│
│  QuotVersion │     │  Renewal     │     │  FraudFlag   │
│  QuotDocument│     │  UWAssessment│     │  Case        │
│  ApprovalLmt │     │  UWFlag      │     │  ClaimDoc    │
│              │     │  PolicyDoc   │     │              │
│              │     │  CreditNote  │     │              │
└──────────────┘     └──────────────┘     └──────┬───────┘
                            │                     │
                            ▼                     ▼
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│   PreAuth    │     │   Billing    │     │ Reinsurance  │
│  ──────────  │     │  ──────────  │     │ ───────────  │
│  PreAuth     │     │  Invoice     │     │ Treaty       │
│              │     │  Payment     │     │ Participant  │
│              │     │  Installment │     │ Layer        │
│              │     │  Schedule    │     │ ProfitCommRl │
│              │     │  Remittance  │     │ Cession      │
│              │     │  Statement   │     │ Recovery     │
│              │     │  LineItem    │     │ Bordereau    │
│              │     │              │     │ ReinStatement│
│              │     │              │     │ TreatyAlert  │
└──────────────┘     └──────────────┘     └──────────────┘

┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│ Notification │     │    Audit     │     │  Analytics   │
│ ───────────  │     │  ──────────  │     │ ───────────  │
│ Notification │     │ AuditEvent   │     │ Dashboard    │
│              │     │              │     │ KPIs         │
└──────────────┘     └──────────────┘     └──────────────┘
```

### Key Relationships

- **Plan → Benefits, Exclusions, PremiumRules, UWRules, ProviderNetworks** (1:N)
- **Benefit → SubBenefits** (1:N hierarchical)
- **Lead → Quotations** (1:N)
- **Quotation → Versions, Documents** (1:N)
- **Quotation → Policy** (1:1 conversion)
- **Policy → Members, Endorsements, Renewals, Documents, CreditNotes, Cases** (1:N)
- **Policy → Plan** (N:1)
- **Member → UnderwritingFlags, Cases** (1:N)
- **Claim → LineItems, ClaimDocuments, FraudFlags, AdjudicationDecision** (1:N)
- **Claim → Policy, Member, Provider** (N:1)
- **PreAuth → Policy, Member, Provider** (N:1)
- **Case → PreAuth** (1:1)
- **Treaty → Participants, Layers, ProfitCommissionRules** (1:N)
- **Treaty → Cessions, Recoveries, Bordereaux, Statements, Alerts** (1:N)
- **Cession → Treaty, Policy** (N:1)
- **Recovery → Claim, Treaty** (N:1)

---

## 4. Authentication & Authorization

### PASETO Token Flow

```
Client                          Server
  │                               │
  │  POST /api/v1/auth/login      │
  │  { email, password }          │
  │──────────────────────────────▶│
  │                               │  Validate credentials
  │                               │  Generate PASETO tokens
  │  { access_token,              │
  │    refresh_token,             │
  │    user }                     │
  │◀──────────────────────────────│
  │                               │
  │  GET /api/v1/policies         │
  │  Authorization: Bearer <AT>   │
  │──────────────────────────────▶│
  │                               │  Verify token
  │                               │  Extract payload
  │                               │  Check role/permissions
  │  { data: [...] }              │
  │◀──────────────────────────────│
  │                               │
  │  POST /api/v1/auth/refresh    │
  │  { refresh_token }            │
  │──────────────────────────────▶│
  │                               │  Validate refresh token
  │                               │  Issue new access token
  │  { access_token }             │
  │◀──────────────────────────────│
```

### Token Payload Structure

```json
{
  "id": "uuid-v4",          // Unique token ID
  "user_id": "uuid-v4",     // User identifier
  "email": "user@email.com",
  "role": "Admin",           // Single role string
  "permissions": [           // Array of permission strings
    "claims:read",
    "claims:approve",
    "*:*"                    // Wildcard = all permissions
  ],
  "issued_at": "2026-01-01T00:00:00Z",
  "expired_at": "2026-01-01T01:00:00Z"
}
```

### Roles (8 total)

| Role | Description | Primary Access |
|------|-------------|----------------|
| `Admin` | Full system access | Everything — bypasses all permission checks |
| `Underwriter` | Risk assessment | Underwriting review, flag resolution, quotation approval |
| `ClaimsOfficer` | Claims processing | Claim vetting, CSV import |
| `Finance` | Financial operations | Payments, remittances, mark-paid |
| `Provider` | Healthcare provider | Self-service portal (limited) |
| `Member` | Insurance member | Self-service portal (limited) |
| `SalesAgent` | Business development | Leads, quotations |
| `Manager` | Supervisory approval | Claims approval/rejection, quotation approval |

### RBAC Middleware

**Two middleware functions:**

1. **`RequireRole(roles ...string)`** — Checks if user's role is in the allowed list
2. **`RequirePermission(resource, action)`** — Checks permission string format `resource:action`

**Permission matching logic:**
- Exact match: `claims:approve` matches `claims:approve`
- Wildcard action: `claims:*` matches any claims action
- Universal: `*:*` matches everything
- Admin role automatically bypasses permission checks (but NOT role checks)

**Auth flow in middleware:**
1. Extract `Authorization: Bearer <token>` header
2. Verify PASETO token signature and expiration
3. Store payload in Gin context (`auth_payload`, `user_id`, `role`, `permissions`)
4. Subsequent middleware/handlers read from context

---

## 5. API Response Format

All API responses follow a consistent JSON envelope.

### Success Response

```json
{
  "status": "success",
  "message": "Resource created successfully",
  "data": { ... }
}
```

### Error Response

```json
{
  "status": "error",
  "message": "Validation failed: email is required"
}
```

### Paginated Response

```json
{
  "status": "success",
  "message": "Records retrieved",
  "data": [ ... ],
  "page": 1,
  "page_size": 20,
  "total_count": 150,
  "total_pages": 8
}
```

**Pagination calculation:** `total_pages = ceil(total_count / page_size)`

### HTTP Status Codes

| Code | Usage |
|------|-------|
| 200 | Successful GET, PUT |
| 201 | Successful POST (resource created) |
| 400 | Validation error, bad request |
| 401 | Missing or invalid/expired token |
| 403 | Insufficient role or permissions |
| 404 | Resource not found |
| 500 | Internal server error |

---

## 6. Conventions

### Money

- Stored as **BIGINT** (PostgreSQL) / **int64** (Go) — values in **cents**
- `8,000 KES` = `800000` in the database
- Default currency: `KES` (Kenyan Shilling)
- Frontend must divide by 100 for display: `800000 → "KES 8,000.00"`

### IDs

- All entity IDs use **UUID v4** (`google/uuid.UUID`)
- Transmitted as strings in JSON: `"550e8400-e29b-41d4-a716-446655440000"`

### Timestamps

- Stored as **TIMESTAMPTZ** (PostgreSQL) / **time.Time** (Go)
- Serialized as ISO 8601: `"2026-01-15T10:30:00Z"`

### Number Formats

| Entity | Format | Example |
|--------|--------|---------|
| Claim | `CLM-YYYY-NNNNNN` | CLM-2026-000042 |
| Policy | `POL-YYYY-NNNNNN` | POL-2026-000015 |
| Treaty | `TRY-YYYY-NNNNNN` | TRY-2026-000003 |
| Cession | `CES-YYYY-NNNNNN` | CES-2026-000100 |
| Recovery | `REC-YYYY-NNNNNN` | REC-2026-000008 |
| Bordereau | `BDX-YYYY-NNNNNN` | BDX-2026-000001 |
| Reinsurer Statement | `RST-YYYY-NNNNNN` | RST-2026-000005 |
| Member | `MBR-YYYY-NNNNNN` | MBR-2026-000200 |
| Credit Note | Generated | CN-2026-000012 |
| Invoice | Generated | INV-2026-000050 |
| Statement | Generated | STM-2026-000010 |
| Case | Generated | CAS-2026-000007 |
| Lead | Generated | LD-2026-000030 |
| Quotation | Generated | QT-2026-000025 |

### Status Enums

All status fields use typed string constants defined in `shared/types.go`. Every status type is a named Go type (e.g., `PolicyStatus string`), ensuring type safety. See the API Reference for complete enum values per entity.

---

## 7. Scheduler Tasks

Background cron jobs run on predefined schedules:

| Task | Cron Schedule | Frequency | Description |
|------|--------------|-----------|-------------|
| **Claim SLA Enforcement** | `0 */4 * * *` | Every 4 hours | Detects SLA-breached claims (>48h) and approaching claims (<24h remaining). Sends IN_APP notifications to claim creators. |
| **PreAuth Expiry** | `0 2 * * *` | Daily at 2am | Expires pre-authorizations past their `validity_end` date. Transitions APPROVED → EXPIRED. |
| **Policy Lapse** | `0 1 * * *` | Daily at 1am | Lapses policies with unpaid invoices >30 days. Terminates policies with unpaid invoices >120 days. |
| **Billing Cycle** | `0 0 1 * *` | 1st of month midnight | Generates invoices for all active policies based on premium amount and billing frequency. |
| **Payment Reminder** | `0 8 * * *` | Daily at 8am | Sends escalating reminders for unpaid invoices (day 7, 14, 21 after due date). |
| **Remittance Cycle** | `0 0 * * 1` | Monday midnight | Generates weekly provider remittances for approved/paid claims. |
| **Payment Retry** | `0 */4 * * *` | Every 4 hours | Retries failed payments up to MaxNotificationRetries (3). |
| **Reconciliation** | `0 2 * * *` | Daily at 2am | Matches confirmed payments to bank statements. |
| **Notification Retry** | `*/30 * * * *` | Every 30 minutes | Retries failed SMS/email notifications up to max retries. |

### Key Constants

| Constant | Value | Description |
|----------|-------|-------------|
| `ClaimSLAHours` | 48 | Hours before a claim breaches SLA |
| `MaxPaymentRetries` | 3 | Maximum payment retry attempts |
| `MaxNotificationRetries` | 3 | Maximum notification delivery attempts |
| `InvoiceDueDays` | 30 | Invoice payment terms (days) |
| `PreAuthValidityDays` | 30 | Pre-authorization validity period (days) |
| `QuotationValidityDays` | 30 | Quotation validity period (days) |
| `FraudAmountThresholdCents` | 50,000,000 | 500K KES fraud detection trigger |
| `CatastropheThresholdCents` | 500,000,000 | 5M KES catastrophe alert threshold |
| `AggregateWarningPercent` | 80 | Alert at 80% of aggregate layer limit |

---

## 8. Notification System

### Channels

| Channel | Code | Delivery |
|---------|------|----------|
| In-App | `IN_APP` | Stored in DB, fetched via API |
| SMS | `SMS` | External SMS gateway |
| Email | `EMAIL` | External email service |
| Push | `PUSH` | Mobile push notifications |

### Notification Types

| Type | Triggered By |
|------|-------------|
| `QUOTATION` | Quotation issued, sent to client, accepted/declined |
| `APPROVAL` | Quotation/endorsement/renewal approval or rejection |
| `CLAIM` | Claim status changes, SLA warnings/breaches |
| `POLICY` | Policy activation, lapse, termination, renewal |
| `DOCUMENT` | Document generation (welcome letter, schedule, LOU) |

### Notification Status Flow

```
PENDING → SENT → DELIVERED
              └→ FAILED (→ retry up to 3 times)
                          └→ READ (user marks read)
```

### When Notifications Fire

| Event | Channel | Recipient |
|-------|---------|-----------|
| Claim SLA approaching (< 24h) | IN_APP | Claim creator |
| Claim SLA breached (> 48h) | IN_APP | Claim creator |
| Quotation sent to client | SMS/EMAIL | Client |
| Quotation approval needed | IN_APP | Approver (Manager/Admin) |
| Policy activated | IN_APP | Policyholder |
| Payment reminder (7/14/21 days) | SMS/EMAIL | Policyholder |
| Notification delivery failed | — | Retry via scheduler |

---

## 9. State Machines

### Policy Lifecycle

```
DRAFT ──────────────────────▶ ACTIVE
                                │
                    ┌───────────┼───────────┐
                    ▼           ▼           ▼
                 LAPSED    SUSPENDED   TERMINATED
                    │           │
                    └─────┬─────┘
                          ▼
                       ACTIVE (reinstate)
```

### Claim Pipeline

```
RECEIVED → VALIDATED → ADJUDICATED ─┬─▶ APPROVED → VETTED ──────────────┬─▶ READY_FOR_PAYMENT → PAID
                                     │              │                    │
                                     ├─▶ REJECTED   └─▶ PARTIALLY_VETTED┘
                                     │
                                     └─▶ MANUAL_REVIEW → APPROVED/REJECTED
```

### Pre-Authorization

```
SUBMITTED → UNDER_REVIEW ─┬─▶ APPROVED ──▶ CLAIMED (when claim references it)
                           │            └─▶ EXPIRED (scheduler)
                           ├─▶ DENIED
                           └─▶ INFO_REQUESTED
```

### Quotation

```
DRAFT → ISSUED → PENDING_DECISION ─┬─▶ ACCEPTED → CONVERTED (to policy)
                                    ├─▶ DECLINED
                                    └─▶ EXPIRED (scheduler)
```

### Endorsement

```
PENDING → APPROVED → APPLIED
       └→ REJECTED
```

### Renewal

```
PENDING → APPROVED → COMPLETED
       └→ REJECTED
       └→ EXPIRED (scheduler)
```

### Treaty (Reinsurance)

```
DRAFT → ACTIVE → TERMINATED
              └→ EXPIRED (scheduler)
```

### Cession

```
PENDING → BOOKED → REVERSED
```

### Recovery

```
NOTIFIED → ACKNOWLEDGED ─┬─▶ APPROVED → PAID
                          │
                          └─▶ INFO_REQUESTED → APPROVED → PAID
                                                        └→ WRITTEN_OFF
```

### Case (Inpatient)

```
SCHEDULED → ADMITTED → IN_TREATMENT → DISCHARGED → CLOSED
```

### Provider

```
PENDING → CREDENTIALING → ACTIVE → SUSPENDED → TERMINATED
```

### Credit Note

```
DRAFT → APPROVED → APPLIED
     └→ CANCELLED
```

### Bordereau

```
DRAFT → FINALIZED → SENT
```

### Reinsurer Statement

```
DRAFT → ISSUED → ACKNOWLEDGED → SETTLED
```

---

## 10. Business Rules — Detailed

This section documents every calculation formula, threshold, and business rule used in the system. A frontend implementation must respect these rules for validation and display.

### 10.1 Premium Calculation Engine

Premium is calculated via rules attached to a plan. There are four calculation types:

**per_member** — Rate applied per individual member:
```
For each member:
  1. Find matching rule (priority order):
     a. Relationship + age range match (best)
     b. Relationship-only match
     c. Generic rule (relationship="") + age range
     d. Generic rule (relationship="") without age
     e. Fallback: plan.BasePremium
  2. Add rule.RateAmount to total

Age calculation: age = now.Year - dob.Year - (1 if now.YearDay < dob.YearDay)
Default MaxAge: 150 if not specified (0 means no limit)
```

**per_family** — Flat rate based on family size:
```
Find rule where CalculationType=PER_FAMILY with highest MinMembers ≤ memberCount
Return that rule's RateAmount directly (ignores per-member matching)
```

**Group discount** (applied after per_member or per_family):
```
If any rule has MinMembers > 0 AND memberCount >= MinMembers:
  If DiscountType = PERCENTAGE:
    totalPremium -= totalPremium × DiscountValue / 10000   (basis points: 5000 = 50%)
  If DiscountType = FIXED:
    totalPremium -= DiscountValue
  If totalPremium < 0: totalPremium = 0
```

### 10.2 Quotation Pricing

When creating a quotation version:
```
basePremium = CalculatePremiumWithMembers(planID, memberCount, proposedMembers)

If DiscountType = PERCENTAGE:
  discountAmount = basePremium × discountValue / 10000
If DiscountType = FIXED:
  discountAmount = discountValue

If LoadingType = PERCENTAGE:
  loadingAmount = basePremium × loadingValue / 10000
If LoadingType = FIXED:
  loadingAmount = loadingValue

finalPremium = basePremium - discountAmount + loadingAmount
If finalPremium < 0: finalPremium = 0
```

**Approval limit check** (at version creation time):
```
Load default approval limits for SalesAgent role
If no limits configured: requiresApproval = true
If discountType=PERCENTAGE AND discountValue > limit.MaxDiscountPercentage: requiresApproval = true
If discountType=FIXED AND discountValue > limit.MaxDiscountAmount: requiresApproval = true
If loadingType=PERCENTAGE AND loadingValue > limit.MaxLoadingPercentage: requiresApproval = true
If loadingType=FIXED AND loadingValue > limit.MaxLoadingAmount: requiresApproval = true
```

**Approval enforcement** (at approve time):
```
Load limits for approver's role
If discountValue exceeds approver's limit:
  Error: "Discount of {value} bps exceeds your limit of {limit} bps. Requires escalation to {escalationRole}"
Same for loading values
```

### 10.3 Convert Quotation to Policy

When an ACCEPTED quotation is converted:
1. Create policy (DRAFT) with quotation's client info, plan, and parsed start/end dates (1 year)
2. Enroll proposed members from the latest version's `proposed_members` JSON
3. Create installment schedule using the version's `billing_frequency`
4. Update quotation: set `policy_id`, status → CONVERTED
5. Update lead: status → WON

### 10.4 Claim Adjudication Pipeline

The adjudication engine runs 9 sequential checks. Each check either PASSES, FAILS (reject), or FLAGS (manual review):

| Step | Check | Category | On Fail |
|------|-------|----------|---------|
| 1 | Policy is ACTIVE | eligibility | REJECT |
| 2 | Member exists | eligibility | REJECT |
| 3 | Provider is ACTIVE | eligibility | REJECT |
| 3a | Provider accreditation status | eligibility | FLAG only (not reject) |
| 3b | Provider in plan's network | eligibility | REJECT |
| 3c | Provider has valid contract covering service date | eligibility | REJECT |
| 3d | Pre-auth valid (if provided): APPROVED, not expired, provider matches | eligibility | REJECT |
| 3e | Pre-auth procedure codes match (if provided) | eligibility | FLAG only |
| 3f | Claim amount vs pre-auth approved amount | eligibility | FLAG only |
| 4 | Plan has benefits | coverage | REJECT |
| 5 | Service date past waiting period | eligibility | REJECT |
| 5b | Member age within benefit min_age/max_age | eligibility | REJECT |
| 6 | Diagnosis codes not in plan exclusions | coverage | REJECT |
| 7 | Amount calculation (annual limit, sub-limits, deductible, co-pay) | limits | See below |
| 8 | Duplicate claim detection | fraud | MANUAL_REVIEW |

**Critical details:**
- Waiting period is calculated from **member.CreatedAt** (enrollment date), NOT policy start date
- Age eligibility uses **claim.ServiceDate** for age calculation, NOT current date
- Contract validity uses **strict inequalities**: `startDate < serviceDate < endDate`
- Benefit category: if `claim.AdmissionDate != nil` → INPATIENT, else → OUTPATIENT
- If no exact benefit category match found, falls back to `benefits[0]`

### 10.5 Amount Calculation (Adjudication Step 7)

```
payableAmount = claim.TotalAmount

1. Annual Limit Check:
   used = GetApprovedAmountForBenefitThisYear(memberID, benefitCategory)
   remaining = benefit.AnnualLimit - used
   If remaining ≤ 0: REJECT ("Annual limit exhausted")
   If payableAmount > remaining: payableAmount = remaining

2. Sub-Limit Enforcement:
   If SubLimitType = PER_VISIT AND payableAmount > SubLimitValue:
     payableAmount = SubLimitValue
   If SubLimitType = PER_ITEM AND payableAmount > SubLimitValue:
     payableAmount = SubLimitValue

3. Deductible:
   payableAmount -= benefit.DeductibleAmount
   If payableAmount < 0: payableAmount = 0

4. Co-Pay:
   If CoPayType = PERCENTAGE:
     coPayAmount = payableAmount × CoPayValue / 100
   If CoPayType = FIXED:
     coPayAmount = CoPayValue
   payableAmount -= coPayAmount

5. Pre-Auth Cap (if pre-auth exists):
   If payableAmount > preauth.ApprovedAmount:
     payableAmount = preauth.ApprovedAmount

6. Member Responsibility:
   memberResponsibility = coPayAmount + (claim.TotalAmount - payableAmount - coPayAmount)
   Simplified: memberResponsibility = claim.TotalAmount - payableAmount
```

### 10.6 Fraud Detection

After adjudication, 6 additional fraud checks run. Each creates a FraudFlag entity:

| Check | Flag Type | Severity | Details |
|-------|-----------|----------|---------|
| Duplicate claim detected | DUPLICATE | — | Forces MANUAL_REVIEW in adjudication |
| High frequency of procedure for member | FREQUENCY | MEDIUM | "High frequency of procedure {code} for member" |
| Claim amount exceeds threshold | AMOUNT_THRESHOLD | HIGH | "Claim amount {amount} exceeds threshold for procedure {code}" |
| Provider has no valid contract for service date | EXPIRED_CONTRACT | HIGH | "Provider has no valid contract covering the service date" |
| Provider is SUSPENDED | SUSPENDED_PROVIDER | CRITICAL | "Provider is suspended" |
| Unit price exceeds provider's rate card | RATE_CARD_OVERCHARGE | MEDIUM | "Unit price exceeds rate card for procedure {code}" |
| Repeat visit for same procedure | REPEAT_VISIT | LOW | "Repeat visit detected for procedure {code}" |

**Fraud amount threshold**: 500,000 KES (50,000,000 cents)

### 10.7 Claim Vetting Rules by Type

| Claim Type | Rule | Error Message |
|------------|------|---------------|
| DIRECT | If inpatient (has admission_date), must have pre-auth reference | "Inpatient direct claims require pre-authorization reference" |
| REIMBURSEMENT | Vetted amount ≤ total claimed amount | "Vetted amount cannot exceed total claimed amount for reimbursement claims" |
| EXCEPTION | Vetted amount ≤ 150% of approved amount (`approved × 3/2`) | "Exception claim vetted amount exceeds 150% of approved amount — requires manual override" |

After vetting:
- If vetted_amount == approved_amount → VETTED
- If vetted_amount < approved_amount → PARTIALLY_VETTED

### 10.8 SLA Tracking

```
sla_breach_at = claim.created_at + 48 hours
Breached = now > sla_breach_at
Approaching = sla_breach_at - now < 24 hours
```

Scheduler checks every 4 hours and sends IN_APP notifications:
- Breached: notification to claim creator
- Approaching: warning notification to claim creator

### 10.9 Member Enrollment Underwriting

When enrolling a member, three underwriting checks run:

1. **Double insurance**: If member's `national_id` exists on another ACTIVE policy → Error + DOUBLE_INSURANCE flag (HIGH severity)
2. **Age vs premium rules**: Member age checked against plan's premium rule age ranges for their relationship → Error + MAX_AGE flag (HIGH severity)
3. **Plan underwriting rules**: Each active rule evaluated:
   - MAX_AGE: member age > parameter_value → Flag with rule's severity
   - MIN_AGE: member age < parameter_value → Flag with rule's severity
   - Rules can be relationship-specific (skip if mismatch)

### 10.10 Member Removal & Pro-Rata Credit Notes

When a member is removed:
```
1. Premium recalculated via CalculatePremiumWithMembers() with remaining members
2. If newPremium < oldPremium:
   totalDays = (policy.EndDate - policy.StartDate).Hours / 24
   remainingDays = (policy.EndDate - now).Hours / 24
   premiumDiff = oldPremium - newPremium
   refundAmount = int64(float64(premiumDiff) × remainingDays / totalDays)
3. If refundAmount > 0:
   Create credit note with reason "Pro-rata refund for member removal: {memberName}"
   Credit note is AUTO-APPROVED (reason contains "Pro-rata refund")
```

### 10.11 Policy Activation Side Effects

When a DRAFT policy is activated:
1. Status → ACTIVE
2. Auto-generate welcome letter document (async, non-blocking)
3. Auto-generate member cards for all members (async, non-blocking)
4. Audit event logged

### 10.12 Renewal Premium Calculation

When completing a renewal, premium is adjusted based on claims experience:

```
Loss Ratio = (totalApprovedClaimAmount / originalPremium) × 100

Claims-Experience Loading Tiers:
  Loss Ratio > 100%:  +25% loading ("Claims loading +25%")
  Loss Ratio > 75%:   +15% loading ("Claims loading +15%")
  Loss Ratio > 50%:   +10% loading ("Claims loading +10%")
  Loss Ratio < 30%:   -5% discount ("Good claims discount -5%")
  Otherwise:           no adjustment

Applied: premium += int64(float64(premium) × loadingPercentage)
```

Then premium rules are recalculated with current members (takes precedence if result > 0).

**Member re-validation during renewal:**
- Each active member's age is re-checked against plan rules
- If age out of range: member SKIPPED, RENEWAL_SKIP flag created (MEDIUM severity)
- Double insurance re-checked: if detected, member SKIPPED, RENEWAL_SKIP flag created (HIGH severity)
- Passing members are copied to new policy with new member numbers

### 10.13 Reinsurance Cession Calculation

**Quota Share:**
```
totalShare = sum(participant.SharePercentage)
avgCommissionRate = sum(participant.CommissionRate) / participantCount
cededAmount = floor(grossAmount × totalShare / 100)
retainedAmount = grossAmount - cededAmount
commissionAmount = floor(cededAmount × avgCommissionRate / 100)

Retention limit override:
  If treaty.RetentionLimit > 0 AND retainedAmount < treaty.RetentionLimit:
    retainedAmount = treaty.RetentionLimit
    cededAmount = grossAmount - retainedAmount
    commissionAmount = floor(cededAmount × avgCommissionRate / 100)
```

**Auto-cede** finds ALL active QUOTA_SHARE treaties and creates one cession per treaty.

### 10.14 Reinsurance Recovery Calculation

**Quota Share recovery:**
```
totalShare = GetTotalShareByTreaty(treatyID)
recoverable = floor(approvedAmount × totalShare / 100)
```

**XOL (Excess of Loss) recovery — per layer, sorted by layer_number ascending:**
```
For each layer:
  excess = approvedAmount - layer.AttachmentPoint
  If excess ≤ 0: skip
  layerExposure = min(excess, layer.LayerLimit)
  recoverable = layerExposure - layer.DeductibleAmount
  If recoverable ≤ 0: skip

  If layer.AggregateLimit exists:
    remaining = layer.AggregateLimit - layer.AggregateUsed
    If remaining ≤ 0: skip
    If recoverable > remaining: recoverable = remaining
    Update: layer.AggregateUsed += recoverable
```

### 10.15 Profit Commission Calculation

```
lossRatio = (claimsRecovered × 100) / premiumCeded  (0 if premiumCeded = 0)
netProfit = premiumCeded - claimsRecovered

If CARRY_FORWARD rule exists with CarryForwardBalance > 0:
  netProfit -= CarryForwardBalance

Match lossRatio to rule where: LossRatioFrom ≤ lossRatio ≤ LossRatioTo
commissionRate = matched rule's CommissionRate

If netProfit > 0:
  commissionAmount = floor(netProfit × commissionRate / 100)
  carryForward = 0
Else:
  commissionAmount = 0
  carryForward = -netProfit  (deficit carried forward)
```

### 10.16 Reinsurer Statement Calculation

```
premiumCeded = totalCeded × participant.SharePercentage / 100
claimsRecovered = totalRecovered × participant.SharePercentage / 100
commissionDue = premiumCeded × participant.CommissionRate / 100
netBalance = premiumCeded - claimsRecovered - commissionDue
```

### 10.17 Treaty Alert Thresholds

| Alert | Trigger | Severity | Threshold |
|-------|---------|----------|-----------|
| LIMIT_BREACH | Layer aggregate usage ≥ 100% | CRITICAL | `AggregateUsed × 100 / AggregateLimit ≥ 100` |
| AGGREGATE_WARNING | Layer aggregate usage ≥ 80% | HIGH | `AggregateUsed × 100 / AggregateLimit ≥ 80` |
| CATASTROPHE_THRESHOLD | Treaty total recoverable > 5M KES | CRITICAL | `totalRecoverable > 500,000,000 cents` |
| EXPIRY_WARNING | Treaty expiring within 30 days | MEDIUM | `ExpiryDate - now ≤ 30 days` |

### 10.18 Lead Status Transitions

```
Valid Pipeline: NEW → CONTACTED → QUALIFIED → PROPOSAL_SENT → NEGOTIATION → WON / LOST
DORMANT can be set from any state
WON and LOST cannot transition back to NEW (error: "Cannot transition from {old} to {new}")
When a quotation is created for a lead, status auto-advances to PROPOSAL_SENT (from NEW/CONTACTED/QUALIFIED)
When a quotation is converted to policy, lead status auto-set to WON
```

### 10.19 Endorsement Application

When an APPROVED endorsement is applied, it dispatches based on type:
- **ADD_MEMBER**: Calls member enrollment with changes as EnrollMemberRequest
- **REMOVE_MEMBER**: Calls member removal with MemberID and Reason from changes
- **UPDATE_MEMBER**: Calls member update with MemberID and Updates from changes
- **PLAN_CHANGE**: Calls policy ChangePlan with new plan ID from changes

After application, `policy.PremiumAmount += endorsement.PremiumAdjustment`

### 10.20 Underwriting Auto-Decision Engine

When an underwriting assessment is submitted, rules are auto-evaluated and a decision is made:

```
Risk Score Thresholds:
  UnderwritingAutoApproveThreshold = 30
  UnderwritingReferThreshold       = 60

Decision logic:
  If ANY rule with is_blocking=true is triggered:
    → DECLINED ("Declined: blocking rule triggered")
  Else if totalRiskScore > 60:
    → DECLINED ("Declined: risk score {N} exceeds threshold 60")
  Else if totalRiskScore > 30:
    → REFER ("Referred: risk score {N} exceeds auto-approve threshold 30")
  Else (totalRiskScore ≤ 30, no blockers):
    → APPROVED ("Auto-approved: risk score within acceptable range")

totalRiskScore = SUM(rule.RiskScoreWeight) for all triggered rules
```

**Underwriting rule types evaluated:**

| Rule Type | Trigger Condition | Details |
|-----------|-------------------|---------|
| MAX_AGE | member age > ParameterValue | `"Member age {age} exceeds max age {maxAge} for {relationship}"` |
| MIN_AGE | member age < ParameterValue | `"Member age {age} below min age {minAge} for {relationship}"` |
| DOUBLE_INSURANCE | same NationalID on another ACTIVE policy | `"Double insurance: NationalID {id} already covered under policy {number}"` |
| PRE_EXISTING_CONDITION | questionnaire[ParameterKey] matches ParameterValue OR equals "yes"/"true" (case-insensitive) | `"Pre-existing condition flagged: {key} = {value}"` |
| BMI_THRESHOLD | questionnaire["bmi"] > ParameterValue (float) | `"BMI {X} exceeds threshold {Y}"` |
| WAITING_PERIOD | questionnaire[ParameterKey] equals "yes"/"true" (case-insensitive) | `"Waiting period applies: {ParameterValue} days for {ParameterKey}"` (informational only) |

### 10.21 Claim Submission Pipeline Details

When a claim is submitted (POST /claims), the full pipeline runs synchronously:
1. Calculate TotalAmount = SUM(lineItem.UnitPrice × lineItem.Quantity)
2. Set SLA breach: sla_breach_at = now + 48 hours
3. Default claim_type to "DIRECT" if empty
4. Generate claim number: CLM-{YEAR}-{COUNTER:06d}
5. Run 8-rule validator (all errors collected, not short-circuit)
6. If validator fails: claim is REJECTED (HTTP 201 returned, NOT 4xx)
7. Update status → VALIDATED
8. Run 9-step adjudicator
9. Store adjudication decision
10. If REJECTED: set rejection_reason from FAIL rule details
11. Run 6 fraud checks (create FraudFlag entities, non-blocking)
12. If pre-auth referenced: set pre-auth status → CLAIMED

**Critical: claim submission returns HTTP 201 even if auto-rejected.** The frontend must check the response `message` field for "Claim submitted but rejected:" to detect auto-rejection.

**Pre-pipeline validator (8 rules, all errors collected):**
1. Policy exists
2. Policy status == ACTIVE
3. Member exists
4. Member belongs to this policy (member.PolicyID == claim.PolicyID)
5. Provider exists
6. Provider status == ACTIVE
7. At least one line item
8. TotalAmount > 0

### 10.22 Claim Status Complete Flow

```
RECEIVED
  ├─ (validator fails) ────────────────────► REJECTED
  └─ (validator passes) ──► VALIDATED
                               └──► (adjudication) ──┬─► ADJUDICATED ──┬─► APPROVED (manual)
                                                      │                 │      └─► VETTED ──────────┬─► READY_FOR_PAYMENT ──┬─► PAID
                                                      │                 │      └─► PARTIALLY_VETTED─┘                      └─► PART_PAID
                                                      │                 └─► REJECTED (manual)
                                                      ├─► REJECTED
                                                      └─► MANUAL_REVIEW ──┬─► APPROVED (manual, same path as ADJUDICATED)
                                                                          └─► REJECTED (manual)
```

**Status gates (which states allow which actions):**
| Action | Required Status | Error |
|--------|----------------|-------|
| Approve | ADJUDICATED or MANUAL_REVIEW | "Cannot approve claim in {status} status" |
| Reject | RECEIVED, VALIDATED, ADJUDICATED, or MANUAL_REVIEW | "Cannot reject claim in {status} status" |
| Vet | ADJUDICATED or APPROVED | "Cannot vet claim in {status} status" |
| Ready for Payment | VETTED or PARTIALLY_VETTED | "Cannot mark claim as ready for payment" |
| Mark Paid | READY_FOR_PAYMENT | "Cannot mark claim as paid" |
| Mark Part Paid | READY_FOR_PAYMENT | "Cannot mark claim as part paid" |

### 10.23 Policy Creation Defaults

When creating a policy:
- `StartDate` defaults to `time.Now()` if not provided
- `EndDate` defaults to `StartDate + 1 year` if not provided
- `PremiumAmount` = `plan.BasePremium` at creation time
- `Status` = DRAFT (always)
- `PolicyNumber` auto-generated: POL-{YEAR}-{COUNTER:06d}

### 10.24 Installment Schedule Details

**Installment count by frequency:**
| Frequency | Installments |
|-----------|-------------|
| monthly | 12 |
| quarterly | 4 |
| semi_annual | 2 |
| annual | 1 |

**Amount per installment:** `policy.PremiumAmount / totalInstallments` (integer division — remainder is lost)

**Due dates:** First installment due immediately (startDate). Subsequent installments spaced by frequency interval.

### 10.25 Treaty Activation Requirements

When activating a treaty (DRAFT → ACTIVE):
- Must have at least one participant (totalShare > 0)
- Total participant share must NOT exceed 100% (totalShare ≤ 100)
- Partial share allowed (e.g., 60% total is valid — insurer retains the rest)

### 10.26 Recovery Payment Details

- `RecordPayment`: partial payment still transitions to PAID status (no PARTIALLY_PAID state)
- `WriteOff`: allowed from ANY status except PAID
- After payment: `recoveredAmount += paymentAmount`, `outstandingAmount = recoverableAmount - recoveredAmount` (floored at 0)

### 10.27 Bordereau Data Sources

**Premium bordereau:** includes only BOOKED cessions for the treaty/period (PENDING and REVERSED excluded)
**Claim bordereau:** filtered by `recovery.CreatedAt` within period range in application code (not SQL). Hard limit: 10,000 recovery records.
**Claim bordereau commission:** always 0 (commission not applicable to claim recoveries)

### 10.28 Reinsurer Statement — Claims Data Scope

`claimsRecovered` in statement generation uses **all-time total** recovered for the treaty, NOT filtered by the statement period. This means the claims figure includes all historical recoveries.

### 10.29 Notification System Details

- Notifications are queue-based: published to `"notification-events"` topic asynchronously
- Retry: up to 50 failed notifications per run, max 3 retries per notification
- Retry only increments counter — actual re-dispatch handled by queue workers
- If queue unavailable: notification record saved but delivery silently fails

### 10.30 Scheduler Status

**Fully active scheduler tasks:**
| Task | Cron | Status |
|------|------|--------|
| Claim SLA Enforcement | every 4 hours | **ACTIVE** — fully wired |

**Stub scheduler tasks (logic defined but service calls commented out):**
| Task | Cron | Intended Action |
|------|------|-----------------|
| Billing Cycle | 1st of month midnight | Generate invoices for active policies |
| Payment Reminder | Daily 08:00 | Send reminders at day 7, 14, 21 |
| Policy Lapse | Daily 01:00 | Lapse unpaid >30d, terminate lapsed >120d |
| Pre-Auth Expiry | Daily 02:00 | Expire approved pre-auths past validity_end |
| Payment Retry | Every 4 hours | Retry failed payments (max 3 retries) |
| Reconciliation | Daily 02:00 | Match payments to bank statements |
| Notification Retry | Every 30 min | Retry failed notifications (max 50/run) |
| Remittance Cycle | Monday midnight | Pay approved claims to providers |

**Manual expiry endpoints (NOT scheduled — must be called via API):**
- `POST /treaties/expire` — expire overdue treaties
- `POST /quotations/expire` — expire quotations past valid_until
- `POST /renewals/expire` — expire pending renewals past expires_at

### 10.31 Quotation Document Permissions

Default `CanEditRoles` and `CanDeleteRoles` when not provided: `["Admin"]`
Edit/delete operations are blocked unless the authenticated user's role is in the respective permission array.

### 10.32 Credit Note Auto-Approval

Auto-approval triggers when the `reason` field contains the **case-sensitive** substring `"Pro-rata refund"`. The creating user is recorded as the approver. Currency is always hardcoded to `"KES"`.

### 10.33 Pre-Auth Approval Details

When a pre-auth is approved:
- `ApprovedAmount` = `EstimatedCost` (copies the estimated cost exactly — no manual override at approval)
- `AuthCode` = `"AUTH-{YEAR}-{6-digit}"` format
- `ValidityStart` = now
- `ValidityEnd` = now + 30 days (PreAuthValidityDays = 30)

### 10.34 CSV Import Formats

**Member CSV columns:**
- Required: `name`, `date_of_birth` (YYYY-MM-DD), `gender`, `relationship`
- Optional: `national_id`, `phone`, `email`, `kra_pin`, `county`, `address`
- Column names are **case-insensitive** and **trimmed** (`strings.TrimSpace(strings.ToLower(col))`)

**Claim CSV columns:**
- Required: `policy_id`, `member_id`, `provider_id`, `service_date` (YYYY-MM-DD), `procedure_code`, `procedure_name`, `quantity`, `unit_price`
- Optional: `claim_type` (defaults to "DIRECT"), `diagnosis_code` (defaults to "UNSPECIFIED"), `notes`, `preauth_id`
- Quantity defaults to 1 if ≤ 0

### 10.35 Endorsement Changes JSON Payload Structure

The `changes` field in endorsements is a JSON blob whose structure depends on the `endorsement_type`:

**ADD_MEMBER:** `EnrollMemberRequest` payload:
```json
{ "name": "string", "date_of_birth": "YYYY-MM-DD", "gender": "string", "relationship": "string", "national_id": "string", "phone": "string", "email": "string" }
```

**REMOVE_MEMBER:**
```json
{ "member_id": "uuid", "reason": "string" }
```

**UPDATE_MEMBER:**
```json
{ "member_id": "uuid", "updates": { "name": "string", "phone": "string", "email": "string", ... } }
```

**PLAN_CHANGE:** `ChangePlanRequest` payload:
```json
{ "new_plan_id": "uuid" }
```

After the dispatched action, if `endorsement.PremiumAdjustment != 0`, it is applied **additively** on top of whatever premium change the action itself made: `policy.PremiumAmount += endorsement.PremiumAdjustment`.

### 10.36 Renewal — Clarifications

- **Premium rules overwrite claims loading**: Claims experience loading (the +25%/+15%/+10%/-5% adjustments) is applied first. Then premium rules are recalculated with current members. If premium rules return a valid result (error == nil AND result > 0), the result **completely replaces** the loaded premium — it does NOT stack on top of the loading.
- **RejectRenewal stores reason in `PremiumChangeReason`**: The rejection reason is stored in the `PremiumChangeReason` field (a repurposed field), not a separate `Reason` column.
- **Bulk renewal auto-sets renewal date**: `BulkInitiateRenewals` auto-sets `RenewalDate` to **30 days from now** (`time.Now().AddDate(0, 0, 30)`).
- **New policy is DRAFT**: The renewal creates a new policy in `DRAFT` status with a new policy number. It must be separately activated via `PUT /policies/:id/activate`.

### 10.37 Per-Family Premium Short-Circuit

When calculating premiums, if **any** premium rule for the plan has `calculation_type == "per_family"`, the engine:
1. Calls `findFamilyRule(rules, memberCount)` — picks the rule with the **largest** `MinMembers` that is ≤ `memberCount` (best-fit match)
2. If matched, returns `matchedRule.RateAmount` immediately as the total premium
3. **No per-member calculation is performed** — per_family short-circuits the entire flow
4. If no per_family rule matches the member count, falls back to `plan.BasePremium`

### 10.38 Pre-Auth Status Guards

Pre-authorization **has no status guards** on `ApprovePreAuth` and `DenyPreAuth`. This means:
- A pre-auth can be approved from ANY status (SUBMITTED, DENIED, INFO_REQUESTED, even EXPIRED)
- A pre-auth can be denied from ANY status
- Only `ExpirePreAuths` (batch job) restricts to APPROVED status
- The documented state machine (SUBMITTED → UNDER_REVIEW → APPROVED/DENIED) represents intended flow, not enforced constraints

### 10.39 Invoice Generation Rules

- **Amount** = `policy.PremiumAmount` (full annual premium, in cents)
- **Due date** = `time.Now().AddDate(0, 0, 30)` (now + 30 days)
- **Billing period** = `PeriodStart: time.Now()`, `PeriodEnd: time.Now().AddDate(0, 1, 0)` (now to now + 1 month)
- **Invoice number format** = `INV-{YYYY}-{NNNNNN}` where NNNNNN = `time.Now().UnixNano() % 1000000` — has collision risk for concurrent calls
- **Status** = `PENDING` initially
- Note: `RunBillingCycle` is a **stub scheduler** — manual invoice generation is via `POST /billing/invoices/:policyId`

### 10.40 Remittance Rules

- **Period** = `PeriodStart: now - 1 month`, `PeriodEnd: now` (auto-set, not user-specified)
- **Provider filter** = Only `ACTIVE` providers; up to **1000** providers per cycle
- **Claim aggregation** = sums `ApprovedAmount` from all approved claims for the provider
- **No claims guard** = returns error `"No approved claims for remittance"` if no approved claims exist
- **RunRemittanceCycle** is a **stub scheduler** — manual creation is via `POST /remittances/:providerId`

### 10.41 Pagination Defaults

All paginated endpoints use `shared/utils/PaginationParams`:
```
Page     = 1       (min: 1)
PageSize = 20      (min: 1, max: 100)
Sort     = "created_at"
Order    = "desc"
```
Offset formula: `(page - 1) * pageSize`. Negative page values produce negative offsets (no guard).

### 10.42 Kenyan Validation Patterns

Backend validation regex patterns for input fields:
- **Phone**: `^(?:\+254|254|0)?([17]\d{8})$` — Kenyan format, normalized to `+254XXXXXXXXX`
- **Email**: `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
- **National ID**: `^\d{7,8}$` — 7 or 8 digits

### 10.43 Document S3 Patterns & Notification Triggers

**S3 key patterns** (all under `policies/{policyID}/documents/`):

| Document Type | S3 Key | File Name |
|---|---|---|
| WELCOME_LETTER | `welcome_letter_{uuid}.pdf` | `Welcome_Letter_{policyNumber}.pdf` |
| MEMBER_CARD | `member_card_{uuid}.pdf` | `Member_Card_{memberNumber}.pdf` |
| POLICY_SCHEDULE | `policy_schedule_{uuid}.pdf` | `Policy_Schedule_{policyNumber}.pdf` |
| RENEWAL_NOTICE | `renewal_notice_{uuid}.pdf` | `Renewal_Notice_{policyNumber}.pdf` |
| LOU | `lou_{uuid}.pdf` | `LOU_{authCode}_{policyNumber}.pdf` |
| DECLINE_LETTER | `decline_letter_{uuid}.pdf` | `Decline_Letter_{claimNumber}_{policyNumber}.pdf` |

**Documents that send IN_APP notifications:** Welcome Letter, Policy Schedule, LOU, Decline Letter
**Documents that do NOT send notifications:** Member Card, Renewal Notice

### 10.44 LOU Idempotency

`GenerateLOU` is idempotent: if a LOU has already been generated for a pre-auth (detected by `LOU_{authCode}_` filename prefix), it returns the existing document with message `"Existing LOU found for this pre-authorization (generated on {date})"` and HTTP 200 instead of generating a duplicate.

### 10.45 Claim SLA Task Limits

The Claim SLA Enforcement task processes a maximum of **100 claims per cycle** per phase:
- Phase 1: Up to 100 SLA-breached claims (sends "SLA Breach Alert" IN_APP notification)
- Phase 2: Up to 100 approaching-SLA claims within 24 hours (sends "SLA Warning" notification)
- Notification send failures are silently ignored

### 10.46 Statement & Cession Integer Truncation

In `GenerateStatement` and `CalculateProfitCommission`, participant `SharePercentage` and `CommissionRate` are `float64` but cast to `int64` before division:
```go
premiumCeded := totalCeded * int64(participant.SharePercentage) / 100
```
This means a 12.5% share is truncated to 12%. Fractional percentages lose precision.

### 10.47 Event Topics

The system defines domain event topics (published via Watermill) but most are not yet wired to consumers:

| Topic | Events |
|---|---|
| Claims | `claim.submitted`, `claim.approved`, `claim.rejected`, `claim.paid` |
| Policies | `policy.activated`, `policy.lapsed`, `policy.terminated`, `policy.reinstated`, `policy.suspended`, `policy.renewed`, `policy.upgraded`, `policy.downgraded` |
| Members | `member.enrolled`, `member.removed`, `member.suspended` |
| Endorsements | `endorsement.created`, `endorsement.approved`, `endorsement.applied` |
| Renewals | `renewal.initiated`, `renewal.completed` |
| Pre-Auth | `preauth.submitted`, `preauth.approved`, `preauth.denied` |
| Sales | `lead.created`, `lead.status_changed`, `quotation.created`, `quotation.issued`, `quotation.accepted`, `quotation.converted` |
| Approvals | `approval.requested`, `approval.granted`, `approval.rejected` |
| Payments | `payment.initiated`, `payment.confirmed`, `payment.failed` |
| Documents | `document.uploaded`, `extraction.completed`, `document.generated` |

### 10.48 ExpireOverdue Treaty Behavior

`POST /treaties/expire` fetches up to **1000** ACTIVE treaties and the expiry filter condition checks `t.ExpiryDate.Before(t.CreatedAt) || t.Status == "ACTIVE"`. Since all fetched treaties are ACTIVE, the second condition always passes — effectively **all fetched active treaties are expired** regardless of their actual expiry date. Use this endpoint with caution; it is intended to be called manually by Admin users.

### 10.49 Co-Pay Has No Floor Check (Payable Can Go Negative)

In the amount calculation (Section 10.5, Step 4), after co-pay subtraction, `payableAmount` is **NOT** floored at 0. Unlike the deductible step (which has `if payableAmount < 0: payableAmount = 0`), the co-pay step simply does:
```
payableAmount -= coPayAmount
```
This means `payableAmount` can theoretically become negative if the co-pay exceeds the post-deductible payable amount. The frontend should guard against displaying negative payable amounts.

### 10.50 Fraud Pipeline — Only CheckDuplicate Affects Adjudication

The claim submission pipeline has **two distinct fraud phases**:

1. **During adjudication** (Step 8 of `RunAdjudication`): Only `CheckDuplicate` runs. If a duplicate is found, the adjudication decision is set to `MANUAL_REVIEW`. This is the **only** fraud check that influences the adjudication outcome.

2. **After adjudication** (`RunFraudChecks` in `SubmitClaim`): Six additional checks run **after** the adjudication decision has already been stored:
   - `CheckFrequency` — high procedure frequency for member
   - `CheckAmountThreshold` — amount exceeds 500K KES
   - `CheckExpiredContract` — provider contract not covering service date
   - `CheckSuspendedProvider` — provider is SUSPENDED
   - `CheckRateCardOvercharge` — unit price exceeds rate card
   - `CheckRepeatVisit` — repeat visit for same procedure

These 6 checks **only create FraudFlag records** for manual investigation. They do NOT change the claim status, adjudication decision, or payable amount. They are informational flags.

### 10.51 Waiting Period & Age Checks Against ALL Benefits

During adjudication steps 5 and 5b, the checks iterate over **ALL** plan benefits — not just the matched benefit for the claim's category:

```
for _, benefit := range allPlanBenefits:
  if serviceDateDaysSinceEnrollment < benefit.WaitingPeriodDays:
    → REJECT ("Service date within waiting period")
  if memberAge < benefit.MinAge OR memberAge > benefit.MaxAge:
    → REJECT ("Member age outside benefit age range")
```

**Implication**: A claim can be rejected due to a waiting period or age restriction defined on a benefit category that is **not relevant** to the claim. For example, an outpatient claim could be rejected because an inpatient benefit has a 90-day waiting period that hasn't elapsed yet.

### 10.52 Provider Statement Reconciliation Algorithm

`ReconcileStatement` uses a two-phase matching algorithm:

```
For each statement line item:
  Phase 1 — Match by claim number:
    If item.ClaimNumber is not empty:
      Look up claim by claim_number
      If found → MATCHED

  Phase 2 — Fallback match by provider + service date + amount:
    If Phase 1 fails:
      Search claims by provider_id + service_date + amount
      If found → MATCHED

  Phase 3 — Unmatched:
    If both phases fail → mark as UNMATCHED
```

**Amount tolerance**: 1 KES (100 cents). Discrepancies within this tolerance are treated as zero.
```
tolerance = 100 cents (1 KES)
discrepancy = item.ClaimedAmount - claim.ApprovedAmount
if abs(discrepancy) <= tolerance: discrepancy = 0
```

**Claim status update on match**:
- If `discrepancy <= 0` (provider claimed ≤ approved): claim status → `PAID`
- If `discrepancy > 0` (provider claimed > approved): claim status → `PART_PAID`

### 10.53 Case Management Transition Preconditions

Each case state transition has specific preconditions enforced by the service:

| Transition | Required Current Status | Additional Preconditions |
|---|---|---|
| Create Case | N/A (new) | Pre-auth must exist and be APPROVED |
| SCHEDULED → ADMITTED | `SCHEDULED` | `admission_date` provided in request |
| ADMITTED → IN_TREATMENT | `ADMITTED` | No additional checks (empty request body) |
| IN_TREATMENT → DISCHARGED | `ADMITTED` or `IN_TREATMENT` | `actual_discharge` datetime and `actual_cost` provided |
| DISCHARGED → CLOSED | `DISCHARGED` | No additional checks in code (manual closure) |

**Note**: Discharge is allowed from both `ADMITTED` and `IN_TREATMENT` states (e.g., early discharge before treatment starts). The close-case endpoint does NOT verify that all linked claims are in terminal status — it only checks the case is `DISCHARGED`.

### 10.54 Authentication Implementation Details

**Register** (`POST /auth/register`):
- Default role: `Member` (if `role_name` not provided)
- Default status: `ACTIVE`
- Generates member number: `MBR-YYYY-NNNNNN`
- Password is hashed via `bcrypt.GenerateFromPassword` with default cost
- Login uses `bcrypt.CompareHashAndPassword` for verification

**RefreshToken** (`POST /auth/refresh`):
- **WARNING**: Does NOT check if the user is active/suspended. A suspended or terminated user with a valid refresh token can still obtain new access tokens.
- Only validates the refresh token signature and expiration

**Logout** (`POST /auth/logout`):
- **WARNING**: This is a no-op. Returns `{"status":"success","message":"Logged out successfully"}` but does NOT invalidate the token.
- Since PASETO tokens are stateless, there is no server-side session to invalidate. The token remains valid until its `expired_at` timestamp.

**CreateUser** (`POST /users` — Admin endpoint):
- **WARNING**: Does NOT hash the password. The password is stored as provided. Only `Register` hashes passwords via bcrypt.
- Admin-created users must have their password manually hashed or use the Register endpoint instead.

### 10.55 RequirePermission Middleware Is Unused

`RequirePermission(resource, action)` middleware is defined but is **NOT wired to any route** in `routes.go`. All route-level RBAC is enforced exclusively via `RequireRole(roles...)`.

**Implications**:
- The permission-based access control system (including Admin's automatic bypass) is effectively dead code
- Admin must be **explicitly listed** in every `RequireRole()` call — there is no automatic Admin bypass for role checks
- Roles `Provider`, `Member`, and `SalesAgent` are never used in any `RequireRole()` call, meaning these roles have access to all non-role-restricted endpoints (authenticated-only)

### 10.56 Backend Route-to-Role Mapping

Complete mapping of which routes have `RequireRole` restrictions. Routes NOT listed here are accessible to **any authenticated user** (all 8 roles):

| Endpoint | Required Roles |
|---|---|
| `PUT /underwriting/:id/review` | Admin, Underwriter |
| `PUT /underwriting-flags/:id/resolve` | Admin, Underwriter |
| `PUT /underwriting-flags/:id/override` | Admin, Underwriter |
| `PUT /credit-notes/:id/approve` | Admin |
| `PUT /credit-notes/:id/apply` | Admin |
| `POST /claims/import-csv` | Admin, ClaimsOfficer |
| `PUT /claims/:id/vet` | Admin, ClaimsOfficer |
| `PUT /claims/:id/approve` | Admin, Manager |
| `PUT /claims/:id/reject` | Admin, Manager |
| `PUT /claims/:id/ready-for-payment` | Admin, Finance |
| `PUT /claims/:id/mark-paid` | Admin, Finance |
| `PUT /claims/:id/mark-part-paid` | Admin, Finance |
| `PUT /quotations/:id/versions/:v/approve` | Admin, Underwriter, Manager |
| `PUT /quotations/:id/versions/:v/reject` | Admin, Underwriter, Manager |
| `POST /renewals/expire` | Admin |
| `POST /quotations/expire` | Admin |
| All `/approval-limits/*` routes | Admin |

**All other endpoints** (plans, benefits, policies, members, providers, leads, pre-auth, billing, reinsurance, analytics, notifications, audit, etc.) require only authentication — any role can access them.

### 10.57 Benefit CheckCoverage — procedureCode Ignored

The `CheckCoverage(planID, procedureCode)` method accepts a `procedureCode` parameter but **ignores it entirely**. Coverage lookup is performed solely by:
1. Fetching all active benefits for the plan
2. Matching by benefit category (INPATIENT if `claim.AdmissionDate != nil`, else OUTPATIENT)
3. If no category match, falling back to `benefits[0]` (first active benefit)

The `procedureCode` parameter is dead code. The frontend should NOT rely on procedure-code-level coverage validation from this endpoint.

**CheckCoverage returns sub-benefits**: When a benefit is matched, its child sub-benefits (linked via `parent_benefit_id`) are included in the response. Sub-benefits provide granular breakdown (e.g., Outpatient → Lab Tests, Consultation, Pharmacy) but are NOT independently evaluated during adjudication coverage decisions.

### 10.58 Analytics Implementation Details

**parsePeriod helper** — maps period query parameter to day ranges:

| Period Value | Days | Example Range |
|---|---|---|
| `week` | 7 | last 7 days |
| `month` | 30 | last 30 days |
| `quarter` | 90 | last 90 days |
| `year` | 365 | last 365 days |
| (default/unknown) | 30 | falls back to month |

**ExportCSV** (`GET /analytics/export`):
- **STUB**: Returns a CSV file with only header columns (`report,data\n`) and no data rows
- The `reportType` and `period` parameters are accepted but ignored
- Frontend should either hide this feature or display a "coming soon" indicator

**GetReinsuranceAnalytics** (`GET /analytics/reinsurance`):
- Uses a **hardcoded** `"last_year"` (365 days) period
- The period cannot be overridden via query parameters
- Always returns data for the trailing 365-day window

### 10.59 Notification Fire-and-Forget Pattern

Notification dispatch uses a **fire-and-forget goroutine** pattern:
```go
go func() {
    ctx := context.Background()  // detached from request context
    // ... send notification ...
}()
```

**Implications for the system**:
- The goroutine is completely detached from the HTTP request context
- Notification failures or panics within the goroutine will NOT propagate to the calling request
- The caller always receives a successful response even if notification dispatch fails
- There is no backpressure mechanism — a burst of requests creates unbounded goroutines
- The frontend should NOT assume notifications are delivered just because the triggering API call succeeded
