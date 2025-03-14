-- name: CreateUser :one
INSERT INTO users (
  email, 
  password_hash,
  account_type,
  personal_account_type,
  first_name,
  last_name,
  nationality,
  residential_country,
  job_role,
  company_website,
  employment_type
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at
LIMIT $1
OFFSET $2;

-- name: UpdateUser :one
UPDATE users
SET 
  first_name = COALESCE($2, first_name),
  last_name = COALESCE($3, last_name),
  nationality = COALESCE($4, nationality),
  residential_country = COALESCE($5, residential_country),
  job_role = COALESCE($6, job_role),
  company_website = COALESCE($7, company_website),
  employment_type = COALESCE($8, employment_type),
  updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;