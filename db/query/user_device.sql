-- name: CreateUserDeviceToken :one
INSERT INTO user_device_tokens (
    id,
    user_id,
    device_token,
    platform,
    device_type,
    device_model,
    os_name,
    os_version,
    push_notification_token,
    is_active,
    is_verified,
    app_version,
    client_ip,
    expires_at
) VALUES (
    $1, 
    $2, 
    $3, 
    COALESCE($4, 'unknown'),
    COALESCE($5, 'unknown'),
    COALESCE($6, 'unknown'),
    COALESCE($7, 'unknown'),
    COALESCE($8, 'unknown'),
    $9,
    COALESCE($10, true),
    COALESCE($11, false),
    $12,
    $13,
    $14
) RETURNING *;

-- name: UpsertUserDeviceToken :one
INSERT INTO user_device_tokens (
    id,
    user_id,
    device_token,
    platform,
    device_type,
    device_model,
    os_name,
    os_version,
    push_notification_token,
    is_active,
    is_verified,
    app_version,
    client_ip,
    expires_at
) VALUES (
    $1, 
    $2, 
    $3, 
    COALESCE($4, 'unknown'),
    COALESCE($5, 'unknown'),
    COALESCE($6, 'unknown'),
    COALESCE($7, 'unknown'),
    COALESCE($8, 'unknown'),
    $9,
    COALESCE($10, true),
    COALESCE($11, false),
    $12,
    $13,
    $14
) ON CONFLICT (device_token) DO UPDATE SET
    platform = COALESCE(EXCLUDED.platform, user_device_tokens.platform),
    device_type = COALESCE(EXCLUDED.device_type, user_device_tokens.device_type),
    device_model = COALESCE(EXCLUDED.device_model, user_device_tokens.device_model),
    os_name = COALESCE(EXCLUDED.os_name, user_device_tokens.os_name),
    os_version = COALESCE(EXCLUDED.os_version, user_device_tokens.os_version),
    push_notification_token = COALESCE(EXCLUDED.push_notification_token, user_device_tokens.push_notification_token),
    is_active = COALESCE(EXCLUDED.is_active, user_device_tokens.is_active),
    is_verified = COALESCE(EXCLUDED.is_verified, user_device_tokens.is_verified),
    app_version = COALESCE(EXCLUDED.app_version, user_device_tokens.app_version),
    client_ip = COALESCE(EXCLUDED.client_ip, user_device_tokens.client_ip),
    expires_at = COALESCE(EXCLUDED.expires_at, user_device_tokens.expires_at),
    last_used_at = NOW()
RETURNING *;

-- name: GetUserDeviceTokenByID :one
SELECT * FROM user_device_tokens
WHERE id = $1;

-- name: GetUserDeviceTokenByDeviceToken :one
SELECT * FROM user_device_tokens
WHERE device_token = $1;

-- name: GetActiveDeviceTokensForUser :many
SELECT * FROM user_device_tokens
WHERE user_id = $1 
  AND is_active = true 
  AND (expires_at IS NULL OR expires_at > NOW())
ORDER BY first_registered_at DESC;

-- name: UpdateDeviceTokenDetails :one
UPDATE user_device_tokens
SET 
    platform = COALESCE($2, platform),
    device_type = COALESCE($3, device_type),
    device_model = COALESCE($4, device_model),
    os_name = COALESCE($5, os_name),
    os_version = COALESCE($6, os_version),
    app_version = COALESCE($7, app_version),
    client_ip = COALESCE($8, client_ip),
    last_used_at = NOW(),
    is_verified = COALESCE($9, is_verified)
WHERE id = $1
RETURNING *;

-- name: UpdateDeviceTokenLastUsed :one
UPDATE user_device_tokens
SET 
    last_used_at = NOW(),
    client_ip = COALESCE($2, client_ip)
WHERE id = $1
RETURNING *;

-- name: RevokeDeviceToken :one
UPDATE user_device_tokens
SET 
    is_active = false,
    is_revoked = true
WHERE id = $1
RETURNING *;

-- name: DeleteExpiredDeviceTokens :exec
DELETE FROM user_device_tokens
WHERE 
    expires_at < NOW() 
    AND is_active = false;

-- name: CountActiveDeviceTokensForUser :one
SELECT COUNT(*) 
FROM user_device_tokens
WHERE 
    user_id = $1 
    AND is_active = true 
    AND (expires_at IS NULL OR expires_at > NOW());

-- name: UpdateDeviceTokenPushNotificationToken :one
UPDATE user_device_tokens
SET 
    push_notification_token = $2,
    is_verified = true
WHERE id = $1
RETURNING *;

-- name: GetDeviceTokensByPlatform :many
SELECT * FROM user_device_tokens
WHERE 
    user_id = $1 
    AND platform = $2 
    AND is_active = true 
    AND (expires_at IS NULL OR expires_at > NOW())
ORDER BY first_registered_at DESC;

-- name: SearchDeviceTokens :many
SELECT * FROM user_device_tokens
WHERE 
    user_id = $1 
    AND (
        COALESCE(device_type, '') ILIKE '%' || $2 || '%' 
        OR COALESCE(device_model, '') ILIKE '%' || $2 || '%' 
        OR platform ILIKE '%' || $2 || '%'
    )
ORDER BY first_registered_at DESC
LIMIT $3 OFFSET $4;