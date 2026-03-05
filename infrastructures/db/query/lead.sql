-- name: CreateLead :one
INSERT INTO leads (lead_number, contact_name, contact_email, contact_phone, company_name, source, segment, plan_type, estimated_members, expected_premium, closure_probability, currency, status, assigned_to, next_follow_up_date, notes, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17) RETURNING *;

-- name: GetLeadByID :one
SELECT * FROM leads WHERE id = $1;

-- name: GetLeadByNumber :one
SELECT * FROM leads WHERE lead_number = $1;

-- name: ListLeads :many
SELECT * FROM leads ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: ListLeadsByStatus :many
SELECT * FROM leads WHERE status = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: ListLeadsByAssignedTo :many
SELECT * FROM leads WHERE assigned_to = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: ListDueFollowUps :many
SELECT * FROM leads
WHERE next_follow_up_date <= NOW()
  AND status NOT IN ('WON', 'LOST', 'DORMANT')
ORDER BY next_follow_up_date ASC
LIMIT $1 OFFSET $2;

-- name: UpdateLead :one
UPDATE leads SET
    contact_name = COALESCE(NULLIF(sqlc.arg('contact_name')::text, ''), contact_name),
    contact_email = COALESCE(NULLIF(sqlc.arg('contact_email')::text, ''), contact_email),
    contact_phone = COALESCE(NULLIF(sqlc.arg('contact_phone')::text, ''), contact_phone),
    company_name = COALESCE(NULLIF(sqlc.arg('company_name')::text, ''), company_name),
    source = COALESCE(NULLIF(sqlc.arg('source')::text, ''), source),
    segment = COALESCE(NULLIF(sqlc.arg('segment')::text, ''), segment),
    plan_type = COALESCE(NULLIF(sqlc.arg('plan_type')::text, ''), plan_type),
    estimated_members = COALESCE(sqlc.narg('estimated_members'), estimated_members),
    expected_premium = COALESCE(sqlc.narg('expected_premium'), expected_premium),
    closure_probability = COALESCE(sqlc.narg('closure_probability'), closure_probability),
    assigned_to = COALESCE(sqlc.narg('assigned_to'), assigned_to),
    next_follow_up_date = COALESCE(sqlc.narg('next_follow_up_date'), next_follow_up_date),
    notes = COALESCE(NULLIF(sqlc.arg('notes')::text, ''), notes),
    updated_at = NOW()
WHERE id = sqlc.arg('id') RETURNING *;

-- name: UpdateLeadStatus :one
UPDATE leads SET status = $2, updated_at = NOW() WHERE id = $1 RETURNING *;

-- name: CountLeads :one
SELECT COUNT(*) FROM leads;

-- name: CountLeadsByStatus :one
SELECT COUNT(*) FROM leads WHERE status = $1;
