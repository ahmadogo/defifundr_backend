-- name: AddPayrollEmployee :one
INSERT INTO payroll_employees (
  payroll_id,
  employee_id
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetPayrollEmployee :one
SELECT * FROM payroll_employees
WHERE payroll_id = $1 AND employee_id = $2
LIMIT 1;

-- name: ListPayrollEmployees :many
SELECT pe.*, u.first_name, u.last_name, u.email
FROM payroll_employees pe
JOIN users u ON pe.employee_id = u.id
WHERE pe.payroll_id = $1
ORDER BY pe.created_at DESC;

-- name: ListEmployeePayrolls :many
SELECT p.*
FROM payrolls p
JOIN payroll_employees pe ON p.id = pe.payroll_id
WHERE pe.employee_id = $1
ORDER BY p.created_at DESC;

-- name: RemovePayrollEmployee :exec
DELETE FROM payroll_employees
WHERE payroll_id = $1 AND employee_id = $2;