-- name: CreateInvoice :one
INSERT INTO invoices (
  freelancer_id,
  employer_id,
  amount,
  currency,
  status
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetInvoice :one
SELECT * FROM invoices
WHERE id = $1 LIMIT 1;

-- name: ListFreelancerInvoices :many
SELECT i.*, u.first_name as employer_first_name, u.last_name as employer_last_name
FROM invoices i
JOIN users u ON i.employer_id = u.id
WHERE i.freelancer_id = $1
ORDER BY i.created_at DESC;

-- name: ListEmployerInvoices :many
SELECT i.*, u.first_name as freelancer_first_name, u.last_name as freelancer_last_name
FROM invoices i
JOIN users u ON i.freelancer_id = u.id
WHERE i.employer_id = $1
ORDER BY i.created_at DESC;

-- name: UpdateInvoiceStatus :one
UPDATE invoices
SET 
  status = $2,
  updated_at = now()
WHERE id = $1
RETURNING *;

-- name: UpdateInvoiceContractAddress :one
UPDATE invoices
SET 
  contract_address = $2,
  updated_at = now()
WHERE id = $1
RETURNING *;