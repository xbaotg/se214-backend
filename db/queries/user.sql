-- name: CreateUser :one
INSERT INTO users (
    user_id, username, user_fullname, password, user_role, year, created_at, updated_at, user_email
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1;