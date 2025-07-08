-- name: GetPaymentOptions :many
SELECT * FROM payment_options
ORDER BY sort_order ASC;

-- name: GetActivePaymentOptions :many
SELECT * FROM payment_options
WHERE is_active = true
ORDER BY sort_order ASC;

-- name: CreatePaymentOption :one
INSERT INTO payment_options (name, description, is_active, sort_order)
VALUES(
    sqlc.arg(name),
    sqlc.arg(description),
    sqlc.arg(is_active),
    sqlc.arg(sort_order)
)
RETURNING *;

-- name: GetPaymentOptionById :one
SELECT * FROM payment_options WHERE id = sqlc.arg(id);

-- name: UpdatePaymentOption :exec
UPDATE payment_options
SET name = sqlc.arg(name),
description = sqlc.arg(description),
sort_order = sqlc.arg(sort_order),
is_active = sqlc.arg(is_active)
where id = sqlc.arg(id);

-- name: DeletePaymentOption :exec
DELETE FROM payment_options
WHERE id = sqlc.arg(id);