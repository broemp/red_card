-- name: CreateComment :one
INSERT INTO "comment" (
  message, author, card
) VALUES ( 
 $1, $2, $3
 ) RETURNING *;
