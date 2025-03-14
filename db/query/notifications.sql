-- name: CreateNotification :one
INSERT INTO notifications (
  user_id,
  message,
  type,
  is_read
) VALUES (
  $1, $2, $3, false
) RETURNING *;

-- name: GetNotification :one
SELECT * FROM notifications
WHERE id = $1 LIMIT 1;

-- name: ListUserNotifications :many
SELECT * FROM notifications
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2
OFFSET $3;

-- name: MarkNotificationRead :one
UPDATE notifications
SET is_read = true
WHERE id = $1
RETURNING *;

-- name: MarkAllUserNotificationsRead :exec
UPDATE notifications
SET is_read = true
WHERE user_id = $1 AND is_read = false;

-- name: GetUnreadNotificationCount :one
SELECT COUNT(*) FROM notifications
WHERE user_id = $1 AND is_read = false;

-- name: DeleteNotification :exec
DELETE FROM notifications
WHERE id = $1 AND user_id = $2;