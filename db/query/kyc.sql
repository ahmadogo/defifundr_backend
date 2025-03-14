-- name: CreateKYC :one
INSERT INTO kyc (
  id,
  user_id,
  face_verification,
  identity_verification
) VALUES (
  uuid_generate_v4(), $1, $2, $3
) RETURNING *;

-- name: GetKYC :one
SELECT * FROM kyc
WHERE user_id = $1 LIMIT 1;

-- name: UpdateKYC :one
UPDATE kyc
SET 
  face_verification = $2,
  identity_verification = $3,
  updated_at = now()
WHERE user_id = $1
RETURNING *;