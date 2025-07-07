-- name: CreateUser :one
INSERT INTO users (email, full_name, password_hash)
VALUES ($1, $2, $3)
RETURNING id, email, full_name, created_at, updated_at, is_admin, is_active, disabled_at;


-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserAccountById :one
SELECT id, email, full_name, created_at, updated_at, is_admin, is_active, disabled_at
FROM users
WHERE id = sqlc.arg(id);

-- name: ListUsers :many
SELECT
  id,
  full_name,
  email,
  created_at,
  updated_at,
  is_admin,
  is_active,
  disabled_at
FROM users;


-- name: UpdateUserById :one
UPDATE users
SET
  full_name = sqlc.arg(full_name),
  email = sqlc.arg(email),
  is_admin = sqlc.arg(is_admin),
  updated_at = NOW()
WHERE id = sqlc.arg(id)
RETURNING
  id,
  full_name,
  email,
  created_at,
  updated_at,
  is_admin,
  is_active,
  disabled_at;

-- name: UpdateUserPassword :execrows
UPDATE users
SET password_hash = sqlc.arg(password)
WHERE id = sqlc.arg(id);

-- name: DeleteUserById :execrows
DELETE FROM users
WHERE id = sqlc.arg(id);

-- name: CountFilteredUsers :one
SELECT COUNT(*)
FROM users
WHERE full_name ILIKE $1 OR email ILIKE $1;



-- name: DisableUser :one
UPDATE users
SET
  is_active = FALSE,
  disabled_at = CASE
    WHEN disabled_at IS NULL THEN CURRENT_TIMESTAMP
    ELSE disabled_at
  END,
  updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, full_name, email, created_at, updated_at, is_admin, is_active, disabled_at;


-- name: EnableUser :one
UPDATE users
SET is_active = TRUE,
    disabled_at = NULL,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, full_name, email, created_at, is_admin, updated_at, is_active;


-- name: CountFilteredUsersWithStatus :one
SELECT COUNT(*)
FROM users
WHERE
  (
    email ILIKE $1
    OR full_name ILIKE $1
  )
  AND (
    CASE
      WHEN sqlc.arg(status) = 'active' THEN is_active = TRUE
      WHEN sqlc.arg(status) = 'disabled' THEN is_active = FALSE
      ELSE TRUE
    END
  );
