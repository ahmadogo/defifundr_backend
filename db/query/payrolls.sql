-- name: CreatePayroll :one
INSERT INTO payrolls (
  employer_id,
  organization_id,
  payment_frequency,
  salary_amount,
  currency,
  status
) VALUES (
  $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetPayroll :one
SELECT * FROM payrolls
WHERE id = $1 LIMIT 1;

-- name: ListEmployerPayrolls :many
SELECT * FROM payrolls
WHERE employer_id = $1
ORDER BY created_at DESC;

-- name: ListOrganizationPayrolls :many
SELECT * FROM payrolls
WHERE organization_id = $1
ORDER BY created_at DESC;

-- name: UpdatePayrollStatus :one
UPDATE payrolls
SET 
  status = $2,
  updated_at = now()
WHERE id = $1
RETURNING *;

-- name: UpdatePayrollContractAddress :one
UPDATE payrolls
SET 
  contract_address = $2,
  updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeletePayroll :exec
DELETE FROM payrolls
WHERE id = $1 AND employer_id = $2 AND status = 'pending';