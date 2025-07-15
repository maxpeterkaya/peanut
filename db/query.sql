-- name: GetRepo :one
SELECT * FROM repository
WHERE id = $1 LIMIT 1;

-- name: ListRepositories :many
SELECT * FROM repository
ORDER BY name;
