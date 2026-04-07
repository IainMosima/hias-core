# Member Domain — Frontend API Integration Guide

All endpoints require `Authorization: Bearer <token>`. Base URL: `/api/v1`.
Money = BIGINT cents (8000 KES = 800000). IDs = UUIDs. Dates = ISO 8601.

---

## Member Creation (2 ways)

### Standalone (Members page)
```
POST /api/v1/members
```
```json
{
  "policy_id": "uuid",
  "name": "Jane Doe",
  "date_of_birth": "1990-05-15",
  "gender": "female",
  "relationship": "principal",
  "national_id": "12345678",
  "phone": "+254712345678",
  "email": "jane@example.com",
  "kra_pin": "A012345678B",
  "county": "Nairobi",
  "city": "Nairobi",
  "country": "Kenya",
  "address": "123 Kenyatta Ave"
}
```

### Nested under policy (Policy detail page)
```
POST /api/v1/policies/:policyId/members
```
Same body **minus** `policy_id`.

---

## All Endpoints

### Policy-nested
```
POST   /policies/:policyId/members              — Enroll member
GET    /policies/:policyId/members               — List members of policy
POST   /policies/:policyId/members/bulk          — Bulk enroll
POST   /policies/:policyId/members/import        — CSV import (multipart/form-data, field: "file")
POST   /policies/:policyId/members/bulk-remove   — Bulk remove
```

### Standalone
```
GET    /members                        — List all (paginated + search)
POST   /members                        — Create member (policy_id in body)
GET    /members/:id                    — Get member
PUT    /members/:id                    — Update member
DELETE /members/:id                    — Remove member (soft delete)
PUT    /members/:id/verify             — Mark verified
PUT    /members/:id/suspend            — Suspend
PUT    /members/:id/reactivate         — Reactivate
GET    /members/:id/eligibility        — Check claim eligibility
POST   /members/:id/card              — Generate member card
GET    /members/:id/underwriting-flags — Underwriting flags
GET    /members/:id/cases             — Cases
GET    /members/:id/documents         — Documents
```

---

## Request Schemas

### Create/Enroll Member
| Field | Type | Required | Notes |
|-------|------|----------|-------|
| policy_id | string (uuid) | YES (standalone only) | Not needed on nested route |
| name | string | YES | |
| date_of_birth | string | YES | YYYY-MM-DD |
| gender | string | YES | "male" / "female" |
| relationship | string | YES | "principal" / "spouse" / "child" / "parent" |
| national_id | string | no | |
| phone | string | no | |
| email | string | no | |
| kra_pin | string | no | |
| county | string | no | Dropdown — Kenyan counties |
| city | string | no | Dropdown |
| country | string | no | Dropdown, default "Kenya" |
| address | string | no | Free text |

### Update Member (PUT /members/:id) — all optional
```json
{
  "name": "Jane Smith",
  "phone": "+254700000000",
  "email": "jane.smith@example.com",
  "kra_pin": "B987654321A",
  "county": "Mombasa",
  "city": "Mombasa",
  "country": "Kenya",
  "address": "456 Moi Ave"
}
```

### Remove Member (DELETE /members/:id) — optional body
```json
{ "reason": "Left the company" }
```

### Bulk Enroll (POST /policies/:id/members/bulk)
```json
{
  "members": [
    { "name": "John", "date_of_birth": "1985-03-20", "gender": "male", "relationship": "principal", "county": "Nairobi", "city": "Nairobi", "country": "Kenya" },
    { "name": "Mary", "date_of_birth": "1988-07-10", "gender": "female", "relationship": "spouse" }
  ]
}
```

### Bulk Remove (POST /policies/:id/members/bulk-remove)
```json
{
  "member_ids": ["uuid1", "uuid2"],
  "reason": "Policy cancelled"
}
```

### CSV Import (POST /policies/:id/members/import)
- Content-Type: `multipart/form-data`
- Field: `file` (CSV)
- Required columns: `name`, `date_of_birth`, `gender`, `relationship`
- Optional columns: `national_id`, `phone`, `email`, `kra_pin`, `county`, `city`, `country`, `address`

---

## Response Schemas

### Member Response
```json
{
  "id": "uuid",
  "policy_id": "uuid",
  "national_id": "12345678",
  "name": "Jane Doe",
  "date_of_birth": "1990-05-15T00:00:00Z",
  "gender": "female",
  "relationship": "principal",
  "member_number": "MBR-2026-000001",
  "phone": "+254712345678",
  "email": "jane@example.com",
  "kra_pin": "A012345678B",
  "county": "Nairobi",
  "city": "Nairobi",
  "country": "Kenya",
  "address": "123 Kenyatta Ave",
  "status": "ACTIVE",
  "verified": false,
  "verified_at": null,
  "coverage_start_date": "2026-01-01T00:00:00Z",
  "coverage_end_date": "2026-12-31T00:00:00Z",
  "created_at": "2026-04-07T10:00:00Z"
}
```

### Bulk Result Response
```json
{
  "succeeded": 2,
  "failed": 1,
  "members": [ ...MemberResponse[] ],
  "errors": ["Member 3: Invalid date of birth"]
}
```

### Eligibility Response (GET /members/:id/eligibility)
```json
{ "status": "success", "message": "Eligibility checked", "data": true }
```

---

## List / Pagination

### List All Members (GET /members)
Query params: `?search=jane&page=1&page_size=10`

Search matches: name, national_id, email, phone.

```json
{
  "status": "success",
  "message": "Members retrieved",
  "data": [ ...MemberResponse[] ],
  "page": 1,
  "page_size": 10,
  "total_count": 150,
  "total_pages": 15
}
```

### List by Policy (GET /policies/:id/members)
Returns all members for that policy (no pagination).

---

## Member Status Workflow

```
PENDING  ──(policy activated)──>  ACTIVE
ACTIVE   ──(suspend)──────────>  SUSPENDED
SUSPENDED ──(reactivate)───────>  ACTIVE
ACTIVE/SUSPENDED ──(remove)────>  REMOVED (terminal)
```

- New members start as `PENDING` if policy is `DRAFT`, `ACTIVE` if policy is `ACTIVE`
- `REMOVED` members cannot be updated
- Removal triggers premium recalculation + pro-rata credit note

---

## Address Form Fields (Recommended Order)

The backend supports these address fields. Recommended form order:

1. **Country** — dropdown, default "Kenya"
2. **County** — dropdown (filter by country)
3. **City** — dropdown (filter by county)
4. **Address** — free text input

---

## PII Masking

Non-admin users receive masked PII:
- `national_id`: `"****5678"`
- `email`: `"j***@example.com"`
- `phone`: `"+254****5678"`

Admin users see full values.

---

## Underwriting Validations (on create)

The backend automatically runs these checks on enrollment:
1. **Double insurance** — rejects if national_id already covered under another active policy
2. **Age validation** — checks against plan premium rules (min/max age per relationship)
3. **Underwriting rules** — evaluates plan-specific rules, creates flags for violations

Errors return 400 with descriptive message.

---

## API Response Envelope

**Success:**
```json
{ "status": "success", "message": "Member enrolled", "data": { ...MemberResponse } }
```

**Error:**
```json
{ "status": "error", "message": "Double insurance detected: member already covered under another active policy" }
```

**Paginated:**
```json
{ "status": "success", "message": "Members retrieved", "data": [...], "page": 1, "page_size": 10, "total_count": 50, "total_pages": 5 }
```
