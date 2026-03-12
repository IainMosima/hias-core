# Document Generation — Single Source of Truth

> **Last updated:** 2026-03-11
> **Base URL:** `/api/v1`
> **Auth:** Bearer token (PASETO) required on all endpoints
> **Content-Type:** `application/json`

---

## Table of Contents

1. [Architecture & Backend Changes](#1-architecture--backend-changes)
2. [Generation Flow (Internal)](#2-generation-flow-internal)
3. [API Endpoints](#3-api-endpoints)
4. [Data Types & Enums](#4-data-types--enums)
5. [Response Shapes](#5-response-shapes)
6. [Frontend Integration Patterns](#6-frontend-integration-patterns)

---

## 1. Architecture & Backend Changes

### What Changed

The document generation pipeline had a race condition: it deleted the PENDING record then tried to create a new GENERATED record. If the create failed, the record was permanently lost.

**Before (unsafe):**
```
DELETE pending record → CREATE new generated record
                        ↑ if this fails, data is gone
```

**After (atomic):**
```
UPDATE pending record → set file_name, file_size, s3_key, status=GENERATED
                        ↑ if this fails, PENDING record survives for retry
```

### Files Changed

| Layer | File | What |
|-------|------|------|
| SQL | `infrastructures/db/query/policy_document.sql` | Added `UpdatePolicyDocumentGenerated` query |
| SQLC | `infrastructures/db/sqlc/*.go` | Auto-regenerated |
| Domain | `domains/policy/repository/policy_document_repository.go` | Added `UpdateGenerated` interface method |
| Infra | `infrastructures/repository/policy_document_repository.go` | Implemented `UpdateGenerated` |
| Service | `services/policy/policy_document_service_impl.go` | Replaced delete+create with atomic update |

### New SQL Query

```sql
UPDATE policy_documents
SET file_name = $2, file_size = $3, s3_key = $4, status = $5, updated_at = NOW()
WHERE id = $1
RETURNING *;
```

---

## 2. Generation Flow (Internal)

```
Step 1  GetNextVersion()         → determine version number
Step 2  CreateV2(PENDING)        → insert record with status=PENDING
Step 3  Generate PDF             → render the document bytes
Step 4  Upload to S3             → store file in cloud storage
Step 5  UpdateGenerated()        → atomically set file info + status=GENERATED
Step 6  Supersede previous       → mark old version superseded_by (if version > 1)
```

**Failure handling:**
- Step 3 or 4 fails → record stays PENDING, status updated to FAILED with error_message
- Step 5 fails → PENDING record survives, can be retried or cleaned up
- No data loss at any step

---

## 3. API Endpoints

### Endpoints at a Glance

| Method | Endpoint | Description |
|--------|----------|-------------|
| **`POST`** | **`/documents/generate`** | **Generate any document (primary endpoint)** |
| `GET` | `/documents/can-generate` | Pre-flight: can this document be generated? |
| `GET` | `/documents/availability` | All document types + status for an entity |
| `GET` | `/policies/:id/documents` | List all documents for a policy |
| `GET` | `/policy-documents/:id` | Get single document by ID |
| `DELETE` | `/policy-documents/:id` | Delete document by ID |
| `DELETE` | `/policies/:id/documents/:docId` | Delete document (nested path) |

**Shortcut endpoints** (no request body needed):

| Method | Endpoint | Equivalent |
|--------|----------|------------|
| `POST` | `/policies/:id/documents/welcome-letter` | `entity_type=policy, doc=WELCOME_LETTER` |
| `POST` | `/policies/:id/documents/policy-schedule` | `entity_type=policy, doc=POLICY_SCHEDULE` |
| `POST` | `/policies/:id/documents/member-cards` | Bulk `MEMBER_CARD` for all policy members |
| `POST` | `/members/:id/card` | `entity_type=member, doc=MEMBER_CARD` |
| `POST` | `/preauths/:id/lou` | `entity_type=preauth, doc=LOU` |
| `POST` | `/claims/:id/decline-letter` | `DECLINE_LETTER` (claim must have rejection_reason) |

---

### 3.1 Generate Document

```
POST /api/v1/documents/generate
```

**Request:**
```json
{
  "entity_type": "policy",
  "entity_id": "550e8400-e29b-41d4-a716-446655440000",
  "document_type": "WELCOME_LETTER"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `entity_type` | string | Yes | See [Entity Types](#entity-types) |
| `entity_id` | string (UUID) | Yes | ID of the target entity |
| `document_type` | string | Yes | See [Document Types](#document-types) |

**Response — `201 Created`:**
```json
{
  "status": "success",
  "message": "Document generated successfully",
  "data": {
    "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "policy_id": "550e8400-e29b-41d4-a716-446655440000",
    "member_id": null,
    "document_type": "WELCOME_LETTER",
    "file_name": "WELCOME_LETTER_POL-2026-000001_v1.pdf",
    "file_size": 45230,
    "s3_key": "documents/policy/550e8400.../WELCOME_LETTER_v1.pdf",
    "generated_by": "user-uuid-here",
    "version": 1,
    "status": "GENERATED",
    "generation_mode": "MANUAL",
    "entity_type": "policy",
    "entity_id": "550e8400-e29b-41d4-a716-446655440000",
    "superseded_by": null,
    "error_message": "",
    "created_at": "2026-03-11T10:30:00Z",
    "updated_at": "2026-03-11T10:30:02Z"
  }
}
```

**Errors:**

| Status | When |
|--------|------|
| `400` | Invalid body, missing fields, bad UUID, wrong entity_type + document_type combo |
| `404` | Entity not found |
| `409` | Document already PENDING (generation in progress) |
| `500` | PDF render or S3 upload failure |

---

### 3.2 Can Generate (Pre-flight)

```
GET /api/v1/documents/can-generate?entity_type=policy&entity_id=UUID&document_type=WELCOME_LETTER
```

| Param | Type | Required |
|-------|------|----------|
| `entity_type` | string | Yes |
| `entity_id` | UUID | Yes |
| `document_type` | string | Yes |

**Response — `200 OK`:**
```json
{
  "status": "success",
  "message": "Document readiness check complete",
  "data": {
    "can_generate": true,
    "errors": []
  }
}
```

When blocked:
```json
{
  "data": {
    "can_generate": false,
    "errors": ["Policy must be ACTIVE to generate welcome letter"]
  }
}
```

---

### 3.3 Document Availability

```
GET /api/v1/documents/availability?entity_type=policy&entity_id=UUID
```

| Param | Type | Required |
|-------|------|----------|
| `entity_type` | string | Yes |
| `entity_id` | UUID | Yes |

**Response — `200 OK`:**
```json
{
  "status": "success",
  "message": "Document availability retrieved",
  "data": [
    {
      "document_type": "WELCOME_LETTER",
      "exists": true,
      "can_generate": true,
      "latest_status": "GENERATED",
      "latest_version": 1,
      "latest_file_url": "",
      "generated_at": "2026-03-11T10:30:02Z",
      "errors": []
    },
    {
      "document_type": "POLICY_SCHEDULE",
      "exists": false,
      "can_generate": true,
      "latest_status": "",
      "latest_version": 0,
      "latest_file_url": "",
      "generated_at": null,
      "errors": []
    }
  ]
}
```

---

### 3.4 List Policy Documents

```
GET /api/v1/policies/:id/documents
```

Returns array of document objects (same shape as generate response `data`).

---

### 3.5 Get Single Document

```
GET /api/v1/policy-documents/:id
```

Returns single document object.

---

### 3.6 Delete Document

```
DELETE /api/v1/policy-documents/:id
DELETE /api/v1/policies/:policyId/documents/:docId
```

**Response — `200 OK`:**
```json
{
  "status": "success",
  "message": "Document deleted"
}
```

---

## 4. Data Types & Enums

### Entity Types

| Value | Use for |
|-------|---------|
| `policy` | Welcome letters, policy schedules |
| `member` | Member cards |
| `endorsement` | Endorsement confirmations |
| `renewal` | Renewal notices |
| `preauth` | Letters of Undertaking (LOU) |
| `claim` | Decline letters |

### Document Types

| Value | Valid Entity Type | Description |
|-------|-------------------|-------------|
| `WELCOME_LETTER` | `policy` | Welcome letter for new/activated policies |
| `MEMBER_CARD` | `member` | Member insurance card PDF |
| `POLICY_SCHEDULE` | `policy` | Full policy benefits schedule |
| `RENEWAL_NOTICE` | `renewal` | Renewal notification letter |
| `ENDORSEMENT` | `endorsement` | Endorsement confirmation document |
| `LOU` | `preauth` | Letter of Undertaking for providers |
| `DECLINE_LETTER` | `claim` | Claim rejection/decline letter |

### Document Statuses

| Status | Meaning | Frontend Action |
|--------|---------|-----------------|
| `PENDING` | Generation in progress | Show spinner |
| `GENERATED` | Ready, file available in S3 | Show download button |
| `FAILED` | Generation failed | Show error + retry button, check `error_message` |

### Generation Modes

| Value | Meaning |
|-------|---------|
| `MANUAL` | User clicked generate in UI |
| `AUTO` | System-triggered (e.g., on policy activation) |

---

## 5. Response Shapes

### Response Envelope (All Endpoints)

```json
{
  "status": "success",
  "message": "Human-readable message",
  "data": { ... }
}
```

Error:
```json
{
  "status": "error",
  "message": "Policy not found"
}
```

### PolicyDocumentResponse (Full Shape)

```typescript
interface PolicyDocumentResponse {
  id: string;              // UUID
  policy_id: string;       // UUID — always present
  member_id?: string;      // UUID — only for member-level docs
  document_type: string;   // e.g. "WELCOME_LETTER"
  file_name: string;       // e.g. "WELCOME_LETTER_POL-2026-000001_v1.pdf"
  file_size: number;       // bytes (0 while PENDING)
  s3_key: string;          // S3 object key (empty while PENDING)
  generated_by: string;    // UUID of user who triggered generation
  version: number;         // starts at 1, increments on regeneration
  status: string;          // "PENDING" | "GENERATED" | "FAILED"
  generation_mode: string; // "MANUAL" | "AUTO"
  entity_type: string;     // "policy" | "member" | etc.
  entity_id: string;       // UUID of the source entity
  superseded_by?: string;  // UUID of newer version (if superseded)
  error_message?: string;  // populated when status=FAILED
  created_at: string;      // ISO 8601
  updated_at: string;      // ISO 8601
}
```

### DocumentAvailabilityItem

```typescript
interface DocumentAvailabilityItem {
  document_type: string;
  exists: boolean;
  can_generate: boolean;
  latest_status: string;      // "" if no document exists
  latest_version: number;     // 0 if no document exists
  latest_file_url: string;
  generated_at?: string;      // ISO 8601 or null
  errors: string[];
}
```

### DocumentReadinessResponse

```typescript
interface DocumentReadinessResponse {
  can_generate: boolean;
  errors: string[];
}
```

### BulkResultResponse (Member Cards)

```typescript
interface BulkResultResponse {
  succeeded: number;
  failed: number;
  errors?: string[];
}
```

---

## 6. Frontend Integration Patterns

### Pattern 1: Documents Tab on Policy Detail

```
Page Load:
  GET /documents/availability?entity_type=policy&entity_id={policyId}

Render each item:
  - exists=false, can_generate=true  → "Generate" button (enabled)
  - exists=false, can_generate=false → "Generate" button (disabled) + show errors tooltip
  - exists=true,  status=GENERATED   → "Download" link + "Regenerate" button
  - exists=true,  status=PENDING     → spinner
  - exists=true,  status=FAILED      → "Retry" button + show error_message

On Generate click:
  POST /documents/generate { entity_type: "policy", entity_id, document_type }
  → On success: refetch availability
  → On error: show error message
```

### Pattern 2: Quick Action Button (e.g., Welcome Letter on Policy Card)

```
Mount:
  GET /documents/can-generate?entity_type=policy&entity_id={id}&document_type=WELCOME_LETTER
  → can_generate=true  → enable button
  → can_generate=false → disable + tooltip with errors

Click:
  POST /documents/generate { entity_type: "policy", entity_id: "{id}", document_type: "WELCOME_LETTER" }
  → Show success toast with file_name
```

### Pattern 3: Member Cards Bulk Generation

```
Click "Generate All Member Cards":
  POST /policies/{policyId}/documents/member-cards
  → Response: { succeeded: 5, failed: 0, errors: [] }
  → Show toast: "5 member cards generated"
  → Refetch /policies/{policyId}/documents to list them
```

### Pattern 4: Document Version History

```
GET /policies/{policyId}/documents
  → Filter by document_type in frontend
  → Sort by version DESC
  → Current = no superseded_by
  → Old versions = has superseded_by (show as "v1 (superseded)" etc.)
```

---

## Versioning Behavior

When regenerating a document that already exists:

1. New record created with `version = previous + 1`
2. Previous record gets `superseded_by = new_record.id`
3. Only the latest `GENERATED` version is the "current" document
4. Old versions remain in the database for audit trail

```
v1 (GENERATED, superseded_by: v2.id)  ← old
v2 (GENERATED, superseded_by: null)   ← current
```