// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: card.sql

package db

import (
	"context"
	"database/sql"
	"time"
)

const createCard = `-- name: CreateCard :one
INSERT INTO "card" (
    author, accused, color, event, description
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING id, author, accused, color, description, event, created_at
`

type CreateCardParams struct {
	Author      int64          `json:"author"`
	Accused     int64          `json:"accused"`
	Color       Color          `json:"color"`
	Event       sql.NullInt64  `json:"event"`
	Description sql.NullString `json:"description"`
}

func (q *Queries) CreateCard(ctx context.Context, arg CreateCardParams) (Card, error) {
	row := q.db.QueryRowContext(ctx, createCard,
		arg.Author,
		arg.Accused,
		arg.Color,
		arg.Event,
		arg.Description,
	)
	var i Card
	err := row.Scan(
		&i.ID,
		&i.Author,
		&i.Accused,
		&i.Color,
		&i.Description,
		&i.Event,
		&i.CreatedAt,
	)
	return i, err
}

const getCard = `-- name: GetCard :one
SELECT c.id, acc.username as accused_username, acc.name as accused_name, acc.id as accused_id,
aut.username as author_username, aut.name as author_name, aut.id as author_id,
color, event, description, c.created_at 
FROM "card" as c
JOIN "user" as acc on acc.id=c.accused
JOIN "user" as aut on aut.id=c.author
WHERE c.id = $1
`

type GetCardRow struct {
	ID              int64          `json:"id"`
	AccusedUsername string         `json:"accused_username"`
	AccusedName     string         `json:"accused_name"`
	AccusedID       int64          `json:"accused_id"`
	AuthorUsername  string         `json:"author_username"`
	AuthorName      string         `json:"author_name"`
	AuthorID        int64          `json:"author_id"`
	Color           Color          `json:"color"`
	Event           sql.NullInt64  `json:"event"`
	Description     sql.NullString `json:"description"`
	CreatedAt       time.Time      `json:"created_at"`
}

func (q *Queries) GetCard(ctx context.Context, id int64) (GetCardRow, error) {
	row := q.db.QueryRowContext(ctx, getCard, id)
	var i GetCardRow
	err := row.Scan(
		&i.ID,
		&i.AccusedUsername,
		&i.AccusedName,
		&i.AccusedID,
		&i.AuthorUsername,
		&i.AuthorName,
		&i.AuthorID,
		&i.Color,
		&i.Event,
		&i.Description,
		&i.CreatedAt,
	)
	return i, err
}

const getCardColorCountByUserID = `-- name: GetCardColorCountByUserID :many
SELECT color, count(color)
FROM "card"
WHERE accused=$1
GROUP BY color
`

type GetCardColorCountByUserIDRow struct {
	Color Color `json:"color"`
	Count int64 `json:"count"`
}

func (q *Queries) GetCardColorCountByUserID(ctx context.Context, accused int64) ([]GetCardColorCountByUserIDRow, error) {
	rows, err := q.db.QueryContext(ctx, getCardColorCountByUserID, accused)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetCardColorCountByUserIDRow{}
	for rows.Next() {
		var i GetCardColorCountByUserIDRow
		if err := rows.Scan(&i.Color, &i.Count); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listCardsByUserID = `-- name: ListCardsByUserID :many
SELECT c.id, c.author, c.accused, c.color, c.event, c.created_at , u.username as author_username, u.name as author_name
FROM "card" as c
JOIN "user" as u 
on u.id=c.author
WHERE c.accused=$1
`

type ListCardsByUserIDRow struct {
	ID             int64         `json:"id"`
	Author         int64         `json:"author"`
	Accused        int64         `json:"accused"`
	Color          Color         `json:"color"`
	Event          sql.NullInt64 `json:"event"`
	CreatedAt      time.Time     `json:"created_at"`
	AuthorUsername string        `json:"author_username"`
	AuthorName     string        `json:"author_name"`
}

func (q *Queries) ListCardsByUserID(ctx context.Context, accused int64) ([]ListCardsByUserIDRow, error) {
	rows, err := q.db.QueryContext(ctx, listCardsByUserID, accused)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListCardsByUserIDRow{}
	for rows.Next() {
		var i ListCardsByUserIDRow
		if err := rows.Scan(
			&i.ID,
			&i.Author,
			&i.Accused,
			&i.Color,
			&i.Event,
			&i.CreatedAt,
			&i.AuthorUsername,
			&i.AuthorName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listMostRecentCard = `-- name: ListMostRecentCard :many
SELECT c.id, acc.username as accused_username, acc.name as accused_name,
aut.username as author_username, aut.name as author_name,
color, event, c.created_at 
FROM "card" as c
JOIN "user" as acc on acc.id=c.accused
JOIN "user" as aut on aut.id=c.author
ORDER BY c.created_at DESC
LIMIT $1
OFFSET $2
`

type ListMostRecentCardParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type ListMostRecentCardRow struct {
	ID              int64         `json:"id"`
	AccusedUsername string        `json:"accused_username"`
	AccusedName     string        `json:"accused_name"`
	AuthorUsername  string        `json:"author_username"`
	AuthorName      string        `json:"author_name"`
	Color           Color         `json:"color"`
	Event           sql.NullInt64 `json:"event"`
	CreatedAt       time.Time     `json:"created_at"`
}

func (q *Queries) ListMostRecentCard(ctx context.Context, arg ListMostRecentCardParams) ([]ListMostRecentCardRow, error) {
	rows, err := q.db.QueryContext(ctx, listMostRecentCard, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListMostRecentCardRow{}
	for rows.Next() {
		var i ListMostRecentCardRow
		if err := rows.Scan(
			&i.ID,
			&i.AccusedUsername,
			&i.AccusedName,
			&i.AuthorUsername,
			&i.AuthorName,
			&i.Color,
			&i.Event,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
