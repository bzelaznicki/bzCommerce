-- name: CreateShippingOption :one
INSERT INTO shipping_options (name, description, price, estimated_days, sort_order, is_active)
VALUES (sqlc.arg(name), sqlc.arg(description), sqlc.arg(price), sqlc.arg(estimated_days), sqlc.arg(sort_order), sqlc.arg(is_active))
RETURNING *;

-- name: GetShippingOptions :many
SELECT * FROM shipping_options
ORDER BY sort_order ASC;

-- name: SelectShippingOptionById :one
SELECT * FROM shipping_options WHERE id = sqlc.arg(id);