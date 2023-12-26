-- name: CreateEvent :one 
INSERT INTO "event" (
    name, date
) VALUES (
    $1, $2
) RETURNING *;
