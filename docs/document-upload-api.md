# Document Upload API — Frontend Integration Guide

## Overview

S3 presigned URL upload system for KYC documents (ID copies, photos, medical declarations) attached to policies, members, claims, or quotations. Uses a **three-step flow**: request URL, upload directly to S3, then confirm.

All endpoints require `Authorization: Bearer <token>`.

**Base URL:** `/api/v1`

---

## Upload Flow

```
Frontend                        Backend                         S3
   |                               |                            |
   |  1. POST /documents/upload-url|                            |
   |  (file metadata)              |                            |
   |------------------------------>|                            |
   |                               |  Creates DB record         |
   |                               |  (status=PENDING_UPLOAD)   |
   |                               |  Generates presigned PUT   |
   |  { document_id, upload_url }  |                            |
   |<------------------------------|                            |
   |                               |                            |
   |  2. PUT upload_url            |                            |
   |  (raw file bytes)             |                            |
   |------------------------------------------------------>    |
   |                          200 OK                            |
   |<------------------------------------------------------|   |
   |                               |                            |
   |  3. POST /documents/{id}/confirm-upload                    |
   |------------------------------>|                            |
   |                               |  HeadObject (verify file)  |
   |                               |--------------------------->|
   |                               |  status -> ACTIVE          |
   |  { document details }         |                            |
   |<------------------------------|                            |
```

---

## Endpoints

### 1. Request Upload URL

Creates a document record and returns a presigned S3 PUT URL.

```
POST /api/v1/documents/upload-url
```

**Request Body:**

| Field           | Type   | Required | Description                                              |
|-----------------|--------|----------|----------------------------------------------------------|
| `entity_type`   | string | Yes      | One of: `policy`, `member`, `claim`, `quotation`         |
| `entity_id`     | string | Yes      | UUID of the parent entity                                |
| `document_type` | string | Yes      | e.g. `KYC_ID`, `KYC_PASSPORT`, `KYC_PHOTO`, `SUPPORTING`|
| `file_name`     | string | Yes      | Original file name (validated server-side)               |
| `file_size`     | int    | Yes      | File size in bytes, must be > 0                          |
| `mime_type`     | string | Yes      | MIME type (validated against server allowlist)            |

```json
{
  "entity_type": "policy",
  "entity_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "document_type": "KYC_ID",
  "file_name": "national-id-front.pdf",
  "file_size": 245000,
  "mime_type": "application/pdf"
}
```

**Success Response (200):**

```json
{
  "status": "success",
  "message": "Upload URL generated",
  "data": {
    "document_id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
    "upload_url": "https://s3.amazonaws.com/bucket/policy/a1b2.../KYC_ID/1710000000-abc12345-national-id-front.pdf?X-Amz-...",
    "s3_key": "policy/a1b2c3d4-e5f6-7890-abcd-ef1234567890/KYC_ID/1710000000-abc12345-national-id-front.pdf",
    "expires_in": 900
  }
}
```

**Error Responses:**

| Code | Reason                                                  |
|------|---------------------------------------------------------|
| 400  | Invalid file name, size exceeds limit, blocked MIME type, blocked extension, missing fields |

---

### 2. Upload File to S3 (Frontend Direct)

Use the `upload_url` from step 1 to PUT the file directly to S3. **This request goes to S3, not the backend.**

```
PUT <upload_url>
Content-Type: <same mime_type from step 1>
Body: <raw file bytes>
```

**Important:**
- The `Content-Type` header MUST match the `mime_type` sent in step 1 — S3 will reject the request otherwise.
- The presigned URL expires in **900 seconds (15 minutes)**.
- S3 returns `200 OK` on success.
- Do NOT send `Authorization` header to S3 — the presigned URL contains its own auth.

**Example (fetch):**

```javascript
const response = await fetch(uploadUrl, {
  method: 'PUT',
  headers: { 'Content-Type': mimeType },
  body: file, // File or Blob object
});

if (response.ok) {
  // Proceed to step 3
}
```

---

### 3. Confirm Upload

Verifies the file exists in S3 (HeadObject) and marks the document as `ACTIVE`.

```
POST /api/v1/documents/:id/confirm-upload
```

**Path Params:** `id` — the `document_id` from step 1.

No request body required.

**Success Response (200):**

```json
{
  "status": "success",
  "message": "Upload confirmed",
  "data": {
    "id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
    "entity_type": "policy",
    "entity_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "document_type": "KYC_ID",
    "status": "ACTIVE",
    "file_name": "national-id-front.pdf",
    "file_size": 245000,
    "mime_type": "application/pdf",
    "s3_key": "policy/a1b2.../KYC_ID/1710000000-abc12345-national-id-front.pdf",
    "uploaded_by": "user-uuid-here",
    "confirmed_at": "2026-03-11T10:30:00Z",
    "created_at": "2026-03-11T10:28:00Z"
  }
}
```

**Error Responses:**

| Code | Reason                                                        |
|------|---------------------------------------------------------------|
| 400  | Document status is not `PENDING_UPLOAD`, or file not found in S3 (upload failed/expired) |
| 404  | Document ID not found                                         |

---

### 4. Bulk Request Upload URLs

Same as step 1 but for multiple files at once (max 10). Fails entirely if any single file fails validation.

```
POST /api/v1/documents/bulk-upload-urls
```

**Request Body:**

```json
{
  "files": [
    {
      "entity_type": "member",
      "entity_id": "member-uuid-1",
      "document_type": "KYC_PHOTO",
      "file_name": "passport-photo.jpg",
      "file_size": 120000,
      "mime_type": "image/jpeg"
    },
    {
      "entity_type": "member",
      "entity_id": "member-uuid-1",
      "document_type": "KYC_ID",
      "file_name": "id-scan.pdf",
      "file_size": 350000,
      "mime_type": "application/pdf"
    }
  ]
}
```

**Success Response (200):**

```json
{
  "status": "success",
  "message": "Upload URLs generated",
  "data": [
    {
      "document_id": "uuid-1",
      "upload_url": "https://s3...",
      "s3_key": "member/member-uuid-1/KYC_PHOTO/...",
      "expires_in": 900
    },
    {
      "document_id": "uuid-2",
      "upload_url": "https://s3...",
      "s3_key": "member/member-uuid-1/KYC_ID/...",
      "expires_in": 900
    }
  ]
}
```

After receiving the array, upload each file to its respective `upload_url`, then call `POST /documents/:id/confirm-upload` for each.

---

### 5. Get Download URL

Returns a presigned GET URL (15 min expiry) to download a confirmed document.

```
POST /api/v1/documents/:id/download-url
```

No request body required.

**Success Response (200):**

```json
{
  "status": "success",
  "message": "Download URL generated",
  "data": {
    "download_url": "https://s3.amazonaws.com/bucket/policy/...?X-Amz-..."
  }
}
```

| Code | Reason                            |
|------|-----------------------------------|
| 400  | Document is not `ACTIVE`          |
| 404  | Document not found                |

---

### 6. Delete Document

Soft deletes the document record and removes the file from S3.

```
DELETE /api/v1/documents/:id
```

**Success Response (200):**

```json
{
  "status": "success",
  "message": "Document deleted successfully",
  "data": "Document deleted"
}
```

---

### 7. List Policy Uploads

Returns all uploaded documents attached to a specific policy.

```
GET /api/v1/policies/:id/uploads
```

**Success Response (200):**

```json
{
  "status": "success",
  "message": "Documents retrieved",
  "data": [
    {
      "id": "doc-uuid",
      "entity_type": "policy",
      "entity_id": "policy-uuid",
      "document_type": "KYC_ID",
      "status": "ACTIVE",
      "file_name": "national-id.pdf",
      "file_size": 245000,
      "mime_type": "application/pdf",
      "s3_key": "policy/.../KYC_ID/...",
      "uploaded_by": "user-uuid",
      "confirmed_at": "2026-03-11T10:30:00Z",
      "created_at": "2026-03-11T10:28:00Z"
    }
  ]
}
```

---

### 8. List Member Documents

Returns all uploaded documents attached to a specific member.

```
GET /api/v1/members/:id/documents
```

Same response shape as List Policy Uploads but with `entity_type: "member"`.

---

## Document Statuses

| Status            | Meaning                                          |
|-------------------|--------------------------------------------------|
| `PENDING_UPLOAD`  | Record created, waiting for S3 upload + confirm   |
| `ACTIVE`          | Upload confirmed, file verified in S3             |
| `UPLOAD_FAILED`   | Reserved for future cleanup jobs                  |
| `DELETED`         | Soft deleted                                      |

---

## Document Types

Use these values for the `document_type` field:

| Value            | Use Case                        |
|------------------|---------------------------------|
| `KYC_ID`         | National ID / government ID     |
| `KYC_PASSPORT`   | Passport copy                   |
| `KYC_PHOTO`      | Passport-size photo             |
| `MEDICAL_DECLARATION` | Medical declaration form   |
| `SUPPORTING`     | Any other supporting document   |

These are not enforced as an enum on the backend — you can send any string. The above are the recommended conventions.

---

## Validation Rules

**File Name:**
- Must not be empty
- Must not contain `/` or `\` (path separators)
- Must not contain `..` (path traversal)
- Must not contain null bytes
- Blocked extensions: `.exe`, `.sh`, `.bat`, `.cmd`, `.php`, `.jsp`, `.com`, `.msi`, `.ps1`, `.vbs`, `.js`, `.wsf`

**File Size:**
- Must be greater than 0
- Must not exceed server-configured max (set via `AWS_S3_MAX_FILE_SIZE` env var)

**MIME Type:**
- Must not be empty
- Must be in the server-configured allowlist (set via `AWS_S3_ALLOWED_TYPES` env var)
- Common allowed types: `application/pdf`, `image/jpeg`, `image/png`, `image/webp`

---

## Frontend Integration Example

```typescript
// Full upload flow for a single document
async function uploadDocument(
  entityType: string,
  entityId: string,
  documentType: string,
  file: File
) {
  // Step 1: Request upload URL
  const urlResponse = await api.post('/documents/upload-url', {
    entity_type: entityType,
    entity_id: entityId,
    document_type: documentType,
    file_name: file.name,
    file_size: file.size,
    mime_type: file.type,
  });

  const { document_id, upload_url } = urlResponse.data.data;

  // Step 2: Upload directly to S3
  await fetch(upload_url, {
    method: 'PUT',
    headers: { 'Content-Type': file.type },
    body: file,
  });

  // Step 3: Confirm upload
  const confirmResponse = await api.post(
    `/documents/${document_id}/confirm-upload`
  );

  return confirmResponse.data.data; // DocumentResponse
}

// Bulk upload for KYC onboarding
async function uploadKYCDocuments(
  memberId: string,
  files: { type: string; file: File }[]
) {
  // Step 1: Bulk request URLs
  const urlResponse = await api.post('/documents/bulk-upload-urls', {
    files: files.map(f => ({
      entity_type: 'member',
      entity_id: memberId,
      document_type: f.type,
      file_name: f.file.name,
      file_size: f.file.size,
      mime_type: f.file.type,
    })),
  });

  const uploadInfos = urlResponse.data.data;

  // Step 2: Upload all to S3 in parallel
  await Promise.all(
    uploadInfos.map((info, i) =>
      fetch(info.upload_url, {
        method: 'PUT',
        headers: { 'Content-Type': files[i].file.type },
        body: files[i].file,
      })
    )
  );

  // Step 3: Confirm all uploads
  const confirmed = await Promise.all(
    uploadInfos.map(info =>
      api.post(`/documents/${info.document_id}/confirm-upload`)
    )
  );

  return confirmed.map(r => r.data.data);
}
```

---

## Existing Endpoints (Unchanged)

These pre-existing endpoints continue to work as before and now also include documents from the new `documents` table:

| Method | Path                              | Description                                |
|--------|-----------------------------------|--------------------------------------------|
| GET    | `/documents/standalone`           | Paginated list of ALL documents system-wide (now includes uploaded docs) |
| GET    | `/documents/:id/download`         | Legacy download — looks up across all tables |

---

## Error Response Format

All errors follow the standard format:

```json
{
  "status": "error",
  "message": "Human-readable error description"
}
```
