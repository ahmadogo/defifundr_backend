-- name: CreateTransaction :one
INSERT INTO transactions (
  user_id,
  tx_hash,
  amount,
  currency,
  type,
  status
) VALUES (
  $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetTransaction :one
SELECT * FROM transactions
WHERE id = $1 LIMIT 1;

-- name: GetTransactionByHash :one
SELECT * FROM transactions
WHERE tx_hash = $1 LIMIT 1;

-- name: ListUserTransactions :many
SELECT * FROM transactions
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2
OFFSET $3;

-- name: UpdateTransactionStatus :one
UPDATE transactions
SET status = $2
WHERE id = $1
RETURNING *;