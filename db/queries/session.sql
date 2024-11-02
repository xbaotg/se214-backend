-- name: CreateSession :one
INSERT INTO sessions (
    session_id, user_id, refresh_token, created_at, updated_at, expires_in, is_active
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetSessionByUserId :one
SELECT * FROM sessions WHERE user_id = $1;

-- name: GetSessionBySessionId :one
SELECT * FROM sessions WHERE session_id = $1;

-- name: UpdateSession :exec
UPDATE sessions SET
    refresh_token = $2,
    updated_at = $3
WHERE user_id = $1;

-- name: RevolveSession :exec
UPDATE sessions SET
    updated_at = $2,
    is_active = $3
WHERE session_id = $1;

-- name: DeleteSession :exec
DELETE FROM sessions WHERE user_id = $1;

-- name: DeleteSessionBySessionId :exec
DELETE FROM sessions WHERE session_id = $1;