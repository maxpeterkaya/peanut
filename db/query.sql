-- name: GetRepository :one
SELECT * FROM repository
WHERE id = $1 LIMIT 1;

-- name: ListRepositories :many
SELECT * FROM repository
ORDER BY name;

-- name: CreateRepository :one
INSERT INTO repository (owner,
                        name,
                        token,
                        is_private,
                        created_at,
                        updated_at)
VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING *;

-- name: UpdateRepository :exec
UPDATE repository
set owner      = $2,
    name       = $3,
    token      = $4,
    is_private = $5
WHERE id = $1;

-- name: DeleteRepository :exec
DELETE
FROM repository
WHERE id = $1;

-- name: GetRelease :one
SELECT *
FROM release
WHERE id = $1 LIMIT 1;

-- name: ListReleases :many
SELECT *
FROM release
ORDER BY name;

-- name: CreateRelease :one
INSERT INTO release (name,
                     tag_name,
                     body,
                     is_draft,
                     is_prerelease,
                     created_at,
                     published_at,
                     author_name,
                     author_id,
                     author_avatar_url,
                     repository_id)
VALUES ($1, $2, $3, $4, $5, NOW(), NOW(), $6, $7, $8, $9) RETURNING *;

-- name: UpdateRelease :exec
UPDATE release
set name          = $2,
    tag_name      = $3,
    body          = $4,
    is_draft      = $5,
    is_prerelease = $6
WHERE id = $1;

-- name: DeleteRelease :exec
DELETE
FROM release
WHERE id = $1;

-- name: GetAsset :one
SELECT *
FROM asset
WHERE id = $1 LIMIT 1;

-- name: ListAssets :many
SELECT *
FROM asset
ORDER BY name;

-- name: CreateAsset :one
INSERT INTO asset (api_url,
                   url,
                   name,
                   content_length,
                   download_count,
                   view_count,
                   created_at,
                   updated_at,
                   uploaded_at,
                   release_id)
VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW(), NOW(), $7) RETURNING *;

-- name: UpdateAsset :exec
UPDATE asset
set api_url        = $2,
    url            = $3,
    name           = $4,
    content_length = $5,
    download_count = $6
WHERE id = $1;

-- name: DeleteAsset :exec
DELETE
FROM asset
WHERE id = $1;

-- name: GetUser :one
SELECT *
FROM "user"
WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT *
FROM "user"
ORDER BY username;

-- name: CreateUser :one
INSERT INTO "user" (username,
                    display_name,
                    pass_hash,
                    created_at,
                    updated_at)
VALUES ($1, $2, $3, NOW(), NOW()) RETURNING *;

-- name: UpdateUser :exec
UPDATE "user"
set username     = $2,
    display_name = $3,
    pass_hash    = $4
WHERE id = $1;

-- name: DeleteUser :exec
DELETE
FROM "user"
WHERE id = $1;