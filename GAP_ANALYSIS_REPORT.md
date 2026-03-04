# HIAS Core — Technical Evaluation Demo Gap Analysis

> Generated: 2026-03-03
> Compared: Current codebase vs `Technical Evaluation Demo Script.xlsx`

---

## Executive Summary

The **domain layer** (entities, interfaces, schemas, DTOs) and **infrastructure layer** (migrations, SQLC, repos, cache, queue, workers, scheduler) are solidly built. What's **completely missing** is the **API Gateway layer** — there is no `main.go`, no HTTP handlers, no routes, no middleware, and no gRPC server. Without this, nothing is callable.

Beyond the API layer, the demo script reveals **significant functional gaps** in product configuration depth, business development/quotation workflow, reinsurance, reporting, and several billing/payment features that have no domain modeling at all.

---

## BLOCKER: API Gateway Layer (Not Started)

| Component | Status |
|---|---|
| `services/api-gateway/main.go` | **MISSING** — referenced in Makefile but directory doesn't exist |
| Gin HTTP route handlers | **MISSING** |
| Auth middleware (PASETO + Cognito) | **MISSING** — PASETO maker exists in `shared/auth/` but no middleware |
| RBAC middleware | **MISSING** |
| gRPC server + `.proto` files | **MISSING** |
| Swagger/OpenAPI docs | **MISSING** — `make swagger` target exists but no source |
| Service wiring / DI | **MISSING** — workers & scheduler have commented-out service calls |
| Service implementations | **MISSING** — only interfaces exist in `domains/*/service/` |

**This must be built first. Every demo scenario requires working API endpoints.**

---

## F1: Core Functionality Demo

### Scenario 1: Product Configuration — `PARTIAL`

| Demo Requirement | Status | Details |
|---|---|---|
| Create medical insurance product | ✅ Have | `Plan` entity + `CreatePlan` service interface |
| Define benefit categories (inpatient, outpatient, dental, optical, maternity) | ✅ Have | `Benefit` entity with category CHECK constraint |
| **Product segmentation (retail, corporate, SME)** | ❌ Missing | `Plan.Type` only supports `individual` / `group`. No retail/corporate/SME segmentation |
| **Capture KYC requirements (KRA Pin, location, contacts)** | ❌ Missing | No KYC entity or config. User has `NationalID` but no KRA Pin, location fields |
| Configure benefit limits and sub-limits | ⚠️ Partial | `Benefit.AnnualLimit` exists but **no sub-limits** (parent-child-grandchild benefit hierarchy) |
| **Waiting periods for maternity, chronic, surgical** | ⚠️ Partial | `Benefit.WaitingPeriodDays` is a flat integer — no condition-specific waiting periods |
| **Upload and configure rate sheets / age-banded pricing** | ❌ Missing | No age-banded pricing model. `Plan.BasePremium` is a single flat value |
| **Premium rules (per member, per family, per relation, unit rates, discounts)** | ❌ Missing | No premium calculation engine, no discount rules, no per-member/family pricing |
| **Provider restrictions per product/category/member** | ❌ Missing | No plan-provider network/restriction mapping |
| **Prorating of premiums for partial periods** | ❌ Missing | No proration logic |
| **Installment premium setup** | ❌ Missing | No installment schedule entity or logic |
| **Co-payment rules** | ⚠️ Partial | `Benefit.CoPayType` + `CoPayValue` exist. `BenefitService.CalculateCoPay()` interface exists but no implementation |
| **Age-limit validations** | ❌ Missing | No age validation rules on plans or benefits |
| **Exclusion configuration** | ✅ Have | `Exclusion` entity with types (pre_existing, cosmetic, experimental) + ICD codes |
| **Audit trail of config changes** | ✅ Have | `AuditEvent` entity + append-only DB table with triggers |

**Gap Score: ~40% covered**

### Scenario 2: Business Development — `NOT STARTED`

| Demo Requirement | Status | Details |
|---|---|---|
| **Lead capture** (probability, expected premium, follow-up) | ❌ Missing | No Lead/Opportunity entity |
| **Quotation generation** (standard or tailor-made) | ❌ Missing | No Quotation entity, no pricing engine |
| **Restricted discounting + escalation for approval** | ❌ Missing | No discount rules, no approval workflow engine |
| **Quotation versioning** (version codes, compare versions) | ❌ Missing | No versioning system |
| **Upload quotation documents** | ❌ Missing | S3 service exists but no document management entity |
| **Send quotation via SMS/email** | ⚠️ Partial | SMS (Africa's Talking) + Email (SES) adapters exist, but no quotation-specific flow |
| **Quote lifecycle statuses** (draft→issued→pending→accepted→declined→expired) | ❌ Missing | No quotation status machine |
| **Audit trail of pricing adjustments** | ⚠️ Partial | Generic audit exists but no quotation-specific audit |

**Gap Score: ~5% covered**

### Scenario 3: Policy Lifecycle Management — `PARTIAL`

| Demo Requirement | Status | Details |
|---|---|---|
| Convert quotation to policy | ❌ Missing | No quotation system to convert from |
| **Underwriting assessment + medical questionnaire** | ❌ Missing | No underwriting entity, no questionnaire model |
| Issue policy with member cards + welcome docs | ❌ Missing | No document generation |
| **Mid-term endorsement (add members)** | ⚠️ Partial | `MemberService.EnrollMember()` exists but no endorsement entity, no premium recalculation |
| **Member deletion with pro-rata refund** | ❌ Missing | `MemberRepository.Delete()` exists but no refund calculation |
| **Renewal workflow with claims experience loading** | ❌ Missing | No renewal entity, no claims-experience premium loading |
| **Policy document generation + digital distribution** | ❌ Missing | No document/template engine |
| Policy status dashboard | ⚠️ Partial | `PolicyRepository.CountByStatus()` exists |
| **Policy/member suspension/unsuspension/reinstatement** | ⚠️ Partial | `PolicyService.ReinstatePolicy()` exists. Policy has LAPSED/TERMINATED but **no SUSPENDED status**. No member suspension |
| **Mass transactions** (bulk upload, additions, deletions) | ❌ Missing | No batch/bulk endpoints |
| **Upgrade/downgrade with premium recalculation** | ❌ Missing | No plan change logic |
| **Underwriting flags** (overage, double insurance) | ❌ Missing | No underwriting rules engine |

**Gap Score: ~20% covered**

### Scenario 4: Claims Processing & Case Management — `PARTIAL`

| Demo Requirement | Status | Details |
|---|---|---|
| Provider validation (registered, accredited, active, tier, contract) | ⚠️ Partial | Provider entity has status + contract. **No tier field, no accreditation entity, no network eligibility check** |
| **Letter of Undertaking (LOU) + decline letters** | ❌ Missing | No LOU entity or generation |
| **Case manager reviews** (admitted, scheduled, discharged) | ❌ Missing | No case management entity |
| Claim data capture + auto claim number | ✅ Have | `ClaimService.SubmitClaim()` + `claim_number_generator.go` |
| Eligibility & benefit verification | ⚠️ Partial | `ValidatorService` + `AdjudicatorService` interfaces exist. `MemberService.GetMemberEligibility()` exists. Benefit limit check via `GetApprovedAmountForBenefitThisYear()` |
| Match claim to preauth | ✅ Have | `Claim.PreAuthID` FK exists |
| **Vetting of reimbursements, credit claims, exception claims** | ❌ Missing | No vetting workflow, no claim type differentiation (reimbursement vs credit vs exception) |
| Fraud detection | ✅ Have | `FraudService` with duplicate, frequency, amount threshold checks. `FraudFlag` entity with severity levels |
| **Approver reviews with prepopulated deductions** | ⚠️ Partial | `ClaimService.ApproveClaim()` / `RejectClaim()` exist but no deduction breakdown UI/API |
| **Claim status workflow** (captured→vetted→approved→ready for payment) | ⚠️ Partial | Status CHECK includes RECEIVED→VALIDATED→ADJUDICATED→APPROVED→REJECTED→MANUAL_REVIEW→PAID. **Missing VETTED status** |
| **Remittance note + payment file generation** | ⚠️ Partial | `Remittance` entity exists, `RemittanceService.SendRemittanceAdvice()` interface exists, but **no file generation** |
| **Provider statement upload + reconciliation** | ❌ Missing | No provider statement entity, no reconciliation matching |

**Gap Score: ~35% covered**

### Scenario 5: Reinsurance Processing — `NOT STARTED`

| Demo Requirement | Status | Details |
|---|---|---|
| **Treaty configuration** (Quota Share, XOL, layers, retention) | ❌ Missing | No reinsurance domain at all |
| **Cession calculation** | ❌ Missing | |
| **Profit commission rules** | ❌ Missing | |
| **Reinsurer statements** | ❌ Missing | |
| **Recovery management + reconciliation** | ❌ Missing | |
| **Treaty limit alerts** | ❌ Missing | |
| **Bordereaux generation** | ❌ Missing | |
| **Claims analytics dashboard** | ⚠️ Partial | `AnalyticsService.GetDashboard()` exists with basic KPIs |

**Gap Score: ~5% covered**

---

## F2: Reporting & Analytics Demo

### Scenario 6: Operational & Management Reporting — `MINIMAL`

| Demo Requirement | Status | Details |
|---|---|---|
| Claims experience report (loss ratio) | ⚠️ Partial | `AnalyticsRepository.GetLossRatio()` exists |
| **Claims Register, Premium Debtors Ageing, Premium Register** | ❌ Missing | No report generation system |
| **Ad-hoc report builder** | ❌ Missing | |
| Management dashboard with KPIs | ⚠️ Partial | `AnalyticsService.GetDashboard()` + `GetKPIs()` interfaces |
| **Export to Excel and PDF** | ⚠️ Partial | `AnalyticsService.ExportCSV()` interface for CSV only. No Excel/PDF |
| **Schedule recurring reports** | ❌ Missing | |
| **Drill-down from summary to detail** | ❌ Missing | |
| **Report access controls** | ❌ Missing | |
| **Membership reports with KYCs** | ❌ Missing | |

**Gap Score: ~15% covered**

---

## F3: Billing & Payment

### Scenario 7: Billing & Payment — `PARTIAL`

| Demo Requirement | Status | Details |
|---|---|---|
| M-Pesa push for premium collection | ⚠️ Partial | `MpesaAdapter` exists in `shared/integrations/mpesa/`. `PaymentService.InitiatePayment()` interface supports MPESA |
| Reconcile M-Pesa payment | ⚠️ Partial | `PaymentService.ReconcilePayment()` + `reconciliation_task.go` scheduler |
| **Card payment integration** | ❌ Missing | Only MPESA + BANK_TRANSFER methods |
| Payment reconciliation dashboard | ❌ Missing | No dashboard API |
| Payment advice with remittance details | ⚠️ Partial | `RemittanceService.SendRemittanceAdvice()` interface |
| Bulk payment processing for provider settlements | ⚠️ Partial | `RemittanceService.RunRemittanceCycle()` batches claims |
| **Detailed premium register** (credits + debits per policy) | ❌ Missing | No premium register/ledger |
| **Compute baseline commissions to intermediaries** | ❌ Missing | No broker/intermediary/commission entity |
| **Cash flow statements** | ❌ Missing | |
| **Multiple premium payments** | ❌ Missing | Invoice is 1:1 with policy billing period |
| **Double receipting detection** | ❌ Missing | `payments.reference_number` is UNIQUE but no cross-check logic |
| **Production statements** | ❌ Missing | |
| **Refunds processing** (credit balances, overpayments) | ❌ Missing | No refund entity or logic |
| **Payment frequencies** (payroll, weekly) | ❌ Missing | No payment frequency configuration |
| **Withholding taxes + regulatory computation** | ❌ Missing | No tax entity or calculation |

**Gap Score: ~20% covered**

---

## F4: Business Rules & Workflow Configuration

### Scenario 8: Business Rules & Workflow Configuration — `NOT STARTED`

| Demo Requirement | Status | Details |
|---|---|---|
| **Auto-adjudication rules (configurable)** | ❌ Missing | `AdjudicatorService` interface exists but no configurable rule engine |
| **Claims escalation workflow** (amount thresholds) | ❌ Missing | No workflow engine |
| **Automated member notifications** (claim status changes) | ⚠️ Partial | `NotificationService.Send()` + notification workers exist, but not wired to claim events |
| **Document template editor** | ❌ Missing | |
| **Approval hierarchy configuration** | ❌ Missing | RBAC exists but no approval chain/hierarchy |
| **Version control for config changes** | ❌ Missing | Audit trail logs changes but no versioning/rollback |
| **Config testing before deployment** | ❌ Missing | |

**Gap Score: ~5% covered**

---

## F6: Usability & Design

### Scenario 9: UI Walkthrough — `N/A (Backend)`

This is a frontend concern. Backend needs to **provide the APIs** for:
- Role-based dashboard data (partial — analytics exists)
- Search/filtering endpoints (missing — no search APIs)
- Personalization storage (missing)

---

## T1: Performance & Reliability

### Scenario 10: System Performance — `PARTIAL`

| Demo Requirement | Status | Details |
|---|---|---|
| **Performance metrics dashboard** | ❌ Missing | No Prometheus/metrics endpoint |
| Response times for key transactions | ❌ Missing | No API to test against |
| **Batch processing** (member uploads, payments, claim vetting) | ❌ Missing | No batch endpoints |
| **Queue management for pre-auth** | ⚠️ Partial | SQS queue + `preauth_submitted_handler.go` exists |
| **DB query performance for large datasets** | ⚠️ Partial | SQLC queries exist with pagination |
| **System alerting and monitoring** | ❌ Missing | No health check, no metrics |

**Gap Score: ~15% covered**

---

## T2: Security & Access Control

### Scenario 11: Security Features — `PARTIAL`

| Demo Requirement | Status | Details |
|---|---|---|
| Role-based access control | ✅ Have | `roles`, `permissions`, `role_permissions` tables + entities + repos |
| **Privilege-based access** | ⚠️ Partial | Permission model exists (resource + action) but no middleware enforcement |
| **Multi-factor authentication** | ⚠️ Partial | Cognito supports MFA but not configured/enforced in code |
| Audit trail for sensitive data access | ✅ Have | `AuditEvent` entity + append-only table |
| **Data masking for PII** | ❌ Missing | No data masking logic |
| **Encryption at rest and in transit** | ⚠️ Partial | TLS is infrastructure-level. DB encryption depends on deployment. No application-level encryption |
| **Session management and timeout** | ⚠️ Partial | PASETO tokens exist with expiry. No session invalidation/revocation |
| **User access review workflow** | ❌ Missing | |

**Gap Score: ~35% covered**

---

## T3: Integration & Data Flow

### Scenario 12: System Integration — `PARTIAL`

| Demo Requirement | Status | Details |
|---|---|---|
| **Provider portal** | ❌ Missing | No provider-facing API or portal |
| Pre-auth submission via portal | ❌ Missing | |
| Real-time eligibility check via API | ⚠️ Partial | `MemberService.GetMemberEligibility()` interface exists |
| Electronic claim submission from provider | ⚠️ Partial | `ClaimService.SubmitClaim()` interface |
| **IPRS integration** | ⚠️ Partial | `IPRSAdapter` in `shared/integrations/iprs/` exists but likely stubbed |
| **API documentation (Swagger/OpenAPI)** | ❌ Missing | No endpoints to document |
| **Live API demo via Postman** | ❌ Missing | No running API |
| **API architecture docs** (patterns, retry, circuit breaker, versioning) | ❌ Missing | |
| **Rate limiting** | ❌ Missing | No rate limiter middleware |
| **Error handling with consistent error bodies** | ⚠️ Partial | `ServiceResponse` pattern exists but no HTTP error envelope |

**Gap Score: ~15% covered**

---

## T4: Microservice Architecture

### Scenario 13: Microservice Check — `DOES NOT APPLY`

| Demo Requirement | Status | Details |
|---|---|---|
| **Service failure isolation** | ❌ N/A | Currently a **monolith** — single Go binary |
| **Independent data ownership per service** | ❌ N/A | Single `hias_db` database |
| **Independent scaling** | ❌ N/A | |
| **Service recovery** | ❌ N/A | |

The architecture is DDD + Clean Architecture within a **monolith**. The demo asks to prove microservice isolation. You'll need to decide: refactor into microservices or demonstrate graceful degradation within the monolith (e.g., notifications failing doesn't crash claims).

---

## T5: System Administration & Health Check

### Scenario 14: System Administration — `PARTIAL`

| Demo Requirement | Status | Details |
|---|---|---|
| User & access management (create, modify, disable) | ✅ Have | `UserService` with full CRUD + status management |
| **System health dashboard** | ❌ Missing | No `/health` endpoint, no metrics |
| **Alerts for failed logins, system errors** | ❌ Missing | |
| **Usage metrics** | ❌ Missing | |
| **Disaster recovery / failover** | ❌ Missing | Infrastructure concern, but no health checks |
| Audit trails (logins, data changes, approvals) | ✅ Have | `AuditEvent` with tamper-proof triggers |
| **Resource monitoring** (CPU, memory, storage, DB health) | ❌ Missing | No observability setup |
| **CI/CD pipeline demo** | ❌ Missing | No GitHub Actions, no pipeline config |
| **Infrastructure-as-Code** | ⚠️ Partial | CloudFormation YAML for 5 SQS queues exists. No full IaC for DB, Redis, etc. |

**Gap Score: ~25% covered**

---

## T6: Data Management

### Scenario 15: Data Management — `MINIMAL`

| Demo Requirement | Status | Details |
|---|---|---|
| **Golden record / deduplication** | ❌ Missing | `users.email` and `members.national_id` are unique, but no merge/dedup workflow |
| **Data classification** (public, internal, confidential) | ❌ Missing | |
| **Retention and archival** | ❌ Missing | |
| **Backup and restore** | ❌ Missing | Infrastructure concern, no application support |
| **Data quality validation** | ⚠️ Partial | `shared/utils/validators.go` exists. DB CHECK constraints exist |
| **Consent management** | ❌ Missing | No consent entity |
| **Document management** | ❌ Missing | S3 service exists but no document entity/metadata |

**Gap Score: ~10% covered**

---

## T7: Innovation — AI & RPA

### Scenario 16: AI/ML & Automation — `MINIMAL`

| Demo Requirement | Status | Details |
|---|---|---|
| AI-powered fraud detection | ⚠️ Partial | Rule-based `FraudService` (duplicate, frequency, threshold). **Not AI/ML** |
| **Intelligent claims routing** (complexity scoring) | ❌ Missing | |
| **OCR/ICR for claim documents** | ⚠️ Partial | `extraction_result_handler.go` worker exists for document extraction results. Implies external OCR service |
| **Predictive analytics** (high-risk members) | ❌ Missing | |
| **Chatbot / virtual assistant** | ❌ Missing | |
| **RPA automations** | ❌ Missing | |
| **Recommendation engine** | ❌ Missing | |
| **NLP for unstructured data** | ❌ Missing | |

**Gap Score: ~10% covered**

---

## Missing Domains (Entirely New)

These are **entirely new domain aggregates** that need to be created from scratch:

| Domain | Why Needed | Entities to Create |
|---|---|---|
| **Quotation** | Scenario 2 (Business Development) | `Lead`, `Quotation`, `QuotationVersion`, `QuotationDocument` |
| **Reinsurance** | Scenario 5 | `Treaty`, `TreatyParticipant`, `Cession`, `Recovery`, `Bordereaux` |
| **Underwriting** | Scenarios 1, 3, 12 | `UnderwritingRule`, `MedicalQuestionnaire`, `UnderwritingDecision` |
| **Document** | Scenarios 2, 3, 4, 6 | `Document`, `DocumentTemplate`, `GeneratedDocument` |
| **Commission** | Scenario 7 | `Intermediary`, `CommissionRule`, `CommissionPayment` |
| **Endorsement** | Scenario 3 | `Endorsement`, `EndorsementLineItem` |
| **Workflow/Rules Engine** | Scenario 8 | `WorkflowDefinition`, `WorkflowInstance`, `ApprovalStep`, `BusinessRule` |
| **Consent** | Scenario 15 (Data Management) | `ConsentRecord`, `ConsentPurpose` |

---

## Priority Implementation Roadmap

### Phase 0 — MUST HAVE (Demo is impossible without these)
1. **API Gateway** — `main.go`, Gin routes, middleware (auth, RBAC, audit, error handling), service wiring/DI, health check endpoint
2. **Service Implementations** — Concrete implementations of all 19 service interfaces
3. **Swagger/OpenAPI** — Auto-generated API docs

### Phase 1 — Core Demo Scenarios (F1)
4. **Enhanced Product Configuration** — age-banded pricing, sub-limits, provider networks, premium rules, product segmentation
5. **Policy Lifecycle Enhancements** — suspension status, endorsements, renewal, bulk operations, pro-rata calculations
6. **Claims Processing Completion** — vetting workflow, LOU generation, case management, claim type differentiation

### Phase 2 — Revenue & Finance (F3)
7. **Billing Enhancements** — premium register, installments, payment frequencies, tax calculations, refunds, double-receipt detection
8. **Commission Management** — intermediary/broker entity, commission rules, commission payments

### Phase 3 — New Domains
9. **Quotation/Business Development** — leads, quotations, versioning, approval workflows
10. **Reinsurance** — treaties, cessions, recoveries, bordereaux
11. **Document Management** — templates, generation (PDF), S3 storage with metadata

### Phase 4 — Technical Demo Requirements (T1-T7)
12. **Reporting Engine** — report builder, Excel/PDF export, scheduling, drill-down
13. **Workflow/Rules Engine** — configurable business rules, approval hierarchies
14. **Observability** — health checks, Prometheus metrics, structured logging
15. **CI/CD Pipeline** — GitHub Actions, deployment gates, IaC
16. **AI/ML Features** — ML fraud scoring, claims routing, OCR integration, predictive analytics

---

## Coverage Summary

| Scenario | Category | Coverage |
|---|---|---|
| S1: Product Configuration | F1 | ~40% |
| S2: Business Development | F1 | ~5% |
| S3: Policy Lifecycle | F1 | ~20% |
| S4: Claims Processing | F1 | ~35% |
| S5: Reinsurance | F1 | ~5% |
| S6: Reporting & Analytics | F2 | ~15% |
| S7: Billing & Payment | F3 | ~20% |
| S8: Business Rules & Workflow | F4 | ~5% |
| S9: UI Walkthrough | F6 | N/A (frontend) |
| S10: Performance | T1 | ~15% |
| S11: Security & Access | T2 | ~35% |
| S12: Integration | T3 | ~15% |
| S13: Microservice Check | T4 | ~0% |
| S14: System Admin | T5 | ~25% |
| S15: Data Management | T6 | ~10% |
| S16: AI/ML & Automation | T7 | ~10% |
| **Overall Estimated Coverage** | | **~17%** |

---

## What You DID Right

The foundation is genuinely solid:

- **Clean DDD layering** — domain interfaces are pure, infrastructure is separate
- **21 entities** with proper relationships and FK constraints
- **22 repository implementations** backed by SQLC-generated type-safe queries
- **11 well-structured migrations** with proper CHECK constraints and indexes
- **Event-driven architecture** — SQS queues, Watermill publishers/consumers, 6 worker handlers
- **8 scheduled tasks** covering billing cycles, payment retries, policy lapsing, etc.
- **Real-time infrastructure** — both SSE and WebSocket hubs ready
- **External integrations** — M-Pesa, IPRS, SMART/Slade360, Bank, Africa's Talking SMS, SES email
- **Security foundation** — PASETO tokens, Cognito integration, RBAC tables, append-only audit
- **Proper conventions** — money in cents, UUID PKs, typed constants, ServiceResponse[T] generic pattern

The architecture supports rapid feature development. The gap is primarily in **wiring it together** (API gateway) and **depth of business logic** (the demo wants insurance-domain-specific features that go beyond basic CRUD).