-- name: AddRefreshToken :one
INSERT INTO refresh_tokens (user_id, token, expires_at)
VALUES (sqlc.arg(user_id), sqlc.arg(token), sqlc.arg(expires_at))
RETURNING *;

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens WHERE token = sqlc.arg(token);

-- name: InvalidateRefreshToken :exec
UPDATE refresh_tokens SET revoked_at = NOW() WHERE token = sqlc.arg(token);