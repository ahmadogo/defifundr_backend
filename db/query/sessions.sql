-- name: CreateSession :one
INSERT INTO sessions (
  id,
  user_id,
  refresh_token,
  user_agent,
  client_ip,
  is_blocked,
  expires_at
) VALUES (
  uuid_generate_v4(), $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetSession :one
SELECT * FROM sessions
WHERE id = $1 LIMIT 1;

-- name: GetSessionByRefreshToken :one
SELECT * FROM sessions
WHERE refresh_token = $1 LIMIT 1;

-- name: ListUserSessions :many
SELECT * FROM sessions
WHERE user_id = $1 AND is_blocked = false AND expires_at > now()
ORDER BY created_at DESC;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE id = $1;

-- name: BlockSession :exec
UPDATE sessions
SET is_blocked = true
WHERE id = $1;

-- name: DeleteExpiredSessions :exec
DELETE FROM sessions
WHERE expires_at < now();