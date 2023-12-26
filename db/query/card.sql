-- name: CreateCard :one
INSERT INTO "card" (
    author, accused, color, event
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetCard :one
SELECT * FROM "card"
WHERE id = $1;

-- name: ListMostRecentCard :many
SELECT * FROM "card"
ORDER BY created_at DESC
LIMIT $1
OFFSET $2;

-- name: ListCardsFromUserByUsername :many
SELECT c.id, c.author, c.accused, c.color, c.event , u.username 
FROM "card" as c
JOIN "user" as u 
on u.id=c.author
WHERE username=$1;
