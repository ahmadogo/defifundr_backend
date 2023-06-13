-- name: CreateWallet :one

INSERT INTO
    wallet (
        owner,
        balance,
        address,
        pem_encoded,
        pem_encoded_pub
    )
VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetWallet :one

SELECT * FROM wallet WHERE owner = $1 LIMIT 1;

-- name: UpdateWallet :one

UPDATE wallet SET balance = $2 WHERE owner = $1 RETURNING *;

-- name: AddAccountBalance :one

UPDATE wallet
SET
    balance = balance + sqlc.arg(amount)
WHERE
    id = sqlc.arg(id) RETURNING *;

-- name: DeleteAccount :exec

DELETE FROM wallet WHERE id = $1;