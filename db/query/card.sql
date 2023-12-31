-- name: CreateCard :one
INSERT INTO "card" (
    author, accused, color, event
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetCard :one
SELECT c.id, acc.username as accused_username, acc.first_name as accused_first_name, acc.last_name as accused_last_name,
aut.username as author_username, aut.first_name as author_first_name, aut.last_name as author_last_name,
color, event, c.created_at 
FROM "card" as c
JOIN "user" as acc on acc.id=c.accused
JOIN "user" as aut on aut.id=c.author
WHERE c.id = $1;

-- name: ListMostRecentCard :many
SELECT c.id, acc.username as accused_username, acc.first_name as accused_first_name, acc.last_name as accused_last_name,
aut.username as author_username, aut.first_name as author_first_name, aut.last_name as author_last_name,
color, event, c.created_at 
FROM "card" as c
JOIN "user" as acc on acc.id=c.accused
JOIN "user" as aut on aut.id=c.author
ORDER BY c.created_at DESC
LIMIT $1
OFFSET $2;

-- name: ListCardsByUserID :many
SELECT c.id, c.author, c.accused, c.color, c.event, c.created_at , u.username as author_username
FROM "card" as c
JOIN "user" as u 
on u.id=c.author
WHERE c.accused=$1;
