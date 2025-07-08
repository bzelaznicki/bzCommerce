-- name: CreateOrder :one
INSERT INTO orders (
    user_id, 
    total_price, 
    customer_email, 
    shipping_name, 
    shipping_address, 
    shipping_city, 
    shipping_postal_code, 
    shipping_country_id, 
    shipping_phone, 
    billing_name, 
    billing_address, 
    billing_city, 
    billing_postal_code, 
    billing_country_id,
    shipping_option_id,
    shipping_price,
    payment_option_id
)
VALUES (
    sqlc.arg(user_id), 
    sqlc.arg(total_price), 
    sqlc.arg(customer_email), 
    sqlc.arg(shipping_name), 
    sqlc.arg(shipping_address), 
    sqlc.arg(shipping_city), 
    sqlc.arg(shipping_postal_code), 
    sqlc.arg(shipping_country_id), 
    sqlc.arg(shipping_phone), 
    sqlc.arg(billing_name), 
    sqlc.arg(billing_address), 
    sqlc.arg(billing_city), 
    sqlc.arg(billing_postal_code), 
    sqlc.arg(billing_country_id),
    sqlc.arg(shipping_option_id),
    sqlc.arg(shipping_price),
    sqlc.arg(payment_option_id)
)
RETURNING *;

-- name: GetOrderById :one
SELECT * FROM orders
WHERE id = sqlc.arg(id);

-- name: GetOrders :many
SELECT * FROM orders
ORDER BY created_at DESC;

-- name: GetOrdersByStatus :many
SELECT * FROM orders
WHERE status IN (sqlc.arg(status))
ORDER BY created_at DESC;

-- name: GetOrdersByOwnerUserId :many
SELECT * FROM orders
WHERE user_id = sqlc.arg(user_id)
ORDER BY created_at DESC;

-- name: CopyCartDataIntoOrder :many
INSERT INTO orders_variants (order_id, product_variant_id, quantity, price_per_item, total_price)
SELECT 
  sqlc.arg(order_id), 
  cv.product_variant_id, 
  cv.quantity, 
  cv.price_per_item,
  (cv.quantity * cv.price_per_item)
FROM carts_variants cv
WHERE cv.cart_id = sqlc.arg(cart_id)
RETURNING *;


-- name: GetOrderItemsByOrderId :many
SELECT * FROM orders_variants
WHERE order_id = sqlc.arg(order_id);

-- name: UpdateOrderStatus :one
UPDATE orders
SET status = sqlc.arg(status)
WHERE id = sqlc.arg(id)
RETURNING *;