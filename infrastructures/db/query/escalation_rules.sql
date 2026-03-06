-- name: CreateEscalationRule :one
INSERT INTO escalation_rules (name, condition_type, threshold_amount, escalation_role, is_active)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetEscalationRuleByID :one
SELECT * FROM escalation_rules WHERE id = $1;

-- name: ListEscalationRules :many
SELECT * FROM escalation_rules ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: ListActiveEscalationRules :many
SELECT * FROM escalation_rules WHERE is_active = TRUE ORDER BY threshold_amount DESC;

-- name: UpdateEscalationRule :one
UPDATE escalation_rules
SET name = $2, condition_type = $3, threshold_amount = $4, escalation_role = $5, is_active = $6, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteEscalationRule :exec
DELETE FROM escalation_rules WHERE id = $1;
