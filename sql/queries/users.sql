-- name: CreateUser :one
INSERT INTO users (email, full_name, password_hash)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1;

-- name: ListUsers :many
SELECT id, full_name, email, created_at, updated_at, is_admin FROM users;