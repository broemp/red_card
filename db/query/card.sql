-- name: CreateCard :one
INSERT INTO "card" (
    author, accused, color, event, description
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetCard :one
SELECT c.id, acc.username as accused_username, acc.name as accused_name, acc.id as accused_id,
aut.username as author_username, aut.name as author_name, aut.id as author_id,
color, event, description, c.created_at 
FROM "card" as c
JOIN "user" as acc on acc.id=c.accused
JOIN "user" as aut on aut.id=c.author
WHERE c.id = $1;

-- name: ListMostRecentCard :many
SELECT c.id, acc.username as accused_username, acc.name as accused_name,
aut.username as author_username, aut.name as author_name,
color, event, c.created_at 
FROM "card" as c
JOIN "user" as acc on acc.id=c.accused
JOIN "user" as aut on aut.id=c.author
ORDER BY c.created_at DESC
LIMIT $1
OFFSET $2;

-- name: ListCardsByUserID :many
SELECT c.id, c.author, c.accused, c.color, c.event, c.created_at , u.username as author_username, u.name as author_name
FROM "card" as c
JOIN "user" as u 
on u.id=c.author
WHERE c.accused=$1;

-- name: GetCardColorCountByUserID :many
SELECT color, count(color)
FROM "card"
WHERE accused=$1
GROUP BY color;
