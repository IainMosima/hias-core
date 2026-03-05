-- name: CreateUnderwritingAssessment :one
INSERT INTO underwriting_assessments (policy_id, member_id, status, questionnaire, medical_declarations, created_by)
VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: GetUnderwritingByID :one
SELECT * FROM underwriting_assessments WHERE id = $1;

-- name: ListUnderwritingByPolicy :many
SELECT * FROM underwriting_assessments WHERE policy_id = $1 ORDER BY created_at DESC;

-- name: GetUnderwritingByMember :one
SELECT * FROM underwriting_assessments WHERE member_id = $1 ORDER BY created_at DESC LIMIT 1;

-- name: UpdateUnderwritingStatus :one
UPDATE underwriting_assessments SET status = $2, updated_at = NOW() WHERE id = $1 RETURNING *;

-- name: UpdateUnderwriting :one
UPDATE underwriting_assessments SET
    status = COALESCE(sqlc.narg('status'), status),
    risk_score = COALESCE(sqlc.narg('risk_score'), risk_score),
    risk_flags = COALESCE(sqlc.narg('risk_flags'), risk_flags),
    decision_reason = COALESCE(sqlc.narg('decision_reason'), decision_reason),
    assessed_by = COALESCE(sqlc.narg('assessed_by'), assessed_by),
    assessed_at = COALESCE(sqlc.narg('assessed_at'), assessed_at),
    updated_at = NOW()
WHERE id = sqlc.arg('id') RETURNING *;
