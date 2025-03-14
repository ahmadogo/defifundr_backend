-- name: CreateOrganization :one
INSERT INTO organizations (
  name,
  employer_id
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetOrganization :one
SELECT * FROM organizations
WHERE id = $1 LIMIT 1;

-- name: ListEmployerOrganizations :many
SELECT * FROM organizations
WHERE employer_id = $1
ORDER BY created_at DESC;

-- name: UpdateOrganization :one
UPDATE organizations
SET 
  name = $2,
  updated_at = now()
WHERE id = $1 AND employer_id = $3
RETURNING *;

-- name: DeleteOrganization :exec
DELETE FROM organizations
WHERE id = $1 AND employer_id = $2;