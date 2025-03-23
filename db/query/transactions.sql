-- name: GetTransactionByID :one
-- Retrieves a single transaction by its ID
SELECT * FROM transactions
WHERE id = $1
LIMIT 1;

-- name: GetTransactionByTxHash :one
-- Retrieves a single transaction by its transaction hash
SELECT * FROM transactions
WHERE tx_hash = $1
LIMIT 1;

-- name: GetTransactionsByUserID :many
-- Retrieves all transactions for a specific user
SELECT * FROM transactions
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: GetTransactionsByStatus :many
-- Retrieves transactions by status
SELECT * FROM transactions
WHERE status = $1
ORDER BY created_at DESC
LIMIT $2
OFFSET $3;

-- name: GetTransactionsByUserIDAndStatus :many
-- Retrieves transactions for a specific user with a specific status
SELECT * FROM transactions
WHERE user_id = $1 AND status = $2
ORDER BY created_at DESC;

-- name: UpdateTransactionStatus :one
-- Updates the status of a transaction and returns the updated transaction
UPDATE transactions
SET
  status = $2,
  updated_at = now()
WHERE id = $1
RETURNING *;

-- name: UpdateTransaction :one
-- Updates transaction details and returns the updated transaction
UPDATE transactions
SET
  status = COALESCE($2, status),
  transaction_pin_hash = COALESCE($3, transaction_pin_hash),
  updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteTransaction :exec
-- Permanently deletes a transaction record
DELETE FROM transactions
WHERE id = $1;

-- name: DeleteTransactionsByUserID :exec
-- Deletes all transactions for a specific user
DELETE FROM transactions
WHERE user_id = $1;