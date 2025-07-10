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
-- name: ListOrders :many
SELECT
  o.id,
  o.user_id,
  o.status,
  o.payment_status,
  o.total_price,
  o.created_at,
  o.updated_at,
  o.customer_email,
  o.shipping_name,
  o.shipping_address,
  o.shipping_city,
  o.shipping_postal_code,
  o.shipping_phone,
  o.billing_name,
  o.billing_address,
  o.billing_city,
  o.billing_postal_code,
  o.shipping_option_id,
  o.shipping_price,
  o.payment_option_id,
  o.shipping_country_id,
  o.billing_country_id,
  s.name AS shipping_method_name,
  p.name AS payment_method_name
FROM
  orders o
LEFT JOIN shipping_options s ON s.id = o.shipping_option_id
LEFT JOIN payment_options p ON p.id = o.payment_option_id
WHERE
  (
    sqlc.arg(search)::text IS NULL
    OR o.customer_email ILIKE '%' || sqlc.arg(search) || '%'
    OR o.shipping_name ILIKE '%' || sqlc.arg(search) || '%'
    OR o.billing_name ILIKE '%' || sqlc.arg(search) || '%'
  )
  AND (
    sqlc.arg(status)::order_status IS NULL
    OR o.status = sqlc.arg(status)
  )
  AND (
    sqlc.arg(payment_status)::payment_status IS NULL
    OR o.payment_status = sqlc.arg(payment_status)
  )
  AND (
    sqlc.arg(date_from)::timestamp IS NULL
    OR o.created_at >= sqlc.arg(date_from)
  )
  AND (
    sqlc.arg(date_to)::timestamp IS NULL
    OR o.created_at <= sqlc.arg(date_to)
  )
ORDER BY o.created_at DESC
LIMIT $1
OFFSET $2;


-- name: CountOrders :one
SELECT COUNT(*)
FROM orders o
WHERE
  (
    sqlc.arg(search)::text IS NULL
    OR o.customer_email ILIKE '%' || sqlc.arg(search) || '%'
    OR o.shipping_name ILIKE '%' || sqlc.arg(search) || '%'
    OR o.billing_name ILIKE '%' || sqlc.arg(search) || '%'
  )
  AND (
    sqlc.arg(status)::order_status IS NULL
    OR o.status = sqlc.arg(status)
  )
  AND (
    sqlc.arg(payment_status)::payment_status IS NULL
    OR o.payment_status = sqlc.arg(payment_status)
  )
  AND (
    sqlc.arg(date_from)::timestamp IS NULL
    OR o.created_at >= sqlc.arg(date_from)
  )
  AND (
    sqlc.arg(date_to)::timestamp IS NULL
    OR o.created_at <= sqlc.arg(date_to)
  );

