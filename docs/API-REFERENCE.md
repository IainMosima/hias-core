# HIAS Core — API Reference

**Base URL:** `/api/v1`
**Auth:** All endpoints require `Authorization: Bearer <PASETO_TOKEN>` unless marked **Public**.
**Money:** All monetary values are in **cents** (int64). `8000 KES = 800000`.

**Pagination defaults** (all paginated endpoints): `page=1` (min 1), `page_size=20` (min 1, max 100), `sort=created_at`, `order=desc`. Offset = `(page - 1) * page_size`.

**Validation patterns** (backend-enforced):
- Phone: `^(?:\+254|254|0)?([17]\d{8})$` — Kenyan format, normalized to `+254XXXXXXXXX`
- Email: `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
- National ID: `^\d{7,8}$` — 7 or 8 digits

---

## Table of Contents

1. [Auth](#1-auth)
2. [Users](#2-users)
3. [Plans](#3-plans)
4. [Benefits](#4-benefits)
5. [Exclusions](#5-exclusions)
6. [Premium Rules](#6-premium-rules)
7. [Underwriting Rules](#7-underwriting-rules)
8. [Provider Networks](#8-provider-networks)
9. [Providers](#9-providers)
10. [Leads](#10-leads)
11. [Quotations](#11-quotations)
12. [Approval Limits](#12-approval-limits)
13. [Policies](#13-policies)
14. [Members](#14-members)
15. [Endorsements](#15-endorsements)
16. [Renewals](#16-renewals)
17. [Underwriting Assessments](#17-underwriting-assessments)
18. [Underwriting Flags](#18-underwriting-flags)
19. [Policy Documents](#19-policy-documents)
20. [Credit Notes](#20-credit-notes)
21. [Pre-Authorization](#21-pre-authorization)
22. [Claims](#22-claims)
23. [Cases](#23-cases)
24. [Claim Documents](#24-claim-documents)
25. [Provider Statements](#25-provider-statements)
26. [Installments](#26-installments)
27. [Invoices](#27-invoices)
28. [Payments](#28-payments)
29. [Remittances](#29-remittances)
30. [Notifications](#30-notifications)
31. [Audit](#31-audit)
32. [Analytics](#32-analytics)
33. [Treaties](#33-treaties)
34. [Cessions](#34-cessions)
35. [Recoveries](#35-recoveries)
36. [Bordereaux](#36-bordereaux)
37. [Reinsurer Statements](#37-reinsurer-statements)
38. [Treaty Alerts](#38-treaty-alerts)

---

## 1. Auth

### POST /auth/login
- **Auth:** Public
- **Request Body:**
  ```json
  {
    "email": "string (required, valid email)",
    "password": "string (required, min 8 chars)"
  }
  ```
- **Response (200):**
  ```json
  {
    "status": "success",
    "data": {
      "access_token": "string",
      "access_token_expires_at": "datetime",
      "refresh_token": "string",
      "user": {
        "id": "uuid",
        "email": "string",
        "name": "string",
        "phone": "string",
        "national_id": "string",
        "role_id": "uuid",
        "role_name": "string",
        "status": "ACTIVE|INACTIVE|SUSPENDED",
        "created_at": "datetime",
        "updated_at": "datetime"
      }
    }
  }
  ```
- **Business Rules:** Validates credentials against stored hash. Returns PASETO access + refresh tokens. Access token contains user_id, email, role, and permissions array.
- **Status Codes:** 200 success, 400 invalid input, 401 invalid credentials

### POST /auth/register
- **Auth:** Public
- **Request Body:**
  ```json
  {
    "email": "string (required, valid email)",
    "password": "string (required, min 8 chars)",
    "name": "string (required)",
    "phone": "string (required)",
    "national_id": "string (optional)",
    "role_name": "string (optional)"
  }
  ```
- **Response (201):**
  ```json
  {
    "status": "success",
    "data": {
      "user_id": "string",
      "email": "string",
      "message": "string"
    }
  }
  ```
- **Business Rules:** Creates user with hashed password (bcrypt). Default role: **`Member`** if role_name not specified. Default status: **`ACTIVE`**. Generates member number `MBR-YYYY-NNNNNN`. Email must be unique.
- **Status Codes:** 201 created, 400 validation error, 409 email already exists

### POST /auth/refresh
- **Auth:** Public
- **Request Body:**
  ```json
  {
    "refresh_token": "string (required)"
  }
  ```
- **Response (200):** New access_token + expiry
- **Business Rules:** Validates refresh token. Issues new access token with same payload. Refresh token remains valid until its own expiry. **WARNING: Does NOT check if user is active/suspended.** A suspended or terminated user with a valid refresh token can still obtain new access tokens.
- **Status Codes:** 200 success, 401 invalid/expired refresh token

### POST /auth/logout
- **Auth:** Required
- **Request Body:** None
- **Response (200):** `{"status":"success","message":"Logged out successfully"}`
- **Business Rules:** **WARNING: This is a no-op.** Returns success but does NOT invalidate the token. Since PASETO tokens are stateless, there is no server-side session to invalidate. The token remains valid until its `expired_at` timestamp. Frontend should clear stored tokens client-side on logout.
- **Status Codes:** 200 success, 401 unauthorized

---

## 2. Users

### GET /users
- **Auth:** Required
- **Query Params:** `page` (int, min 1), `page_size` (int, min 1, max 100)
- **Response (200):** Paginated list of UserResponse
- **Status Codes:** 200 success, 401 unauthorized

### GET /users/:id
- **Auth:** Required
- **Response (200):** Single UserResponse
- **Status Codes:** 200 success, 404 not found

### POST /users
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "email": "string (required, valid email)",
    "name": "string (required)",
    "phone": "string (required)",
    "national_id": "string (optional)",
    "role_name": "string (required)",
    "password": "string (required, min 8 chars)"
  }
  ```
- **Response (201):** UserResponse
- **Business Rules:** Admin creates users with specific roles. Email must be unique. **WARNING: Does NOT hash the password** — the password is stored as provided. Only `POST /auth/register` hashes passwords via bcrypt. Admin-created users should use the Register endpoint or have their password pre-hashed.
- **Status Codes:** 201 created, 400 validation error

### PUT /users/:id
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "name": "string (optional)",
    "phone": "string (optional)",
    "national_id": "string (optional)"
  }
  ```
- **Response (200):** Updated UserResponse
- **Status Codes:** 200 success, 404 not found

### PUT /users/:id/role
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "role_id": "uuid (required)"
  }
  ```
- **Response (200):** Updated UserResponse
- **Business Rules:** Changes user's role. Affects all future token payloads.
- **Status Codes:** 200 success, 404 not found

### PUT /users/:id/status
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "status": "ACTIVE|INACTIVE|SUSPENDED (required)"
  }
  ```
- **Response (200):** Updated UserResponse
- **Business Rules:** SUSPENDED users cannot log in. INACTIVE users are soft-disabled.
- **Status Codes:** 200 success, 400 invalid status, 404 not found

---

## 3. Plans

### GET /plans
- **Auth:** Required
- **Query Params:** `page`, `page_size`
- **Response (200):** Paginated list of PlanResponse
  ```json
  {
    "id": "uuid",
    "name": "string",
    "type": "individual|group",
    "segment": "retail|corporate|sme",
    "base_premium": 800000,
    "currency": "KES",
    "status": "ACTIVE|INACTIVE",
    "description": "string",
    "created_at": "datetime",
    "updated_at": "datetime"
  }
  ```

### GET /plans/:id
- **Auth:** Required
- **Response (200):** Single PlanResponse

### POST /plans
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "name": "string (required)",
    "type": "individual|group (required)",
    "segment": "retail|corporate|sme (optional)",
    "base_premium": 800000,
    "currency": "KES (optional, defaults to KES)",
    "description": "string (optional)"
  }
  ```
- **Response (201):** PlanResponse
- **Business Rules:** base_premium must be > 0 (in cents). Default status is ACTIVE. Currency defaults to KES.

### PUT /plans/:id
- **Auth:** Required
- **Request Body:** All fields optional (partial update)
  ```json
  {
    "name": "string",
    "type": "individual|group",
    "segment": "retail|corporate|sme",
    "base_premium": 900000,
    "description": "string",
    "status": "ACTIVE|INACTIVE"
  }
  ```
- **Response (200):** Updated PlanResponse

---

## 4. Benefits

### GET /plans/:id/benefits
- **Auth:** Required
- **Response (200):** List of BenefitResponse (top-level benefits for plan)
  ```json
  {
    "id": "uuid",
    "plan_id": "uuid",
    "parent_benefit_id": null,
    "name": "Outpatient Cover",
    "category": "outpatient|inpatient|dental|optical|maternity",
    "annual_limit": 50000000,
    "co_pay_type": "percentage|fixed",
    "co_pay_value": 2000,
    "waiting_period_days": 30,
    "sub_limit_type": "none|per_visit|per_item",
    "sub_limit_value": 500000,
    "min_age": 0,
    "max_age": 65,
    "waiting_period_type": "general|maternity|pre_existing|chronic|surgical",
    "deductible_amount": 100000,
    "created_at": "datetime"
  }
  ```

### POST /plans/:id/benefits
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "parent_benefit_id": "uuid (optional, for sub-benefits)",
    "name": "string (required)",
    "category": "outpatient|inpatient|dental|optical|maternity (required)",
    "annual_limit": 50000000,
    "co_pay_type": "percentage|fixed (required)",
    "co_pay_value": 2000,
    "waiting_period_days": 30,
    "sub_limit_type": "none|per_visit|per_item",
    "sub_limit_value": 500000,
    "min_age": 0,
    "max_age": 65,
    "waiting_period_type": "general|maternity|pre_existing|chronic|surgical",
    "deductible_amount": 100000
  }
  ```
- **Business Rules:**
  - `annual_limit`: Maximum amount covered per year (cents). Hard stop in adjudication.
  - `co_pay_type=percentage`: Member pays `co_pay_value`% of approved amount. 2000 = 20%.
  - `co_pay_type=fixed`: Member pays fixed `co_pay_value` cents per claim.
  - `deductible_amount`: Amount deducted before co-pay calculation.
  - `waiting_period_days`: Days after enrollment before benefit becomes active.
  - `sub_limit_type=per_visit`: Caps each claim at `sub_limit_value` for this benefit.
  - `sub_limit_type=per_item`: Caps each line item at `sub_limit_value`.
  - `min_age/max_age`: Member age eligibility (0/0 means no restriction).

### GET /benefits/:id/sub-benefits
- **Auth:** Required
- **Response (200):** List of BenefitResponse (children of parent benefit)

### POST /benefits/:id/sub-benefits
- **Auth:** Required
- **Request Body:** Same as POST /plans/:id/benefits (parent_benefit_id auto-set)
- **Business Rules:** Creates a child benefit under the specified parent. Inherits plan_id from parent. Sub-benefits form a hierarchical tree for detailed coverage breakdown.

---

## 5. Exclusions

### GET /plans/:id/exclusions
- **Auth:** Required
- **Response (200):** List of ExclusionResponse
  ```json
  {
    "id": "uuid",
    "plan_id": "uuid",
    "description": "Pre-existing conditions diagnosed before enrollment",
    "type": "pre_existing|cosmetic|experimental",
    "icd_codes": ["Z00", "Z01"],
    "created_at": "datetime"
  }
  ```

### POST /plans/:id/exclusions
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "description": "string (required)",
    "type": "pre_existing|cosmetic|experimental (required)",
    "icd_codes": ["string array (optional)"]
  }
  ```
- **Business Rules:** During adjudication, claim diagnosis codes are matched against exclusion ICD codes. If a match is found, the claim is **rejected** with reason "Excluded condition".

### PUT /exclusions/:id
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "description": "string (optional)",
    "type": "pre_existing|cosmetic|experimental (optional)",
    "icd_codes": ["string array (optional)"]
  }
  ```

### DELETE /exclusions/:id
- **Auth:** Required
- **Response (200):** Success message

---

## 6. Premium Rules

### GET /plans/:id/premium-rules
- **Auth:** Required
- **Response (200):** List of PremiumRuleResponse
  ```json
  {
    "id": "uuid",
    "plan_id": "uuid",
    "calculation_type": "per_member|per_family|tiered|flat",
    "relationship": "principal|spouse|child|parent",
    "rate_amount": 250000,
    "discount_type": "percentage|fixed",
    "discount_value": 5000,
    "min_members": 5,
    "min_age": 0,
    "max_age": 18,
    "created_at": "datetime",
    "updated_at": "datetime"
  }
  ```

### POST /plans/:id/premium-rules
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "calculation_type": "per_member|per_family|tiered|flat (required)",
    "relationship": "principal|spouse|child|parent (optional, empty=generic)",
    "rate_amount": 250000,
    "discount_type": "percentage|fixed (optional)",
    "discount_value": 5000,
    "min_members": 5,
    "min_age": 0,
    "max_age": 18
  }
  ```
- **Business Rules:**
  - **per_member**: Rate applied per individual member. Matching priority:
    1. Relationship + age range match
    2. Relationship-only match
    3. Generic rule (no relationship) + age range
    4. Generic rule (no relationship, no age)
    5. Fallback to plan's base_premium
  - **per_family**: Flat rate based on family size. Rule with highest `min_members ≤ memberCount` wins. **Short-circuit:** If any rule for the plan has `calculation_type == "per_family"`, per-member calculation is skipped entirely — the matched family rule's `rate_amount` is returned as the total premium.
  - **tiered**: Tier-based pricing by member count.
  - **Discounts**: Applied when `memberCount >= min_members`
    - `percentage`: `discount = total × discount_value / 10000` (basis points: 5000 = 50%)
    - `fixed`: `discount = discount_value` (cents)
  - **Age bands**: `min_age` and `max_age` define the age range. Default max_age is 150 if not specified.

### POST /plans/:id/calculate-premium
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "relationships": ["principal", "spouse", "child"],
    "proposed_members": [
      {
        "relationship": "principal",
        "date_of_birth": "1990-01-15"
      },
      {
        "relationship": "child",
        "date_of_birth": "2015-06-20"
      }
    ]
  }
  ```
- **Response (200):**
  ```json
  {
    "status": "success",
    "data": {
      "total_premium": 1500000,
      "breakdown": { ... }
    }
  }
  ```
- **Business Rules:**
  - Accepts either `relationships` array (simple) or `proposed_members` with DOB (age-band matching)
  - Calculates age from DOB: `years = now.Year() - dob.Year()` with YearDay adjustment
  - Returns total premium after discount/loading application
  - Premium can never be negative (floored at 0)

### GET /plans/:id/rate-sheet
- **Auth:** Required
- **Response (200):** Complete rate sheet with all rules for the plan, organized by relationship and age band

### DELETE /premium-rules/:id
- **Auth:** Required
- **Response (200):** Success message

---

## 7. Underwriting Rules

### GET /plans/:id/underwriting-rules
- **Auth:** Required
- **Response (200):** List of UnderwritingRuleResponse
  ```json
  {
    "id": "uuid",
    "plan_id": "uuid",
    "rule_type": "MAX_AGE|MIN_AGE|DOUBLE_INSURANCE|PRE_EXISTING_CONDITION|BMI_THRESHOLD|WAITING_PERIOD",
    "relationship": "principal|spouse|child|parent",
    "parameter_key": "max_age",
    "parameter_value": "65",
    "severity": "LOW|MEDIUM|HIGH",
    "risk_score_weight": 10,
    "is_blocking": false,
    "is_active": true,
    "description": "Maximum entry age",
    "created_at": "datetime",
    "updated_at": "datetime"
  }
  ```

### POST /plans/:id/underwriting-rules
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "rule_type": "MAX_AGE|MIN_AGE|... (required)",
    "relationship": "string (optional, empty=all relationships)",
    "parameter_key": "string (required)",
    "parameter_value": "string (required)",
    "severity": "LOW|MEDIUM|HIGH",
    "risk_score_weight": 10,
    "is_blocking": false,
    "is_active": true,
    "description": "string"
  }
  ```
- **Business Rules:** Rules are evaluated during member enrollment. If `is_blocking=true`, enrollment is rejected. Otherwise, an underwriting flag is created for manual review. Rules can be relationship-specific (e.g., MAX_AGE only for "child").

### PUT /underwriting-rules/:id
- **Auth:** Required
- **Request Body:** Partial update, all fields optional

### DELETE /underwriting-rules/:id
- **Auth:** Required

---

## 8. Provider Networks

### GET /plans/:id/provider-networks
- **Auth:** Required
- **Response (200):** List of ProviderNetworkResponse
  ```json
  {
    "id": "uuid",
    "plan_id": "uuid",
    "provider_id": "uuid",
    "benefit_category": "outpatient|inpatient|dental|optical|maternity",
    "status": "ACTIVE|INACTIVE",
    "created_at": "datetime",
    "updated_at": "datetime"
  }
  ```

### POST /plans/:id/provider-networks
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "provider_id": "uuid (required)",
    "benefit_category": "string (optional)"
  }
  ```
- **Business Rules:** Associates a provider with a plan. During adjudication, the system checks if the claim's provider is in the plan's network. Claims from out-of-network providers are **rejected**.

### PUT /provider-networks/:id/status
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "status": "ACTIVE|INACTIVE (required)"
  }
  ```

### DELETE /provider-networks/:id
- **Auth:** Required

---

## 9. Providers

### GET /providers
- **Auth:** Required
- **Query Params:** `page`, `page_size`
- **Response (200):** Paginated list of ProviderResponse
  ```json
  {
    "id": "uuid",
    "name": "string",
    "type": "hospital|clinic|pharmacy|lab",
    "license_number": "string",
    "status": "PENDING|CREDENTIALING|ACTIVE|SUSPENDED|TERMINATED",
    "tier": "TIER_1|TIER_2|TIER_3",
    "county": "string",
    "phone": "string",
    "email": "string",
    "contact_person": "string",
    "accreditation_status": "NONE|PENDING|ACCREDITED|EXPIRED|REVOKED",
    "accreditation_expiry": "datetime",
    "accreditation_body": "string",
    "created_at": "datetime",
    "updated_at": "datetime"
  }
  ```

### GET /providers/by-tier
- **Auth:** Required
- **Query Params:** `tier` (TIER_1|TIER_2|TIER_3), `page`, `page_size`

### GET /providers/by-accreditation
- **Auth:** Required
- **Query Params:** `status` (accreditation status), `page`, `page_size`

### GET /providers/expiring-accreditations
- **Auth:** Required
- **Query Params:** `days` (expiring within N days), `page`, `page_size`

### GET /providers/:id
- **Auth:** Required
- **Response (200):** Single ProviderResponse

### POST /providers
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "name": "string (required)",
    "type": "hospital|clinic|pharmacy|lab (required)",
    "license_number": "string (required)",
    "county": "string",
    "address": "string",
    "phone": "string (required)",
    "email": "string (required, valid email)",
    "contact_person": "string"
  }
  ```
- **Business Rules:** Initial status is PENDING. Must go through credentialing before activation. Tier defaults to TIER_3.

### PUT /providers/:id
- **Auth:** Required
- **Request Body:** Partial update (name, county, address, phone, email, contact_person)

### PUT /providers/:id/credential
- **Auth:** Required
- **Business Rules:** PENDING → CREDENTIALING. Initiates credentialing process.

### PUT /providers/:id/activate
- **Auth:** Required
- **Business Rules:** CREDENTIALING → ACTIVE. Provider can now receive claims.

### PUT /providers/:id/suspend
- **Auth:** Required
- **Business Rules:** ACTIVE → SUSPENDED. Claims from suspended providers trigger fraud flags (CRITICAL severity).

### PUT /providers/:id/terminate
- **Auth:** Required
- **Business Rules:** Any → TERMINATED. Permanent deactivation.

### PUT /providers/:id/tier
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "tier": "TIER_1|TIER_2|TIER_3 (required)"
  }
  ```
- **Business Rules:** TIER_1 = premium, TIER_2 = standard, TIER_3 = basic. Affects rate card matching.

### PUT /providers/:id/accreditation
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "accreditation_status": "NONE|PENDING|ACCREDITED|EXPIRED|REVOKED (required)",
    "accreditation_expiry": "YYYY-MM-DD (optional)",
    "accreditation_body": "string (optional)"
  }
  ```
- **Business Rules:** During adjudication, non-ACCREDITED providers trigger a FLAG (not automatic rejection). Expired accreditations are tracked for compliance.

### GET /providers/:id/contracts
- **Auth:** Required
- **Response (200):** List of ContractResponse
  ```json
  {
    "id": "uuid",
    "provider_id": "uuid",
    "start_date": "datetime",
    "end_date": "datetime",
    "terms": "string",
    "status": "ACTIVE|EXPIRED|TERMINATED",
    "created_at": "datetime"
  }
  ```

### POST /providers/:id/contracts
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "start_date": "datetime (required)",
    "end_date": "datetime (required)",
    "terms": "string"
  }
  ```
- **Business Rules:** During adjudication, provider must have an active contract covering the service date. No active contract → claim **rejected**.

### GET /providers/:id/rate-cards
- **Auth:** Required
- **Response (200):** List of RateCardResponse
  ```json
  {
    "id": "uuid",
    "provider_id": "uuid",
    "procedure_code": "string",
    "procedure_name": "string",
    "rate_amount": 500000,
    "effective_date": "datetime",
    "age_from": 0,
    "age_to": 65,
    "gender": "male|female",
    "relationship": "string"
  }
  ```

### POST /providers/:id/rate-cards
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "procedure_code": "string (required)",
    "procedure_name": "string (required)",
    "rate_amount": 500000,
    "effective_date": "datetime",
    "age_from": 0,
    "age_to": 65,
    "gender": "string",
    "relationship": "string"
  }
  ```
- **Business Rules:** Rate cards define the expected cost per procedure. During fraud checks, if a claim line item's unit_price exceeds the rate card amount → RATE_CARD_OVERCHARGE fraud flag.

### POST /providers/:id/rate-cards/bulk
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "rate_cards": [{ ...CreateRateCardRequest }]
  }
  ```

---

## 10. Leads

### GET /leads
- **Auth:** Required
- **Query Params:** `page`, `page_size`, `status` (optional filter)
- **Response (200):** Paginated list of LeadResponse
  ```json
  {
    "id": "uuid",
    "lead_number": "string",
    "contact_name": "string",
    "contact_email": "string",
    "contact_phone": "string",
    "company_name": "string",
    "source": "direct|referral|web|agent|broker",
    "segment": "retail|corporate|sme",
    "plan_type": "individual|group",
    "estimated_members": 10,
    "expected_premium": 5000000,
    "closure_probability": 75,
    "currency": "KES",
    "status": "NEW|CONTACTED|QUALIFIED|PROPOSAL_SENT|NEGOTIATION|WON|LOST|DORMANT",
    "assigned_to": "uuid",
    "next_follow_up_date": "datetime",
    "notes": "string",
    "created_by": "uuid",
    "created_at": "datetime",
    "updated_at": "datetime"
  }
  ```

### POST /leads
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "contact_name": "string (required)",
    "contact_email": "string",
    "contact_phone": "string",
    "company_name": "string",
    "source": "direct|referral|web|agent|broker (required)",
    "segment": "retail|corporate|sme (required)",
    "plan_type": "individual|group (required)",
    "estimated_members": 10,
    "expected_premium": 5000000,
    "closure_probability": 75,
    "next_follow_up_date": "datetime",
    "notes": "string"
  }
  ```
- **Business Rules:** Initial status is NEW. Assigned to creating user. Lead number auto-generated.

### GET /leads/due-follow-ups
- **Auth:** Required
- **Response (200):** List of leads with next_follow_up_date <= today

### GET /leads/:id
- **Auth:** Required

### PUT /leads/:id
- **Auth:** Required
- **Request Body:** Partial update of all lead fields

### PUT /leads/:id/status
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "status": "NEW|CONTACTED|QUALIFIED|PROPOSAL_SENT|NEGOTIATION|WON|LOST|DORMANT (required)"
  }
  ```
- **Business Rules:**
  - Lead status pipeline: NEW → CONTACTED → QUALIFIED → PROPOSAL_SENT → NEGOTIATION → WON/LOST
  - DORMANT can be set from any state
  - WON and LOST **cannot** transition back to NEW (error: "Cannot transition from {old} to {new}")
  - Auto-creates a status_change activity log entry

### GET /leads/:id/activities
- **Auth:** Required
- **Response (200):** List of LeadActivityResponse
  ```json
  {
    "id": "uuid",
    "lead_id": "uuid",
    "activity_type": "call|email|meeting|note|follow_up",
    "description": "string",
    "scheduled_at": "datetime",
    "completed_at": "datetime",
    "created_by": "uuid",
    "created_at": "datetime"
  }
  ```

### POST /leads/:id/activities
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "activity_type": "call|email|meeting|note|follow_up (required)",
    "description": "string",
    "scheduled_at": "datetime",
    "completed_at": "datetime"
  }
  ```

### GET /leads/:id/quotations
- **Auth:** Required
- **Response (200):** List of quotations linked to this lead

---

## 11. Quotations

### GET /quotations
- **Auth:** Required
- **Query Params:** `page`, `page_size`, `status` (optional)
- **Response (200):** Paginated list of QuotationResponse
  ```json
  {
    "id": "uuid",
    "quotation_number": "string",
    "lead_id": "uuid",
    "plan_id": "uuid",
    "quotation_type": "standard|tailor_made",
    "status": "DRAFT|ISSUED|PENDING_DECISION|ACCEPTED|DECLINED|EXPIRED|CONVERTED",
    "current_version": 1,
    "policy_id": "uuid (if converted)",
    "valid_from": "datetime",
    "valid_until": "datetime",
    "client_name": "string",
    "client_email": "string",
    "client_phone": "string",
    "currency": "KES",
    "created_by": "uuid",
    "created_at": "datetime",
    "updated_at": "datetime"
  }
  ```

### POST /quotations
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "lead_id": "uuid (required)",
    "plan_id": "uuid (required)",
    "quotation_type": "standard|tailor_made (required)",
    "client_name": "string (required)",
    "client_email": "string",
    "client_phone": "string",
    "member_count": 5,
    "proposed_members": [{"relationship":"principal","date_of_birth":"1990-01-15"}],
    "billing_frequency": "monthly|quarterly|semi_annual|annual (required)",
    "discount_type": "percentage|fixed (optional)",
    "discount_value": 5000,
    "discount_reason": "string",
    "loading_type": "percentage|fixed (optional)",
    "loading_value": 2000,
    "loading_reason": "string"
  }
  ```
- **Response (201):** QuotationDetailResponse (quotation + first version)
- **Business Rules:**
  - Lead must exist and be in valid status (NEW, CONTACTED, QUALIFIED, PROPOSAL_SENT, NEGOTIATION)
  - Creates quotation in DRAFT status with validity = 30 days
  - Automatically creates Version 1 with pricing calculation
  - **Premium calculation**: Calls CalculatePremiumWithMembers (uses age-band matching from proposed_members DOB)
  - **Discount**: Applied to base: `percentage → base × value / 10000` (basis points: 1000=10%), `fixed → base - value`
  - **Loading**: Added to base: `percentage → base × value / 10000`, `fixed → base + value`
  - **Final premium** = basePremium - discountAmount + loadingAmount (floored at 0)
  - **Approval check**: Loads SalesAgent role limits. If no limits configured → `requires_approval = true`. If discount/loading exceeds role limits → `requires_approval = true`, approval_status = PENDING
  - Lead status auto-advances to PROPOSAL_SENT (from NEW/CONTACTED/QUALIFIED)
  - Quotation number auto-generated

### POST /quotations/expire
- **Auth:** Required | Role: Admin
- **Business Rules:** Batch expires quotations past their `valid_until` date. ISSUED/PENDING_DECISION → EXPIRED.

### GET /quotations/:id
- **Auth:** Required
- **Response (200):** QuotationDetailResponse (includes versions and documents)

### PUT /quotations/:id/issue
- **Auth:** Required
- **Business Rules:** DRAFT → ISSUED. Sets `valid_from` and `valid_until` (default 30 days validity).

### PUT /quotations/:id/accept
- **Auth:** Required
- **Business Rules:** PENDING_DECISION → ACCEPTED.

### PUT /quotations/:id/decline
- **Auth:** Required
- **Business Rules:** PENDING_DECISION → DECLINED.

### PUT /quotations/:id/send
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "channel": "SMS|EMAIL (required)",
    "message": "string"
  }
  ```
- **Business Rules:** ISSUED → PENDING_DECISION. Sends notification to client via specified channel.

### POST /quotations/:id/convert
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "start_date": "YYYY-MM-DD (required)",
    "notes": "string"
  }
  ```
- **Response (200):**
  ```json
  {
    "quotation_id": "uuid",
    "policy_id": "uuid",
    "quotation_number": "string",
    "policy_number": "string",
    "message": "string"
  }
  ```
- **Business Rules:**
  - ACCEPTED → CONVERTED
  - **Conversion steps:**
    1. Creates new policy (DRAFT) with quotation's client info and plan
    2. Enrolls all proposed_members from the latest version
    3. Creates installment schedule using version's billing_frequency
    4. Updates quotation: sets policy_id, status → CONVERTED
    5. Updates lead: status → WON
  - Quotation cannot be converted twice

### GET /quotations/:id/versions
- **Auth:** Required
- **Response (200):** List of QuotationVersionResponse
  ```json
  {
    "id": "uuid",
    "quotation_id": "uuid",
    "version_number": 1,
    "base_premium": 1500000,
    "discount_type": "percentage",
    "discount_value": 5000,
    "discount_reason": "string",
    "loading_type": "fixed",
    "loading_value": 200000,
    "loading_reason": "string",
    "final_premium": 1450000,
    "member_count": 5,
    "proposed_members": [],
    "billing_frequency": "monthly",
    "requires_approval": true,
    "approval_status": "NONE|PENDING|APPROVED|REJECTED",
    "approved_by": "uuid",
    "approved_at": "datetime",
    "rejection_reason": "string",
    "pricing_breakdown": {},
    "created_by": "uuid",
    "created_at": "datetime",
    "updated_at": "datetime"
  }
  ```

### POST /quotations/:id/versions
- **Auth:** Required
- **Request Body:** Same as quotation creation (member_count, billing_frequency, discount, loading, etc.)
- **Business Rules:** Creates a new version with updated pricing. Version number auto-incremented. Requires approval if discount/loading exceeds limits.

### GET /quotations/:id/versions/compare
- **Auth:** Required
- **Query Params:** `version_a` (int), `version_b` (int)
- **Response (200):**
  ```json
  {
    "version_a": { ...QuotationVersionResponse },
    "version_b": { ...QuotationVersionResponse },
    "pricing_diff": {
      "base_premium_diff": -50000,
      "discount_diff": 2000,
      "loading_diff": 0,
      "final_premium_diff": -52000,
      "member_count_diff": -1
    }
  }
  ```

### GET /quotations/:id/versions/:version
- **Auth:** Required

### PUT /quotations/:id/versions/:version/submit-approval
- **Auth:** Required
- **Business Rules:** NONE → PENDING. Submits version for approval by authorized role.

### PUT /quotations/:id/versions/:version/approve
- **Auth:** Required | Role: Admin, Underwriter, Manager
- **Request Body:**
  ```json
  {
    "notes": "string"
  }
  ```
- **Business Rules:** PENDING → APPROVED. Records approver and timestamp.

### PUT /quotations/:id/versions/:version/reject
- **Auth:** Required | Role: Admin, Underwriter, Manager
- **Request Body:**
  ```json
  {
    "reason": "string (required)"
  }
  ```
- **Business Rules:** PENDING → REJECTED. Records reason.

### GET /quotations/:id/documents
- **Auth:** Required

### POST /quotations/:id/documents
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "file_name": "string (required)",
    "file_type": "string (required)",
    "file_size": 1024,
    "version_number": 1,
    "can_edit_roles": ["Admin", "SalesAgent"],
    "can_delete_roles": ["Admin"]
  }
  ```

### PUT /quotation-documents/:id
- **Auth:** Required
- **Request Body:** Partial update (file_name, can_edit_roles, can_delete_roles)

### DELETE /quotation-documents/:id
- **Auth:** Required

---

## 12. Approval Limits

### GET /approval-limits
- **Auth:** Required | Role: Admin
- **Response (200):** List of ApprovalLimitResponse
  ```json
  {
    "id": "uuid",
    "role_name": "SalesAgent",
    "max_discount_percentage": 1000,
    "max_discount_amount": 500000,
    "max_loading_percentage": 2000,
    "max_loading_amount": 1000000,
    "escalation_role": "Manager",
    "is_active": true,
    "created_at": "datetime",
    "updated_at": "datetime"
  }
  ```

### POST /approval-limits
- **Auth:** Required | Role: Admin
- **Request Body:**
  ```json
  {
    "role_name": "string (required)",
    "max_discount_percentage": 1000,
    "max_discount_amount": 500000,
    "max_loading_percentage": 2000,
    "max_loading_amount": 1000000,
    "escalation_role": "Manager"
  }
  ```
- **Business Rules:** Defines per-role limits for quotation pricing. If a quotation version's discount/loading exceeds the creating user's role limit, `requires_approval = true` and must be approved by the `escalation_role` or Admin. Percentage values in basis points (1000 = 10%).

### PUT /approval-limits/:id
- **Auth:** Required | Role: Admin

---

## 13. Policies

### GET /policies
- **Auth:** Required
- **Query Params:** `page`, `page_size`
- **Response (200):** Paginated list of PolicyResponse
  ```json
  {
    "id": "uuid",
    "plan_id": "uuid",
    "policyholder_name": "string",
    "policyholder_email": "string",
    "policyholder_phone": "string",
    "policy_number": "POL-2026-000015",
    "status": "DRAFT|ACTIVE|LAPSED|TERMINATED|SUSPENDED",
    "start_date": "datetime",
    "end_date": "datetime",
    "premium_amount": 1500000,
    "currency": "KES",
    "renewed_from_id": "uuid",
    "created_at": "datetime",
    "updated_at": "datetime"
  }
  ```

### GET /policies/by-status
- **Auth:** Required
- **Query Params:** `status` (required), `page`, `page_size`

### GET /policies/:id
- **Auth:** Required

### POST /policies
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "plan_id": "uuid (required)",
    "policyholder_name": "string (required)",
    "policyholder_email": "string (required, valid email)",
    "policyholder_phone": "string (required)",
    "start_date": "datetime",
    "end_date": "datetime"
  }
  ```
- **Business Rules:** Creates policy in DRAFT status. Policy number auto-generated (POL-YYYY-NNNNNN). Premium calculated from plan's base_premium. Must be activated explicitly.

### PUT /policies/:id
- **Auth:** Required
- **Request Body:** Partial update (policyholder info, dates)

### PUT /policies/:id/activate
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "payment_reference": "string (required)"
  }
  ```
- **Business Rules:**
  - DRAFT → ACTIVE
  - Requires payment reference (proof of premium payment)
  - **Side effects (async, non-blocking):**
    - Auto-generates welcome letter document
    - Auto-generates member cards for all active members
  - Audit logged

### PUT /policies/:id/lapse
- **Auth:** Required
- **Business Rules:** ACTIVE → LAPSED. Triggered by non-payment (30+ days overdue). Members remain enrolled but cannot claim. Can be reinstated.

### PUT /policies/:id/terminate
- **Auth:** Required
- **Business Rules:** ACTIVE or LAPSED → TERMINATED. Permanent. Cannot be reversed.

### PUT /policies/:id/reinstate
- **Auth:** Required
- **Business Rules:** LAPSED/SUSPENDED → ACTIVE. Requires outstanding payments to be cleared.

### PUT /policies/:id/suspend
- **Auth:** Required
- **Business Rules:** ACTIVE → SUSPENDED. Temporary hold. Members cannot claim. Can be reinstated.

### PUT /policies/:id/change-plan
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "new_plan_id": "uuid (required)",
    "reason": "string"
  }
  ```
- **Business Rules:**
  - Policy must be ACTIVE, new plan must be ACTIVE
  - Premium recalculated based on new plan and existing member count
  - **Pro-rata adjustment:**
    ```
    remainingDays = (endDate - now) / 24h
    totalDays = (endDate - startDate) / 24h
    premiumDiff = newPremium - oldPremium
    proratedAdjustment = premiumDiff × remainingDays / totalDays
    finalPremium = oldPremium + proratedAdjustment
    ```
  - If newPremium < oldPremium: auto-creates credit note ("Plan downgrade — pro-rata refund")

### GET /policies/:id/prorate
- **Auth:** Required
- **Response (200):** Pro-rata calculation for remaining period
- **Business Rules:** `refund = (oldPremium - newPremium) × remainingDays / totalPolicyDays`

### POST /policies/bulk/activate
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "ids": ["uuid", "uuid", ...]
  }
  ```
- **Response (200):** BulkResultResponse with succeeded/failed counts

### POST /policies/bulk/lapse
- **Auth:** Required
- **Request Body:** Same as bulk/activate

---

## 14. Members

### GET /policies/:id/members
- **Auth:** Required
- **Response (200):** List of MemberResponse
  ```json
  {
    "id": "uuid",
    "policy_id": "uuid",
    "national_id": "string",
    "name": "string",
    "date_of_birth": "datetime",
    "gender": "male|female|other",
    "relationship": "principal|spouse|child|parent",
    "member_number": "MBR-2026-000200",
    "phone": "string",
    "email": "string",
    "kra_pin": "string",
    "county": "string",
    "address": "string",
    "status": "ACTIVE|SUSPENDED|REMOVED",
    "verified": false,
    "verified_at": "datetime",
    "created_at": "datetime"
  }
  ```

### POST /policies/:id/members
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "national_id": "string",
    "name": "string (required)",
    "date_of_birth": "YYYY-MM-DD (required)",
    "gender": "male|female|other (required)",
    "relationship": "principal|spouse|child|parent (required)",
    "phone": "string",
    "email": "string",
    "kra_pin": "string",
    "county": "string",
    "address": "string"
  }
  ```
- **Business Rules:**
  1. Policy must be ACTIVE or DRAFT
  2. **Double insurance check**: If national_id exists on another active policy → Error + DOUBLE_INSURANCE flag (HIGH severity)
  3. **Age validation**: Member age checked against plan's premium rules (min_age/max_age). Violation → Error + MAX_AGE flag
  4. **Plan underwriting rules**: Each active rule evaluated:
     - MAX_AGE: age > parameter → flag created
     - MIN_AGE: age < parameter → flag created
     - Relationship-specific rules skip non-matching relationships
  5. Member created with status ACTIVE, verified=false
  6. Member number auto-generated (MBR-YYYY-NNNNNN)

### POST /policies/:id/members/bulk
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "members": [{ ...EnrollMemberRequest }]
  }
  ```
- **Response (200):** BulkMemberResultResponse

### POST /policies/:id/members/import
- **Auth:** Required
- **Content-Type:** multipart/form-data (CSV file)
- **CSV Columns:** name (required), date_of_birth (required, YYYY-MM-DD), gender (required), relationship (required), national_id, phone, email, kra_pin, county, address
- **Response (200):** BulkMemberResultResponse

### POST /policies/:id/members/bulk-remove
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "member_ids": ["uuid", ...],
    "reason": "string"
  }
  ```
- **Business Rules:** Removes multiple members. For each, triggers pro-rata credit note if premium decreases.

### GET /members/:id
- **Auth:** Required

### PUT /members/:id
- **Auth:** Required
- **Request Body:** Partial update (name, phone, email, kra_pin, county, address)

### PUT /members/:id/verify
- **Auth:** Required
- **Business Rules:** Sets verified=true, verified_at=now. Confirms member identity.

### PUT /members/:id/suspend
- **Auth:** Required
- **Business Rules:** ACTIVE → SUSPENDED. Member cannot claim benefits.

### PUT /members/:id/reactivate
- **Auth:** Required
- **Business Rules:** SUSPENDED → ACTIVE.

### DELETE /members/:id
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "reason": "string"
  }
  ```
- **Business Rules:**
  - Member must not already be REMOVED; policy must be ACTIVE or DRAFT
  - Status → REMOVED
  - Premium recalculated for policy using CalculatePremiumWithMembers with remaining members
  - If newPremium < oldPremium, **pro-rata credit note** auto-created:
    ```
    totalDays = (policy.EndDate - policy.StartDate).Hours / 24
    remainingDays = (policy.EndDate - now).Hours / 24
    premiumDiff = oldPremium - newPremium
    refundAmount = int64(float64(premiumDiff) × remainingDays / totalDays)
    ```
    - Only created if refundAmount > 0
    - Reason: "Pro-rata refund for member removal: {MemberName}"
    - Status immediately APPROVED (auto-approved because reason contains "Pro-rata refund")

### GET /members/:id/eligibility
- **Auth:** Required
- **Response (200):**
  ```json
  {
    "eligible": true,
    "member_status": "ACTIVE",
    "policy_status": "ACTIVE",
    "policy_end_date": "datetime"
  }
  ```
- **Business Rules:** Eligible = member ACTIVE + policy ACTIVE + current date < policy end_date.

### POST /members/:id/card
- **Auth:** Required
- **Response (200):** PolicyDocumentResponse (generated member card)

### GET /members/:id/underwriting-flags
- **Auth:** Required
- **Response (200):** List of UnderwritingFlagResponse for this member

### GET /members/:id/cases
- **Auth:** Required
- **Response (200):** List of CaseRecordResponse for this member

---

## 15. Endorsements

### GET /policies/:id/endorsements
- **Auth:** Required
- **Response (200):** List of EndorsementResponse
  ```json
  {
    "id": "uuid",
    "policy_id": "uuid",
    "endorsement_type": "ADD_MEMBER|REMOVE_MEMBER|UPDATE_MEMBER|PLAN_CHANGE",
    "status": "PENDING|APPROVED|REJECTED|APPLIED",
    "effective_date": "datetime",
    "changes": { ... },
    "reason": "string",
    "premium_adjustment": 50000,
    "requested_by": "uuid",
    "approved_by": "uuid",
    "approved_at": "datetime",
    "applied_at": "datetime",
    "created_at": "datetime",
    "updated_at": "datetime"
  }
  ```

### POST /policies/:id/endorsements
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "endorsement_type": "ADD_MEMBER|REMOVE_MEMBER|UPDATE_MEMBER|PLAN_CHANGE (required)",
    "effective_date": "YYYY-MM-DD (required)",
    "changes": {},
    "reason": "string"
  }
  ```
  **`changes` JSON structure per endorsement type:**
  - **ADD_MEMBER**: `{ "name": "string", "date_of_birth": "YYYY-MM-DD", "gender": "string", "relationship": "string", "national_id": "string", "phone": "string", "email": "string" }` (EnrollMemberRequest)
  - **REMOVE_MEMBER**: `{ "member_id": "uuid", "reason": "string" }`
  - **UPDATE_MEMBER**: `{ "member_id": "uuid", "updates": { "name": "string", "phone": "string", ... } }` (UpdateMemberRequest nested)
  - **PLAN_CHANGE**: `{ "new_plan_id": "uuid" }` (ChangePlanRequest)
- **Business Rules:** Creates endorsement in PENDING status. Changes stored as JSON. Premium adjustment calculated based on endorsement type. Note: `RejectEndorsement` **overwrites** the original `reason` field with the rejection reason (creation reason is lost).

### GET /endorsements/:id
- **Auth:** Required

### PUT /endorsements/:id/approve
- **Auth:** Required
- **Business Rules:** PENDING → APPROVED. Records approver.

### PUT /endorsements/:id/reject
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "reason": "string (required)"
  }
  ```

### PUT /endorsements/:id/apply
- **Auth:** Required
- **Business Rules:**
  - APPROVED → APPLIED. Dispatches to appropriate service based on endorsement type:
    - **ADD_MEMBER**: Calls member enrollment with changes as EnrollMemberRequest
    - **REMOVE_MEMBER**: Calls member removal with MemberID and Reason from changes JSON
    - **UPDATE_MEMBER**: Calls member update with MemberID and Updates from changes JSON
    - **PLAN_CHANGE**: Calls policy ChangePlan with new plan ID from changes JSON
  - After execution: `policy.PremiumAmount += endorsement.PremiumAdjustment`

---

## 16. Renewals

### GET /policies/:id/renewals
- **Auth:** Required
- **Response (200):** List of RenewalResponse
  ```json
  {
    "id": "uuid",
    "policy_id": "uuid",
    "renewed_policy_id": "uuid",
    "status": "PENDING|APPROVED|REJECTED|COMPLETED|EXPIRED",
    "renewal_date": "datetime",
    "new_premium": 1600000,
    "premium_change_reason": "string",
    "new_plan_id": "uuid",
    "approved_by": "uuid",
    "approved_at": "datetime",
    "completed_at": "datetime",
    "expires_at": "datetime",
    "created_by": "uuid",
    "created_at": "datetime",
    "updated_at": "datetime"
  }
  ```

### POST /policies/:id/renewals
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "new_plan_id": "uuid (optional, keep current plan if empty)",
    "renewal_date": "YYYY-MM-DD (required)",
    "expires_at": "YYYY-MM-DD (optional, renewal offer expiry)"
  }
  ```
- **Business Rules:** Creates PENDING renewal. New premium calculated. Can optionally switch plans during renewal.

### GET /renewals/:id
- **Auth:** Required

### PUT /renewals/:id/approve
- **Auth:** Required
- **Business Rules:** PENDING → APPROVED.

### PUT /renewals/:id/reject
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "reason": "string (required)"
  }
  ```
- **Business Rules:** PENDING → REJECTED. **Note:** The rejection reason is stored in the `premium_change_reason` field (repurposed), not a separate column.

### POST /renewals/:id/complete
- **Auth:** Required
- **Business Rules:**
  - APPROVED → COMPLETED
  - **Creates new policy** (DRAFT status):
    - StartDate = originalPolicy.EndDate
    - EndDate = StartDate + 1 year
    - RenewedFromID = originalPolicy.ID
    - Inherits policyholder info
  - **Premium adjusted by claims experience (loss ratio loading):**
    ```
    lossRatio = (totalApprovedClaimsAmount / originalPremium) × 100
    If lossRatio > 100%: +25% loading
    If lossRatio > 75%:  +15% loading
    If lossRatio > 50%:  +10% loading
    If lossRatio < 30%:  -5% discount
    ```
    Then premium rules recalculated with current members. **Important:** If premium rules return a valid result (no error AND result > 0), it **completely replaces** the claims-loaded premium — it does NOT stack on top of the loading.
  - **Member copy with re-validation:**
    - Each active member's age re-checked against plan rules
    - If age out of range → member SKIPPED, RENEWAL_SKIP flag (MEDIUM severity) created
    - Double insurance re-checked → if detected, member SKIPPED, RENEWAL_SKIP flag (HIGH severity)
    - Passing members copied to new policy with new member numbers, preserving verified flag

### POST /renewals/expire
- **Auth:** Required | Role: Admin
- **Business Rules:** Batch expires renewals past their `expires_at` date. PENDING → EXPIRED.

### POST /renewals/bulk
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "policy_ids": ["uuid", ...]
  }
  ```
- **Business Rules:** Initiates renewals for multiple policies. Auto-sets `renewal_date` to **30 days from now** (`time.Now().AddDate(0, 0, 30)`). Returns bulk result with succeeded/failed counts.

---

## 17. Underwriting Assessments

### GET /policies/:id/underwriting
- **Auth:** Required
- **Response (200):** List of UnderwritingResponse
  ```json
  {
    "id": "uuid",
    "policy_id": "uuid",
    "member_id": "uuid",
    "status": "PENDING|APPROVED|DECLINED|REFER",
    "questionnaire": { ... },
    "medical_declarations": { ... },
    "risk_score": 75,
    "risk_flags": ["HIGH_BMI", "SMOKER"],
    "decision_reason": "string",
    "assessed_by": "uuid",
    "assessed_at": "datetime",
    "created_by": "uuid",
    "created_at": "datetime",
    "updated_at": "datetime"
  }
  ```

### POST /policies/:id/underwriting
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "member_id": "uuid (optional, policy-level if empty)",
    "questionnaire": { "smoker": true, "bmi": 32 },
    "medical_declarations": { "pre_existing": ["diabetes"] }
  }
  ```
- **Business Rules:**
  - Creates assessment in PENDING status
  - Immediately runs auto-evaluation against plan's underwriting rules
  - **Auto-decision engine:**
    - For each active rule matching the member's relationship:
      - MAX_AGE: triggers if member age > ParameterValue
      - MIN_AGE: triggers if member age < ParameterValue
      - DOUBLE_INSURANCE: triggers if same NationalID exists on another ACTIVE policy
      - PRE_EXISTING_CONDITION: triggers if questionnaire[ParameterKey] matches ParameterValue or equals "yes"/"true" (case-insensitive)
      - BMI_THRESHOLD: triggers if questionnaire["bmi"] > ParameterValue (float comparison)
      - WAITING_PERIOD: triggers if questionnaire[ParameterKey] equals "yes"/"true" (informational flag only)
    - Each triggered rule creates an UnderwritingFlag with status=OPEN and accumulates `risk_score += rule.RiskScoreWeight`
    - **Decision thresholds:**
      - Any blocking rule triggered → DECLINED (`"Declined: blocking rule triggered"`)
      - Risk score > 60 → DECLINED (`"Declined: risk score {N} exceeds threshold 60"`)
      - Risk score > 30 → REFER (`"Referred: risk score {N} exceeds auto-approve threshold 30"`)
      - Risk score ≤ 30, no blockers → APPROVED (`"Auto-approved: risk score within acceptable range"`)

### GET /underwriting/:id
- **Auth:** Required

### PUT /underwriting/:id/review
- **Auth:** Required | Role: Admin, Underwriter
- **Request Body:**
  ```json
  {
    "status": "APPROVED|DECLINED|REFER (required)",
    "risk_score": 75,
    "decision_reason": "string"
  }
  ```
- **Business Rules:** PENDING → APPROVED/DECLINED/REFER. Underwriter reviews questionnaire/declarations and assigns risk score. REFER sends to senior underwriter for additional review.

---

## 18. Underwriting Flags

### GET /underwriting-flags
- **Auth:** Required
- **Response (200):** List of all open flags across the system

### GET /underwriting-flags/count
- **Auth:** Required
- **Response (200):** `{ "count": 15 }`

### GET /policies/:id/underwriting-flags
- **Auth:** Required

### GET /members/:id/underwriting-flags
- **Auth:** Required

### GET /underwriting-flags/:id
- **Auth:** Required
- **Response (200):**
  ```json
  {
    "id": "uuid",
    "assessment_id": "uuid",
    "policy_id": "uuid",
    "member_id": "uuid",
    "flag_type": "MAX_AGE|MIN_AGE|DOUBLE_INSURANCE|PRE_EXISTING_CONDITION|BMI_THRESHOLD|WAITING_PERIOD|RENEWAL_SKIP",
    "severity": "LOW|MEDIUM|HIGH",
    "details": "Member age 67 exceeds maximum age 65",
    "status": "OPEN|ACKNOWLEDGED|RESOLVED|OVERRIDDEN",
    "resolved_by": "uuid",
    "resolved_at": "datetime",
    "resolution": "string",
    "created_at": "datetime",
    "updated_at": "datetime"
  }
  ```

### PUT /underwriting-flags/:id/resolve
- **Auth:** Required | Role: Admin, Underwriter
- **Request Body:**
  ```json
  {
    "resolution": "string (required)"
  }
  ```
- **Business Rules:** OPEN/ACKNOWLEDGED → RESOLVED. Underwriter provides resolution explanation.

### PUT /underwriting-flags/:id/override
- **Auth:** Required | Role: Admin, Underwriter
- **Request Body:**
  ```json
  {
    "reason": "string (required)"
  }
  ```
- **Business Rules:** OPEN/ACKNOWLEDGED → OVERRIDDEN. Overrides the flag with documented reason. Used when business decision overrules the automated check.

---

## 19. Policy Documents

### GET /policies/:id/documents
- **Auth:** Required
- **Response (200):** List of PolicyDocumentResponse
  ```json
  {
    "id": "uuid",
    "policy_id": "uuid",
    "member_id": "uuid",
    "document_type": "WELCOME_LETTER|MEMBER_CARD|POLICY_SCHEDULE|RENEWAL_NOTICE|ENDORSEMENT|LOU|DECLINE_LETTER",
    "file_name": "string",
    "file_size": 1024,
    "s3_key": "string",
    "generated_by": "uuid",
    "created_at": "datetime"
  }
  ```

### POST /policies/:id/documents/welcome-letter
- **Auth:** Required
- **Business Rules:** Generates welcome letter PDF for policyholder. Contains policy details, coverage summary, member information.

### POST /policies/:id/documents/policy-schedule
- **Auth:** Required
- **Business Rules:** Generates policy schedule PDF. Detailed breakdown of benefits, premiums, exclusions, provider networks.

### POST /policies/:id/documents/member-cards
- **Auth:** Required
- **Business Rules:** Bulk generates member cards for all active members on the policy.

### POST /members/:id/card
- **Auth:** Required
- **Business Rules:** Generates individual member card.

### POST /preauths/:id/lou
- **Auth:** Required
- **Business Rules:** Generates Letter of Undertaking (LOU) for approved pre-auth. Contains approved amount, procedures, validity period.

### POST /claims/:id/decline-letter
- **Auth:** Required
- **Business Rules:** Generates decline letter for rejected claims. Contains rejection reasons.

### GET /policy-documents/:id
- **Auth:** Required

### DELETE /policy-documents/:id
- **Auth:** Required

---

## 20. Credit Notes

### GET /policies/:id/credit-notes
- **Auth:** Required
- **Response (200):** List of CreditNoteResponse
  ```json
  {
    "id": "uuid",
    "policy_id": "uuid",
    "member_id": "uuid",
    "credit_note_number": "string",
    "amount": 125000,
    "currency": "KES",
    "reason": "Pro-rata refund for member removal: John Doe",
    "status": "DRAFT|APPROVED|APPLIED|CANCELLED",
    "applied_to_invoice_id": "uuid",
    "approved_by": "uuid",
    "approved_at": "datetime",
    "applied_at": "datetime",
    "created_by": "uuid",
    "created_at": "datetime",
    "updated_at": "datetime"
  }
  ```

### GET /credit-notes/:id
- **Auth:** Required

### PUT /credit-notes/:id/approve
- **Auth:** Required | Role: Admin
- **Business Rules:** DRAFT → APPROVED.
- **Note:** Pro-rata refunds (from member removal) are **auto-approved** — they skip DRAFT and go directly to APPROVED.

### PUT /credit-notes/:id/apply
- **Auth:** Required | Role: Admin
- **Request Body:**
  ```json
  {
    "invoice_id": "uuid (required)"
  }
  ```
- **Business Rules:** APPROVED → APPLIED. Credits the specified invoice. Invoice must belong to the same policy.

---

## 21. Pre-Authorization

### GET /preauths
- **Auth:** Required
- **Query Params:** `page`, `page_size`
- **Response (200):** Paginated list of PreAuthResponse
  ```json
  {
    "id": "uuid",
    "policy_id": "uuid",
    "member_id": "uuid",
    "provider_id": "uuid",
    "auth_code": "string",
    "procedure_codes": ["PROC001", "PROC002"],
    "diagnosis_codes": ["ICD10-J06"],
    "estimated_cost": 5000000,
    "approved_amount": 4500000,
    "status": "SUBMITTED|UNDER_REVIEW|APPROVED|DENIED|INFO_REQUESTED|EXPIRED|CLAIMED",
    "validity_start": "datetime",
    "validity_end": "datetime",
    "notes": "string",
    "denial_reason": "string",
    "created_at": "datetime",
    "updated_at": "datetime"
  }
  ```

### GET /preauths/:id
- **Auth:** Required

### POST /preauths
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "policy_id": "uuid (required)",
    "member_id": "uuid (required)",
    "provider_id": "uuid (required)",
    "procedure_codes": ["string"] (required),
    "diagnosis_codes": ["string"] (required),
    "estimated_cost": 5000000,
    "notes": "string"
  }
  ```
- **Business Rules:** Creates pre-auth in SUBMITTED status. Auth code auto-generated.

### PUT /preauths/:id/review
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "decision": "APPROVED|DENIED|INFO_REQUESTED (required)",
    "approved_amount": 4500000,
    "denial_reason": "string",
    "validity_days": 30
  }
  ```
- **Business Rules:**
  - SUBMITTED → UNDER_REVIEW (automatic) → APPROVED/DENIED/INFO_REQUESTED
  - If APPROVED: Sets validity_start=now, validity_end=now+validity_days
  - If DENIED: Records denial_reason
  - If INFO_REQUESTED: Requires additional information from provider

### PUT /preauths/:id/approve
- **Auth:** Required
- **Business Rules:**
  - `approved_amount` = `estimated_cost` (copies estimated cost exactly — no manual override at this endpoint)
  - `auth_code` = `"AUTH-{YEAR}-{6-digit}"` (auto-generated)
  - `validity_start` = now
  - `validity_end` = now + 30 days (PreAuthValidityDays = 30)
  - Status → APPROVED
  - **No status guard:** Can approve from ANY status (SUBMITTED, DENIED, INFO_REQUESTED, even EXPIRED). The documented state machine (SUBMITTED → UNDER_REVIEW → APPROVED) represents intended flow, not enforced constraints.

### PUT /preauths/:id/deny
- **Auth:** Required
- **Business Rules:** Direct denial shortcut. Records denial reason.
  - **No status guard:** Can deny from ANY status. Same as approve — no enforcement of source status.

### POST /preauths/:id/lou
- **Auth:** Required
- **Business Rules:** Generates Letter of Undertaking document. Pre-auth must be APPROVED.
  - **Idempotent:** If a LOU has already been generated for this pre-auth (detected by `LOU_{authCode}_` filename prefix), returns the existing document with HTTP 200 and message `"Existing LOU found for this pre-authorization (generated on {date})"` instead of generating a duplicate.

---

## 22. Claims

### GET /claims
- **Auth:** Required
- **Query Params:** `page`, `page_size`, `status` (optional)
- **Response (200):** Paginated list of ClaimResponse
  ```json
  {
    "id": "uuid",
    "claim_number": "CLM-2026-000042",
    "policy_id": "uuid",
    "member_id": "uuid",
    "provider_id": "uuid",
    "status": "RECEIVED|VALIDATED|ADJUDICATED|APPROVED|REJECTED|MANUAL_REVIEW|PAID|VETTED|PARTIALLY_VETTED|READY_FOR_PAYMENT|PART_PAID",
    "total_amount": 5000000,
    "approved_amount": 4200000,
    "co_pay_amount": 420000,
    "member_responsibility": 800000,
    "diagnosis_codes": ["ICD10-J06"],
    "service_date": "datetime",
    "notes": "string",
    "claim_type": "DIRECT|REIMBURSEMENT|CREDIT|EXCEPTION",
    "vetted_amount": 4200000,
    "vetted_by": "uuid",
    "vetted_at": "datetime",
    "sla_breach_at": "datetime",
    "rejection_reason": "string",
    "line_items": [...],
    "decision": { ... },
    "fraud_flags": [...],
    "created_at": "datetime",
    "updated_at": "datetime"
  }
  ```

### GET /claims/sla-breached
- **Auth:** Required
- **Query Params:** `page`, `page_size`
- **Business Rules:** Returns claims where current time > sla_breach_at and status is not terminal (PAID, REJECTED).

### GET /claims/:id
- **Auth:** Required
- **Response (200):** Full ClaimResponse with embedded line_items, decision, and fraud_flags

### POST /claims
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "policy_id": "uuid (required)",
    "member_id": "uuid (required)",
    "provider_id": "uuid (required)",
    "pre_auth_id": "uuid (optional)",
    "diagnosis_codes": ["ICD10-J06"] (required),
    "service_date": "datetime (required)",
    "admission_date": "datetime (optional, for inpatient)",
    "discharge_date": "datetime (optional, for inpatient)",
    "notes": "string",
    "claim_type": "DIRECT|REIMBURSEMENT|CREDIT|EXCEPTION (optional, default DIRECT)",
    "line_items": [
      {
        "procedure_code": "string (required)",
        "procedure_name": "string (required)",
        "diagnosis_code": "string",
        "quantity": 1,
        "unit_price": 250000
      }
    ]
  }
  ```
- **Business Rules — Full Claim Pipeline:**
  1. **Create**: Claim created in RECEIVED status. Total = sum(quantity × unit_price). SLA breach = `now + 48 hours`. Default claim_type = DIRECT if not specified.
  2. **Validate**: Policy ACTIVE, member exists AND belongs to this policy, provider ACTIVE, has line items, total > 0. If validation fails → immediate REJECT with combined error messages.
  3. **Adjudicate** (9-step engine):
     - **Eligibility**: Policy active, member exists, provider active + in-network + has contract covering service date
     - **PreAuth validation** (if provided): Must be APPROVED, not expired, provider matches, procedures match
     - **Coverage**: Plan has benefits, member age within benefit min_age/max_age (age calculated at service_date, not today)
     - **Waiting periods**: Service date must be after **member.created_at** + waiting_period_days (uses enrollment date, not policy start)
     - **Exclusions**: Claim diagnosis codes matched against plan exclusion ICD codes → REJECT on match
     - **Benefit category**: Has admission_date → INPATIENT, else OUTPATIENT. Falls back to first benefit if no category match
     - **Waiting period & age scope**: These checks iterate over **ALL** plan benefits (not just the matched benefit). If ANY benefit's waiting period or age range is violated, the claim is REJECTED.
     - **Amount calculation**:
       - Annual limit check (hard stop if exhausted)
       - Sub-limit application (per_visit or per_item caps — both behave identically)
       - Deductible: `payable -= deductible_amount` (floored at 0)
       - Co-pay: percentage → `copay = payable × rate / 100`, fixed → `copay = fixed_amount`
       - `payable -= copay` (**WARNING: NOT floored at 0 — can go negative**)
       - Cap at pre-auth approved amount if applicable
     - **Fraud check**: Duplicate detection → MANUAL_REVIEW (only fraud check that affects adjudication decision)
     - **Decision**: APPROVE (with amounts), REJECT (payable=0), or MANUAL_REVIEW
  4. **Store decision** with rule results and reasons
  5. **Update status**: Maps decision to claim status
  6. **Update amounts**: approved_amount, co_pay, member_responsibility
  7. **Fraud checks** (6 additional checks post-adjudication — **informational only**, do NOT change claim status or adjudication):
     - High frequency of same procedure for member → FREQUENCY flag
     - Amount exceeds 500K KES threshold → AMOUNT_THRESHOLD flag
     - Provider contract expired for service date → EXPIRED_CONTRACT flag
     - Provider suspended → SUSPENDED_PROVIDER flag
     - Rate card overcharge (unit price > rate card) → RATE_CARD_OVERCHARGE flag
     - Repeat visit for same procedure → REPEAT_VISIT flag
  8. **PreAuth update**: If pre_auth_id provided, set pre-auth status to CLAIMED
- **CRITICAL:** Returns HTTP **201** even when auto-rejected. Check `message` field:
  - `"Claim submitted and processed"` — success
  - `"Claim submitted but rejected: {reasons}"` — auto-rejected by validator/adjudicator
  - `"Claim submitted, adjudication failed"` — error, claim stuck in VALIDATED
- **Pre-pipeline validator (8 rules, all errors collected):**
  1. Policy exists
  2. Policy status == ACTIVE
  3. Member exists
  4. Member belongs to this policy (member.PolicyID == claim.PolicyID)
  5. Provider exists
  6. Provider status == ACTIVE
  7. At least one line item
  8. TotalAmount > 0

### POST /claims/bulk
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "claims": [{ ...SubmitClaimRequest }]
  }
  ```
- **Response (200):** BulkClaimResultResponse

### POST /claims/import-csv
- **Auth:** Required | Role: Admin, ClaimsOfficer
- **Content-Type:** multipart/form-data
- **Business Rules:** Imports claims from CSV file. Each row creates a claim through the full pipeline.

### PUT /claims/:id/vet
- **Auth:** Required | Role: Admin, ClaimsOfficer
- **Request Body:**
  ```json
  {
    "vetted_amount": 4200000,
    "notes": "string"
  }
  ```
- **Business Rules:**
  - **Status gate:** ADJUDICATED or APPROVED only. Error: `"Cannot vet claim in {status} status; must be ADJUDICATED or APPROVED"` (400)
  - **Claim-type-specific rules:**
    - DIRECT: If inpatient (has admission_date) AND no pre-auth reference → Error: `"Inpatient direct claims require pre-authorization reference"` (400)
    - REIMBURSEMENT: If vetted_amount > total_amount → Error: `"Vetted amount cannot exceed total claimed amount for reimbursement claims"` (400)
    - EXCEPTION: If approved_amount > 0 AND vetted_amount > approved_amount × 3/2 → Error: `"Exception claim vetted amount exceeds 150% of approved amount — requires manual override"` (400). Integer formula: `approved_amount * 3 / 2`
    - CREDIT: No additional vetting rules
  - **Vet status determination:**
    - If vetted_amount == approved_amount → VETTED
    - If vetted_amount < approved_amount → PARTIALLY_VETTED

### PUT /claims/:id/approve
- **Auth:** Required | Role: Admin, Manager
- **Request Body:**
  ```json
  {
    "decision": "APPROVED",
    "reason": "string"
  }
  ```
- **Business Rules:**
  - **Status gate:** ADJUDICATED or MANUAL_REVIEW only. Error: `"Cannot approve claim in {status} status; must be ADJUDICATED or MANUAL_REVIEW"` (400)
  - Recalculates co_pay from adjudication decision: `coPayAmount = totalAmount - decision.PayableAmount - decision.MemberResponsibility` (floored at 0)
  - Updates amounts and transitions to APPROVED

### PUT /claims/:id/reject
- **Auth:** Required | Role: Admin, Manager
- **Request Body:**
  ```json
  {
    "decision": "REJECTED",
    "reason": "string"
  }
  ```
- **Business Rules:**
  - **Status gate:** RECEIVED, VALIDATED, ADJUDICATED, or MANUAL_REVIEW. Error: `"Cannot reject claim in {status} status"` (400)
  - Records rejection reason

### PUT /claims/:id/ready-for-payment
- **Auth:** Required | Role: Admin, Finance
- **Business Rules:** VETTED/PARTIALLY_VETTED → READY_FOR_PAYMENT. Claim is cleared for payment processing.

### PUT /claims/:id/mark-paid
- **Auth:** Required | Role: Admin, Finance
- **Business Rules:** READY_FOR_PAYMENT → PAID. Full payment recorded.

### PUT /claims/:id/mark-part-paid
- **Auth:** Required | Role: Admin, Finance
- **Business Rules:** READY_FOR_PAYMENT → PART_PAID. Partial payment recorded.

---

## 23. Cases

### GET /cases
- **Auth:** Required
- **Query Params:** `page`, `page_size`, `status` (optional)

### GET /cases/count
- **Auth:** Required
- **Response (200):** `{ "counts": { "SCHEDULED": 5, "ADMITTED": 3, ... } }`

### GET /cases/:id
- **Auth:** Required
- **Response (200):**
  ```json
  {
    "id": "uuid",
    "case_number": "string",
    "pre_auth_id": "uuid",
    "policy_id": "uuid",
    "member_id": "uuid",
    "provider_id": "uuid",
    "status": "SCHEDULED|ADMITTED|IN_TREATMENT|DISCHARGED|CLOSED",
    "admission_date": "datetime",
    "expected_discharge": "datetime",
    "actual_discharge": "datetime",
    "diagnosis": "string",
    "treating_doctor": "string",
    "room_type": "string",
    "total_estimated_cost": 10000000,
    "total_actual_cost": 9500000,
    "notes": "string",
    "closed_at": "datetime",
    "created_by": "uuid",
    "created_at": "datetime",
    "updated_at": "datetime"
  }
  ```

### POST /cases
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "pre_auth_id": "uuid (required)",
    "expected_discharge": "datetime",
    "diagnosis": "string",
    "treating_doctor": "string",
    "room_type": "string",
    "estimated_cost": 10000000,
    "notes": "string"
  }
  ```
- **Business Rules:** Creates case from approved pre-auth in SCHEDULED status. Links to policy, member, and provider from the pre-auth. Pre-auth must exist and be in APPROVED status.

### PUT /cases/:id
- **Auth:** Required
- **Request Body:** Partial update (diagnosis, treating_doctor, room_type, estimated_cost, notes)

### PUT /cases/:id/admit
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "admission_date": "datetime (required)"
  }
  ```
- **Business Rules:** SCHEDULED → ADMITTED. Status must be SCHEDULED; error if not: `"Cannot admit case in {status} status"`.

### PUT /cases/:id/start-treatment
- **Auth:** Required
- **Business Rules:** ADMITTED → IN_TREATMENT. Status must be ADMITTED; error if not: `"Cannot start treatment for case in {status} status"`. No request body required.

### PUT /cases/:id/discharge
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "actual_discharge": "datetime (required)",
    "actual_cost": 9500000
  }
  ```
- **Business Rules:** ADMITTED or IN_TREATMENT → DISCHARGED. Discharge is allowed from both ADMITTED (early discharge before treatment) and IN_TREATMENT. Records actual discharge date and cost.

### PUT /cases/:id/close
- **Auth:** Required
- **Business Rules:** DISCHARGED → CLOSED. Status must be DISCHARGED; error if not. Final case closure. Sets closed_at timestamp. Does NOT verify that all linked claims are in terminal status.

### GET /policies/:id/cases
- **Auth:** Required

### GET /members/:id/cases
- **Auth:** Required

### GET /providers/:id/cases
- **Auth:** Required

---

## 24. Claim Documents

### GET /claims/:id/documents
- **Auth:** Required
- **Response (200):** List of ClaimDocumentResponse
  ```json
  {
    "id": "uuid",
    "claim_id": "uuid",
    "file_name": "string",
    "file_type": "string",
    "file_size": 1024,
    "s3_key": "string",
    "uploaded_by": "uuid",
    "created_at": "datetime"
  }
  ```

### POST /claims/:id/documents
- **Auth:** Required
- **Content-Type:** multipart/form-data
- **Business Rules:** Uploads document and associates with claim. S3 key generated for storage.

### DELETE /claim-documents/:id
- **Auth:** Required

---

## 25. Provider Statements

### GET /providers/:id/statements
- **Auth:** Required
- **Response (200):** List of ProviderStatementResponse
  ```json
  {
    "id": "uuid",
    "provider_id": "uuid",
    "statement_number": "string",
    "period_start": "datetime",
    "period_end": "datetime",
    "total_claimed": 50000000,
    "total_matched": 48000000,
    "total_discrepancy": 2000000,
    "matched_count": 45,
    "unmatched_count": 5,
    "status": "UPLOADED|RECONCILED",
    "file_name": "string",
    "reconciled_at": "datetime",
    "created_at": "datetime"
  }
  ```

### POST /providers/:id/statements
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "period_start": "datetime (required)",
    "period_end": "datetime (required)",
    "file_name": "string",
    "s3_key": "string",
    "line_items": [
      {
        "claim_number": "CLM-2026-000042",
        "service_date": "datetime",
        "member_name": "string",
        "procedure_code": "string",
        "claimed_amount": 500000
      }
    ]
  }
  ```
- **Business Rules:** Creates statement in UPLOADED status with line items. Each line item starts as UNMATCHED.

### GET /provider-statements/:id
- **Auth:** Required

### GET /provider-statements/:id/line-items
- **Auth:** Required
- **Response (200):** List of StatementLineItemResponse
  ```json
  {
    "id": "uuid",
    "statement_id": "uuid",
    "claim_number": "string",
    "service_date": "datetime",
    "member_name": "string",
    "procedure_code": "string",
    "claimed_amount": 500000,
    "matched_claim_id": "uuid",
    "match_status": "UNMATCHED|MATCHED|DISPUTED",
    "discrepancy_amount": 50000,
    "notes": "string",
    "created_at": "datetime"
  }
  ```

### POST /provider-statements/:id/reconcile
- **Auth:** Required
- **Business Rules:** UPLOADED → RECONCILED. Uses **two-phase matching algorithm**:
  1. **Phase 1 — Match by claim number**: If `item.ClaimNumber` is not empty, look up claim by claim_number
  2. **Phase 2 — Fallback match**: If Phase 1 fails, match by `provider_id + service_date + amount` using `FindByProviderAndDate`
  3. **Phase 3 — Unmatched**: If both phases fail, mark as UNMATCHED
  - **Amount tolerance**: 1 KES (100 cents). `discrepancy = item.ClaimedAmount - claim.ApprovedAmount`. If `abs(discrepancy) <= 100 cents`, discrepancy is reset to 0.
  - **Claim status update on match**: If `discrepancy <= 0` → claim status set to `PAID`. If `discrepancy > 0` → claim status set to `PART_PAID`.
  - Updates matched/unmatched counts and totals on the statement record.

---

## 26. Installments

### GET /policies/:id/installments
- **Auth:** Required
- **Response (200):** List of InstallmentScheduleResponse
  ```json
  {
    "id": "uuid",
    "policy_id": "uuid",
    "frequency": "monthly|quarterly|semi_annual|annual",
    "total_installments": 12,
    "amount_per_installment": 125000,
    "start_date": "datetime",
    "status": "ACTIVE|COMPLETED|CANCELLED",
    "created_at": "datetime",
    "installments": [
      {
        "id": "uuid",
        "schedule_id": "uuid",
        "installment_number": 1,
        "due_date": "datetime",
        "amount": 125000,
        "status": "PENDING|PAID|OVERDUE",
        "paid_at": "datetime",
        "invoice_id": "uuid",
        "created_at": "datetime"
      }
    ]
  }
  ```

### POST /policies/:id/installments
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "frequency": "monthly|quarterly|semi_annual|annual (required)",
    "start_date": "datetime"
  }
  ```
- **Business Rules:** Creates installment schedule. Total = policy premium. Divides into installments based on frequency (monthly=12, quarterly=4, semi_annual=2, annual=1). Each installment has a due date calculated from start_date.

### GET /installments/schedule/:id
- **Auth:** Required
- **Response (200):** List of InstallmentResponse for this schedule

### PUT /installments/:id/pay
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "invoice_id": "string (optional)"
  }
  ```
- **Business Rules:** PENDING/OVERDUE → PAID. Records paid_at timestamp. Links to invoice if provided.

---

## 27. Invoices

### GET /invoices
- **Auth:** Required
- **Query Params:** `page`, `page_size`, `status` (optional)
- **Response (200):** Paginated list of InvoiceResponse
  ```json
  {
    "id": "uuid",
    "policy_id": "uuid",
    "invoice_number": "string",
    "amount": 1500000,
    "currency": "KES",
    "due_date": "datetime",
    "status": "PENDING|PAID|OVERDUE|CANCELLED",
    "billing_period_start": "datetime",
    "billing_period_end": "datetime",
    "created_at": "datetime"
  }
  ```

### POST /invoices/:policyId
- **Auth:** Required
- **Business Rules:**
  - **Amount** = `policy.PremiumAmount` (full premium in cents)
  - **Due date** = now + 30 days (`InvoiceDueDays = 30`)
  - **Billing period** = `billing_period_start: now`, `billing_period_end: now + 1 month`
  - **Invoice number** = `INV-{YYYY}-{NNNNNN}` (6-digit nanosecond modulo — collision risk for concurrent calls)
  - **Currency** = copied from `policy.Currency`
  - **Status** = `PENDING`
  - Note: `RunBillingCycle` (batch invoice generation) is a **stub** — use this endpoint for manual generation.

### GET /invoices/:id
- **Auth:** Required

---

## 28. Payments

### GET /payments
- **Auth:** Required
- **Query Params:** `page`, `page_size`
- **Response (200):** Paginated list of PaymentResponse
  ```json
  {
    "id": "uuid",
    "type": "PREMIUM|REMITTANCE",
    "amount": 1500000,
    "currency": "KES",
    "method": "MPESA|BANK_TRANSFER",
    "reference_number": "string",
    "status": "INITIATED|PROCESSING|CONFIRMED|FAILED|RECONCILED|CANCELLED",
    "retry_count": 0,
    "paid_at": "datetime",
    "created_at": "datetime"
  }
  ```

### GET /payments/:id
- **Auth:** Required

### POST /payments
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "invoice_id": "uuid (optional)",
    "claim_id": "uuid (optional)",
    "amount": 1500000,
    "method": "MPESA|BANK_TRANSFER (required)",
    "phone": "string (required for MPESA)"
  }
  ```
- **Business Rules:** Creates payment in INITIATED status. For MPESA, triggers STK push. Reference number auto-generated.

### PUT /payments/:id/retry
- **Auth:** Required
- **Business Rules:** FAILED → INITIATED. Increments retry_count. Max retries = 3.

### PUT /payments/:id/reconcile
- **Auth:** Required
- **Business Rules:** CONFIRMED → RECONCILED. Matches payment to bank statement.

### POST /webhooks/mpesa (Public)
- **Auth:** Public (webhook)
- **Business Rules:** Receives M-Pesa callback. Updates payment status (PROCESSING → CONFIRMED or FAILED).

---

## 29. Remittances

### GET /remittances
- **Auth:** Required
- **Query Params:** `page`, `page_size`
- **Response (200):** Paginated list of RemittanceResponse
  ```json
  {
    "id": "uuid",
    "provider_id": "uuid",
    "total_amount": 25000000,
    "currency": "KES",
    "status": "PENDING|PROCESSING|SENT|CONFIRMED|FAILED",
    "remittance_advice_sent": false,
    "period_start": "datetime",
    "period_end": "datetime",
    "created_at": "datetime"
  }
  ```

### GET /remittances/:id
- **Auth:** Required

### POST /remittances/:providerId
- **Auth:** Required
- **Business Rules:**
  - Fetches all approved claims for the provider via `GetApprovedForRemittance`
  - If no approved claims: Error `"No approved claims for remittance"` (HTTP 400)
  - **Total amount** = sum of all `ApprovedAmount` from approved claims
  - **Period** = `period_start: now - 1 month`, `period_end: now` (auto-set, not user-specified)
  - **Status** = `PENDING`
  - **Currency** = `"KES"` (hardcoded)
  - Note: `RunRemittanceCycle` (batch) processes only **ACTIVE** providers, up to **1000** per cycle. It is a **stub** scheduler.

### GET /remittances/:id/export
- **Auth:** Required
- **Response (200):**
  ```json
  {
    "remittance_id": "uuid",
    "provider_id": "uuid",
    "provider_name": "string",
    "total_amount": 25000000,
    "currency": "KES",
    "period_start": "datetime",
    "period_end": "datetime",
    "claims": [
      {
        "claim_number": "CLM-2026-000042",
        "amount": 500000,
        "service_date": "datetime"
      }
    ]
  }
  ```
- **Business Rules:** Exports payment file with claim-level detail for bank processing.

---

## 30. Notifications

### GET /notifications
- **Auth:** Required
- **Query Params:** `page`, `page_size`
- **Response (200):** Paginated list of NotificationResponse
  ```json
  {
    "id": "uuid",
    "user_id": "uuid",
    "channel": "SMS|EMAIL|IN_APP|PUSH",
    "type": "QUOTATION|APPROVAL|CLAIM|POLICY|DOCUMENT",
    "subject": "string",
    "body": "string",
    "metadata": {},
    "status": "PENDING|SENT|DELIVERED|FAILED|READ",
    "retry_count": 0,
    "max_retries": 3,
    "sent_at": "datetime",
    "read_at": "datetime",
    "created_at": "datetime",
    "updated_at": "datetime"
  }
  ```
- **Business Rules:** Returns notifications for the authenticated user. Ordered by created_at desc.

### PUT /notifications/:id/read
- **Auth:** Required
- **Business Rules:** Marks notification as READ. Sets read_at timestamp.

### GET /notifications/unread-count
- **Auth:** Required
- **Response (200):** `{ "count": 7 }`

---

## 31. Audit

### GET /audit
- **Auth:** Required
- **Query Params:** `page`, `page_size`
- **Response (200):** Paginated list of AuditEventResponse
  ```json
  {
    "id": "uuid",
    "user_id": "uuid",
    "entity_type": "CLAIM|POLICY|MEMBER|PLAN|...",
    "entity_id": "uuid",
    "action": "CREATE|UPDATE|DELETE|STATE_CHANGE",
    "old_value": {},
    "new_value": {},
    "ip_address": "string",
    "user_agent": "string",
    "created_at": "datetime"
  }
  ```

### GET /audit/entity/:type/:id
- **Auth:** Required
- **Business Rules:** Lists all audit events for a specific entity (e.g., all changes to policy POL-2026-000015).

### GET /audit/user/:id
- **Auth:** Required
- **Business Rules:** Lists all audit events performed by a specific user.

---

## 32. Analytics

**parsePeriod helper** — All analytics endpoints that accept a `period` query parameter map it to day ranges:
| Period Value | Days |
|---|---|
| `week` | 7 |
| `month` | 30 (default if unspecified or unknown) |
| `quarter` | 90 |
| `year` | 365 |

### GET /analytics/dashboard
- **Auth:** Required
- **Query Params:** `period` (string: week/month/quarter/year, default: month)
- **Response (200):**
  ```json
  {
    "claims_volume": {
      "total_claims": 500,
      "approved_claims": 350,
      "rejected_claims": 80,
      "manual_review_claims": 20,
      "paid_claims": 300
    },
    "approval_rate": 70.0,
    "average_tat": 24.5,
    "loss_ratio": 65.0,
    "fraud_rate": 2.5,
    "total_premium": 500000000,
    "total_claims_paid": 325000000,
    "top_providers": [
      {
        "id": "uuid",
        "name": "string",
        "claim_count": 50,
        "total_amount": 25000000,
        "total_approved": 22000000
      }
    ]
  }
  ```

### GET /analytics/kpis
- **Auth:** Required
- **Response (200):**
  ```json
  {
    "approval_rate": 70.0,
    "average_tat": 24.5,
    "loss_ratio": 65.0,
    "fraud_rate": 2.5,
    "total_premium": 500000000,
    "total_claims_paid": 325000000
  }
  ```
- **Business Rules:**
  - `approval_rate` = approved / total × 100
  - `average_tat` = average time from RECEIVED to terminal status (hours)
  - `loss_ratio` = total_claims_paid / total_premium × 100
  - `fraud_rate` = claims_with_fraud_flags / total × 100

### GET /analytics/export
- **Auth:** Required
- **Query Params:** `report_type` (string), `period` (string)
- **Response:** CSV file download
- **Business Rules:** **STUB — returns empty CSV.** The response contains only header columns (`report,data\n`) with no data rows. The `report_type` and `period` parameters are accepted but ignored. Frontend should hide this feature or display a "coming soon" indicator.

### GET /analytics/reinsurance
- **Auth:** Required
- **Business Rules:** Uses a **hardcoded** `"last_year"` (365 days) period. The period cannot be overridden via query parameters. Always returns data for the trailing 365-day window.
- **Response (200):**
  ```json
  {
    "active_treaty_count": 3,
    "total_ceded_premiums": 150000000,
    "total_recoverable": 45000000,
    "total_recovered": 38000000,
    "total_outstanding": 7000000,
    "cession_ratio": 0.30,
    "recovery_success_rate": 0.844,
    "unacknowledged_alerts": 2
  }
  ```
- **Business Rules:**
  - `total_outstanding` = total_recoverable - total_recovered
  - `cession_ratio` = total_ceded / total_gross
  - `recovery_success_rate` = total_recovered / total_recoverable

---

## 33. Treaties

### GET /treaties
- **Auth:** Required
- **Query Params:** `page`, `page_size`, `status` (optional)
- **Response (200):** Paginated list of TreatyResponse
  ```json
  {
    "id": "uuid",
    "treaty_number": "TRY-2026-000003",
    "name": "string",
    "treaty_type": "QUOTA_SHARE|XOL",
    "status": "DRAFT|ACTIVE|EXPIRED|TERMINATED",
    "effective_date": "datetime",
    "expiry_date": "datetime",
    "retention_limit": 100000000,
    "currency": "KES",
    "notes": "string",
    "created_by": "uuid",
    "created_at": "datetime",
    "updated_at": "datetime"
  }
  ```

### POST /treaties
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "name": "string (required)",
    "treaty_type": "QUOTA_SHARE|XOL (required)",
    "effective_date": "datetime (required)",
    "expiry_date": "datetime (required)",
    "retention_limit": 100000000,
    "currency": "KES",
    "notes": "string"
  }
  ```
- **Business Rules:** Creates treaty in DRAFT status. Treaty number auto-generated (TRY-YYYY-NNNNNN). Default currency KES.

### POST /treaties/expire
- **Auth:** Required
- **Business Rules:** Fetches up to **1000** ACTIVE treaties. **Caution:** Due to a filter condition issue, this endpoint effectively expires **all** fetched active treaties regardless of their actual expiry date. Use with care — intended for Admin manual expiry checks.

### GET /treaties/:id
- **Auth:** Required
- **Response (200):** TreatyDetailResponse (includes participants, layers, profit commission rules)

### PUT /treaties/:id
- **Auth:** Required
- **Request Body:** Partial update (name, dates, retention_limit, currency, notes)

### PUT /treaties/:id/activate
- **Auth:** Required
- **Business Rules:**
  - DRAFT → ACTIVE
  - **Validation**:
    - Must have at least one participant (totalShare > 0). Error: `"Treaty must have at least one participant"` (400)
    - Total participant share must NOT exceed 100%. Error: `"Total participant share exceeds 100%"` (400)
    - Partial share allowed (e.g., 60% total is valid — insurer retains 40%)

### PUT /treaties/:id/terminate
- **Auth:** Required
- **Business Rules:** ACTIVE → TERMINATED.

### GET /treaties/:id/participants
- **Auth:** Required
- **Response (200):** List of TreatyParticipantResponse
  ```json
  {
    "id": "uuid",
    "treaty_id": "uuid",
    "reinsurer_name": "string",
    "share_percentage": 30.0,
    "commission_rate": 25.0,
    "is_lead": true,
    "created_at": "datetime",
    "updated_at": "datetime"
  }
  ```

### POST /treaties/:id/participants
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "reinsurer_name": "string (required)",
    "share_percentage": 30.0,
    "commission_rate": 25.0,
    "is_lead": true
  }
  ```
- **Business Rules:** share_percentage must be > 0 and ≤ 100. Total across all participants must not exceed 100%.

### PUT /treaties/:id/participants/:participantId
- **Auth:** Required
- **Request Body:** Partial update (reinsurer_name, share_percentage, commission_rate, is_lead)

### DELETE /treaties/:id/participants/:participantId
- **Auth:** Required

### GET /treaties/:id/layers
- **Auth:** Required
- **Response (200):** List of TreatyLayerResponse
  ```json
  {
    "id": "uuid",
    "treaty_id": "uuid",
    "layer_number": 1,
    "attachment_point": 10000000,
    "layer_limit": 50000000,
    "deductible_amount": 500000,
    "premium_rate": 5.0,
    "aggregate_limit": 200000000,
    "aggregate_used": 45000000,
    "created_at": "datetime",
    "updated_at": "datetime"
  }
  ```

### POST /treaties/:id/layers
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "layer_number": 1,
    "attachment_point": 10000000,
    "layer_limit": 50000000,
    "deductible_amount": 500000,
    "premium_rate": 5.0,
    "aggregate_limit": 200000000
  }
  ```
- **Business Rules:**
  - **Attachment point**: Loss amount at which the layer starts responding
  - **Layer limit**: Maximum the layer will pay per occurrence
  - **Aggregate limit**: Maximum total payout for the layer across all occurrences
  - **aggregate_used** tracks cumulative usage — alerts triggered at 80% and 100%

### PUT /treaties/:id/layers/:layerId
- **Auth:** Required

### DELETE /treaties/:id/layers/:layerId
- **Auth:** Required

### GET /treaties/:id/profit-commission-rules
- **Auth:** Required
- **Response (200):** List of ProfitCommissionResponse
  ```json
  {
    "id": "uuid",
    "treaty_id": "uuid",
    "commission_type": "SLIDING_SCALE|FLAT|CARRY_FORWARD",
    "loss_ratio_from": 0.0,
    "loss_ratio_to": 50.0,
    "commission_rate": 25.0,
    "carry_forward_years": 3,
    "carry_forward_balance": 0,
    "period_start": "datetime",
    "period_end": "datetime",
    "calculated_amount": 0,
    "created_at": "datetime",
    "updated_at": "datetime"
  }
  ```

### POST /treaties/:id/profit-commission-rules
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "commission_type": "SLIDING_SCALE|FLAT|CARRY_FORWARD (required)",
    "loss_ratio_from": 0.0,
    "loss_ratio_to": 50.0,
    "commission_rate": 25.0,
    "carry_forward_years": 3
  }
  ```
- **Business Rules:**
  - **SLIDING_SCALE**: Commission rate varies by loss ratio band
  - **FLAT**: Fixed commission rate regardless of loss ratio
  - **CARRY_FORWARD**: Losses carried forward to offset future profits (carry_forward_years defines the lookback)

### DELETE /treaties/:id/profit-commission-rules/:ruleId
- **Auth:** Required

### GET /treaties/:id/cessions
- **Auth:** Required (see [Cessions](#34-cessions))

### GET /treaties/:id/recoveries
- **Auth:** Required (see [Recoveries](#35-recoveries))

### GET /treaties/:id/bordereaux
- **Auth:** Required (see [Bordereaux](#36-bordereaux))

### GET /treaties/:id/statements
- **Auth:** Required (see [Reinsurer Statements](#37-reinsurer-statements))

### GET /treaties/:id/alerts
- **Auth:** Required (see [Treaty Alerts](#38-treaty-alerts))

---

## 34. Cessions

### POST /cessions
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "treaty_id": "uuid (required)",
    "policy_id": "uuid (required)",
    "amount": 1500000
  }
  ```
- **Response (201):** CessionResponse
  ```json
  {
    "id": "uuid",
    "cession_number": "CES-2026-000100",
    "treaty_id": "uuid",
    "policy_id": "uuid",
    "treaty_layer_id": "uuid",
    "cession_type": "PREMIUM|CLAIM",
    "gross_amount": 1500000,
    "ceded_amount": 450000,
    "retained_amount": 1050000,
    "commission_amount": 112500,
    "share_percentage": 30.0,
    "status": "PENDING|BOOKED|REVERSED",
    "created_by": "uuid",
    "created_at": "datetime",
    "updated_at": "datetime"
  }
  ```
- **Business Rules:**
  - Treaty must be ACTIVE; must have at least one participant
  - **Calculation (Quota Share)**:
    ```
    totalShare = sum(participant.SharePercentage)
    avgCommissionRate = sum(participant.CommissionRate) / participantCount
    cededAmount = floor(grossAmount × totalShare / 100)
    retainedAmount = grossAmount - cededAmount
    commissionAmount = floor(cededAmount × avgCommissionRate / 100)
    ```
  - **Retention limit override**: If `treaty.RetentionLimit > 0 AND retainedAmount < retentionLimit`:
    ```
    retainedAmount = retentionLimit
    cededAmount = grossAmount - retainedAmount
    commissionAmount = floor(cededAmount × avgCommissionRate / 100)
    ```
  - Status: PENDING

### POST /cessions/auto-cede
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "policy_id": "uuid (required)",
    "amount": 1500000
  }
  ```
- **Business Rules:** Automatically cedes to ALL active QUOTA_SHARE treaties. Creates one cession per qualifying treaty. Fails if no active quota share treaties exist.

### GET /cessions/:id
- **Auth:** Required

### PUT /cessions/:id/book
- **Auth:** Required
- **Business Rules:** PENDING → BOOKED. Confirms the cession.

### PUT /cessions/:id/reverse
- **Auth:** Required
- **Business Rules:** BOOKED → REVERSED. Only booked cessions can be reversed.

---

## 35. Recoveries

### POST /recoveries
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "claim_id": "uuid (required)",
    "treaty_id": "uuid (required)",
    "treaty_layer_id": "uuid (optional)",
    "cession_id": "uuid (optional)",
    "gross_amount": 5000000,
    "recoverable_amount": 1500000,
    "notes": "string"
  }
  ```
- **Response (201):** RecoveryResponse
  ```json
  {
    "id": "uuid",
    "recovery_number": "REC-2026-000008",
    "claim_id": "uuid",
    "treaty_id": "uuid",
    "treaty_layer_id": "uuid",
    "cession_id": "uuid",
    "gross_claim_amount": 5000000,
    "recoverable_amount": 1500000,
    "recovered_amount": 0,
    "outstanding_amount": 1500000,
    "status": "NOTIFIED|ACKNOWLEDGED|INFO_REQUESTED|APPROVED|PAID|WRITTEN_OFF",
    "workflow_status": "NOTIFICATION|ACKNOWLEDGMENT|INFO_REQUEST|APPROVAL|PAYMENT",
    "notes": "string",
    "created_by": "uuid",
    "created_at": "datetime",
    "updated_at": "datetime"
  }
  ```
- **Business Rules:** Creates recovery in NOTIFIED status, workflow_status=NOTIFICATION. recovered_amount=0, outstanding_amount=recoverable_amount.

### GET /recoveries/outstanding
- **Auth:** Required
- **Response (200):** List of recoveries with outstanding_amount > 0

### GET /recoveries/aged-analysis
- **Auth:** Required
- **Response (200):**
  ```json
  [
    { "bucket": "0-30 days", "count": 5, "total_outstanding": 7500000 },
    { "bucket": "31-60 days", "count": 3, "total_outstanding": 4500000 },
    { "bucket": "61-90 days", "count": 1, "total_outstanding": 2000000 },
    { "bucket": "90+ days", "count": 1, "total_outstanding": 3000000 }
  ]
  ```

### POST /recoveries/apply-for-claim/:claimId
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "approved_amount": 5000000
  }
  ```
- **Business Rules:**
  - Finds active treaties and auto-creates recoveries:
  - Finds ALL active treaties (both QUOTA_SHARE and XOL) and creates recoveries
  - **Quota Share**: `recoverable = floor(approvedAmount × totalShare / 100)`
  - **XOL** — per layer, sorted by layer_number ascending:
    ```
    excess = approvedAmount - layer.AttachmentPoint
    If excess ≤ 0: skip layer
    layerExposure = min(excess, layer.LayerLimit)
    recoverable = layerExposure - layer.DeductibleAmount
    If recoverable ≤ 0: skip layer
    If layer.AggregateLimit exists:
      remaining = AggregateLimit - AggregateUsed
      If remaining ≤ 0: skip layer (exhausted)
      If recoverable > remaining: recoverable = remaining
      AggregateUsed += recoverable  (tracked for alert thresholds)
    ```

### GET /recoveries/:id
- **Auth:** Required
- **Response (200):** RecoveryDetailResponse (includes workflow events)

### PUT /recoveries/:id/acknowledge
- **Auth:** Required
- **Request Body:** `{ "notes": "string" }`
- **Business Rules:** NOTIFIED → ACKNOWLEDGED. Workflow event: ACKNOWLEDGMENT.

### PUT /recoveries/:id/request-info
- **Auth:** Required
- **Request Body:** `{ "notes": "string" }`
- **Business Rules:** ACKNOWLEDGED → INFO_REQUESTED. Workflow event: INFO_REQUEST.

### PUT /recoveries/:id/approve
- **Auth:** Required
- **Request Body:** `{ "notes": "string" }`
- **Business Rules:** ACKNOWLEDGED/INFO_REQUESTED → APPROVED. Workflow event: APPROVAL.

### PUT /recoveries/:id/record-payment
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "amount": 1500000,
    "notes": "string"
  }
  ```
- **Business Rules:**
  - APPROVED → PAID
  - `recovered_amount += amount`
  - `outstanding_amount = recoverable_amount - recovered_amount` (min 0)
  - **IMPORTANT:** Partial payment still transitions to PAID status — there is no PARTIALLY_PAID state
  - Workflow event: PAYMENT

### PUT /recoveries/:id/write-off
- **Auth:** Required
- **Request Body:** `{ "notes": "string" }`
- **Business Rules:**
  - Allowed from ANY status except PAID. Error: `"Cannot write off a PAID recovery"` (400)
  - Allowed statuses: NOTIFIED, ACKNOWLEDGED, INFO_REQUESTED, APPROVED → WRITTEN_OFF
  - Preserves existing `workflow_status` (does not change it)
  - Workflow event: WRITE_OFF

### GET /recoveries/:id/workflow
- **Auth:** Required
- **Response (200):** List of RecoveryWorkflowEventResponse
  ```json
  [
    {
      "id": "uuid",
      "recovery_id": "uuid",
      "from_status": "NOTIFIED",
      "to_status": "ACKNOWLEDGED",
      "event_type": "ACKNOWLEDGMENT",
      "notes": "string",
      "performed_by": "uuid",
      "created_at": "datetime"
    }
  ]
  ```

---

## 36. Bordereaux

### POST /bordereaux/premium
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "treaty_id": "uuid (required)",
    "period_start": "datetime (required)",
    "period_end": "datetime (required)"
  }
  ```
- **Response (201):** BordereauResponse
  ```json
  {
    "id": "uuid",
    "bordereau_number": "BDX-2026-000001",
    "treaty_id": "uuid",
    "bordereau_type": "PREMIUM|CLAIM",
    "period_start": "datetime",
    "period_end": "datetime",
    "total_gross": 50000000,
    "total_ceded": 15000000,
    "total_commission": 3750000,
    "item_count": 25,
    "status": "DRAFT|FINALIZED|SENT",
    "created_by": "uuid",
    "created_at": "datetime",
    "updated_at": "datetime"
  }
  ```
- **Business Rules:**
  - Fetches only **BOOKED** cessions for the treaty within the period (PENDING and REVERSED excluded)
  - Creates one bordereau item per cession with: CessionID, GrossAmount, CededAmount, CommissionAmount
  - Aggregates: TotalGross = SUM(GrossAmount), TotalCeded = SUM(CededAmount), TotalCommission = SUM(CommissionAmount)
  - Initial status: DRAFT

### POST /bordereaux/claim
- **Auth:** Required
- **Request Body:** Same as premium bordereau
- **Business Rules:**
  - Fetches up to **10,000** recoveries for the treaty (hard limit)
  - Filters by `recovery.CreatedAt` within period range in application code (not SQL)
  - Creates one bordereau item per recovery: GrossAmount = GrossClaimAmount, CededAmount = RecoverableAmount, CommissionAmount = **0** (always zero for claim bordereaux)
  - Aggregates: TotalGross = SUM(GrossClaimAmount), TotalCeded = SUM(RecoverableAmount), TotalCommission = 0

### GET /bordereaux/:id
- **Auth:** Required

### PUT /bordereaux/:id/finalize
- **Auth:** Required
- **Business Rules:** DRAFT → FINALIZED. Locks the bordereau for submission.

### PUT /bordereaux/:id/mark-sent
- **Auth:** Required
- **Business Rules:** FINALIZED → SENT. Confirms sent to reinsurer.

### GET /bordereaux/:id/items
- **Auth:** Required
- **Response (200):** List of BordereauItemResponse
  ```json
  {
    "id": "uuid",
    "bordereau_id": "uuid",
    "cession_id": "uuid",
    "recovery_id": "uuid",
    "policy_number": "string",
    "claim_number": "string",
    "gross_amount": 1500000,
    "ceded_amount": 450000,
    "commission_amount": 112500,
    "created_at": "datetime"
  }
  ```

---

## 37. Reinsurer Statements

### POST /reinsurer-statements
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "treaty_id": "uuid (required)",
    "participant_id": "uuid (required)",
    "period_start": "datetime (required)",
    "period_end": "datetime (required)"
  }
  ```
- **Response (201):** ReinsurerStatementResponse
  ```json
  {
    "id": "uuid",
    "statement_number": "RST-2026-000005",
    "treaty_id": "uuid",
    "participant_id": "uuid",
    "period_start": "datetime",
    "period_end": "datetime",
    "premium_ceded": 15000000,
    "claims_recovered": 4500000,
    "commission_due": 3750000,
    "profit_commission": 500000,
    "net_balance": 6250000,
    "status": "DRAFT|ISSUED|ACKNOWLEDGED|SETTLED",
    "created_by": "uuid",
    "created_at": "datetime",
    "updated_at": "datetime"
  }
  ```
- **Business Rules:**
  - Fetches participant's share and commission rate
  - `premium_ceded = totalCeded × sharePercentage / 100`
  - `claims_recovered = totalRecovered × sharePercentage / 100`
  - **Important:** `totalRecovered` is the **all-time cumulative total** of recoveries for this treaty — it is NOT filtered by the statement's `period_start`/`period_end`. The period parameters affect only the statement metadata, not the claims figure.
  - `commission_due = premium_ceded × commissionRate / 100`
  - `net_balance = premium_ceded - claims_recovered - commission_due`
  - **Integer truncation warning:** `sharePercentage` and `commissionRate` are `float64` but cast to `int64` before division. A 12.5% share truncates to 12%. Fractional percentages lose precision.

### POST /reinsurer-statements/profit-commission
- **Auth:** Required
- **Request Body:**
  ```json
  {
    "treaty_id": "uuid (required)",
    "period_start": "datetime (required)",
    "period_end": "datetime (required)"
  }
  ```
- **Response (200):**
  ```json
  {
    "treaty_id": "uuid",
    "premium_ceded": 15000000,
    "claims_recovered": 4500000,
    "loss_ratio": 30.0,
    "net_profit": 10500000,
    "commission_rate": 25.0,
    "commission_amount": 2625000,
    "carry_forward": 0
  }
  ```
- **Business Rules:**
  - `loss_ratio = (claims_recovered × 100) / premium_ceded` (0 if premiumCeded = 0)
  - Match loss ratio to profit commission rule where `LossRatioFrom ≤ lossRatio ≤ LossRatioTo`
  - `net_profit = premium_ceded - claims_recovered`
  - If any CARRY_FORWARD rule has CarryForwardBalance > 0: `net_profit -= CarryForwardBalance`
  - If net_profit > 0: `commission_amount = floor(net_profit × commission_rate / 100)`, carryForward = 0
  - If net_profit ≤ 0: `commission_amount = 0`, `carry_forward = -net_profit` (deficit carried forward to next period)

### GET /reinsurer-statements/:id
- **Auth:** Required

### PUT /reinsurer-statements/:id/issue
- **Auth:** Required
- **Business Rules:** DRAFT → ISSUED.

### PUT /reinsurer-statements/:id/acknowledge
- **Auth:** Required
- **Business Rules:** ISSUED → ACKNOWLEDGED.

### PUT /reinsurer-statements/:id/settle
- **Auth:** Required
- **Business Rules:** ACKNOWLEDGED → SETTLED. Final settlement confirmation.

---

## 38. Treaty Alerts

### GET /treaty-alerts
- **Auth:** Required
- **Query Params:** `page`, `page_size`
- **Response (200):** Paginated list of TreatyAlertResponse
  ```json
  {
    "id": "uuid",
    "treaty_id": "uuid",
    "treaty_layer_id": "uuid",
    "alert_type": "LIMIT_BREACH|AGGREGATE_WARNING|CATASTROPHE_THRESHOLD|EXPIRY_WARNING",
    "severity": "LOW|MEDIUM|HIGH|CRITICAL",
    "message": "string",
    "threshold_value": 200000000,
    "current_value": 180000000,
    "is_acknowledged": false,
    "acknowledged_by": "uuid",
    "acknowledged_at": "datetime",
    "created_at": "datetime"
  }
  ```

### GET /treaty-alerts/unacknowledged
- **Auth:** Required

### GET /treaty-alerts/count
- **Auth:** Required
- **Response (200):** `{ "count": 3 }`

### PUT /treaty-alerts/:id/acknowledge
- **Auth:** Required
- **Business Rules:** Marks alert as acknowledged. Records user and timestamp.

### POST /treaty-alerts/check-limits/:treatyId
- **Auth:** Required
- **Business Rules:** Checks all layers with aggregate limits:
  - `usagePercent = aggregate_used × 100 / aggregate_limit`
  - If ≥ 100%: Creates LIMIT_BREACH alert (CRITICAL severity)
  - If ≥ 80%: Creates AGGREGATE_WARNING alert (HIGH severity)

### POST /treaty-alerts/check-catastrophe/:treatyId
- **Auth:** Required
- **Business Rules:** Checks if total recoverable amount exceeds catastrophe threshold (5,000,000 KES / 500,000,000 cents). Creates CATASTROPHE_THRESHOLD alert (CRITICAL severity).

### POST /treaty-alerts/check-expiry
- **Auth:** Required
- **Business Rules:** Checks all ACTIVE treaties expiring within 30 days. Creates EXPIRY_WARNING alert (MEDIUM severity) with days until expiry.

---

## Endpoint Count Summary

| Domain | Count |
|--------|-------|
| Auth | 4 |
| Users | 6 |
| Plans | 4 |
| Benefits | 4 |
| Exclusions | 4 |
| Premium Rules | 5 |
| Underwriting Rules | 4 |
| Provider Networks | 4 |
| Providers | 18 |
| Leads | 9 |
| Quotations | 18 |
| Approval Limits | 3 |
| Policies | 16 |
| Members | 14 |
| Endorsements | 4 |
| Renewals | 8 |
| Underwriting Assessments | 3 |
| Underwriting Flags | 7 |
| Policy Documents | 7 |
| Credit Notes | 3 |
| Pre-Authorization | 7 |
| Claims | 16 |
| Cases | 11 |
| Claim Documents | 3 |
| Provider Statements | 4 |
| Installments | 4 |
| Invoices | 2 |
| Payments | 6 |
| Remittances | 4 |
| Notifications | 3 |
| Audit | 3 |
| Analytics | 4 |
| Treaties | 20 |
| Cessions | 5 |
| Recoveries | 11 |
| Bordereaux | 6 |
| Reinsurer Statements | 6 |
| Treaty Alerts | 7 |
| **TOTAL** | **~246** |

---

## Appendix B: Route-to-Role Mapping (RequireRole Enforcement)

Only the following routes have `RequireRole` middleware restrictions. **All other routes** are accessible to any authenticated user (all 8 roles).

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

**Notes:**
- `RequirePermission` middleware is defined but **never wired to any route** — all RBAC is via `RequireRole` only
- `RequireRole` does NOT have Admin bypass — Admin must be explicitly listed
- Roles `Provider`, `Member`, and `SalesAgent` never appear in any `RequireRole()` call — these roles can access all non-restricted endpoints
- The frontend should implement UI-level restrictions even for unrestricted endpoints (hide admin-only pages from Members, etc.)
