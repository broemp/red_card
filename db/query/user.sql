-- name: CreateUser :one
INSERT INTO "user" (
    username, hashed_password, name
) VALUES (
    $1, $2, $3 
) RETURNING *;

-- name: ListUser :many
SELECT id, username
FROM "user"
LIMIT $1
OFFSET $2;

-- name: ListUserFilter :many
SELECT id, username, name
FROM "user"
WHERE username LIKE CONCAT(sqlc.arg(filter)::text, '%')
OR name LIKE CONCAT(sqlc.arg(filter)::text, '%')
LIMIT $1
OFFSET $2;

-- name: GetUserAuth :one
SELECT * FROM "user"
WHERE username = $1
LIMIT 1;

-- name: GetUserByID :one
SELECT id, username, name FROM "user"
WHERE id = $1
LIMIT 1;

-- name: DeleteUser :exec
DELETE FROM "user"
WHERE id $1;

-- name: UpdatePassword :exec
UPDATE "user"
SET hashed_password = $1
WHERE id = $2;

-- name: GetUserID :one
SELECT id FROM "user"
WHERE username = $1
LIMIT 1;

