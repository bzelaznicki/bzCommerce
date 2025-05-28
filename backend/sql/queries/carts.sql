-- name: CreateCart :one
INSERT INTO carts (user_id)
VALUES (sqlc.arg(user_id))
RETURNING id, user_id, created_at, updated_at;

-- name: GetCartById :one
SELECT * FROM carts
WHERE id = sqlc.arg(id);

-- name: UpdateCartOwner :one
UPDATE carts
SET user_id = sqlc.arg(user_id)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: GetCartsByUserId :many
SELECT * FROM carts
WHERE user_id = sqlc.arg(user_id)
ORDER BY created_at DESC;

-- name: DeleteCart :exec
DELETE FROM carts
WHERE id = sqlc.arg(id);

-- name: AddVariantToCart :one
INSERT INTO carts_variants (cart_id, product_variant_id, quantity, price_per_item)
VALUES (sqlc.arg(cart_id), sqlc.arg(product_variant_id), sqlc.arg(quantity), sqlc.arg(price_per_item))
RETURNING *;

-- name: GetCartProductsByCartId :many
SELECT * FROM carts_variants
WHERE cart_id = sqlc.arg(cart_id);

-- name: UpdateCartVariantQuantity :one
UPDATE carts_variants
SET quantity = sqlc.arg(quantity)
WHERE cart_id = sqlc.arg(cart_id)
  AND product_variant_id = sqlc.arg(product_variant_id)
RETURNING *;

-- name: MarkCartsAsAbandoned :exec
UPDATE carts
SET status = 'abandoned', updated_at = NOW()
WHERE status = 'new' AND updated_at < sqlc.arg(threshold);


-- name: DeleteCartVariant :exec
DELETE FROM carts_variants
WHERE cart_id = sqlc.arg(cart_id)
  AND product_variant_id = sqlc.arg(product_variant_id);

-- name: ClearCartItems :exec
DELETE FROM carts_variants
WHERE cart_id = sqlc.arg(cart_id);

-- name: GetCartItemsCount :one
SELECT COUNT(*)
FROM carts_variants
WHERE cart_id = sqlc.arg(cart_id);

-- name: GetTotalCartPrice :one
SELECT SUM(quantity * price_per_item) AS total
FROM carts_variants 
WHERE cart_id = sqlc.arg(cart_id);

-- name: GetCartDetailsByCartId :many
SELECT
  cv.cart_id,
  cv.product_variant_id,
  cv.quantity,
  cv.price_per_item,
  v.sku,
  v.price AS variant_price,
  v.stock_quantity,
  v.image_url AS variant_image,
  v.variant_name,
  p.id AS product_id,
  p.name AS product_name,
  p.slug AS product_slug,
  p.image_url AS product_image
FROM carts_variants cv
JOIN product_variants v ON cv.product_variant_id = v.id
JOIN products p ON v.product_id = p.id
WHERE cv.cart_id = sqlc.arg(cart_id);

-- name: GetCartDetailsWithSnapshotPrice :many
SELECT
  cv.cart_id,
  cv.product_variant_id,
  cv.quantity,
  cv.price_per_item,
  v.sku,
  v.price AS variant_price,
  v.stock_quantity,
  v.image_url AS variant_image,
  v.variant_name,
  p.id AS product_id,
  p.name AS product_name,
  p.slug AS product_slug,
  p.image_url AS product_image
FROM carts_variants cv
JOIN product_variants v ON cv.product_variant_id = v.id
JOIN products p ON v.product_id = p.id
WHERE cv.cart_id = sqlc.arg(cart_id);

-- name: GetCartDetailsWithLivePrice :many
SELECT
  cv.cart_id,
  cv.product_variant_id,
  cv.quantity,
  v.price AS current_price,
  v.sku,
  v.stock_quantity,
  v.image_url AS variant_image,
  v.variant_name,
  p.id AS product_id,
  p.name AS product_name,
  p.slug AS product_slug,
  p.image_url AS product_image
FROM carts_variants cv
JOIN product_variants v ON cv.product_variant_id = v.id
JOIN products p ON v.product_id = p.id
WHERE cv.cart_id = sqlc.arg(cart_id);

-- name: UpsertVariantToCart :one
INSERT INTO carts_variants (cart_id, product_variant_id, quantity, price_per_item)
VALUES (
  sqlc.arg(cart_id),
  sqlc.arg(product_variant_id),
  sqlc.arg(quantity),
  sqlc.arg(price_per_item)
)
ON CONFLICT (cart_id, product_variant_id) DO UPDATE
SET 
  quantity = carts_variants.quantity + EXCLUDED.quantity,
  price_per_item = EXCLUDED.price_per_item,
  updated_at = CURRENT_TIMESTAMP
RETURNING *;
