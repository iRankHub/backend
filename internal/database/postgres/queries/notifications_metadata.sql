-- name: CreateNotificationMetadata :one
INSERT INTO notification_metadata (
    notification_id,
    user_id,
    category,
    type,
    status,
    priority,
    delivery_methods,
    delivery_status,
    metadata,
    expires_at,
    file_size
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;

-- name: GetNotificationMetadata :one
SELECT * FROM notification_metadata
WHERE notification_id = $1 LIMIT 1;

-- name: GetUnreadNotificationsForUser :many
SELECT * FROM notification_metadata
WHERE user_id = $1
  AND is_read = FALSE
  AND expires_at > CURRENT_TIMESTAMP
ORDER BY created_at DESC;

-- name: GetNotificationsByStatus :many
SELECT * FROM notification_metadata
WHERE status = $1
  AND expires_at > CURRENT_TIMESTAMP
ORDER BY created_at DESC;

-- name: UpdateNotificationStatus :exec
UPDATE notification_metadata
SET status = $2,
    delivery_status = $3,
    updated_at = CURRENT_TIMESTAMP
WHERE notification_id = $1;

-- name: MarkNotificationAsRead :exec
UPDATE notification_metadata
SET is_read = TRUE,
    read_at = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
WHERE notification_id = $1;

-- name: DeleteExpiredNotifications :exec
DELETE FROM notification_metadata
WHERE expires_at < CURRENT_TIMESTAMP;

-- name: UpdateNotificationDeliveryStatus :exec
UPDATE notification_metadata
SET delivery_status = delivery_status || $2::jsonb,
    updated_at = CURRENT_TIMESTAMP
WHERE notification_id = $1;

-- name: GetNotificationsToRetry :many
SELECT * FROM notification_metadata
WHERE status = 'failed'
  AND expires_at > CURRENT_TIMESTAMP
  AND updated_at < CURRENT_TIMESTAMP - INTERVAL '30 minutes'
ORDER BY created_at DESC;