-- name: GetAllActiveDonations :many

SELECT *
FROM donations
WHERE deadline > now()
ORDER BY created_at DESC
LIMIT 10
OFFSET 0;