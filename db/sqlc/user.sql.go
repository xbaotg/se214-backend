// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: user.sql

package sqlc

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
    user_id, username, user_fullname, password, user_role, year, created_at, updated_at, user_email
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING user_id, username, password, user_email, user_fullname, user_role, year, created_at, updated_at
`

type CreateUserParams struct {
	UserID       uuid.UUID
	Username     string
	UserFullname string
	Password     string
	UserRole     Role
	Year         int32
	CreatedAt    time.Time
	UpdatedAt    time.Time
	UserEmail    string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.UserID,
		arg.Username,
		arg.UserFullname,
		arg.Password,
		arg.UserRole,
		arg.Year,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.UserEmail,
	)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.Username,
		&i.Password,
		&i.UserEmail,
		&i.UserFullname,
		&i.UserRole,
		&i.Year,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserByUsername = `-- name: GetUserByUsername :one
SELECT user_id, username, password, user_email, user_fullname, user_role, year, created_at, updated_at FROM users WHERE username = $1
`

func (q *Queries) GetUserByUsername(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByUsername, username)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.Username,
		&i.Password,
		&i.UserEmail,
		&i.UserFullname,
		&i.UserRole,
		&i.Year,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
