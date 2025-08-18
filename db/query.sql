-- name: GetRepository :one
SELECT *
FROM repository
WHERE id = $1 LIMIT 1;

-- name: GetGithubRepository :one
SELECT *
FROM repository
WHERE github_id = $1 LIMIT 1;

-- name: SearchRepository :one
SELECT *
FROM repository
WHERE to_tsvector('english', name) @@ plainto_tsquery($1);

-- name: ListRepositories :many
SELECT *
FROM repository
ORDER BY name;

-- name: ListUserRepository :many
SELECT *
FROM repository
WHERE user_id = $1
ORDER BY updated_at DESC;

-- name: CreateRepository :one
INSERT INTO repository (user_id,
                        owner,
                        name,
                        token,
                        is_private,
                        github_id,
                        created_at,
                        updated_at)
VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW()) RETURNING *;

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

-- name: GetLatestRelease :one
SELECT *
FROM release
WHERE repository_id = $1
ORDER BY published_at DESC LIMIT 1;

-- name: GetReleaseWithAssets :one
SELECT b.*,
       a.id                as release_id,
       a.github_id         as release_github_id,
       a.name              as release_name,
       a.tag_name          as release_tag_name,
       a.body              as release_body,
       a.is_draft          as release_is_draft,
       a.is_prerelease     as release_is_prerelease,
       a.created_at        as release_created_at,
       a.published_at      as release_published_at,
       a.author_name       as release_author_name,
       a.author_id         as release_author_id,
       a.author_avatar_url as release_author_avatar_url,
       a.repository_id     as release_repository_id
FROM asset b
         JOIN release a ON a.id = b.release_id
WHERE a.id = $1;

-- name: GetGithubRelease :one
SELECT *
FROM release
WHERE github_id = $1 LIMIT 1;

-- name: GetGithubReleaseWithAssets :one
SELECT b.*,
       a.id                as release_id,
       a.github_id         as release_github_id,
       a.name              as release_name,
       a.tag_name          as release_tag_name,
       a.body              as release_body,
       a.is_draft          as release_is_draft,
       a.is_prerelease     as release_is_prerelease,
       a.created_at        as release_created_at,
       a.published_at      as release_published_at,
       a.author_name       as release_author_name,
       a.author_id         as release_author_id,
       a.author_avatar_url as release_author_avatar_url,
       a.repository_id     as release_repository_id
FROM asset b
         JOIN release a ON a.id = b.release_id
WHERE a.id = $1;

-- name: SearchRelease :one
SELECT *
FROM release
WHERE name = $1 LIMIT 1;

-- name: GetReleaseVersion :one
SELECT *
FROM release
WHERE tag_name = $1
  AND repository_id = $2 LIMIT 1;

-- name: ListReleases :many
SELECT *
FROM release
WHERE repository_id = $1
ORDER BY published_at;

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
                     repository_id,
                     github_id)
VALUES ($1, $2, $3, $4, $5, NOW(), $6, $7, $8, $9, $10, $11) RETURNING *;

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

-- name: GetGithubAsset :one
SELECT *
FROM asset
WHERE github_id = $1 LIMIT 1;

-- name: SearchAsset :one
SELECT *
FROM asset
WHERE name = $1 LIMIT 1;

-- name: ListAssets :many
SELECT *
FROM asset
ORDER BY name;

-- name: ListReleaseAssets :many
SELECT *
FROM asset
WHERE release_id = $1
ORDER BY uploaded_at;

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
                   release_id,
                   github_id)
VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW(), $7, $8, $9) RETURNING *;

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