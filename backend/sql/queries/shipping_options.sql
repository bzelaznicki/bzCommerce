-- name: CreateShippingOption :one
INSERT INTO shipping_options (name, description, price, estimated_days, is_active)
VALUES (sqlc.arg(name), sqlc.arg(description), sqlc.arg(price), sqlc.arg(estimated_days), sqlc.arg(is_active))
RETURNING *;

-- name: GetShippingOptions :many
SELECT * FROM shipping_options
ORDER BY sort_order ASC;

-- name: GetActiveShippingOptions :many
SELECT * FROM shipping_options
WHERE is_active = true
ORDER BY sort_order ASC;

-- name: SelectShippingOptionById :one
SELECT * FROM shipping_options WHERE id = sqlc.arg(id);

-- name: UpdateShippingOption :one
UPDATE shipping_options
SET name = sqlc.arg(name),
    description = sqlc.arg(description),
    price = sqlc.arg(price),
    estimated_days = sqlc.arg(estimated_days),
    sort_order = sqlc.arg(sort_order),
    is_active = sqlc.arg(is_active)
WHERE id = sqlc.arg(id)
RETURNING *;
-- name: DeleteShippingOption :execrows
DELETE FROM shipping_options
WHERE id = sqlc.arg(id);

-- name: ToggleShippingOptionStatus :one
UPDATE shipping_options
SET is_active = sqlc.arg(is_active)
WHERE id = sqlc.arg(id)
RETURNING id, is_active;