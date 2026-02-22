-- name: CreateNotification :one
INSERT INTO notifications (user_id, channel, type, subject, body, metadata, status)
VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: GetNotificationByID :one
SELECT * FROM notifications WHERE id = $1;

-- name: ListNotificationsByUser :many
SELECT * FROM notifications WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: ListUnreadNotificationsByUser :many
SELECT * FROM notifications WHERE user_id = $1 AND status != 'READ' ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: CountUnreadNotificationsByUser :one
SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND status != 'READ';

-- name: MarkNotificationRead :one
UPDATE notifications SET status = 'READ', read_at = NOW(), updated_at = NOW() WHERE id = $1 RETURNING *;

-- name: UpdateNotificationStatus :one
UPDATE notifications SET status = $2, updated_at = NOW() WHERE id = $1 RETURNING *;

-- name: MarkNotificationSent :one
UPDATE notifications SET status = 'SENT', sent_at = NOW(), updated_at = NOW() WHERE id = $1 RETURNING *;

-- name: IncrementNotificationRetry :one
UPDATE notifications SET retry_count = retry_count + 1, updated_at = NOW() WHERE id = $1 RETURNING *;

-- name: GetFailedNotificationsForRetry :many
SELECT * FROM notifications WHERE status = 'FAILED' AND retry_count < max_retries ORDER BY created_at ASC LIMIT $1;
