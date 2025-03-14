-- name: AddOrganizationMember :one
INSERT INTO organization_members (
  organization_id,
  employee_id,
  role
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetOrganizationMember :one
SELECT * FROM organization_members
WHERE organization_id = $1 AND employee_id = $2
LIMIT 1;

-- name: ListOrganizationMembers :many
SELECT om.*, u.first_name, u.last_name, u.email
FROM organization_members om
JOIN users u ON om.employee_id = u.id
WHERE om.organization_id = $1
ORDER BY om.created_at DESC;

-- name: ListUserOrganizationMemberships :many
SELECT om.*, o.name as organization_name
FROM organization_members om
JOIN organizations o ON om.organization_id = o.id
WHERE om.employee_id = $1
ORDER BY om.created_at DESC;

-- name: UpdateOrganizationMemberRole :one
UPDATE organization_members
SET role = $3
WHERE organization_id = $1 AND employee_id = $2
RETURNING *;

-- name: RemoveOrganizationMember :exec
DELETE FROM organization_members
WHERE organization_id = $1 AND employee_id = $2;