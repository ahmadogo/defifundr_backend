-- name: CreateWallet :one
INSERT INTO wallets (
  user_id,
  wallet_address,
  chain,
  is_primary,
  pin_hash
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetWallet :one
SELECT * FROM wallets
WHERE id = $1 LIMIT 1;

-- name: GetWalletByAddress :one
SELECT * FROM wallets
WHERE wallet_address = $1 LIMIT 1;

-- name: ListUserWallets :many
SELECT * FROM wallets
WHERE user_id = $1
ORDER BY is_primary DESC, created_at DESC;

-- name: GetPrimaryWallet :one
SELECT * FROM wallets
WHERE user_id = $1 AND is_primary = true
LIMIT 1;

-- name: SetPrimaryWallet :exec
UPDATE wallets
SET is_primary = false
WHERE user_id = $1;

-- name: UpdateWalletPrimary :one
UPDATE wallets
SET is_primary = true
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: DeleteWallet :exec
DELETE FROM wallets
WHERE id = $1 AND user_id = $2;