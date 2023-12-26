// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: user.sql

package db

import (
	"context"
	"database/sql"
)

const createUser = `-- name: CreateUser :one
INSERT INTO "user" (
    username, hashed_password, first_name, last_name
) VALUES (
    $1, $2, $3, $4
) RETURNING id, username, first_name, last_name, hashed_password, password_changed_at, deleted_at, created_at
`

type CreateUserParams struct {
	Username       string         `json:"username"`
	HashedPassword string         `json:"hashed_password"`
	FirstName      sql.NullString `json:"first_name"`
	LastName       sql.NullString `json:"last_name"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.Username,
		arg.HashedPassword,
		arg.FirstName,
		arg.LastName,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.FirstName,
		&i.LastName,
		&i.HashedPassword,
		&i.PasswordChangedAt,
		&i.DeletedAt,
		&i.CreatedAt,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM "user"
WHERE id $1
`

func (q *Queries) DeleteUser(ctx context.Context, dollar_1 interface{}) error {
	_, err := q.db.ExecContext(ctx, deleteUser, dollar_1)
	return err
}

const getUser = `-- name: GetUser :one
SELECT id, username, first_name, last_name, hashed_password, password_changed_at, deleted_at, created_at FROM "user"
WHERE username = $1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.FirstName,
		&i.LastName,
		&i.HashedPassword,
		&i.PasswordChangedAt,
		&i.DeletedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getUserID = `-- name: GetUserID :one
SELECT id FROM "user"
WHERE username = $1
LIMIT 1
`

func (q *Queries) GetUserID(ctx context.Context, username string) (int64, error) {
	row := q.db.QueryRowContext(ctx, getUserID, username)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const listUser = `-- name: ListUser :many
SELECT username
FROM "user"
LIMIT $1
OFFSET $2
`

type ListUserParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListUser(ctx context.Context, arg ListUserParams) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, listUser, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []string{}
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			return nil, err
		}
		items = append(items, username)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updatePassword = `-- name: UpdatePassword :exec
UPDATE "user"
SET hashed_password = $1
WHERE id = $2
`

type UpdatePasswordParams struct {
	HashedPassword string `json:"hashed_password"`
	ID             int64  `json:"id"`
}

func (q *Queries) UpdatePassword(ctx context.Context, arg UpdatePasswordParams) error {
	_, err := q.db.ExecContext(ctx, updatePassword, arg.HashedPassword, arg.ID)
	return err
}
