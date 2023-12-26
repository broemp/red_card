-- name: CreateUser :one
INSERT INTO "user" (
    username, hashed_password, first_name, last_name
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: ListUser :many
SELECT username
FROM "user"
LIMIT $1
OFFSET $2;

-- name: GetUser :one
SELECT * FROM "user"
WHERE username = $1 LIMIT 1;

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

