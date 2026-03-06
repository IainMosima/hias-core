-- name: CreateAdjudicationRule :one
INSERT INTO adjudication_rules (name, rule_type, parameters, priority, is_active)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetAdjudicationRuleByID :one
SELECT * FROM adjudication_rules WHERE id = $1;

-- name: ListAdjudicationRules :many
SELECT * FROM adjudication_rules ORDER BY priority ASC LIMIT $1 OFFSET $2;

-- name: ListActiveAdjudicationRules :many
SELECT * FROM adjudication_rules WHERE is_active = TRUE ORDER BY priority ASC;

-- name: UpdateAdjudicationRule :one
UPDATE adjudication_rules
SET name = $2, rule_type = $3, parameters = $4, priority = $5, is_active = $6, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteAdjudicationRule :exec
DELETE FROM adjudication_rules WHERE id = $1;
