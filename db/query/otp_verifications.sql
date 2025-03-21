-- name: CreateOTPVerification :one
INSERT INTO otp_verifications (
    user_id,
    otp_code,
    hashed_otp,
    purpose,
    contact_method,
    attempts_made,
    max_attempts,
    expires_at,
    ip_address,
    user_agent,
    device_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
) RETURNING *;

-- name: GetOTPVerificationByID :one
SELECT * FROM otp_verifications
WHERE id = $1;

-- name: GetOTPVerificationByUserAndPurpose :one
SELECT * FROM otp_verifications
WHERE user_id = $1 
  AND purpose = $2 
  AND is_verified = false 
  AND expires_at > NOW()
ORDER BY created_at DESC
LIMIT 1;

-- name: UpdateOTPAttempts :one
UPDATE otp_verifications
SET attempts_made = attempts_made + 1
WHERE id = $1
RETURNING *;

-- name: VerifyOTP :one
UPDATE otp_verifications
SET 
    is_verified = true,
    verified_at = NOW()
WHERE id = $1
  AND otp_code = $2
  AND expires_at > NOW()
  AND attempts_made <= max_attempts
RETURNING *;

-- name: InValidateOTP :exec
UPDATE otp_verifications
SET is_verified = false
WHERE id = $1;

-- name: DeleteExpiredOTPs :exec
DELETE FROM otp_verifications
WHERE expires_at < NOW();

-- name: CountActiveOTPsForUser :one
SELECT COUNT(*) 
FROM otp_verifications
WHERE user_id = $1 
  AND purpose = $2 
  AND is_verified = false 
  AND expires_at > NOW();

-- name: GetUnverifiedOTPsForUser :many
SELECT * FROM otp_verifications
WHERE user_id = $1 
  AND is_verified = false 
  AND expires_at > NOW()
ORDER BY created_at DESC;