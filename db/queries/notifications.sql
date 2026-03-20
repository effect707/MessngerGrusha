-- name: CreateNotification :exec
INSERT INTO notifications (id, user_id, type, payload, is_read, created_at)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetUserNotifications :many
SELECT id, user_id, type, payload, is_read, created_at
FROM notifications
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2;

-- name: GetUnreadNotifications :many
SELECT id, user_id, type, payload, is_read, created_at
FROM notifications
WHERE user_id = $1 AND is_read = false
ORDER BY created_at DESC
LIMIT $2;

-- name: MarkNotificationRead :exec
UPDATE notifications SET is_read = true WHERE id = $1 AND user_id = $2;

-- name: MarkAllNotificationsRead :exec
UPDATE notifications SET is_read = true WHERE user_id = $1 AND is_read = false;

-- name: CountUnreadNotifications :one
SELECT COUNT(*) FROM notifications
WHERE user_id = $1 AND is_read = false;
