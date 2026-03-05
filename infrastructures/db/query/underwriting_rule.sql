-- name: CreateUnderwritingRule :one
INSERT INTO underwriting_rules (plan_id, rule_type, relationship, parameter_key, parameter_value, severity, risk_score_weight, is_blocking, is_active, description)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING *;

-- name: GetUnderwritingRuleByID :one
SELECT * FROM underwriting_rules WHERE id = $1;

-- name: ListUnderwritingRulesByPlan :many
SELECT * FROM underwriting_rules WHERE plan_id = $1 ORDER BY created_at DESC;

-- name: ListActiveUnderwritingRulesByPlan :many
SELECT * FROM underwriting_rules WHERE plan_id = $1 AND is_active = true ORDER BY rule_type;

-- name: UpdateUnderwritingRule :one
UPDATE underwriting_rules SET
    rule_type = COALESCE(sqlc.narg('rule_type'), rule_type),
    relationship = COALESCE(sqlc.narg('relationship'), relationship),
    parameter_key = COALESCE(sqlc.narg('parameter_key'), parameter_key),
    parameter_value = COALESCE(sqlc.narg('parameter_value'), parameter_value),
    severity = COALESCE(sqlc.narg('severity'), severity),
    risk_score_weight = COALESCE(sqlc.narg('risk_score_weight'), risk_score_weight),
    is_blocking = COALESCE(sqlc.narg('is_blocking'), is_blocking),
    is_active = COALESCE(sqlc.narg('is_active'), is_active),
    description = COALESCE(sqlc.narg('description'), description),
    updated_at = NOW()
WHERE id = sqlc.arg('id') RETURNING *;

-- name: DeleteUnderwritingRule :exec
DELETE FROM underwriting_rules WHERE id = $1;
