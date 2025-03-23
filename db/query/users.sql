-- name: CreateUser :one
-- Creates a new user record and returns the created user
INSERT INTO users (
  id,
  email,
  password_hash,
  profile_picture,
  account_type,
  gender,
  personal_account_type,
  first_name,
  last_name,
  nationality,
  residential_country,
  job_role,
  company_website,
  employment_type,
  created_at,
  updated_at
) VALUES (
  COALESCE($1, uuid_generate_v4()),
  $2,
  $3,
 COALESCE($4, ''),
  $5,
  $6,
  $7,
  $8,
  $9,
  $10,
  $11,
  $12,
  $13,
  $14,
  COALESCE($15, now()),
  COALESCE($16, now())
) RETURNING *;


-- name: GetUser :one
SELECT * FROM users WHERE id = $1 OR id::text = $1 LIMIT 1;


-- name: GetUserByEmail :one
-- Retrieves a single user by their email address
SELECT * FROM users
WHERE email = $1
LIMIT 1;

-- name: CheckEmailExists :one
SELECT EXISTS (
        SELECT 1
        FROM users
        WHERE
            email = $1
        LIMIT 1
    ) AS exists;

-- name: ListUsers :many
-- Lists users with pagination support
SELECT *
FROM users
ORDER BY 
  CASE WHEN $3::text = 'ASC' THEN created_at END ASC,
  CASE WHEN $3::text = 'DESC' OR $3::text IS NULL THEN created_at END DESC
LIMIT $1
OFFSET $2;

-- name: ListUsersByAccountType :many
-- Lists users filtered by account type with pagination
SELECT *
FROM users
WHERE account_type = $3
ORDER BY 
  CASE WHEN $4::text = 'ASC' THEN created_at END ASC,
  CASE WHEN $4::text = 'DESC' OR $4::text IS NULL THEN created_at END DESC
LIMIT $1
OFFSET $2;

-- name: CountUsers :one
-- Counts the total number of users (useful for pagination)
SELECT COUNT(*) FROM users;

-- name: CountUsersByAccountType :one
-- Counts users filtered by account type
SELECT COUNT(*) FROM users
WHERE account_type = $1;

-- name: UpdateUser :one
-- Updates user details and returns the updated user
UPDATE users
SET
  email = COALESCE($2, email),
  profile_picture = $3,
  account_type = COALESCE($4, account_type),
  gender = $5,
  personal_account_type = COALESCE($6, personal_account_type),
  first_name = COALESCE($7, first_name),
  last_name = COALESCE($8, last_name),
  nationality = COALESCE($9, nationality),
  residential_country = $10,
  job_role = $11,
  company_website = $12,
  employment_type = $13,
  updated_at = now()
WHERE id = $1
RETURNING *;

-- name: UpdateUserPassword :exec
-- Updates a user's password
UPDATE users
SET
  password_hash = $2,
  updated_at = now()
WHERE id = $1;

-- name: UpdateUserEmail :one
-- Updates a user's email address with validation that the new email is unique
UPDATE users
SET
  email = $2,
  updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
-- Permanently deletes a user record
DELETE FROM users
WHERE id = $1;

-- name: SearchUsers :many
-- Searches for users by name, email, or nationality with pagination
SELECT *
FROM users
WHERE 
  (
    first_name ILIKE '%' || $3 || '%' OR
    last_name ILIKE '%' || $3 || '%' OR
    email ILIKE '%' || $3 || '%' OR
    nationality ILIKE '%' || $3 || '%'
  )
ORDER BY
  CASE WHEN $4::text = 'ASC' THEN created_at END ASC,
  CASE WHEN $4::text = 'DESC' OR $4::text IS NULL THEN created_at END DESC
LIMIT $1
OFFSET $2;

-- name: CountSearchUsers :one
-- Counts the number of users matching a search query
SELECT COUNT(*)
FROM users
WHERE 
  (
    first_name ILIKE '%' || $1 || '%' OR
    last_name ILIKE '%' || $1 || '%' OR
    email ILIKE '%' || $1 || '%' OR
    nationality ILIKE '%' || $1 || '%'
  );