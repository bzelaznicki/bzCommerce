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

-- name: GetUserAccountById :one
SELECT id, email, full_name, created_at, updated_at, is_admin
FROM users
WHERE id = sqlc.arg(id);

-- name: ListUsers :many
SELECT id, full_name, email, created_at, updated_at, is_admin FROM users;

-- name: UpdateUserById :exec
UPDATE users
SET full_name = sqlc.arg(full_name),
    email = sqlc.arg(email),
    is_admin = sqlc.arg(is_admin),
    updated_at = NOW()
WHERE id = sqlc.arg(id);

-- name: UpdateUserPassword :exec
UPDATE users
SET password_hash = sqlc.arg(password)
WHERE id = sqlc.arg(id);

-- name: DeleteUserById :exec
DELETE FROM users
WHERE id = sqlc.arg(id);

-- name: CountFilteredUsers :one
SELECT COUNT(*)
FROM users
WHERE full_name ILIKE $1 OR email ILIKE $1;