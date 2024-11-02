-- name: CreateUser :one
INSERT INTO users (
    user_id, username, user_fullname, password, user_role, year, created_at, updated_at, user_email
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1;

-- name: GetUserById :one
SELECT * FROM users WHERE user_id = $1;

-- name: ValidateNewUser :one
SELECT * FROM users WHERE username = $1 OR user_email = $2;

-- name: UpdatPassword :one
UPDATE users SET password = $1, updated_at = $2 WHERE user_id = $3 RETURNING *;