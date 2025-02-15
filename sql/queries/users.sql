-- name: CreateUser :one
INSERT INTO users (email, full_name, password_hash)
VALUES ($1, $2, $3)
RETURNING *;